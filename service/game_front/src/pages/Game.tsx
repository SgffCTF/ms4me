import { useNavigate, useParams } from "react-router"
import { CreatorGame } from "./CreatorGame"
import { useEffect, useRef, useState } from "react";
import { enterGame, getGameByID } from "../api/games";

import { GameDetails, Message } from "../models/models";
import { useAuth } from "../context/AuthProvider";
import { ParticipantGame } from "./ParticipantGame";
import { ClickGameEvent, DeleteRoomEvent, DeleteRoomEventType, ExitRoomEvent, ExitRoomEventType, JoinRoomEvent, JoinRoomEventType, LoseGameEvent, LoseGameEventType, NewMessageEventType, OpenCellEventType, RoomParticipant, StartGameEventType, UpdateRoomEvent, UpdateRoomEventType, WinGameEvent, WinGameEventType, WSEvent } from "../models/events";
import { toast } from "react-toastify";
import { gameContainsUserID, getCookie } from "../utils/utils";
import { WS_URI } from "../api/api";
import { getGameInfo, getMessages } from "../api/ingame";

export const GameDetail = () => {
    const { id } = useParams<{ id: string }>();
    const { user } = useAuth();
    const [game, setGame] = useState<GameDetails>();
    const [messages, setMessages] = useState<Message[]>();
    const navigate = useNavigate();
    const wsRef = useRef<WebSocket | null>(null);
    const reconnectRef = useRef<number | null>(null);
    const isActiveRef = useRef(true);
    const [isStart, setIsStart] = useState(false);
    const [roomParticipants, setRoomParticipants] = useState<Array<RoomParticipant> | null>(null);

    const eventHandler = (event: WSEvent) => {
        if (!event.payload) return;
        let eventData: any;
        switch (event.event_type) {
        case JoinRoomEventType:
            eventData = event.payload as JoinRoomEvent;
            if (eventData.user_id != user?.id) {
                toast("К игре присоединился " + eventData.username);
                setGame((prev) => {
                    if (!prev) return prev;

                    // Проверим, есть ли уже такой игрок
                    const alreadyExists = prev.players.some(p => p.id === eventData.user_id);
                    if (alreadyExists) return prev;

                    // Добавим нового игрока
                    return {
                        ...prev,
                        players: [...prev.players, {
                            id: eventData.user_id,
                            username: eventData.username
                        }]
                    };
                });
                setRoomParticipants((prev) => {
                    if (!prev) return prev;

                    // Проверим, есть ли уже такой игрок
                    const alreadyExists = prev.some(p => p.id === eventData.user_id);
                    if (alreadyExists) return prev;

                    // Добавим нового игрока
                    return [...prev, {id: eventData.user_id, username: eventData.username, is_owner: false, field: null} as RoomParticipant]
                })
            }
            break;
        case ExitRoomEventType:
            eventData = event.payload as ExitRoomEvent;
            if (eventData.user_id != user?.id) {
                toast(eventData.username + " вышел из игры");
                setGame((prev) => {
                    if (!prev) return prev;

                    const players = prev.players.filter(p => p.id != eventData.user_id);
                    return {
                        ...prev,
                        players: players,
                    };
                });
            }
            break;
        case UpdateRoomEventType: 
            eventData = event.payload as UpdateRoomEvent;
            setGame((prev) => {
                if (!prev) return prev;
                return {
                ...prev,
                title: eventData.title,
                is_public: eventData.is_public ?? prev.is_public,
                };
            });
            break;
        case DeleteRoomEventType:
            eventData = event.payload as DeleteRoomEvent;
            if (eventData.id === id) {
                toast("Игра удалена");
                navigate("/");
            }
            break;
        case StartGameEventType:
            setIsStart(true);
            toast("Игра началась!", {autoClose: 5000});
            break;
        case OpenCellEventType:
            eventData = event.payload as ClickGameEvent;
            setRoomParticipants(eventData.participants);
            break;
        case LoseGameEventType:
            eventData = event.payload as LoseGameEvent;
            if (eventData.loser_id != user?.id) {
                toast.success("🎉 Победа! Вы выиграли игру!", {
                    position: "top-center",
                    autoClose: 1000,
                    hideProgressBar: false,
                    closeOnClick: true,
                    pauseOnHover: true,
                    draggable: false,
                    progress: undefined,
                    theme: "colored",
                    onClose() {
                        navigate("/");
                    },
                });
            } else {
                toast.warn(`🙁 Поражение! Игру выиграл ${roomParticipants?.find((v) => {
                    if (v.id != eventData.loser_id) {
                        return v;
                    }
                })?.username}!`, {
                    position: "top-center",
                    autoClose: 1000,
                    hideProgressBar: false,
                    closeOnClick: true,
                    pauseOnHover: true,
                    draggable: false,
                    progress: undefined,
                    theme: "colored",
                    onClose() {
                        navigate("/");
                    },
                });
            }
            break;
        case WinGameEventType:
            eventData = event.payload as WinGameEvent;
            if (eventData.winner_id == user?.id) {
                toast.success("🎉 Победа! Вы выиграли игру!", {
                    position: "top-center",
                    autoClose: 1000,
                    hideProgressBar: false,
                    closeOnClick: true,
                    pauseOnHover: true,
                    draggable: false,
                    progress: undefined,
                    theme: "colored",
                    onClose() {
                        navigate("/");
                    },
                });
            } else {
                toast.warn(`🙁 Поражение! Игру выиграл ${eventData.winner_username}!`, {
                    position: "top-center",
                    autoClose: 1000,
                    hideProgressBar: false,
                    closeOnClick: true,
                    pauseOnHover: true,
                    draggable: false,
                    progress: undefined,
                    theme: "colored",
                    onClose() {
                        navigate("/");
                    },
                });
            }
            break;
        case NewMessageEventType:
            eventData = event.payload as Message;
            setMessages(prevMessages => {
                if (!prevMessages) return [eventData];
                    return [...prevMessages, eventData];
            });
    break;
        default:
            console.error("Неизвестный event_type: " + event.event_type);
            break;
        }
    };

    const disconnect = () => {
            if (wsRef.current) {
                wsRef.current.close();
                wsRef.current = null;
            }
            if (reconnectRef.current) {
                clearTimeout(reconnectRef.current);
                reconnectRef.current = null;
            }
        }
    
    const connectWS = () => {
        const token = getCookie("token");
        if (!user || !token || !isActiveRef.current) return;

        disconnect();

        const ws = new WebSocket(`${WS_URI}/${id}`);
        wsRef.current = ws;

        ws.onopen = () => {
            console.log("WebSocket открыт");
            ws.send(JSON.stringify({ token }));
        };

        ws.onmessage = (ev) => {
            if (ev.data == "ping") return;
            try {
                const data: WSEvent = JSON.parse(ev.data);
                eventHandler(data);
            } catch (err) {
                console.error("Ошибка парсинга:", err);
            }
        };

        ws.onclose = () => {
            console.warn("WebSocket закрыт");
            if (isActiveRef.current && reconnectRef.current === null) {
                reconnectRef.current = window.setTimeout(() => {
                    reconnectRef.current = null;
                    connectWS(); // не переподключаем, если компонент размонтирован
                }, 2000);
            }
        };

        ws.onerror = (err) => {
            console.error("WebSocket ошибка:", err);
            ws.close();
        };
    };

    
    useEffect(() => {
        if (user === null) return;
        isActiveRef.current = true;

        const load = async () => {
            if (!id || !user) return;

            try {
                const game = await getGameByID(id);
                if (!gameContainsUserID(game, user.id)) {
                    throw Error("Пользователь отсутствует в данной игре");
                }
                setGame(game);
            } catch (e: any) {
                try {
                    await enterGame(id);
                    setGame(await getGameByID(id));
                } catch (e: any) {
                    toast.error(e.message);
                    navigate("/");
                    return;
                }
            }

            try {
                setRoomParticipants(await getGameInfo(id));
            } catch (e: any) {
                toast.error(e.message);
            }

            try {
                setMessages(await getMessages(id));
            } catch (e: any) {
                toast.error(e.message);
            }

            connectWS();
        };

        load();

        return () => {
            isActiveRef.current = false;
            disconnect();
        };
    }, []);


    return (
        <>
        {
            (id && game && user && messages && (game.owner_id === user.id &&
            <CreatorGame // Интерфейс владельца игры
            id={id}
            gameInfo={game}
            wsRef={wsRef}
            isStart={isStart}
            roomParticipants={roomParticipants}
            messages={messages}
            ></CreatorGame>
            ||
            <ParticipantGame // Интерфейс участника игры
            isStart={isStart}
            id={id}
            gameInfo={game}
            wsRef={wsRef}
            roomParticipants={roomParticipants}
            messages={messages}
            ></ParticipantGame>))
        }
        </>
    )
}
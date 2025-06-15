import { useNavigate } from "react-router";
import { Game, Message } from "../../models/models";
import { useEffect, useRef, useState } from "react";
import { getGames, getMyGames } from "../../api/games";
import { toast } from "react-toastify";
import { CreateRoomEventType, DeleteRoomEventType, ExitRoomEvent, ExitRoomEventType, JoinRoomEvent, JoinRoomEventType, StartGameEventType, WSEvent } from "../../models/events";
import { useAuth } from "../../context/AuthProvider";
import { getCookie } from "../../utils/utils";
import { WS_URI } from "../../api/api";
import { getMessages } from "../../api/ingame";
import { ChatModal } from "./ChatModal";

interface Props {
    searchQuery: string;
    showMyGames: boolean;
}

export const GameList = (props: Props) => {
    const navigate = useNavigate();
    const { user } = useAuth();
    const [games, setGames] = useState<Array<Game>>([]);
    const [newGameIds, setNewGameIds] = useState<Record<string, boolean>>({});
    const [isLoading, setIsLoading] = useState(true);
    const wsRef = useRef<WebSocket | null>(null);
    const reconnectRef = useRef<number | null>(null);
    const isActiveRef = useRef(true);
    const [messages, setMessages] = useState<Message[]>([]);
    const [clickedID, setClickedID] = useState<string | null>(null);
    const [showModal, setShowModal] = useState<boolean>(false);

    useEffect(() => {
        const load = async () => {
            setIsLoading(true);
            try {
                let games: Array<Game>;
                if (props.showMyGames) {
                    games = await getMyGames();
                } else {
                    games = await getGames(props.searchQuery, "open");
                }
                setGames(games);
            } catch (e: any) {
                toast.error(e.message);
            }
            setIsLoading(false);
        }

        if (user === null) return;
        load();
    }, [props.searchQuery, props.showMyGames]);

    const eventHandler = (event: WSEvent) => {
        if (!event.payload) return;
        switch (event.event_type) {
            case CreateRoomEventType:
                const game = event.payload as Game;
                setGames((prev) => [game, ...prev]);
                setNewGameIds((prev) => ({ ...prev, [game.id]: true }));

                setTimeout(() => {
                    setNewGameIds((prev) => {
                        const copy = { ...prev };
                        delete copy[game.id];
                        return copy;
                    });
                }, 5000);
                break;
            case StartGameEventType:
            case DeleteRoomEventType:
                if (!(event.payload && event.payload.id)) return;
                const gameID = event.payload.id;
                setGames((prev) => prev.filter((game) => game.id !== gameID));
                break;
            case JoinRoomEventType:
                var data = event.payload as JoinRoomEvent;
                setGames((prev) => prev.map((game) => {
                    if (game.id === data.id) game.players_count++;
                    return game;
                }));
                break;
            case ExitRoomEventType:
                var data = event.payload as ExitRoomEvent;
                setGames((prev) => prev.map((game) => {
                    if (game.id === data.id) game.players_count--;
                    return game;
                }));
                break;
            default:
                console.error("Неизвестный event_type: " + event.event_type);
                break;
        }
    }

    const disconnect = () => {
        if (wsRef.current) {
            if (wsRef.current.readyState === WebSocket.OPEN) {
                wsRef.current.close(1000, "User navigated away");
            }
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

        const ws = new WebSocket(WS_URI);
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
                }, 5000);
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

        connectWS();

        return () => {
            isActiveRef.current = false;
            disconnect();
        };
    }, []);

    const handleClick = async (game: Game) => {
        if (game.status == "closed") {
            setMessages(await getMessages(game.id));
            setClickedID(game.id);
            setShowModal(true);
            return;
        }

        disconnect();
        navigate("/game/" + game.id);
    }

    if (isLoading) {
        return <div className="text-center mt-5">Загрузка...</div>;
    }

    return (
        <>
            {games?.map((game) => (
                <div
                    key={game.id}
                    className="border rounded hover p-3 m-1 relative"
                    role="button"
                    onClick={() => handleClick(game)}
                >
                    <div className="container">
                        <div className="row">
                            <div className="col">
                                <p>Создатель: {game.owner_name}</p>
                                <p>Название: {game.title}</p>
                                {!props.showMyGames &&
                                    <p>{game.players_count}/{game.max_players}</p>
                                }
                            </div>
                            <div className="col d-flex justify-content-end">
                                {user && ((game.status == "started" && (
                                    <span className="text-game-in-progress">
                                        В процессе
                                    </span>
                                )) || (game.status == "closed" && game.winner_id == user.id && (
                                    <span className="text-game-win">
                                        Победа!
                                    </span>
                                )) || (game.status == "closed" && game.winner_id != user.id && (
                                    <span className="text-game-lose">
                                        Поражение!
                                    </span>
                                )) || (newGameIds[game.id] && (
                                    <span className="text-game-new">
                                        Новая
                                    </span>
                                ))
                                )}
                            </div>
                        </div>
                    </div>
                </div>
            ))}
            <ChatModal id={clickedID} messages={messages} show={showModal} setShow={setShowModal}></ChatModal>
        </>
    );
}

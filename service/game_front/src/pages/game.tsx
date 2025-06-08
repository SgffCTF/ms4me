import { useNavigate, useParams } from "react-router"
import { CreatorGame } from "./CreatorGame"
import { useEffect, useState } from "react";
import { enterGame, getGameByID } from "../api/games";

import { GameDetails } from "../models/models";
import { useAuth } from "../context/AuthProvider";
import { ParticipantGame } from "./ParticipantGame";
import { DeleteRoomEvent, DeleteRoomEventType, JoinRoomEvent, JoinRoomEventType, UpdateRoomEvent, UpdateRoomEventType, WSEvent } from "../models/events";
import { toast } from "react-toastify";
import { useWS } from "../context/GameWSProvider";

export const GameDetail = () => {
    const { id } = useParams<{ id: string }>();
    const { user } = useAuth();
    const [game, setGame] = useState<GameDetails>();
    const navigate = useNavigate();
    const { addMessageListener, wsRef } = useWS();

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
        default:
            console.error("Неизвестный event_type: " + event.event_type);
            break;
        }
    };
    
    useEffect(() => {
        const load = async () => {
            if (!id) return;

            try {
                setGame(await getGameByID(id));
                addMessageListener(eventHandler);
                return;
            } catch (e: any) {
                try {
                    await enterGame(id);
                    setGame(await getGameByID(id));
                    addMessageListener(eventHandler);
                } catch (e: any) {
                    toast.error(e.message);
                    navigate("/");
                    wsRef.current?.close();
                }
            }
        }

        load();
    }, []);

    return (
        <>
        {
            (id && game && user && (game.owner_id === user.id && <CreatorGame id={id} gameInfo={game}></CreatorGame> || <ParticipantGame id={id} gameInfo={game}></ParticipantGame>))
        }
        </>
    )
}
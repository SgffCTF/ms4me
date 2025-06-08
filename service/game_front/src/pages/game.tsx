import { useNavigate, useParams } from "react-router"
import { CreatorGame } from "./CreatorGame"
import { useEffect, useState } from "react";
import { getGameByID } from "../api/games";

import { GameDetails } from "../models/models";
import { useAuth } from "../context/AuthProvider";
import { ParticipantGame } from "./ParticipantGame";
import { DeleteRoomEvent, DeleteRoomEventType, EnterRoomEventType, UpdateRoomEvent, UpdateRoomEventType, WSEvent } from "../models/events";
import { toast } from "react-toastify";
import { useWS } from "../context/GameWSProvider";

export const GameDetail = () => {
    const { id } = useParams<{ id: string }>();
    const { user } = useAuth();
    const [game, setGame] = useState<GameDetails>();
    const navigate = useNavigate();
    const { addMessageListener } = useWS();

    useEffect(() => {
        if (id) getGameByID(id).then(game => setGame(game));
    }, []);

    const eventHandler = (event: WSEvent) => {
        if (!event.payload) return;
        let eventData: any;
        switch (event.event_type) {
        case EnterRoomEventType:
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
    
    useEffect(() => addMessageListener(eventHandler), []);

    return (
        <>
        {
            (id && game && user && game.owner_id === user.id &&
            <CreatorGame id={id} gameInfo={game}></CreatorGame>) || <ParticipantGame></ParticipantGame>
        }
        </>
    )
}
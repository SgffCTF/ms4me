import { useParams } from "react-router"
import { CreatorGame } from "./CreatorGame"
import { useEffect, useState } from "react";
import { getGameByID } from "../api/games";

import { Game } from "../models/models";
import { useAuth } from "../context/AuthProvider";
import { ParticipantGame } from "./ParticipantGame";

export const GameDetail = () => {
    const { id } = useParams<{ id: string }>();
    const { user } = useAuth();
    const [game, setGame] = useState<Game | null>(null);

    useEffect(() => {
        if (id) getGameByID(id).then(game => setGame(game));
    }, []);

    return (
        <>
        {
            (id && game && user && game.owner_id === user?.id && <CreatorGame id={id}></CreatorGame>) || <ParticipantGame></ParticipantGame>
        }
        </>
    )
}
import { useNavigate } from "react-router";
import { Game } from "../models/models";
import { useWS } from "../context/WebsocketProvider";
import { useEffect, useState } from "react";
import { getGames } from "../api/games";
import { toast } from "react-toastify";
import { CreateRoomEventType, DeleteRoomEventType, WSEvent } from "../models/events";
import { useAuth } from "../context/AuthProvider";

interface Props {
    searchQuery: string;
}

export const GameList = (props: Props) => {
    const navigate = useNavigate();
    const { user } = useAuth();
    const { addMessageListener } = useWS();
    const [games, setGames] = useState<Array<Game>>([]);
    const [newGameIds, setNewGameIds] = useState<Record<string, boolean>>({});

    useEffect(() => {
        const load = async () => {
            try {
                setGames(await getGames(props.searchQuery));
            } catch (e: any) {
                toast.error(e.message);
            }
        }

        load();
    }, [props.searchQuery]);

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
            case DeleteRoomEventType:
                if (!(event.payload && event.payload.id)) return;
                const gameID = event.payload.id;
                setGames((prev) => prev.filter((game) => game.id !== gameID));
                break;
            default:
                console.error("Неизвестный event_type: " + event.event_type);
                break;
        }
    }

    useEffect(() => addMessageListener(eventHandler), []);

    return (
        <>
            {games?.map((game) => (
                <div
                    key={game.id}
                    className="border rounded hover p-3 m-1 relative"
                    role="button"
                    onClick={() => navigate("/game/" + game.id)}
                >
                    <div className="container">
                        <div className="row">
                            <div className="col">
                                <p>Создатель: {game.owner_name}</p>
                                <p>Название: {game.title}</p>
                                <p>{game.players}/{game.max_players}</p>
                            </div>
                            <div className="col d-flex justify-content-end">
                                {user && user.id == game.owner_id && (
                                    <span className="text-yellow">
                                        Твоя игра
                                    </span>
                                ) || newGameIds[game.id] && (
                                    <span className="text-important">
                                        Новая
                                    </span>
                                )
                                }
                            </div>
                        </div>
                    </div>
                </div>
            ))}
        </>
    );
}

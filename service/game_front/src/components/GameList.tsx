import { useNavigate } from "react-router";
import { Game } from "../models/models";
import { useWS } from "../context/WebsocketProvider";
import { useEffect, useState } from "react";
import { getGames } from "../api/games";
import { toast } from "react-toastify";
import { CreateRoomEventType, WSEvent } from "../models/events";

interface Props {
    searchQuery: string;
}

export const GameList = (props: Props) => {
    const navigate = useNavigate();
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

    const newGameHandler = (event: WSEvent) => {
        if (!event.payload) return;
        switch (event.event_type) {
            case CreateRoomEventType:
                const game = event.payload as Game;
                setGames((prev) => [game, ...prev]);
                setNewGameIds((prev) => ({ ...prev, [game.id]: true }));

                // Удаляем метку "новое" через 5 секунд
                setTimeout(() => {
                    setNewGameIds((prev) => {
                        const copy = { ...prev };
                        delete copy[game.id];
                        return copy;
                    });
                }, 5000);
                break;
            default:
                console.error("Неизвестный event_type: " + event.event_type);
                break;
        }
    }

    useEffect(() => addMessageListener(newGameHandler), []);

    return (
        <>
            {games?.map((game) => (
                <div
                    key={game.id}
                    className="border rounded hover p-3 m-1 relative"
                    role="button"
                    onClick={() => navigate("/game/" + game.id)}
                >
                    {newGameIds[game.id] && (
                        <span className="absolute top-0 right-0 bg-green-500 text-important text-xs px-2 py-1 rounded-bl">
                            Новое
                        </span>
                    )}
                    <p>Создатель: {game.owner_name}</p>
                    <p>Название: {game.title}</p>
                    <p>{game.players}/{game.max_players}</p>
                </div>
            ))}
        </>
    );
}

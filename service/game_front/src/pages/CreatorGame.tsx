import { useEffect, useState } from "react";
import { Field } from "../components/Field";
import { Chat } from "../components/Chat";
import { useWS } from "../context/WebsocketProvider";
import { EnterRoomEventType, UpdateRoomEvent, UpdateRoomEventType, WSEvent } from "../models/events";
import { GameDetails } from "../models/models";
import { getGameByID } from "../api/games";
import { formatDate } from "../utils/utils";
import { UpdateGameModal } from "../components/UpdateGameModal";

interface Props {
    id: string;
}

export const CreatorGame = (props: Props) => {
    const [updateModalShow, setUpdateModalShow] = useState(false);
    const [infoOpen, setInfoOpen] = useState(false);
    const { addMessageListener } = useWS();
    const [gameInfo, setGameInfo] = useState<GameDetails>();

    const eventHandler = (event: WSEvent) => {
    if (!event.payload) return;
        switch (event.event_type) {
        case EnterRoomEventType:
            break;
        case UpdateRoomEventType: 
            const eventData = event.payload as UpdateRoomEvent;
            setGameInfo((prev) => {
                if (!prev) return prev;
                return {
                ...prev,
                title: eventData.title,
                is_public: eventData.is_public ?? prev.is_public,
                };
            });
            break;
        default:
            console.error("Неизвестный event_type: " + event.event_type);
            break;
        }
    };

    useEffect(() => addMessageListener(eventHandler), []);
    useEffect(() => {
        getGameByID(props.id).then((info) => setGameInfo(info));
    }, []);

    return (
        <div className="container-fluid d-flex flex-column min-vh-100">
            <h1 className="text-center mt-4 mb-3">Игра</h1>

            { gameInfo &&
            <div className="mb-3">
                <button
                    className="btn btn-primary"
                    onClick={() => setInfoOpen(!infoOpen)}
                    aria-expanded={infoOpen}
                    aria-controls="gameInfoCollapse"
                >
                    {infoOpen ? "Скрыть информацию об игре" : "Показать информацию об игре"}
                </button>
                <button className="btn btn-primary ms-2" onClick={() => setUpdateModalShow(true)}>✏️</button>

                <div
                    className={`collapse mt-2 ${infoOpen ? "show" : ""}`}
                    id="gameInfoCollapse"
                >
                    <div className="card card-body">
                    <h5>Название: {gameInfo.title}</h5>
                    <p>Дата создания: {formatDate(gameInfo.created_at)}</p>
                    <strong>Участники:</strong>
                    <ul>
                        {gameInfo.players.map((p) => (
                        <li key={p.id}>{p.username}</li>
                        ))}
                    </ul>
                    </div>
                </div>
            </div>
            }

            <div className="row flex-grow-1">
            <div className="col-4">
                <Field />
            </div>
            <div className="col-4">
                <Field />
            </div>
            <div className="col-4">
                <Chat />
            </div>
            </div>
            {gameInfo && <UpdateGameModal id={props.id} show={updateModalShow} setShow={setUpdateModalShow} gameInfo={gameInfo}></UpdateGameModal>}
        </div>
    );
};

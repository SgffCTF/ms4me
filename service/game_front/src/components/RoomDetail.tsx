import { useState } from "react";
import { GameDetails } from "../models/models";
import { formatDate } from "../utils/utils"

interface Props {
    gameInfo: GameDetails;
}

export const RoomDetail = (props: Props) => {
    const [infoOpen, setInfoOpen] = useState(false);
    return (
        <>
        <button
                    className="btn btn-primary"
                    onClick={() => setInfoOpen(!infoOpen)}
                    aria-expanded={infoOpen}
                    aria-controls="gameInfoCollapse"
                >
                    {infoOpen ? "Скрыть информацию об игре" : "Показать информацию об игре"}
                </button>
        <div
            className={`collapse mt-2 ${infoOpen ? "show" : ""}`}
            id="gameInfoCollapse"
        >
            <div className="card card-body">
            <h5>Название: {props.gameInfo.title}</h5>
            <p>Дата создания: {formatDate(props.gameInfo.created_at)}</p>
            <strong>Участники:</strong>
            <ul>
                {props.gameInfo.players.map((p) => (
                <li key={p.id}>{p.username}</li>
                ))}
            </ul>
            </div>
        </div>
        </>
    )
}
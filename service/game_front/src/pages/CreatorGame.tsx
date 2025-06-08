import { useState } from "react";
import { Field } from "../components/Field";
import { Chat } from "../components/Chat";
import { GameDetails } from "../models/models";
import { deleteGame } from "../api/games";
import { UpdateGameModal } from "../components/UpdateGameModal";
import { toast } from "react-toastify";
import { RoomDetail } from "../components/RoomDetail";
import { useNavigate } from "react-router";
import { useWS } from "../context/GameWSProvider";

interface Props {
    id: string;
    gameInfo: GameDetails;
}

export const CreatorGame = (props: Props) => {
    const [updateModalShow, setUpdateModalShow] = useState(false);
    const navigate = useNavigate();
    const { wsRef } = useWS();

    const deleteGameHandler = async () => {
        try {
            await deleteGame(props.id);
            toast("Игра удалена");
            navigate("/");
        } catch (e: any) {
            toast.error(e.message);
        }
    }

    return (
        <div className="container-fluid d-flex flex-column min-vh-100">
            <div className="d-flex justify-content-between align-items-center mt-4 mb-3">
                <h1 className="mb-0">Игра</h1>
                <button className="btn btn-outline-secondary" onClick={() => {
                    wsRef.current?.close();
                    navigate("/")}
                }>
                    ← Назад
                </button>
            </div>

            { props.gameInfo &&
            <div className="mb-3">
                <button className="btn btn-primary ms-2" onClick={() => setUpdateModalShow(true)}>✏️</button>
                <button className="btn btn-red ms-2 me-2" onClick={deleteGameHandler}>❌</button>
                <RoomDetail gameInfo={props.gameInfo}></RoomDetail>
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
            <UpdateGameModal id={props.id} show={updateModalShow} setShow={setUpdateModalShow} gameInfo={props.gameInfo}></UpdateGameModal>
        </div>
    );
};

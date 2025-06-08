import { useState } from "react";
import { Field } from "../components/Field";
import { Chat } from "../components/Chat";
import { GameDetails } from "../models/models";
import { deleteGame } from "../api/games";
import { UpdateGameModal } from "../components/UpdateGameModal";
import { toast } from "react-toastify";
import { RoomDetail } from "../components/RoomDetail";

interface Props {
    id: string;
    gameInfo: GameDetails;
}

export const CreatorGame = (props: Props) => {
    const [updateModalShow, setUpdateModalShow] = useState(false);

    const deleteGameHandler = async () => {
        try {
            await deleteGame(props.id);
        } catch (e: any) {
            toast.error(e.message);
        }
    }

    return (
        <div className="container-fluid d-flex flex-column min-vh-100">
            <h1 className="text-center mt-4 mb-3">Игра</h1>

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

import { useState } from "react";
import { Field } from "../components/Field/Field";
import { Chat } from "../components/Chat";
import { GameDetails, Message } from "../models/models";
import { deleteGame, startGame } from "../api/games";
import { UpdateGameModal } from "../components/UpdateGameModal";
import { toast } from "react-toastify";
import { RoomDetail } from "../components/RoomDetail";
import { useNavigate } from "react-router";
import { RoomParticipant } from "../models/events";
import { useAuth } from "../context/AuthProvider";

interface Props {
    id: string;
    gameInfo: GameDetails;
    wsRef: React.RefObject<WebSocket | null>;
    isStart: boolean;
    roomParticipants: Array<RoomParticipant> | null;
    messages: Message[];
}

export const CreatorGame = (props: Props) => {
    const [updateModalShow, setUpdateModalShow] = useState(false);
    const navigate = useNavigate();
    const { user } = useAuth();

    const deleteGameHandler = async () => {
        if (props.wsRef) {
            try {
                await deleteGame(props.id);
            } catch (e: any) {
                toast.error(e.message);
            }
        }
    }

    const startGameHandler = async () => {
        if (props.wsRef) {
            try {
                await startGame(props.id);
            } catch (e: any) {
                toast.error(e.message);
            }
        }
    }

    return (
        <div className="container-fluid d-flex flex-column min-vh-100">
            <div className="d-flex justify-content-between align-items-center mt-4 mb-3">
                <h1 className="mb-0">Игра</h1>
                <button className="btn btn-outline-secondary" onClick={() => {
                    props.wsRef.current?.close();
                    navigate("/")}
                }>
                    ← Назад
                </button>
            </div>

            { props.gameInfo &&
            <div className="mb-3">
                <button className="btn btn-primary me-2" onClick={startGameHandler}>▶️</button>
                <button className="btn btn-orange me-2" onClick={() => setUpdateModalShow(true)}>✏️</button>
                <button className="btn btn-red me-2" onClick={deleteGameHandler}>❌</button>

                <RoomDetail gameInfo={props.gameInfo}></RoomDetail>
            </div>
            }

            <div className="row flex-grow-1">
                <div className="col-4">
                    <Field roomParticipants={props.roomParticipants} gameID={props.gameInfo.id} fieldOwnerID={user ? user.id : null}/>
                </div>
                <div className="col-4">
                    {
                        (props.gameInfo.players.length > 1 &&
                        <Field roomParticipants={props.roomParticipants} gameID={props.gameInfo.id} fieldOwnerID={props.gameInfo.players[1].id}/>) ||
                        <Field roomParticipants={props.roomParticipants} gameID={props.gameInfo.id} fieldOwnerID={null}/>
                    }
                </div>
                <div className="col-4">
                    <Chat messages={props.messages} id={props.id} withInput={true}/>
                </div>
            </div>
            <UpdateGameModal id={props.id} show={updateModalShow} setShow={setUpdateModalShow} gameInfo={props.gameInfo}></UpdateGameModal>
        </div>
    );
};

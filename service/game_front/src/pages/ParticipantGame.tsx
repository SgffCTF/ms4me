import { Field } from "../components/Field";
import { Chat } from "../components/Chat";
import { GameDetails } from "../models/models";
import { RoomDetail } from "../components/RoomDetail";
import { exitGame } from "../api/games";
import { toast } from "react-toastify";
import { useNavigate } from "react-router";

interface Props {
    id: string;
    gameInfo: GameDetails;
}

export const ParticipantGame = (props: Props) => {
    const navigate = useNavigate();

    const exitHandler = async () => {
        try {
            await exitGame(props.id);
            toast("Вы вышли из игры");
            navigate("/");
        } catch (e: any) {
            toast.error(e.message);
        }
    }

    return (
        <div className="container-fluid d-flex flex-column min-vh-100">
            <h1 className="text-center mt-4 mb-3">Игра</h1>

            { props.gameInfo &&
            <div className="mb-3">
                <button className="btn btn-red ms-2 me-2" onClick={exitHandler}>Выйти</button>
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
        </div>
    );
};

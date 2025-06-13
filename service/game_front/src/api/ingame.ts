import { RoomParticipant } from "../models/events";
import { API_GAME_URI, BaseResponse, STATUS_ERROR } from "./api"

interface GetGameInfoResponse extends BaseResponse {
    participants: Array<RoomParticipant>;
}

export const openCell = async (id: string, row: number, col: number) => {
    const res = await fetch(`${API_GAME_URI}/api/v1/game/${id}/cell/open`, {
        method: "PATCH",
        credentials: "include",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({"row": row, "col": col})
    });
    const data: BaseResponse = await res.json();

    if (data.status == STATUS_ERROR) {
        throw Error(data.error);
    }
}

export const setFlag = async (id: string, row: number, col: number) => {
    const res = await fetch(`${API_GAME_URI}/api/v1/game/${id}/cell/flag`, {
        method: "PATCH",
        credentials: "include",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({"row": row, "col": col})
    });
    const data: BaseResponse = await res.json();

    if (data.status == STATUS_ERROR) {
        throw Error(data.error);
    }
}

export const getGameInfo = async (id: string) => {
    const res = await fetch(`${API_GAME_URI}/api/v1/game/${id}/info`, {
        credentials: "include"
    })
    const data: GetGameInfoResponse = await res.json();
    if (data.status == STATUS_ERROR) {
        throw Error(data.error);
    }
    return data.participants;
}
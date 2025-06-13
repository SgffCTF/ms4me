import { RoomParticipant } from "../models/events";
import { Message } from "../models/models";
import { API_GAME_URI, BaseResponse, STATUS_ERROR } from "./api"

interface GetGameInfoResponse extends BaseResponse {
    participants: Array<RoomParticipant>;
}

interface GetMessagesResponse extends BaseResponse {
    messages: Array<Message> | null;
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

export const getMessages = async (id: string) => {
    const res = await fetch(`${API_GAME_URI}/api/v1/game/${id}/chat`, {
        credentials: "include"
    })
    const data: GetMessagesResponse = await res.json();
    if (data.status == STATUS_ERROR) {
        throw Error(data.error);
    }
    return data.messages ? data.messages : [];
}

export const sendMessage = async (id: string, text: string) => {
    const res = await fetch(`${API_GAME_URI}/api/v1/game/${id}/chat`, {
        method: "POST",
        credentials: "include",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({"text": text})
    })
    const data: BaseResponse = await res.json();
    if (data.status == STATUS_ERROR) {
        throw Error(data.error);
    }
}
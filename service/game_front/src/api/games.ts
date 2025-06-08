import { Game } from "../models/models";
import { API_URI, BaseResponse, STATUS_ERROR } from "./api"

export interface GamesResponse extends BaseResponse {
    games: Array<Game>;
}

export const getGames = async (query: string) => {
    const res = await fetch(`${API_URI}/api/v1/game?query=${query}`, {
        credentials: "include"
    });
    const data: GamesResponse = await res.json();

    if (data.status == STATUS_ERROR) {
        throw Error(data.error);
    }

    return data.games;
}

export const createGame = async (name: string, isPublic: boolean) => {
    const res = await fetch(`${API_URI}/api/v1/game`, {
        method: "POST",
        credentials: "include",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({"title": name, "is_public": isPublic})
    })
    const data: BaseResponse = await res.json();

    if (data.status == STATUS_ERROR) {
        throw Error(data.error);
    }
}

export const deleteGame = async (id: string) => {
    const res = await fetch(`${API_URI}/api/v1/game/${id}`, {
        method: "DELETE",
        credentials: "include"
    })
    const data: BaseResponse = await res.json();
    if (data.status == STATUS_ERROR) {
        throw Error(data.error);
    }
}
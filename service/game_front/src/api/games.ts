import { Game, GameDetails } from "../models/models";
import { API_URI, BaseResponse, STATUS_ERROR } from "./api"

export interface GamesResponse extends BaseResponse {
    games: Array<Game>;
}

export interface GameResponse extends BaseResponse {
    game: GameDetails;
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

export const getMyGames = async () => {
    const res = await fetch(`${API_URI}/api/v1/user/game`, {
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

export const updateGame = async (id: string, name: string, isPublic: boolean) => {
    const res = await fetch(`${API_URI}/api/v1/game/${id}`, {
        method: "PUT",
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

export const getGameByID = async (id: string) => {
    const res = await fetch(`${API_URI}/api/v1/game/${id}`, {
        credentials: "include"
    })
    const data: GameResponse = await res.json();
    if (data.status == STATUS_ERROR) {
        throw Error(data.error);
    }
    return data.game;
}

export const enterGame = async (id: string) => {
    const res = await fetch(`${API_URI}/api/v1/game/${id}/enter`, {
        method: "POST",
        credentials: "include"
    })
    const data: BaseResponse = await res.json();
    if (data.status == STATUS_ERROR) {
        throw Error(data.error);
    }
}
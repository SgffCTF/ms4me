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
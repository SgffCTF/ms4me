import { Game } from "./models";

export const CreateRoomEventType = "CREATE_ROOM";
export const DeleteRoomEventType = "DELETE_ROOM";

export interface WSEvent {
    status: string;
    error?: string;
    event_type: string;
    payload?: any;
    message?: string;
}

export interface CreateRoomEvent {
    game: Game;
}
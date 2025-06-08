import { Game } from "./models";

export const CreateRoomEventType = "CREATE_ROOM";
export const DeleteRoomEventType = "DELETE_ROOM";
export const UpdateRoomEventType = "UPDATE_ROOM";
export const ExitRoomEventType = "EXIT_ROOM";
export const JoinRoomEventType = "JOIN_ROOM";

export const StartGameEventType = "START_GAME";
export const OpenCellEventType = "OPEN_CELL";
export const ResultGameEventType = "RESULT_GAME";

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

export interface UpdateRoomEvent {
    title: string;
    is_public?: boolean;
}

export interface DeleteRoomEvent {
    id: string;
    user_id: number;
}

export interface JoinRoomEvent {
    id: string;
    user_id: number;
    username: string;
}
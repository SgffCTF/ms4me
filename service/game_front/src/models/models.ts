export interface User {
    id: number;
    username: string;
}

export interface Game {
    id: string;
    title: string;
    mines: number;
    rows: number;
    cols: number;
    owner_id: number;
    owner_name: string;
    is_public: boolean;
    created_at: string;
    status: string;
    players: number;
    max_players: number;
}
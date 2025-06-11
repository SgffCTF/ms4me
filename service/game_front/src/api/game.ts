import { API_GAME_URI, BaseResponse, STATUS_ERROR } from "./api"

export const openCell = async (id: string, row: number, col: number) => {
    const res = await fetch(`${API_GAME_URI}/api/v1/game/${id}/cell`, {
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
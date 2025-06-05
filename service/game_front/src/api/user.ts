import { User } from "../models/models";
import { API_URI, BaseResponse, STATUS_OK } from "./api"

interface RegisterResponse extends BaseResponse {
    id: number;
}

export const fetchLogin = async (username: string, password: string) => {
    const res = await fetch(`${API_URI}/api/v1/user/login`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            "username": username,
            "password": password
        }),
        credentials: "include"
    })

    const data: BaseResponse = await res.json();

    if (data.status != "OK") {
        throw Error(data.error);
    }
}

export const fetchRegister = async (username: string, password: string) => {
    const res = await fetch(`${API_URI}/api/v1/user`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            "username": username,
            "password": password
        })
    })

    const data: RegisterResponse = await res.json();

    if (data.status != STATUS_OK) {
        throw Error(data.error);
    }
}

interface UserResponse extends BaseResponse {
    user: User;
}

export const fetchUser = async () => {
    const res = await fetch(`${API_URI}/api/v1/user`, {
        credentials: "include"
    });
    const data: UserResponse = await res.json();
    if (data.status != STATUS_OK) {
        throw Error(data.error);
    }
    return data.user;
}

export const fetchLogout = async () => {
    const res = await fetch(`${API_URI}/api/v1/user/logout`, {
        credentials: "include"
    })
    const data: BaseResponse = await res.json();
    if (data.status == "Error") {
        throw Error(data.error);
    }
}
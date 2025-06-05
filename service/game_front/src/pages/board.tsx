import { useState } from "react";
import { GameList } from "../components/GameList";

export const Board = () => {
    const [searchQuery, setSearchQuery] = useState("");

    return (
        <div className="container-sm">
            <h1 className="text-center mt-5">Поиск игр</h1>

            <div className="row justify-content-start">
                <input
                    className="form-control form-control-lg m-1"
                    type="text"
                    placeholder="Поиск"
                    aria-label=".form-control-lg example"
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}>
                </input>
                <button className={`w-1 btn btn-primary`}>+</button>
                <GameList searchQuery={searchQuery}></GameList>
            </div>
        </div>
    )
}
import { useState } from "react";
import { GameList } from "../components/GameList/GameList";
import { CreateGameModal } from "../components/GameList/CreateGameModal";
import { useAuth } from "../context/AuthProvider";

export const List = () => {
  const [createModalShow, setCreateModalShow] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");
  const [showMyGames, setShowMyGames] = useState(false);
  const { logout } = useAuth();

  return (
    <div className="container-sm">
      <h1 className="text-center mt-5">Поиск игр</h1>

      <div className="mb-3">
        <button
          className={`btn me-3 ${showMyGames ? "btn-primary" : "btn-outline-primary"}`}
          onClick={() => setShowMyGames(true)}
        >
          Мои игры
        </button>
        <button
          className={`btn me-3 ${!showMyGames ? "btn-primary" : "btn-outline-primary"}`}
          onClick={() => setShowMyGames(false)}
        >
          Все игры
        </button>
        <button className="btn btn-red float-end" onClick={logout}>
          Выйти
        </button>
      </div>

      <div className="row justify-content-start align-items-center mb-3">
        {!showMyGames &&
        <input
          className="form-control form-control-lg m-1"
          type="text"
          placeholder="Поиск"
          aria-label="Поиск"
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
        />
        }
        <button className="btn btn-primary w-auto m-1" onClick={() => setCreateModalShow(true)}>
          +
        </button>
      </div>

      <GameList searchQuery={searchQuery} showMyGames={showMyGames} />

      <CreateGameModal show={createModalShow} setShow={setCreateModalShow} />
    </div>
  );
};

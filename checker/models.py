from dataclasses import dataclass
from typing import Optional, Any
from enum import StrEnum

import json


@dataclass
class User:
    id: int
    username: str

@dataclass
class Game:
    id: str
    title: str
    is_public: bool
    mines: int
    rows: int
    cols: int
    owner_id: int
    owner_name: str
    created_at: str
    status: str
    winner_id: Optional[int]
    players_count: int
    max_players: int
    players: list[User]

@dataclass
class Cell:
    value: str
    is_open: bool

@dataclass
class Field:
    rows: int
    cols: int
    mines: int
    cells_open: int
    mine_is_open: bool
    grid: list[list[Cell]]

@dataclass
class Participant:
    id: int
    username: str
    is_owner: bool
    field: Optional[Field] = None

class EventType(StrEnum):
    TYPE_AUTH = "AUTH"
    TYPE_START_GAME = "START_GAME"
    TYPE_CREATE_GAME = "CREATE_ROOM"
    TYPE_JOIN_GAME = "JOIN_ROOM"
    TYPE_DELETE_GAME = "DELETE_ROOM"
    TYPE_UPDATE_GAME = "UPDATE_ROOM"
    TYPE_EXIT_GAME = "EXIT_ROOM"
    TYPE_CLICK_GAME = "OPEN_CELL"
    TYPE_LOSE_GAME = "LOSE_GAME"
    TYPE_WIN_GAME = "WIN_GAME"
    TYPE_NEW_MESSAGE = "NEW_MESSAGE"

@dataclass
class Event:
    event_type: EventType
    status: str
    message: Optional[str] = None
    payload: Optional[dict] = None
    error: Optional[str] = None

@dataclass
class Message:
    id: str
    creator_id: int
    creator_username: str
    text: str
    created_at: str
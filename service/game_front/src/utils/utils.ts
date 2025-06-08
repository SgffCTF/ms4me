import { getGameByID } from "../api/games";
import { GameDetails } from "../models/models";

export function getCookie(name: string) {
  const cookies = document.cookie.split('; ');
  for (let cookie of cookies) {
    const [key, value] = cookie.split('=');
    if (key === name) return decodeURIComponent(value);
  }
  return null;
}

export function formatDate(isoString: string) {
  const date = new Date(isoString);
  return date.toLocaleString('ru-RU', {
    day: 'numeric',
    month: 'long',  // полное название месяца
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false,  // 24-часовой формат
    timeZone: 'UTC' // по умолчанию в UTC, можно изменить при необходимости
  });
}

export async function gameContainsUserID(game: GameDetails, userID: number) {
    try {
        game.players.forEach((player) => {
            if (player.id === userID) {
                return true
            }
        })
    } catch (e: any) {
        console.error("error getting game: " + e.message);
    }
    return false
}
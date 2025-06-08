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

export function gameContainsUserID(game: GameDetails, userID: number) {
    for (let i = 0; i < game.players.length; i++) {
        if (game.players[i].id === userID) {
            return true
        }
    }
    return false
}
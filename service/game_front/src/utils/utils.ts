export function getCookie(name: string) {
  const cookies = document.cookie.split('; ');
  for (let cookie of cookies) {
    const [key, value] = cookie.split('=');
    if (key === name) return decodeURIComponent(value);
  }
  return null;
}
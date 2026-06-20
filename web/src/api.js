// Тонкий клиент к Habitica Monitor API (/api/v1).
// База берётся из VITE_API_BASE; по умолчанию "" — значит относительные пути
// (в dev их проксирует Vite, в prod — отдаёт тот же origin).
const BASE = import.meta.env.VITE_API_BASE || "";

async function request(path, options = {}) {
  const res = await fetch(`${BASE}${path}`, {
    headers: { "Content-Type": "application/json" },
    ...options,
  });
  const text = await res.text();
  const data = text ? JSON.parse(text) : null;
  if (!res.ok) {
    throw new Error((data && data.error) || `HTTP ${res.status}`);
  }
  return data;
}

// POST /api/v1/user — добавить пользователя.
export function createUser({ id, api_key, name }) {
  return request("/api/v1/user", {
    method: "POST",
    body: JSON.stringify({ id, api_key, name }),
  });
}

// GET /api/v1/users — список пользователей со всеми данными (профиль + последний снапшот).
export function getUsers() {
  return request("/api/v1/users");
}

// GET /api/v1/user/stats/history — история статов.
export function getStatsHistory(userId, limit) {
  const params = new URLSearchParams({ user_id: userId });
  if (limit) params.set("limit", String(limit));
  return request(`/api/v1/user/stats/history?${params.toString()}`);
}

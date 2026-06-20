# Habitica Monitor — Web UI

React-интерфейс (Vite + React 18) поверх Go API `/api/v1`.

Две панели:
- Добавление пользователя — `POST /api/v1/user`
- История статов — `GET /api/v1/user/stats/history` (таблица + линейный график HP/MP/Exp/Gold через Recharts)

Выбор пользователя в панели «История статов» — выпадающий список,
который наполняется из `GET /api/v1/users`. При выборе пользователя
история статов подгружается автоматически. Список обновляется после
добавления нового пользователя в левой панели.

## Требования

- Node.js 18+ и npm
- Запущенный Go-сервис (бэкенд) из корня проекта

## Запуск в режиме разработки

1. Запустите бэкенд (из корня репозитория):

   ```bash
   HABITICA_X_CLIENT="<ваш-user-id>-Monitor" go run ./cmd/monitor
   # API поднимется на http://localhost:8080
   ```

2. Запустите фронтенд:

   ```bash
   cd web
   npm install
   npm run dev
   # UI: http://localhost:5173
   ```

В режиме разработки Vite проксирует все запросы `/api/*` на бэкенд, поэтому
CORS не нужен. Если бэкенд на другом адресе — задайте цель прокси:

```bash
VITE_API_TARGET=http://localhost:9000 npm run dev
```

## Production-сборка

```bash
cd web
npm run build      # результат в web/dist
npm run preview    # локальный предпросмотр собранной версии
```

Готовую статику из `web/dist` можно раздавать любым статик-сервером
(nginx, Caddy и т.п.).

### Если UI и API на разных origin (prod)

В этом случае прокси Vite уже не работает. Два варианта:

1. **Указать адрес API при сборке** через переменную `VITE_API_BASE`:

   ```bash
   VITE_API_BASE=https://api.example.com npm run build
   ```

   Тогда клиент будет ходить напрямую на этот адрес.

2. **Разрешить CORS на бэкенде** — Go-сервис уже поддерживает это через
   переменную окружения `CORS_ORIGIN`:

   ```bash
   CORS_ORIGIN=https://ui.example.com go run ./cmd/monitor
   ```

   (по умолчанию `*`).

## Переменные окружения фронтенда

| Переменная | Назначение | По умолчанию |
|---|---|---|
| `VITE_API_TARGET` | цель dev-прокси Vite | `http://localhost:8080` |
| `VITE_API_BASE` | базовый URL API в собранной версии | `""` (тот же origin) |

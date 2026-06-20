# Habitica Monitor

Сервис мониторинга команды в Habitica на Go. Раз в час (настраиваемо) снимает статы
каждого пользователя через `GET /user` и сохраняет каждый замер отдельной записью.

Реализован по [golang-standards/project-layout](https://github.com/golang-standards/project-layout)
с разделением на слои в духе **DDD**.

## Структура

```
cmd/monitor/             — точка входа (composition root)
api/openapi.yaml         — спецификация публичного API v1
configs/                 — конфигурация из ENV
internal/
  domain/                — чистый домен, без внешних зависимостей
    user/                — агрегат User + порт Repository
    stats/               — агрегат StatsHistory (Snapshot) + порт Repository
  application/           — сценарии (use cases): UserService, StatsService, SnapshotService
  infrastructure/        — адаптеры: SQLite-репозитории, Habitica API-клиент
  interfaces/http/       — входной адаптер: REST-хендлеры
```

Зависимости направлены внутрь: `interfaces → application → domain`, инфраструктура
реализует доменные порты. Домен ничего не знает об HTTP и SQLite.

## База данных (SQLite)

- **user**: `id` (x-api-user), `api_key` (x-api-key), `name` (задаётся вручную)
- **stats_history**: `user_id, hp, mp, exp, gp, lvl, timestamp` — каждый снапшот = новая строка

Драйвер `modernc.org/sqlite` — чистый Go, без CGO. Схема создаётся автоматически.

## API (`/api/v1`)

| Метод | Путь | Назначение |
|---|---|---|
| POST | `/api/v1/user` | добавить пользователя в таблицу `user` |
| GET  | `/api/v1/users` | список пользователей со всеми данными (профиль + последний снапшот) |
| GET  | `/api/v1/user/stats/history?user_id=...&limit=...` | история статов из `stats_history` |

Примеры:

```bash
curl -X POST localhost:8080/api/v1/user \
  -H 'Content-Type: application/json' \
  -d '{"id":"abc-123","api_key":"secret-token","name":"alice"}'

curl "localhost:8080/api/v1/user/stats/history?user_id=abc-123&limit=10"
```

## Снапшоты

`SnapshotService` запускается фоном из `main.go`. На каждой итерации тикера он берёт всех
пользователей и обрабатывает **каждого в отдельной горутине** (`sync.WaitGroup` ждёт
завершения цикла). Каждый успешный замер пишется новой записью в `stats_history`.

## Конфигурация (ENV)

| Переменная | Назначение | По умолчанию |
|---|---|---|
| `HTTP_ADDR` | адрес HTTP-сервера | `:8080` |
| `DB_PATH` | путь к файлу SQLite | `monitor.db` |
| `HABITICA_X_CLIENT` | заголовок `x-client` (`<your-user-id>-<appname>`) | нужно задать |
| `SNAPSHOT_INTERVAL` | период снапшотов (Go-duration: `1h`, `30m`, `15m`) | `1h` |

## Запуск

```bash
go mod tidy
HABITICA_X_CLIENT="<ваш-user-id>-Monitor" SNAPSHOT_INTERVAL=1h go run ./cmd/monitor
```

## Web UI

React-интерфейс лежит в каталоге `web/` (Vite + React). Он добавляет пользователей
и показывает историю статов (таблица + график). Инструкция запуска — в `web/README.md`.

Коротко: запустите бэкенд (`go run ./cmd/monitor`), затем `cd web && npm install && npm run dev`
и откройте http://localhost:5173.

## Тесты

```bash
go test ./...
```

`internal/application/snapshot_service_test.go` проверяет, что за один прогон
сохраняется по снапшоту на каждого пользователя (с моками портов).

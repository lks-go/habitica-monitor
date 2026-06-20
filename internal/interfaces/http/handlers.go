// Package http — входной адаптер: REST-хендлеры поверх application-сервисов.
package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/example/habitica-monitor/internal/application"
	"github.com/example/habitica-monitor/internal/domain/stats"
	"github.com/example/habitica-monitor/internal/domain/user"
)

// Handler держит зависимости HTTP-слоя.
type Handler struct {
	users *application.UserService
	stats *application.StatsService
}

func NewHandler(users *application.UserService, stats *application.StatsService) *Handler {
	return &Handler{users: users, stats: stats}
}

// Routes регистрирует маршруты под префиксом /api/v1 (Go 1.22+ метод+путь).
func (h *Handler) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/user", h.createUser)
	mux.HandleFunc("GET /api/v1/users", h.listUsers)
	mux.HandleFunc("GET /api/v1/user/stats/history", h.statsHistory)
	return mux
}

// --- DTO ---

type createUserRequest struct {
	ID     string `json:"id"`      // x-api-user
	APIKey string `json:"api_key"` // x-api-key
	Name   string `json:"name"`
}

type userResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// userWithStatsResponse — пользователь со всеми данными для GET /api/v1/users.
// latest_stats == null, если снапшотов ещё нет.
type userWithStatsResponse struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	LatestStats *snapshotResponse `json:"latest_stats"`
}

type snapshotResponse struct {
	UserID    string    `json:"user_id"`
	HP        float64   `json:"hp"`
	MP        float64   `json:"mp"`
	Exp       float64   `json:"exp"`
	GP        float64   `json:"gp"`
	Lvl       int       `json:"lvl"`
	Timestamp time.Time `json:"timestamp"`
}

// POST /api/v1/user — добавить пользователя в таблицу user.
func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	u, err := h.users.AddUser(r.Context(), req.ID, req.APIKey, req.Name)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrEmptyID),
			errors.Is(err, user.ErrEmptyAPIKey),
			errors.Is(err, user.ErrEmptyName):
			writeError(w, http.StatusBadRequest, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "failed to save user")
		}
		return
	}
	writeJSON(w, http.StatusCreated, userResponse{ID: u.ID, Name: u.Name})
}

// GET /api/v1/users — список пользователей со всеми данными (профиль + последний снапшот).
func (h *Handler) listUsers(w http.ResponseWriter, r *http.Request) {
	items, err := h.users.ListUsers(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list users")
		return
	}

	out := make([]userWithStatsResponse, 0, len(items))
	for _, it := range items {
		resp := userWithStatsResponse{ID: it.User.ID, Name: it.User.Name}
		if it.LatestStats != nil {
			s := it.LatestStats
			resp.LatestStats = &snapshotResponse{
				UserID: s.UserID, HP: s.HP, MP: s.MP, Exp: s.Exp,
				GP: s.GP, Lvl: s.Lvl, Timestamp: s.Timestamp,
			}
		}
		out = append(out, resp)
	}
	writeJSON(w, http.StatusOK, out)
}

// GET /api/v1/user/stats/history?user_id=...&limit=...
func (h *Handler) statsHistory(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "query param user_id is required")
		return
	}
	limit := parseLimit(r.URL.Query().Get("limit"))

	items, err := h.stats.History(r.Context(), userID, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to read history")
		return
	}
	writeJSON(w, http.StatusOK, toSnapshotResponses(items))
}

// --- helpers ---

func toSnapshotResponses(items []*stats.Snapshot) []snapshotResponse {
	out := make([]snapshotResponse, 0, len(items))
	for _, s := range items {
		out = append(out, snapshotResponse{
			UserID: s.UserID, HP: s.HP, MP: s.MP, Exp: s.Exp,
			GP: s.GP, Lvl: s.Lvl, Timestamp: s.Timestamp,
		})
	}
	return out
}

func parseLimit(s string) int {
	if s == "" {
		return 0
	}
	n, err := strconv.Atoi(s)
	if err != nil || n < 0 {
		return 0
	}
	return n
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

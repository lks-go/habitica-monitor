// Package stats — доменный слой агрегата StatsHistory.
// Каждый снапшот статов пользователя — отдельная неизменяемая запись.
package stats

import (
	"context"
	"errors"
	"time"
)

// ErrEmptyUserID — попытка создать снапшот без привязки к пользователю.
var ErrEmptyUserID = errors.New("stats: user_id is required")

// Snapshot — одна точка истории статов пользователя в момент времени.
// Таблица stats_history: (user_id, hp, mp, exp, gp, lvl, timestamp).
type Snapshot struct {
	UserID    string
	HP        float64
	MP        float64
	Exp       float64
	GP        float64
	Lvl       int
	Timestamp time.Time
}

// NewSnapshot создаёт валидный снапшот с меткой времени now (UTC).
func NewSnapshot(userID string, hp, mp, exp, gp float64, lvl int, now time.Time) (*Snapshot, error) {
	if userID == "" {
		return nil, ErrEmptyUserID
	}
	return &Snapshot{
		UserID:    userID,
		HP:        hp,
		MP:        mp,
		Exp:       exp,
		GP:        gp,
		Lvl:       lvl,
		Timestamp: now.UTC(),
	}, nil
}

// Repository — контракт хранилища истории статов (порт).
type Repository interface {
	Add(ctx context.Context, s *Snapshot) error
	// ListByUser возвращает снапшоты пользователя, новые первыми; limit<=0 — без ограничения.
	ListByUser(ctx context.Context, userID string, limit int) ([]*Snapshot, error)
}

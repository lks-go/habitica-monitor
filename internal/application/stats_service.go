package application

import (
	"context"

	"github.com/example/habitica-monitor/internal/domain/stats"
)

// StatsService — сценарии чтения истории статов.
type StatsService struct {
	history stats.Repository
}

func NewStatsService(history stats.Repository) *StatsService {
	return &StatsService{history: history}
}

// History возвращает историю статов пользователя (GET /api/v1/user/stats/history).
func (s *StatsService) History(ctx context.Context, userID string, limit int) ([]*stats.Snapshot, error) {
	return s.history.ListByUser(ctx, userID, limit)
}

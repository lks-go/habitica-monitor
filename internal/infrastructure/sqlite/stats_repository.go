package sqlite

import (
	"context"
	"database/sql"
	"time"

	"github.com/example/habitica-monitor/internal/domain/stats"
)

// StatsRepository реализует stats.Repository поверх SQLite.
type StatsRepository struct {
	db *sql.DB
}

func NewStatsRepository(db *sql.DB) *StatsRepository {
	return &StatsRepository{db: db}
}

// Add сохраняет снапшот как новую запись (история не перезаписывается).
func (r *StatsRepository) Add(ctx context.Context, s *stats.Snapshot) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO stats_history (user_id, hp, mp, exp, gp, lvl, timestamp)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		s.UserID, s.HP, s.MP, s.Exp, s.GP, s.Lvl, s.Timestamp.UTC().Format(time.RFC3339),
	)
	return err
}

// ListByUser возвращает снапшоты пользователя, новые первыми.
func (r *StatsRepository) ListByUser(ctx context.Context, userID string, limit int) ([]*stats.Snapshot, error) {
	query := `SELECT user_id, hp, mp, exp, gp, lvl, timestamp
	          FROM stats_history WHERE user_id = ? ORDER BY timestamp DESC`
	args := []any{userID}
	if limit > 0 {
		query += ` LIMIT ?`
		args = append(args, limit)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*stats.Snapshot
	for rows.Next() {
		var (
			s  stats.Snapshot
			ts string
		)
		if err := rows.Scan(&s.UserID, &s.HP, &s.MP, &s.Exp, &s.GP, &s.Lvl, &ts); err != nil {
			return nil, err
		}
		s.Timestamp, _ = time.Parse(time.RFC3339, ts)
		out = append(out, &s)
	}
	return out, rows.Err()
}

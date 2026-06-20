package application

import (
	"context"

	"github.com/example/habitica-monitor/internal/domain/stats"
	"github.com/example/habitica-monitor/internal/domain/user"
)

// UserWithStats — пользователь со всеми данными: профиль + последний снапшот статов.
// LatestStats == nil, если для пользователя ещё нет ни одного снапшота.
type UserWithStats struct {
	User        *user.User
	LatestStats *stats.Snapshot
}

// UserService — сценарии работы с пользователями.
type UserService struct {
	users user.Repository
	stats stats.Repository
}

func NewUserService(users user.Repository, statsRepo stats.Repository) *UserService {
	return &UserService{users: users, stats: statsRepo}
}

// AddUser создаёт и сохраняет пользователя (POST /api/v1/user).
func (s *UserService) AddUser(ctx context.Context, id, apiKey, name string) (*user.User, error) {
	u, err := user.NewUser(id, apiKey, name)
	if err != nil {
		return nil, err
	}
	if err := s.users.Save(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

// ListUsers возвращает всех пользователей вместе с их последним снапшотом статов
// (GET /api/v1/users).
func (s *UserService) ListUsers(ctx context.Context) ([]UserWithStats, error) {
	users, err := s.users.List(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]UserWithStats, 0, len(users))
	for _, u := range users {
		item := UserWithStats{User: u}

		// Берём только самый свежий снапшот (новые первыми, limit=1).
		snaps, err := s.stats.ListByUser(ctx, u.ID, 1)
		if err != nil {
			return nil, err
		}
		if len(snaps) > 0 {
			item.LatestStats = snaps[0]
		}
		out = append(out, item)
	}
	return out, nil
}

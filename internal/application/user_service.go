package application

import (
	"context"

	"github.com/example/habitica-monitor/internal/domain/user"
)

// UserService — сценарии работы с пользователями.
type UserService struct {
	users user.Repository
}

func NewUserService(users user.Repository) *UserService {
	return &UserService{users: users}
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

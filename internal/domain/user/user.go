// Package user — доменный слой агрегата User.
// Содержит сущность и контракт репозитория. Зависимостей от инфраструктуры нет.
package user

import (
	"context"
	"errors"
	"strings"
)

// Ошибки домена.
var (
	ErrEmptyID     = errors.New("user: id (x-api-user) is required")
	ErrEmptyAPIKey = errors.New("user: api_key (x-api-key) is required")
	ErrEmptyName   = errors.New("user: name is required")
	ErrNotFound    = errors.New("user: not found")
)

// User — участник команды, за которым ведётся наблюдение.
//   - ID     соответствует заголовку Habitica x-api-user.
//   - APIKey соответствует заголовку Habitica x-api-key (секрет).
//   - Name   задаётся вручную при добавлении.
type User struct {
	ID     string
	APIKey string
	Name   string
}

// NewUser создаёт валидную доменную сущность User.
func NewUser(id, apiKey, name string) (*User, error) {
	id = strings.TrimSpace(id)
	apiKey = strings.TrimSpace(apiKey)
	name = strings.TrimSpace(name)

	switch {
	case id == "":
		return nil, ErrEmptyID
	case apiKey == "":
		return nil, ErrEmptyAPIKey
	case name == "":
		return nil, ErrEmptyName
	}
	return &User{ID: id, APIKey: apiKey, Name: name}, nil
}

// Repository — контракт хранилища пользователей (порт).
// Реализуется в слое infrastructure.
type Repository interface {
	Save(ctx context.Context, u *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	List(ctx context.Context) ([]*User, error)
}

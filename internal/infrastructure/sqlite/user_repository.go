package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/example/habitica-monitor/internal/domain/user"
)

// UserRepository реализует user.Repository поверх SQLite.
type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Save вставляет нового или обновляет существующего пользователя (upsert по id).
func (r *UserRepository) Save(ctx context.Context, u *user.User) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO user (id, api_key, name) VALUES (?, ?, ?)
		 ON CONFLICT(id) DO UPDATE SET api_key=excluded.api_key, name=excluded.name`,
		u.ID, u.APIKey, u.Name,
	)
	return err
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*user.User, error) {
	var u user.User
	err := r.db.QueryRowContext(ctx,
		`SELECT id, api_key, name FROM user WHERE id = ?`, id,
	).Scan(&u.ID, &u.APIKey, &u.Name)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, user.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) List(ctx context.Context) ([]*user.User, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, api_key, name FROM user ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*user.User
	for rows.Next() {
		var u user.User
		if err := rows.Scan(&u.ID, &u.APIKey, &u.Name); err != nil {
			return nil, err
		}
		out = append(out, &u)
	}
	return out, rows.Err()
}

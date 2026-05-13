package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nathanim1919/mezgeb/internal/domain"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (r *UserRepo) Upsert(ctx context.Context, user *domain.User) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO users (id, first_name, username, language_code)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET
			first_name = EXCLUDED.first_name,
			username = EXCLUDED.username,
			updated_at = NOW()
	`, user.ID, user.FirstName, user.Username, user.LanguageCode)
	return err
}

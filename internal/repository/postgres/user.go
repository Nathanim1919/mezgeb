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

func (r *UserRepo) GetLang(ctx context.Context, userID int64) (string, error) {
	var lang string
	err := r.pool.QueryRow(ctx, `SELECT lang FROM users WHERE id = $1`, userID).Scan(&lang)
	if err != nil {
		return "am", nil // default to Amharic
	}
	return lang, nil
}

func (r *UserRepo) SetLang(ctx context.Context, userID int64, lang string) error {
	_, err := r.pool.Exec(ctx, `UPDATE users SET lang = $1, updated_at = NOW() WHERE id = $2`, lang, userID)
	return err
}

// ClearData deletes all transactions, customers, and products for a user in a single transaction.
// The user record itself is kept (so language preference is preserved).
func (r *UserRepo) ClearData(ctx context.Context, userID int64) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Order matters: transactions reference customers and products
	_, err = tx.Exec(ctx, `DELETE FROM transactions WHERE user_id = $1`, userID)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `DELETE FROM customers WHERE user_id = $1`, userID)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `DELETE FROM products WHERE user_id = $1`, userID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

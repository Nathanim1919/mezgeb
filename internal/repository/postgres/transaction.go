package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nathanim1919/mezgeb/internal/domain"
)

type TransactionRepo struct {
	pool *pgxpool.Pool
}

func NewTransactionRepo(pool *pgxpool.Pool) *TransactionRepo {
	return &TransactionRepo{pool: pool}
}

func (r *TransactionRepo) Create(ctx context.Context, tx *domain.Transaction) error {
	return r.pool.QueryRow(ctx, `
		INSERT INTO transactions (user_id, customer_id, product_id, type, amount, note)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`, tx.UserID, tx.CustomerID, tx.ProductID, tx.Type, tx.Amount, tx.Note).Scan(&tx.ID, &tx.CreatedAt)
}

func (r *TransactionRepo) ListByUser(ctx context.Context, userID int64, from, to time.Time) ([]domain.Transaction, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT t.id, t.user_id, t.customer_id, t.product_id, t.type, t.amount, t.note, t.created_at,
		       c.name AS customer_name,
		       COALESCE(p.name, '') AS product_name
		FROM transactions t
		JOIN customers c ON c.id = t.customer_id
		LEFT JOIN products p ON p.id = t.product_id
		WHERE t.user_id = $1 AND t.created_at >= $2 AND t.created_at < $3
		ORDER BY t.created_at DESC
	`, userID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txns []domain.Transaction
	for rows.Next() {
		var t domain.Transaction
		if err := rows.Scan(&t.ID, &t.UserID, &t.CustomerID, &t.ProductID, &t.Type, &t.Amount, &t.Note, &t.CreatedAt, &t.CustomerName, &t.ProductName); err != nil {
			return nil, err
		}
		txns = append(txns, t)
	}
	return txns, nil
}

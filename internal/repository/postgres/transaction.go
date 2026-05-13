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

// CreateWithBalanceUpdate atomically inserts a transaction, updates the customer balance
// (if customer exists), and updates product stock (if product and stockDelta provided).
func (r *TransactionRepo) CreateWithBalanceUpdate(ctx context.Context, tx *domain.Transaction, balanceDelta int64, stockDelta int64) error {
	dbTx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer dbTx.Rollback(ctx)

	err = dbTx.QueryRow(ctx, `
		INSERT INTO transactions (user_id, customer_id, product_id, type, amount, quantity, note)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at
	`, tx.UserID, tx.CustomerID, tx.ProductID, tx.Type, tx.Amount, tx.Quantity, tx.Note).Scan(&tx.ID, &tx.CreatedAt)
	if err != nil {
		return err
	}

	// Update customer balance if customer is set (debt/payment/purchase flows)
	if tx.CustomerID != nil && balanceDelta != 0 {
		_, err = dbTx.Exec(ctx, `
			UPDATE customers SET balance = balance + $1, updated_at = NOW()
			WHERE id = $2 AND user_id = $3
		`, balanceDelta, *tx.CustomerID, tx.UserID)
		if err != nil {
			return err
		}
	}

	// Update product stock if product is set (sell/buy flows)
	if tx.ProductID != nil && stockDelta != 0 {
		_, err = dbTx.Exec(ctx, `
			UPDATE products SET stock = stock + $1, updated_at = NOW()
			WHERE id = $2 AND user_id = $3
		`, stockDelta, *tx.ProductID, tx.UserID)
		if err != nil {
			return err
		}
	}

	return dbTx.Commit(ctx)
}

func (r *TransactionRepo) ListByUser(ctx context.Context, userID int64, from, to time.Time) ([]domain.Transaction, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT t.id, t.user_id, t.customer_id, t.product_id, t.type, t.amount, t.quantity, t.note, t.created_at,
		       COALESCE(c.name, '') AS customer_name,
		       COALESCE(p.name, '') AS product_name
		FROM transactions t
		LEFT JOIN customers c ON c.id = t.customer_id
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
		if err := rows.Scan(&t.ID, &t.UserID, &t.CustomerID, &t.ProductID, &t.Type, &t.Amount, &t.Quantity, &t.Note, &t.CreatedAt, &t.CustomerName, &t.ProductName); err != nil {
			return nil, err
		}
		txns = append(txns, t)
	}
	return txns, nil
}

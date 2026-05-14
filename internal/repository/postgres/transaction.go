package postgres

import (
	"context"
	"fmt"
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

func (r *TransactionRepo) ListByType(ctx context.Context, userID int64, txType domain.TransactionType, limit int) ([]domain.Transaction, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT t.id, t.user_id, t.customer_id, t.product_id, t.type, t.amount, t.quantity, t.note, t.created_at,
		       COALESCE(c.name, '') AS customer_name,
		       COALESCE(p.name, '') AS product_name
		FROM transactions t
		LEFT JOIN customers c ON c.id = t.customer_id
		LEFT JOIN products p ON p.id = t.product_id
		WHERE t.user_id = $1 AND t.type = $2
		ORDER BY t.created_at DESC
		LIMIT $3
	`, userID, txType, limit)
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

func (r *TransactionRepo) GetByID(ctx context.Context, userID, id int64) (*domain.Transaction, error) {
	var t domain.Transaction
	err := r.pool.QueryRow(ctx, `
		SELECT t.id, t.user_id, t.customer_id, t.product_id, t.type, t.amount, t.quantity, t.note, t.created_at,
		       COALESCE(c.name, '') AS customer_name,
		       COALESCE(p.name, '') AS product_name
		FROM transactions t
		LEFT JOIN customers c ON c.id = t.customer_id
		LEFT JOIN products p ON p.id = t.product_id
		WHERE t.id = $1 AND t.user_id = $2
	`, id, userID).Scan(&t.ID, &t.UserID, &t.CustomerID, &t.ProductID, &t.Type, &t.Amount, &t.Quantity, &t.Note, &t.CreatedAt, &t.CustomerName, &t.ProductName)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TransactionRepo) UpdateAmountAndQuantity(ctx context.Context, userID, id int64, amount, quantity int64, balanceDelta, stockDelta int64) error {
	dbTx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer dbTx.Rollback(ctx)

	tag, err := dbTx.Exec(ctx, `
		UPDATE transactions SET amount = $1, quantity = $2
		WHERE id = $3 AND user_id = $4
	`, amount, quantity, id, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("transaction not found")
	}

	// Fetch the transaction to get customer_id and product_id for rollback
	var customerID *int64
	var productID *int64
	err = dbTx.QueryRow(ctx, `SELECT customer_id, product_id FROM transactions WHERE id = $1`, id).Scan(&customerID, &productID)
	if err != nil {
		return err
	}

	if customerID != nil && balanceDelta != 0 {
		_, err = dbTx.Exec(ctx, `
			UPDATE customers SET balance = balance + $1, updated_at = NOW()
			WHERE id = $2 AND user_id = $3
		`, balanceDelta, *customerID, userID)
		if err != nil {
			return err
		}
	}

	if productID != nil && stockDelta != 0 {
		_, err = dbTx.Exec(ctx, `
			UPDATE products SET stock = stock + $1, updated_at = NOW()
			WHERE id = $2 AND user_id = $3
		`, stockDelta, *productID, userID)
		if err != nil {
			return err
		}
	}

	return dbTx.Commit(ctx)
}

func (r *TransactionRepo) UpdateNote(ctx context.Context, userID, id int64, note string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE transactions SET note = $1
		WHERE id = $2 AND user_id = $3
	`, note, id, userID)
	return err
}

func (r *TransactionRepo) DeleteWithRollback(ctx context.Context, userID, id int64, balanceDelta, stockDelta int64) error {
	dbTx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer dbTx.Rollback(ctx)

	// Fetch the transaction to get customer_id and product_id
	var customerID *int64
	var productID *int64
	err = dbTx.QueryRow(ctx, `
		SELECT customer_id, product_id FROM transactions WHERE id = $1 AND user_id = $2
	`, id, userID).Scan(&customerID, &productID)
	if err != nil {
		return err
	}

	// Delete the transaction
	tag, err := dbTx.Exec(ctx, `DELETE FROM transactions WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("transaction not found")
	}

	// Rollback customer balance
	if customerID != nil && balanceDelta != 0 {
		_, err = dbTx.Exec(ctx, `
			UPDATE customers SET balance = balance + $1, updated_at = NOW()
			WHERE id = $2 AND user_id = $3
		`, balanceDelta, *customerID, userID)
		if err != nil {
			return err
		}
	}

	// Rollback product stock
	if productID != nil && stockDelta != 0 {
		_, err = dbTx.Exec(ctx, `
			UPDATE products SET stock = stock + $1, updated_at = NOW()
			WHERE id = $2 AND user_id = $3
		`, stockDelta, *productID, userID)
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

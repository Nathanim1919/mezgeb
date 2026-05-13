package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nathanim1919/mezgeb/internal/domain"
)

type CustomerRepo struct {
	pool *pgxpool.Pool
}

func NewCustomerRepo(pool *pgxpool.Pool) *CustomerRepo {
	return &CustomerRepo{pool: pool}
}

func (r *CustomerRepo) FindOrCreate(ctx context.Context, userID int64, name string) (*domain.Customer, error) {
	var c domain.Customer
	err := r.pool.QueryRow(ctx, `
		INSERT INTO customers (user_id, name)
		VALUES ($1, $2)
		ON CONFLICT (user_id, LOWER(name)) DO UPDATE SET updated_at = NOW()
		RETURNING id, user_id, name, phone, balance, created_at, updated_at
	`, userID, name).Scan(&c.ID, &c.UserID, &c.Name, &c.Phone, &c.Balance, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CustomerRepo) ListByUser(ctx context.Context, userID int64) ([]domain.Customer, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, name, phone, balance, created_at, updated_at
		FROM customers
		WHERE user_id = $1
		ORDER BY updated_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []domain.Customer
	for rows.Next() {
		var c domain.Customer
		if err := rows.Scan(&c.ID, &c.UserID, &c.Name, &c.Phone, &c.Balance, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		customers = append(customers, c)
	}
	return customers, nil
}

func (r *CustomerRepo) GetByID(ctx context.Context, userID, id int64) (*domain.Customer, error) {
	var c domain.Customer
	err := r.pool.QueryRow(ctx, `
		SELECT id, user_id, name, phone, balance, created_at, updated_at
		FROM customers WHERE id = $1 AND user_id = $2
	`, id, userID).Scan(&c.ID, &c.UserID, &c.Name, &c.Phone, &c.Balance, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CustomerRepo) UpdateBalance(ctx context.Context, userID, id int64, delta int64) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE customers SET balance = balance + $1, updated_at = NOW()
		WHERE id = $2 AND user_id = $3
	`, delta, id, userID)
	return err
}

package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nathanim1919/mezgeb/internal/domain"
)

type ProductRepo struct {
	pool *pgxpool.Pool
}

func NewProductRepo(pool *pgxpool.Pool) *ProductRepo {
	return &ProductRepo{pool: pool}
}

func (r *ProductRepo) FindOrCreate(ctx context.Context, userID int64, name string, price int64) (*domain.Product, error) {
	var p domain.Product
	err := r.pool.QueryRow(ctx, `
		INSERT INTO products (user_id, name, price)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, LOWER(name)) DO UPDATE SET name = EXCLUDED.name
		RETURNING id, user_id, name, price, created_at
	`, userID, name, price).Scan(&p.ID, &p.UserID, &p.Name, &p.Price, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepo) ListByUser(ctx context.Context, userID int64) ([]domain.Product, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, name, price, created_at
		FROM products
		WHERE user_id = $1
		ORDER BY name
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(&p.ID, &p.UserID, &p.Name, &p.Price, &p.CreatedAt); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *ProductRepo) GetByID(ctx context.Context, userID, id int64) (*domain.Product, error) {
	var p domain.Product
	err := r.pool.QueryRow(ctx, `
		SELECT id, user_id, name, price, created_at
		FROM products WHERE id = $1 AND user_id = $2
	`, id, userID).Scan(&p.ID, &p.UserID, &p.Name, &p.Price, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

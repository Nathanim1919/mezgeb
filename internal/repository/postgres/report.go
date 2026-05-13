package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nathanim1919/mezgeb/internal/domain"
)

type ReportRepo struct {
	pool *pgxpool.Pool
}

func NewReportRepo(pool *pgxpool.Pool) *ReportRepo {
	return &ReportRepo{pool: pool}
}

func (r *ReportRepo) GetReport(ctx context.Context, userID int64, from, to time.Time) (*domain.ReportData, error) {
	report := &domain.ReportData{}

	err := r.pool.QueryRow(ctx, `
		SELECT
			COUNT(*),
			COALESCE(SUM(CASE WHEN type = 'sell' THEN amount ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN type = 'buy' THEN amount ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN type IN ('payment', 'purchase') THEN amount ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN type = 'debt' THEN amount ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN type = 'sell' THEN quantity ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN type = 'buy' THEN quantity ELSE 0 END), 0)
		FROM transactions
		WHERE user_id = $1 AND created_at >= $2 AND created_at < $3
	`, userID, from, to).Scan(
		&report.TotalTransactions,
		&report.TotalSales,
		&report.TotalExpenses,
		&report.TotalRevenue,
		&report.TotalDebt,
		&report.ItemsSold,
		&report.ItemsBought,
	)
	if err != nil {
		return nil, err
	}

	// Top sold products (by revenue)
	rows, err := r.pool.Query(ctx, `
		SELECT p.name, SUM(t.quantity) AS cnt, SUM(t.amount) AS total
		FROM transactions t
		JOIN products p ON p.id = t.product_id
		WHERE t.user_id = $1 AND t.created_at >= $2 AND t.created_at < $3
			AND t.type = 'sell'
		GROUP BY p.name
		ORDER BY total DESC
		LIMIT 5
	`, userID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ps domain.ProductStat
		if err := rows.Scan(&ps.Name, &ps.Count, &ps.Total); err != nil {
			return nil, err
		}
		report.TopProducts = append(report.TopProducts, ps)
	}

	return report, nil
}

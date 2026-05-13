package domain

import (
	"context"
	"time"
)

type UserRepo interface {
	Upsert(ctx context.Context, user *User) error
	GetLang(ctx context.Context, userID int64) (string, error)
	SetLang(ctx context.Context, userID int64, lang string) error
}

type CustomerRepo interface {
	FindOrCreate(ctx context.Context, userID int64, name string) (*Customer, error)
	ListByUser(ctx context.Context, userID int64) ([]Customer, error)
	GetByID(ctx context.Context, userID, id int64) (*Customer, error)
	UpdateBalance(ctx context.Context, userID, id int64, delta int64) error
}

type ProductRepo interface {
	FindOrCreate(ctx context.Context, userID int64, name string, price int64) (*Product, error)
	ListByUser(ctx context.Context, userID int64) ([]Product, error)
	GetByID(ctx context.Context, userID, id int64) (*Product, error)
}

type TransactionRepo interface {
	// CreateWithBalanceUpdate atomically creates a transaction and updates the customer balance.
	CreateWithBalanceUpdate(ctx context.Context, tx *Transaction, balanceDelta int64) error
	ListByUser(ctx context.Context, userID int64, from, to time.Time) ([]Transaction, error)
}

type ReportData struct {
	TotalTransactions int
	TotalRevenue      int64 // payments + purchases
	TotalDebt         int64 // new debt added
	TopProducts       []ProductStat
}

type ProductStat struct {
	Name  string
	Count int
	Total int64
}

type ReportRepo interface {
	GetReport(ctx context.Context, userID int64, from, to time.Time) (*ReportData, error)
}

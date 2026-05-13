package domain

import (
	"context"
	"time"
)

type UserRepo interface {
	Upsert(ctx context.Context, user *User) error
	GetLang(ctx context.Context, userID int64) (string, error)
	SetLang(ctx context.Context, userID int64, lang string) error
	ClearData(ctx context.Context, userID int64) error
}

type CustomerRepo interface {
	FindOrCreate(ctx context.Context, userID int64, name string) (*Customer, error)
	ListByUser(ctx context.Context, userID int64) ([]Customer, error)
	GetByID(ctx context.Context, userID, id int64) (*Customer, error)
	UpdateBalance(ctx context.Context, userID, id int64, delta int64) error
}

type ProductRepo interface {
	FindOrCreate(ctx context.Context, userID int64, name string, price int64, stock int64) (*Product, error)
	ListByUser(ctx context.Context, userID int64) ([]Product, error)
	GetByID(ctx context.Context, userID, id int64) (*Product, error)
}

type TransactionRepo interface {
	// CreateWithBalanceUpdate atomically creates a transaction, updates customer balance and product stock.
	CreateWithBalanceUpdate(ctx context.Context, tx *Transaction, balanceDelta int64, stockDelta int64) error
	ListByUser(ctx context.Context, userID int64, from, to time.Time) ([]Transaction, error)
}

type ReportData struct {
	TotalTransactions int
	TotalSales        int64 // money in from sell transactions
	TotalExpenses     int64 // money out from buy transactions
	TotalBorrowed     int64 // others owe you (borrow/debt)
	TotalLoaned       int64 // you owe others (loan)
	TotalRevenue      int64 // legacy: payments + purchases
	TotalDebt         int64 // legacy: new debt added
	ItemsSold         int64 // total quantity sold
	ItemsBought       int64 // total quantity bought
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

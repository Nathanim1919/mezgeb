package service

import (
	"context"
	"time"

	"github.com/nathanim1919/mezgeb/internal/domain"
)

type Service struct {
	Users        domain.UserRepo
	Customers    domain.CustomerRepo
	Products     domain.ProductRepo
	Transactions domain.TransactionRepo
	Reports      domain.ReportRepo
}

func (s *Service) EnsureUser(ctx context.Context, user *domain.User) error {
	return s.Users.Upsert(ctx, user)
}

func (s *Service) AddTransaction(ctx context.Context, tx *domain.Transaction) error {
	var delta int64
	switch tx.Type {
	case domain.TxDebt:
		delta = tx.Amount // they owe more
	case domain.TxPayment:
		delta = -tx.Amount // they paid, reduce what they owe
	case domain.TxPurchase:
		delta = tx.Amount // purchase on credit = debt
	}

	// Atomic: transaction insert + balance update in a single DB transaction
	return s.Transactions.CreateWithBalanceUpdate(ctx, tx, delta)
}

func (s *Service) GetReport(ctx context.Context, userID int64, from, to time.Time) (*domain.ReportData, error) {
	return s.Reports.GetReport(ctx, userID, from, to)
}

func (s *Service) FindOrCreateCustomer(ctx context.Context, userID int64, name string) (*domain.Customer, error) {
	return s.Customers.FindOrCreate(ctx, userID, name)
}

func (s *Service) ListCustomers(ctx context.Context, userID int64) ([]domain.Customer, error) {
	return s.Customers.ListByUser(ctx, userID)
}

func (s *Service) FindOrCreateProduct(ctx context.Context, userID int64, name string, price int64) (*domain.Product, error) {
	return s.Products.FindOrCreate(ctx, userID, name, price)
}

func (s *Service) ListProducts(ctx context.Context, userID int64) ([]domain.Product, error) {
	return s.Products.ListByUser(ctx, userID)
}

func (s *Service) GetCustomer(ctx context.Context, userID, id int64) (*domain.Customer, error) {
	return s.Customers.GetByID(ctx, userID, id)
}

func (s *Service) GetLang(ctx context.Context, userID int64) (string, error) {
	return s.Users.GetLang(ctx, userID)
}

func (s *Service) SetLang(ctx context.Context, userID int64, lang string) error {
	return s.Users.SetLang(ctx, userID, lang)
}

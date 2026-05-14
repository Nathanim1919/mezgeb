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
	var balanceDelta int64
	var stockDelta int64

	switch tx.Type {
	case domain.TxDebt:
		balanceDelta = tx.Amount // they owe more
	case domain.TxPayment:
		balanceDelta = -tx.Amount // they paid, reduce what they owe
	case domain.TxPurchase:
		balanceDelta = tx.Amount // purchase on credit = debt
	case domain.TxSell:
		stockDelta = -tx.Quantity // sold items leave stock
	case domain.TxBuy:
		stockDelta = tx.Quantity // bought items enter stock
	case domain.TxLoan:
		balanceDelta = -tx.Amount // you borrowed, you owe them
	}

	return s.Transactions.CreateWithBalanceUpdate(ctx, tx, balanceDelta, stockDelta)
}

func (s *Service) ListTransactionsByType(ctx context.Context, userID int64, txType domain.TransactionType, limit int) ([]domain.Transaction, error) {
	return s.Transactions.ListByType(ctx, userID, txType, limit)
}

func (s *Service) GetTransaction(ctx context.Context, userID, id int64) (*domain.Transaction, error) {
	return s.Transactions.GetByID(ctx, userID, id)
}

func (s *Service) UpdateTransactionAmount(ctx context.Context, userID, txID int64, oldTx *domain.Transaction, newAmount, newQuantity int64) error {
	var balanceDelta, stockDelta int64

	switch oldTx.Type {
	case domain.TxDebt:
		balanceDelta = newAmount - oldTx.Amount // adjust: new debt - old debt
	case domain.TxPayment:
		balanceDelta = -(newAmount - oldTx.Amount) // adjust: payments reduce balance
	case domain.TxPurchase:
		balanceDelta = newAmount - oldTx.Amount
	case domain.TxSell:
		stockDelta = -(newQuantity - oldTx.Quantity) // more sold = less stock
	case domain.TxBuy:
		stockDelta = newQuantity - oldTx.Quantity // more bought = more stock
	case domain.TxLoan:
		balanceDelta = -(newAmount - oldTx.Amount) // loans reduce balance (you owe)
	}

	return s.Transactions.UpdateAmountAndQuantity(ctx, userID, txID, newAmount, newQuantity, balanceDelta, stockDelta)
}

func (s *Service) UpdateTransactionNote(ctx context.Context, userID, txID int64, note string) error {
	return s.Transactions.UpdateNote(ctx, userID, txID, note)
}

func (s *Service) DeleteTransaction(ctx context.Context, userID int64, tx *domain.Transaction) error {
	var balanceDelta, stockDelta int64

	// Reverse the original effect
	switch tx.Type {
	case domain.TxDebt:
		balanceDelta = -tx.Amount // undo: they owed more
	case domain.TxPayment:
		balanceDelta = tx.Amount // undo: they paid, add back
	case domain.TxPurchase:
		balanceDelta = -tx.Amount
	case domain.TxSell:
		stockDelta = tx.Quantity // undo: return items to stock
	case domain.TxBuy:
		stockDelta = -tx.Quantity // undo: remove items from stock
	case domain.TxLoan:
		balanceDelta = tx.Amount // undo: you no longer owe
	}

	return s.Transactions.DeleteWithRollback(ctx, userID, tx.ID, balanceDelta, stockDelta)
}

// RecordPayment creates a payment transaction with the correct balance direction.
// For borrow payments (they pay you back): balance decreases (positive → 0).
// For loan repayments (you pay them back): balance increases (negative → 0).
func (s *Service) RecordPayment(ctx context.Context, tx *domain.Transaction, isLoanRepay bool) error {
	var balanceDelta int64
	if isLoanRepay {
		balanceDelta = tx.Amount // increase balance (make less negative)
	} else {
		balanceDelta = -tx.Amount // decrease balance (make less positive)
	}
	return s.Transactions.CreateWithBalanceUpdate(ctx, tx, balanceDelta, 0)
}

func (s *Service) GetReport(ctx context.Context, userID int64, from, to time.Time) (*domain.ReportData, error) {
	return s.Reports.GetReport(ctx, userID, from, to)
}

func (s *Service) FindOrCreateCustomer(ctx context.Context, userID int64, name string) (*domain.Customer, error) {
	return s.Customers.FindOrCreate(ctx, userID, name)
}

func (s *Service) FindOrCreateProduct(ctx context.Context, userID int64, name string, price int64, stock int64) (*domain.Product, error) {
	return s.Products.FindOrCreate(ctx, userID, name, price, stock)
}

func (s *Service) ListProducts(ctx context.Context, userID int64) ([]domain.Product, error) {
	return s.Products.ListByUser(ctx, userID)
}

func (s *Service) GetProduct(ctx context.Context, userID, id int64) (*domain.Product, error) {
	return s.Products.GetByID(ctx, userID, id)
}

func (s *Service) UpdateProductPrice(ctx context.Context, userID, id int64, price int64) error {
	return s.Products.UpdatePrice(ctx, userID, id, price)
}

func (s *Service) UpdateProductStock(ctx context.Context, userID, id int64, stock int64) error {
	return s.Products.UpdateStock(ctx, userID, id, stock)
}

func (s *Service) DeleteProduct(ctx context.Context, userID, id int64) error {
	return s.Products.Delete(ctx, userID, id)
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

func (s *Service) ClearData(ctx context.Context, userID int64) error {
	return s.Users.ClearData(ctx, userID)
}

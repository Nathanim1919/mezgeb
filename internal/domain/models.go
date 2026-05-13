package domain

import (
	"fmt"
	"time"
)

type User struct {
	ID           int64
	FirstName    string
	Username     string
	LanguageCode string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Customer struct {
	ID        int64
	UserID    int64
	Name      string
	Phone     string
	Balance   int64 // positive = they owe you, negative = you owe them (in cents)
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Product struct {
	ID        int64
	UserID    int64
	Name      string
	Price     int64 // default price in cents
	CreatedAt time.Time
}

type TransactionType string

const (
	TxDebt     TransactionType = "debt"
	TxPayment  TransactionType = "payment"
	TxPurchase TransactionType = "purchase"
)

type Transaction struct {
	ID         int64
	UserID     int64
	CustomerID int64
	ProductID  *int64
	Type       TransactionType
	Amount     int64 // always positive, in cents
	Note       string
	CreatedAt  time.Time

	// Joined fields (not always populated)
	CustomerName string
	ProductName  string
}

// FormatBirr converts cents to a human-readable birr string.
func FormatBirr(cents int64) string {
	whole := cents / 100
	frac := cents % 100
	if frac == 0 {
		return fmt.Sprintf("%d birr", whole)
	}
	return fmt.Sprintf("%d.%02d birr", whole, frac)
}

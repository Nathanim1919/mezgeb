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
	Lang         string // "am" or "en"
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
	Stock     int64 // current stock quantity
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TransactionType string

const (
	TxDebt     TransactionType = "debt"
	TxPayment  TransactionType = "payment"
	TxPurchase TransactionType = "purchase"
	TxSell     TransactionType = "sell"
	TxBuy      TransactionType = "buy"
	TxLoan     TransactionType = "loan"
)

type Transaction struct {
	ID         int64
	UserID     int64
	CustomerID *int64 // nullable — sell/buy don't require a customer
	ProductID  *int64
	Type       TransactionType
	Amount     int64 // always positive, in cents
	Quantity   int64 // number of items
	Note       string
	CreatedAt  time.Time

	// Joined fields (not always populated)
	CustomerName string
	ProductName  string
}

// FormatBirr converts cents to a human-readable birr string with thousands separators.
// Example: 150000000 cents → "1,500,000 ብር", 150050 cents → "1,500.50 ብር"
func FormatBirr(cents int64, label string) string {
	whole := cents / 100
	frac := cents % 100
	if label == "" {
		label = "ብር"
	}
	wholeStr := formatWithCommas(whole)
	if frac == 0 {
		return fmt.Sprintf("%s %s", wholeStr, label)
	}
	return fmt.Sprintf("%s.%02d %s", wholeStr, frac, label)
}

// formatWithCommas adds comma separators to an integer (e.g. 1500000 → "1,500,000").
func formatWithCommas(n int64) string {
	if n < 0 {
		return "-" + formatWithCommas(-n)
	}
	s := fmt.Sprintf("%d", n)
	if len(s) <= 3 {
		return s
	}
	// Insert commas from the right
	var result []byte
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, byte(c))
	}
	return string(result)
}

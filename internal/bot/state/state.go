package state

import (
	"sync"

	"github.com/nathanim1919/mezgeb/internal/domain"
)

// Step represents where the user is in a conversation flow.
type Step int

const (
	StepIdle Step = iota

	// Transaction menu
	StepTxMenu

	// Sell flow: product → quantity → note → confirm
	StepSellProduct        // choose existing or "new product"
	StepSellNewName        // new product: name
	StepSellNewPrice       // new product: price
	StepSellNewStock       // new product: initial stock
	StepSellQuantity       // how many items to sell
	StepSellNote           // optional note
	StepSellConfirm        // confirm

	// Buy flow: product → price → quantity → note → confirm
	StepBuyProduct         // choose existing or "new product"
	StepBuyNewName         // new product: name
	StepBuyNewPrice        // new product: price
	StepBuyPrice           // buy price per unit (existing product)
	StepBuyQuantity        // how many items to buy
	StepBuyNote            // optional note
	StepBuyConfirm         // confirm

	// Borrow flow: customer → amount → product (optional) → note → confirm
	StepBorrowCustomer
	StepBorrowAmount
	StepBorrowProduct
	StepBorrowNote
	StepBorrowConfirm

	// Loan flow: person → amount → note → confirm
	StepLoanPerson
	StepLoanAmount
	StepLoanNote
	StepLoanConfirm

	// Legacy debt/payment flow (kept for compatibility)
	StepTxCustomerName
	StepTxType
	StepTxAmount
	StepTxProduct
	StepTxNote
	StepTxConfirm

	// Product flow
	StepProductMenu
	StepProductName
	StepProductPrice
	StepProductStock

	// Settings flow
	StepSettingsMenu
	StepSettingsLang
	StepClearDataConfirm
)

// Conversation holds the in-progress state for one user.
type Conversation struct {
	Step         Step
	CustomerID   int64
	Customer     string
	TxType       domain.TransactionType
	Amount       int64   // total amount in cents
	Quantity     int64   // number of items
	UnitPrice    int64   // price per unit in cents
	ProductID    *int64
	Product      string
	Note         string
	ProductPrice int64   // used in product-add flow
	ProductStock int64   // used in product-add flow
}

// Manager is a thread-safe in-memory conversation state store.
type Manager struct {
	mu    sync.RWMutex
	convs map[int64]*Conversation
}

func NewManager() *Manager {
	return &Manager{
		convs: make(map[int64]*Conversation),
	}
}

func (m *Manager) Get(userID int64) *Conversation {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if c, ok := m.convs[userID]; ok {
		return c
	}
	return &Conversation{Step: StepIdle}
}

func (m *Manager) Set(userID int64, conv *Conversation) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.convs[userID] = conv
}

func (m *Manager) Reset(userID int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.convs, userID)
}

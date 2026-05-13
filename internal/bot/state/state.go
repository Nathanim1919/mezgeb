package state

import (
	"sync"

	"github.com/nathanim1919/mezgeb/internal/domain"
)

// Step represents where the user is in a conversation flow.
type Step int

const (
	StepIdle Step = iota

	// Add Transaction flow
	StepTxCustomerName
	StepTxType
	StepTxAmount
	StepTxProduct
	StepTxNote
	StepTxConfirm

	// Product flow
	StepProductName
	StepProductPrice

	// Settings flow
	StepSettingsMenu
	StepSettingsLang
)

// Conversation holds the in-progress state for one user.
type Conversation struct {
	Step       Step
	CustomerID int64
	Customer   string
	TxType     domain.TransactionType
	Amount     int64
	ProductID  *int64
	Product    string
	Note       string
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

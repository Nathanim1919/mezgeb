package handler

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nathanim1919/mezgeb/internal/bot/keyboard"
	"github.com/nathanim1919/mezgeb/internal/bot/state"
	"github.com/nathanim1919/mezgeb/internal/domain"
	"github.com/nathanim1919/mezgeb/internal/i18n"
)

func (h *Handler) startTransaction(msg *tgbotapi.Message, m *i18n.Messages) {
	conv := &state.Conversation{Step: state.StepTxCustomerName}
	h.state.Set(msg.From.ID, conv)
	h.sendWithKeyboard(msg.Chat.ID, m.AskCustomerName, keyboard.Cancel(m))
}

const (
	maxNameLength = 100
	maxAmountCents = 10_000_000_00 // 10 million birr in cents
)

func (h *Handler) handleTxCustomerName(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
	name := strings.TrimSpace(msg.Text)
	if name == "" {
		h.send(msg.Chat.ID, m.EnterCustomerName)
		return
	}
	if len([]rune(name)) > maxNameLength {
		h.send(msg.Chat.ID, m.NameTooLong)
		return
	}

	customer, err := h.svc.FindOrCreateCustomer(ctx, msg.From.ID, name)
	if err != nil {
		h.send(msg.Chat.ID, m.ErrorGeneric)
		return
	}

	conv.CustomerID = customer.ID
	conv.Customer = customer.Name
	conv.Step = state.StepTxType
	h.state.Set(msg.From.ID, conv)

	h.sendWithKeyboard(msg.Chat.ID, fmt.Sprintf(m.AskTxType, customer.Name), keyboard.TransactionType(m))
}

func (h *Handler) handleTxType(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
	switch msg.Text {
	case m.BtnOwesMe:
		conv.TxType = domain.TxDebt
	case m.BtnPaidMe:
		conv.TxType = domain.TxPayment
	case m.BtnBoughtProduct:
		conv.TxType = domain.TxPurchase
	default:
		h.send(msg.Chat.ID, m.InvalidChoice)
		return
	}

	conv.Step = state.StepTxAmount
	h.state.Set(msg.From.ID, conv)
	h.sendWithKeyboard(msg.Chat.ID, m.AskAmount, keyboard.Cancel(m))
}

func (h *Handler) handleTxAmount(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
	text := strings.TrimSpace(msg.Text)
	text = strings.ReplaceAll(text, ",", "")

	amount, err := parseAmount(text)
	if err != nil || amount <= 0 {
		h.send(msg.Chat.ID, m.InvalidAmount)
		return
	}
	if amount > maxAmountCents {
		h.send(msg.Chat.ID, m.AmountTooLarge)
		return
	}

	conv.Amount = amount
	conv.Step = state.StepTxProduct
	h.state.Set(msg.From.ID, conv)

	products, _ := h.svc.ListProducts(ctx, msg.From.ID)
	var names []string
	for _, p := range products {
		names = append(names, p.Name)
	}
	h.sendWithKeyboard(msg.Chat.ID, m.AskProduct, keyboard.ProductChoice(m, names))
}

func (h *Handler) handleTxProduct(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
	text := strings.TrimSpace(msg.Text)

	if text != m.BtnSkip && text != "" {
		product, err := h.svc.FindOrCreateProduct(ctx, msg.From.ID, text, conv.Amount)
		if err != nil {
			h.send(msg.Chat.ID, m.ProductError)
		} else {
			conv.ProductID = &product.ID
			conv.Product = product.Name
		}
	}

	conv.Step = state.StepTxNote
	h.state.Set(msg.From.ID, conv)

	h.sendWithKeyboard(msg.Chat.ID, m.AskNote, keyboard.SkipCancel(m))
}

func (h *Handler) handleTxNote(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
	text := strings.TrimSpace(msg.Text)

	if text != m.BtnSkip && text != "" {
		if len(text) > 200 {
			text = text[:200]
		}
		conv.Note = text
	}

	conv.Step = state.StepTxConfirm
	h.state.Set(msg.From.ID, conv)

	summary := buildTxSummary(conv, m)
	h.sendWithKeyboard(msg.Chat.ID, summary, keyboard.Confirm(m))
}

func (h *Handler) handleTxConfirm(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
	if msg.Text != m.BtnConfirm {
		h.send(msg.Chat.ID, m.InvalidChoice)
		return
	}

	tx := &domain.Transaction{
		UserID:     msg.From.ID,
		CustomerID: conv.CustomerID,
		ProductID:  conv.ProductID,
		Type:       conv.TxType,
		Amount:     conv.Amount,
		Note:       conv.Note,
	}

	if err := h.svc.AddTransaction(ctx, tx); err != nil {
		h.sendWithKeyboard(msg.Chat.ID, m.TxFailed, keyboard.MainMenu(m))
		h.state.Reset(msg.From.ID)
		return
	}

	h.state.Reset(msg.From.ID)
	confirmation := buildConfirmation(conv, m)
	h.sendWithKeyboard(msg.Chat.ID, confirmation, keyboard.MainMenu(m))
}

func buildTxSummary(conv *state.Conversation, m *i18n.Messages) string {
	typeLabel := txTypeLabel(conv.TxType, m)
	s := fmt.Sprintf("%s\n\n%s\n%s\n%s",
		m.TxSummaryTitle,
		fmt.Sprintf(m.TxSummaryCustomer, escMD(conv.Customer)),
		fmt.Sprintf(m.TxSummaryType, typeLabel),
		fmt.Sprintf(m.TxSummaryAmount, domain.FormatBirr(conv.Amount, m.Birr)),
	)
	if conv.Product != "" {
		s += "\n" + fmt.Sprintf(m.TxSummaryProduct, escMD(conv.Product))
	}
	if conv.Note != "" {
		s += "\n" + fmt.Sprintf(m.TxSummaryNote, escMD(conv.Note))
	}
	s += "\n\n" + m.TxSummaryConfirm
	return s
}

func buildConfirmation(conv *state.Conversation, m *i18n.Messages) string {
	birr := domain.FormatBirr(conv.Amount, m.Birr)
	customer := escMD(conv.Customer)
	product := escMD(conv.Product)
	switch conv.TxType {
	case domain.TxDebt:
		s := fmt.Sprintf(m.TxConfirmDebt, customer, birr)
		if product != "" {
			s += " — *" + product + "*"
		}
		return s
	case domain.TxPayment:
		return fmt.Sprintf(m.TxConfirmPayment, birr, customer)
	case domain.TxPurchase:
		if product == "" {
			product = "-"
		}
		return fmt.Sprintf(m.TxConfirmPurchase, customer, product, birr)
	default:
		return m.TxConfirmGeneric
	}
}

func txTypeLabel(t domain.TransactionType, m *i18n.Messages) string {
	switch t {
	case domain.TxDebt:
		return m.BtnOwesMe
	case domain.TxPayment:
		return m.BtnPaidMe
	case domain.TxPurchase:
		return m.BtnBoughtProduct
	default:
		return string(t)
	}
}

func parseAmount(s string) (int64, error) {
	if strings.Contains(s, ".") {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return 0, err
		}
		return int64(f * 100), nil
	}
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return n * 100, nil
}

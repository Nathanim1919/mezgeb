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
)

func (h *Handler) startTransaction(msg *tgbotapi.Message) {
	conv := &state.Conversation{Step: state.StepTxCustomerName}
	h.state.Set(msg.From.ID, conv)
	h.sendWithKeyboard(msg.Chat.ID, "👤 Customer name?", keyboard.Cancel())
}

func (h *Handler) handleTxCustomerName(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation) {
	name := strings.TrimSpace(msg.Text)
	if name == "" {
		h.send(msg.Chat.ID, "Please enter a customer name.")
		return
	}

	customer, err := h.svc.FindOrCreateCustomer(ctx, msg.From.ID, name)
	if err != nil {
		h.send(msg.Chat.ID, "❌ Error saving customer. Please try again.")
		return
	}

	conv.CustomerID = customer.ID
	conv.Customer = customer.Name
	conv.Step = state.StepTxType
	h.state.Set(msg.From.ID, conv)

	h.sendWithKeyboard(msg.Chat.ID, fmt.Sprintf("Got it! *%s*\n\nWhat type of transaction?", customer.Name), keyboard.TransactionType())
}

func (h *Handler) handleTxType(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation) {
	switch msg.Text {
	case "💸 Owes Me":
		conv.TxType = domain.TxDebt
	case "💰 Paid Me":
		conv.TxType = domain.TxPayment
	case "🛒 Bought Product":
		conv.TxType = domain.TxPurchase
	default:
		h.send(msg.Chat.ID, "Please choose from the buttons below 👇")
		return
	}

	conv.Step = state.StepTxAmount
	h.state.Set(msg.From.ID, conv)
	h.sendWithKeyboard(msg.Chat.ID, "💰 How much? (in birr)", keyboard.Cancel())
}

func (h *Handler) handleTxAmount(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation) {
	text := strings.TrimSpace(msg.Text)
	text = strings.ReplaceAll(text, ",", "")

	amount, err := parseAmount(text)
	if err != nil || amount <= 0 {
		h.send(msg.Chat.ID, "Please enter a valid amount.\nExamples: `1500`, `250.50`")
		return
	}

	conv.Amount = amount
	conv.Step = state.StepTxProduct
	h.state.Set(msg.From.ID, conv)

	// Show existing products or skip
	products, _ := h.svc.ListProducts(ctx, msg.From.ID)
	var names []string
	for _, p := range products {
		names = append(names, p.Name)
	}
	h.sendWithKeyboard(msg.Chat.ID, "📦 Which product? (or skip)", keyboard.ProductChoice(names))
}

func (h *Handler) handleTxProduct(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation) {
	text := strings.TrimSpace(msg.Text)

	if text != "⏭ Skip" && text != "" {
		product, err := h.svc.FindOrCreateProduct(ctx, msg.From.ID, text, conv.Amount)
		if err != nil {
			h.send(msg.Chat.ID, "❌ Error with product. Skipping.")
		} else {
			conv.ProductID = &product.ID
			conv.Product = product.Name
		}
	}

	conv.Step = state.StepTxConfirm
	h.state.Set(msg.From.ID, conv)

	summary := h.buildTxSummary(conv)
	h.sendWithKeyboard(msg.Chat.ID, summary, keyboard.Confirm())
}

func (h *Handler) handleTxConfirm(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation) {
	if msg.Text != "✅ Confirm" {
		h.send(msg.Chat.ID, "Tap ✅ Confirm or ❌ Cancel")
		return
	}

	tx := &domain.Transaction{
		UserID:     msg.From.ID,
		CustomerID: conv.CustomerID,
		ProductID:  conv.ProductID,
		Type:       conv.TxType,
		Amount:     conv.Amount,
	}

	if err := h.svc.AddTransaction(ctx, tx); err != nil {
		h.sendWithKeyboard(msg.Chat.ID, "❌ Failed to save transaction. Please try again.", keyboard.MainMenu())
		h.state.Reset(msg.From.ID)
		return
	}

	h.state.Reset(msg.From.ID)

	confirmation := h.buildConfirmation(conv)
	h.sendWithKeyboard(msg.Chat.ID, confirmation, keyboard.MainMenu())
}

func (h *Handler) buildTxSummary(conv *state.Conversation) string {
	typeLabel := txTypeLabel(conv.TxType)
	s := fmt.Sprintf("📋 *Transaction Summary*\n\n👤 Customer: *%s*\n📝 Type: *%s*\n💰 Amount: *%s*",
		conv.Customer, typeLabel, domain.FormatBirr(conv.Amount))
	if conv.Product != "" {
		s += fmt.Sprintf("\n📦 Product: *%s*", conv.Product)
	}
	s += "\n\nConfirm?"
	return s
}

func (h *Handler) buildConfirmation(conv *state.Conversation) string {
	switch conv.TxType {
	case domain.TxDebt:
		s := fmt.Sprintf("✅ *%s* now owes you *%s*", conv.Customer, domain.FormatBirr(conv.Amount))
		if conv.Product != "" {
			s += fmt.Sprintf(" for *%s*", conv.Product)
		}
		return s
	case domain.TxPayment:
		return fmt.Sprintf("✅ Recorded *%s* payment from *%s*", domain.FormatBirr(conv.Amount), conv.Customer)
	case domain.TxPurchase:
		s := fmt.Sprintf("✅ *%s* bought", conv.Customer)
		if conv.Product != "" {
			s += fmt.Sprintf(" *%s*", conv.Product)
		}
		s += fmt.Sprintf(" for *%s*", domain.FormatBirr(conv.Amount))
		return s
	default:
		return "✅ Transaction recorded!"
	}
}

func txTypeLabel(t domain.TransactionType) string {
	switch t {
	case domain.TxDebt:
		return "💸 Owes Me"
	case domain.TxPayment:
		return "💰 Paid Me"
	case domain.TxPurchase:
		return "🛒 Bought Product"
	default:
		return string(t)
	}
}

// parseAmount parses "1500" or "250.50" into cents (int64).
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

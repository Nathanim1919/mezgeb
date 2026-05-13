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

const (
	maxNameLength  = 100
	maxAmountCents = 10_000_000_00 // 10 million birr in cents
	maxQuantity    = 1_000_000
)

// ─── Transaction Menu ───────────────────────────────────────────────

func (h *Handler) startTransaction(msg *tgbotapi.Message, m *i18n.Messages) {
	conv := &state.Conversation{Step: state.StepTxMenu}
	h.state.Set(msg.From.ID, conv)
	h.sendWithKeyboard(msg.Chat.ID, m.TxMenuTitle, keyboard.TransactionMenu(m))
}

func (h *Handler) handleTxMenu(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
	switch msg.Text {
	case m.BtnSell:
		conv.TxType = domain.TxSell
		h.startSellBuyProductChoice(ctx, msg, conv, m, m.AskSellProduct, state.StepSellProduct)
	case m.BtnBuy:
		conv.TxType = domain.TxBuy
		h.startSellBuyProductChoice(ctx, msg, conv, m, m.AskBuyProduct, state.StepBuyProduct)
	case m.BtnBorrow, m.BtnLoan:
		h.send(msg.Chat.ID, m.ComingSoon)
	default:
		h.send(msg.Chat.ID, m.InvalidChoice)
	}
}

// ─── Shared: Product Selection ──────────────────────────────────────

func (h *Handler) startSellBuyProductChoice(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages, prompt string, step state.Step) {
	products, _ := h.svc.ListProducts(ctx, msg.From.ID)
	var names []string
	for _, p := range products {
		names = append(names, p.Name)
	}
	conv.Step = step
	h.state.Set(msg.From.ID, conv)
	h.sendWithKeyboard(msg.Chat.ID, prompt, keyboard.ProductChoiceWithNew(m, names))
}

// ─── SELL FLOW ──────────────────────────────────────────────────────

func (h *Handler) handleSellProduct(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
	text := strings.TrimSpace(msg.Text)

	if text == m.BtnNewProduct {
		conv.Step = state.StepSellNewName
		h.state.Set(msg.From.ID, conv)
		h.sendWithKeyboard(msg.Chat.ID, m.ProductAskName, keyboard.Cancel(m))
		return
	}

	// User picked an existing product
	products, _ := h.svc.ListProducts(ctx, msg.From.ID)
	for _, p := range products {
		if p.Name == text {
			conv.ProductID = &p.ID
			conv.Product = p.Name
			conv.UnitPrice = p.Price
			conv.ProductStock = p.Stock
			conv.Step = state.StepSellQuantity
			h.state.Set(msg.From.ID, conv)
			h.sendWithKeyboard(msg.Chat.ID,
				fmt.Sprintf("%s\n💰 %s | 📊 %s: %d", m.AskQuantity, domain.FormatBirr(p.Price, m.Birr), m.ProductStock, p.Stock),
				keyboard.Cancel(m))
			return
		}
	}

	h.send(msg.Chat.ID, m.InvalidChoice)
}

func (h *Handler) handleSellNewName(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
	name := strings.TrimSpace(msg.Text)
	if name == "" {
		h.send(msg.Chat.ID, m.ProductAskName)
		return
	}
	if len([]rune(name)) > maxNameLength {
		h.send(msg.Chat.ID, m.NameTooLong)
		return
	}
	conv.Product = name
	conv.Step = state.StepSellNewPrice
	h.state.Set(msg.From.ID, conv)
	h.sendWithKeyboard(msg.Chat.ID, fmt.Sprintf(m.ProductAskPrice, escMD(name)), keyboard.Cancel(m))
}

func (h *Handler) handleSellNewPrice(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
	price, ok := h.parsePrice(msg, m)
	if !ok {
		return
	}
	conv.UnitPrice = price
	conv.Step = state.StepSellNewStock
	h.state.Set(msg.From.ID, conv)
	h.sendWithKeyboard(msg.Chat.ID, fmt.Sprintf(m.ProductAskStock, escMD(conv.Product)), keyboard.Cancel(m))
}

func (h *Handler) handleSellNewStock(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
	stock, ok := h.parseQuantityInput(msg, m)
	if !ok {
		return
	}

	// Create the product
	product, err := h.svc.FindOrCreateProduct(ctx, msg.From.ID, conv.Product, conv.UnitPrice, stock)
	if err != nil {
		h.sendWithKeyboard(msg.Chat.ID, m.ProductError2, keyboard.MainMenu(m))
		h.state.Reset(msg.From.ID)
		return
	}

	conv.ProductID = &product.ID
	conv.Product = product.Name
	conv.ProductStock = stock
	conv.Step = state.StepSellQuantity
	h.state.Set(msg.From.ID, conv)
	h.sendWithKeyboard(msg.Chat.ID,
		fmt.Sprintf("%s\n📊 %s: %d", m.AskQuantity, m.ProductStock, stock),
		keyboard.Cancel(m))
}

func (h *Handler) handleSellQuantity(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
	text := strings.TrimSpace(msg.Text)

	// Handle "Sell All" button
	if text == m.BtnSellAll {
		if conv.ProductStock <= 0 {
			h.sendWithKeyboard(msg.Chat.ID, fmt.Sprintf(m.NotEnoughStock, conv.ProductStock), keyboard.NotEnoughStock(m))
			return
		}
		conv.Quantity = conv.ProductStock
		conv.Amount = conv.UnitPrice * conv.Quantity
		conv.Step = state.StepSellNote
		h.state.Set(msg.From.ID, conv)
		h.sendWithKeyboard(msg.Chat.ID, m.AskNote, keyboard.SkipCancel(m))
		return
	}

	// Handle "Other Product" button — go back to product selection
	if text == m.BtnChangeProduct {
		h.startSellBuyProductChoice(ctx, msg, conv, m, m.AskSellProduct, state.StepSellProduct)
		return
	}

	qty, ok := h.parseQuantityInput(msg, m)
	if !ok {
		return
	}

	// Check stock — show helpful options if not enough
	if qty > conv.ProductStock {
		h.sendWithKeyboard(msg.Chat.ID, fmt.Sprintf(m.NotEnoughStock, conv.ProductStock), keyboard.NotEnoughStock(m))
		return
	}

	conv.Quantity = qty
	conv.Amount = conv.UnitPrice * qty
	conv.Step = state.StepSellNote
	h.state.Set(msg.From.ID, conv)
	h.sendWithKeyboard(msg.Chat.ID, m.AskNote, keyboard.SkipCancel(m))
}

func (h *Handler) handleSellNote(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
	h.handleNoteStep(msg, conv, m, state.StepSellConfirm)
}

func (h *Handler) handleSellConfirm(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
	if msg.Text != m.BtnConfirm {
		h.send(msg.Chat.ID, m.InvalidChoice)
		return
	}

	tx := &domain.Transaction{
		UserID:    msg.From.ID,
		ProductID: conv.ProductID,
		Type:      domain.TxSell,
		Amount:    conv.Amount,
		Quantity:  conv.Quantity,
		Note:      conv.Note,
	}

	if err := h.svc.AddTransaction(ctx, tx); err != nil {
		h.sendWithKeyboard(msg.Chat.ID, m.TxFailed, keyboard.MainMenu(m))
		h.state.Reset(msg.From.ID)
		return
	}

	h.state.Reset(msg.From.ID)
	h.sendWithKeyboard(msg.Chat.ID,
		fmt.Sprintf(m.SellConfirm, conv.Quantity, escMD(conv.Product), domain.FormatBirr(conv.Amount, m.Birr)),
		keyboard.MainMenu(m))
}

// ─── BUY FLOW ───────────────────────────────────────────────────────

func (h *Handler) handleBuyProduct(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
	text := strings.TrimSpace(msg.Text)

	if text == m.BtnNewProduct {
		conv.Step = state.StepBuyNewName
		h.state.Set(msg.From.ID, conv)
		h.sendWithKeyboard(msg.Chat.ID, m.ProductAskName, keyboard.Cancel(m))
		return
	}

	// User picked an existing product
	products, _ := h.svc.ListProducts(ctx, msg.From.ID)
	for _, p := range products {
		if p.Name == text {
			conv.ProductID = &p.ID
			conv.Product = p.Name
			conv.Step = state.StepBuyPrice
			h.state.Set(msg.From.ID, conv)
			h.sendWithKeyboard(msg.Chat.ID, m.AskBuyPrice, keyboard.Cancel(m))
			return
		}
	}

	h.send(msg.Chat.ID, m.InvalidChoice)
}

func (h *Handler) handleBuyNewName(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
	name := strings.TrimSpace(msg.Text)
	if name == "" {
		h.send(msg.Chat.ID, m.ProductAskName)
		return
	}
	if len([]rune(name)) > maxNameLength {
		h.send(msg.Chat.ID, m.NameTooLong)
		return
	}
	conv.Product = name
	conv.Step = state.StepBuyNewPrice
	h.state.Set(msg.From.ID, conv)
	h.sendWithKeyboard(msg.Chat.ID, m.AskBuyPrice, keyboard.Cancel(m))
}

func (h *Handler) handleBuyNewPrice(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
	price, ok := h.parsePrice(msg, m)
	if !ok {
		return
	}
	conv.UnitPrice = price
	conv.Step = state.StepBuyQuantity
	h.state.Set(msg.From.ID, conv)
	h.sendWithKeyboard(msg.Chat.ID, m.AskQuantity, keyboard.Cancel(m))
}

func (h *Handler) handleBuyPrice(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
	price, ok := h.parsePrice(msg, m)
	if !ok {
		return
	}
	conv.UnitPrice = price
	conv.Step = state.StepBuyQuantity
	h.state.Set(msg.From.ID, conv)
	h.sendWithKeyboard(msg.Chat.ID, m.AskQuantity, keyboard.Cancel(m))
}

func (h *Handler) handleBuyQuantity(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
	qty, ok := h.parseQuantityInput(msg, m)
	if !ok {
		return
	}

	conv.Quantity = qty
	conv.Amount = conv.UnitPrice * qty

	// If new product, create it now with stock = quantity bought
	if conv.ProductID == nil {
		product, err := h.svc.FindOrCreateProduct(ctx, msg.From.ID, conv.Product, conv.UnitPrice, 0)
		if err != nil {
			h.sendWithKeyboard(msg.Chat.ID, m.ProductError2, keyboard.MainMenu(m))
			h.state.Reset(msg.From.ID)
			return
		}
		conv.ProductID = &product.ID
		conv.Product = product.Name
	}

	conv.Step = state.StepBuyNote
	h.state.Set(msg.From.ID, conv)
	h.sendWithKeyboard(msg.Chat.ID, m.AskNote, keyboard.SkipCancel(m))
}

func (h *Handler) handleBuyNote(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
	h.handleNoteStep(msg, conv, m, state.StepBuyConfirm)
}

func (h *Handler) handleBuyConfirm(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
	if msg.Text != m.BtnConfirm {
		h.send(msg.Chat.ID, m.InvalidChoice)
		return
	}

	tx := &domain.Transaction{
		UserID:    msg.From.ID,
		ProductID: conv.ProductID,
		Type:      domain.TxBuy,
		Amount:    conv.Amount,
		Quantity:  conv.Quantity,
		Note:      conv.Note,
	}

	if err := h.svc.AddTransaction(ctx, tx); err != nil {
		h.sendWithKeyboard(msg.Chat.ID, m.TxFailed, keyboard.MainMenu(m))
		h.state.Reset(msg.From.ID)
		return
	}

	h.state.Reset(msg.From.ID)
	h.sendWithKeyboard(msg.Chat.ID,
		fmt.Sprintf(m.BuyConfirm, conv.Quantity, escMD(conv.Product), domain.FormatBirr(conv.Amount, m.Birr)),
		keyboard.MainMenu(m))
}

// ─── LEGACY FLOW (debt/payment/purchase — for borrow/loan later) ────

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
		product, err := h.svc.FindOrCreateProduct(ctx, msg.From.ID, text, conv.Amount, 0)
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
	h.handleNoteStep(msg, conv, m, state.StepTxConfirm)
}

func (h *Handler) handleTxConfirm(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
	if msg.Text != m.BtnConfirm {
		h.send(msg.Chat.ID, m.InvalidChoice)
		return
	}

	custID := conv.CustomerID
	tx := &domain.Transaction{
		UserID:     msg.From.ID,
		CustomerID: &custID,
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
	confirmation := buildLegacyConfirmation(conv, m)
	h.sendWithKeyboard(msg.Chat.ID, confirmation, keyboard.MainMenu(m))
}

// ─── Summary Builders ───────────────────────────────────────────────

func buildSellBuySummary(conv *state.Conversation, m *i18n.Messages) string {
	var typeLabel string
	if conv.TxType == domain.TxSell {
		typeLabel = m.BtnSell
	} else {
		typeLabel = m.BtnBuy
	}

	s := fmt.Sprintf("%s\n\n%s\n%s\n%s\n%s\n%s",
		m.TxSummaryTitle,
		fmt.Sprintf(m.TxSummaryType, typeLabel),
		fmt.Sprintf(m.TxSummaryProduct, escMD(conv.Product)),
		fmt.Sprintf(m.TxSummaryQty, conv.Quantity),
		fmt.Sprintf(m.TxSummaryUnitPrice, domain.FormatBirr(conv.UnitPrice, m.Birr)),
		fmt.Sprintf(m.TxSummaryTotal, domain.FormatBirr(conv.Amount, m.Birr)),
	)
	if conv.Note != "" {
		s += "\n" + fmt.Sprintf(m.TxSummaryNote, escMD(conv.Note))
	}
	s += "\n\n" + m.TxSummaryConfirm
	return s
}

func buildLegacyConfirmation(conv *state.Conversation, m *i18n.Messages) string {
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

// ─── Shared Helpers ─────────────────────────────────────────────────

func (h *Handler) handleNoteStep(msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages, nextStep state.Step) {
	text := strings.TrimSpace(msg.Text)

	if text != m.BtnSkip && text != "" {
		if len(text) > 200 {
			text = text[:200]
		}
		conv.Note = text
	}

	conv.Step = nextStep
	h.state.Set(msg.From.ID, conv)

	var summary string
	if conv.TxType == domain.TxSell || conv.TxType == domain.TxBuy {
		summary = buildSellBuySummary(conv, m)
	} else {
		summary = buildLegacySummary(conv, m)
	}
	h.sendWithKeyboard(msg.Chat.ID, summary, keyboard.Confirm(m))
}

func buildLegacySummary(conv *state.Conversation, m *i18n.Messages) string {
	typeLabel := legacyTypeLabel(conv.TxType, m)
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

func legacyTypeLabel(t domain.TransactionType, m *i18n.Messages) string {
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

func (h *Handler) parsePrice(msg *tgbotapi.Message, m *i18n.Messages) (int64, bool) {
	text := strings.TrimSpace(msg.Text)
	text = strings.ReplaceAll(text, ",", "")
	price, err := parseAmount(text)
	if err != nil || price <= 0 {
		h.send(msg.Chat.ID, m.InvalidPrice)
		return 0, false
	}
	if price > maxAmountCents {
		h.send(msg.Chat.ID, m.AmountTooLarge)
		return 0, false
	}
	return price, true
}

func (h *Handler) parseQuantityInput(msg *tgbotapi.Message, m *i18n.Messages) (int64, bool) {
	text := strings.TrimSpace(msg.Text)
	text = strings.ReplaceAll(text, ",", "")
	qty, err := strconv.ParseInt(text, 10, 64)
	if err != nil || qty <= 0 {
		h.send(msg.Chat.ID, m.InvalidQuantity)
		return 0, false
	}
	if qty > maxQuantity {
		h.send(msg.Chat.ID, m.InvalidQuantity)
		return 0, false
	}
	return qty, true
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

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
		conv.Step = state.StepSellMenu
		h.state.Set(msg.From.ID, conv)
		h.sendWithKeyboard(msg.Chat.ID, m.SellMenuTitle, keyboard.SellMenu(m))
	case m.BtnBuy:
		conv.Step = state.StepBuyMenu
		h.state.Set(msg.From.ID, conv)
		h.sendWithKeyboard(msg.Chat.ID, m.BuyMenuTitle, keyboard.BuyMenu(m))
	case m.BtnBorrow:
		conv.Step = state.StepBorrowMenu
		h.state.Set(msg.From.ID, conv)
		h.sendWithKeyboard(msg.Chat.ID, m.BorrowMenuTitle, keyboard.BorrowMenu(m))
	case m.BtnLoan:
		conv.Step = state.StepLoanMenu
		h.state.Set(msg.From.ID, conv)
		h.sendWithKeyboard(msg.Chat.ID, m.LoanMenuTitle, keyboard.LoanMenu(m))
	default:
		h.send(msg.Chat.ID, m.InvalidChoice)
	}
}

// ─── Sub-menu handlers ─────────────────────────────────────────────

func (h *Handler) handleSellMenu(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
	switch msg.Text {
	case m.BtnNewSell:
		conv.TxType = domain.TxSell
		h.startSellBuyProductChoice(ctx, msg, conv, m, m.AskSellProduct, state.StepSellProduct)
	case m.BtnListSells:
		h.listTransactions(ctx, msg, conv, m, domain.TxSell)
	default:
		h.send(msg.Chat.ID, m.InvalidChoice)
	}
}

func (h *Handler) handleBuyMenu(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
	switch msg.Text {
	case m.BtnNewBuy:
		conv.TxType = domain.TxBuy
		h.startSellBuyProductChoice(ctx, msg, conv, m, m.AskBuyProduct, state.StepBuyProduct)
	case m.BtnListBuys:
		h.listTransactions(ctx, msg, conv, m, domain.TxBuy)
	default:
		h.send(msg.Chat.ID, m.InvalidChoice)
	}
}

func (h *Handler) handleBorrowMenu(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
	switch msg.Text {
	case m.BtnNewBorrow:
		conv.TxType = domain.TxDebt
		conv.Step = state.StepBorrowCustomer
		h.state.Set(msg.From.ID, conv)
		h.sendWithKeyboard(msg.Chat.ID, m.AskBorrowCustomer, keyboard.Cancel(m))
	case m.BtnListBorrows:
		h.listTransactions(ctx, msg, conv, m, domain.TxDebt)
	default:
		h.send(msg.Chat.ID, m.InvalidChoice)
	}
}

func (h *Handler) handleLoanMenu(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
	switch msg.Text {
	case m.BtnNewLoan:
		conv.TxType = domain.TxLoan
		conv.Step = state.StepLoanPerson
		h.state.Set(msg.From.ID, conv)
		h.sendWithKeyboard(msg.Chat.ID, m.AskLoanPerson, keyboard.Cancel(m))
	case m.BtnListLoans:
		h.listTransactions(ctx, msg, conv, m, domain.TxLoan)
	default:
		h.send(msg.Chat.ID, m.InvalidChoice)
	}
}

func (h *Handler) listTransactions(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages, txType domain.TransactionType) {
	txns, err := h.svc.ListTransactionsByType(ctx, msg.From.ID, txType, 10)
	if err != nil {
		h.state.Reset(msg.From.ID)
		h.sendWithKeyboard(msg.Chat.ID, m.ErrorGeneric, keyboard.MainMenu(m))
		return
	}

	if len(txns) == 0 {
		h.state.Reset(msg.From.ID)
		h.sendWithKeyboard(msg.Chat.ID, m.TxListEmpty, keyboard.MainMenu(m))
		return
	}

	// Store IDs so user can pick by number
	var ids []int64
	for _, tx := range txns {
		ids = append(ids, tx.ID)
	}

	conv.ListTxIDs = ids
	conv.ListTxType = txType
	conv.Step = state.StepTxListSelect
	h.state.Set(msg.From.ID, conv)

	text := h.formatTxListTitle(txType, m) + "\n\n"
	for i, tx := range txns {
		text += h.formatTxListItem(i+1, tx, txType, m) + "\n"
	}
	text += fmt.Sprintf(m.TxListTotal, len(txns))
	text += m.TxListSelectHint

	h.sendWithKeyboard(msg.Chat.ID, text, keyboard.Cancel(m))
}

func (h *Handler) formatTxListTitle(txType domain.TransactionType, m *i18n.Messages) string {
	switch txType {
	case domain.TxSell:
		return m.SellMenuTitle
	case domain.TxBuy:
		return m.BuyMenuTitle
	case domain.TxDebt:
		return m.BorrowMenuTitle
	case domain.TxLoan:
		return m.LoanMenuTitle
	default:
		return m.TxMenuTitle
	}
}

func (h *Handler) formatTxListItem(idx int, tx domain.Transaction, txType domain.TransactionType, m *i18n.Messages) string {
	date := tx.CreatedAt.Format("02/01")
	amount := domain.FormatBirr(tx.Amount, m.Birr)

	var detail string
	switch txType {
	case domain.TxSell:
		if tx.ProductName != "" {
			detail = fmt.Sprintf("%d×%s", tx.Quantity, escMD(tx.ProductName))
		} else {
			detail = amount
		}
	case domain.TxBuy:
		if tx.ProductName != "" {
			detail = fmt.Sprintf("%d×%s", tx.Quantity, escMD(tx.ProductName))
		} else {
			detail = amount
		}
	case domain.TxDebt:
		if tx.CustomerName != "" {
			detail = escMD(tx.CustomerName)
		} else {
			detail = amount
		}
	case domain.TxLoan:
		if tx.CustomerName != "" {
			detail = escMD(tx.CustomerName)
		} else {
			detail = amount
		}
	default:
		detail = amount
	}

	return fmt.Sprintf(m.TxListItem, idx, detail, amount) + " _(" + date + ")_"
}

// ─── Transaction Edit/Delete ───────────────────────────────────────

func (h *Handler) handleTxListSelect(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
	text := strings.TrimSpace(msg.Text)
	idx, err := strconv.Atoi(text)
	if err != nil || idx < 1 || idx > len(conv.ListTxIDs) {
		h.send(msg.Chat.ID, m.InvalidChoice)
		return
	}

	txID := conv.ListTxIDs[idx-1]
	tx, err := h.svc.GetTransaction(ctx, msg.From.ID, txID)
	if err != nil {
		h.state.Reset(msg.From.ID)
		h.sendWithKeyboard(msg.Chat.ID, m.TxNotFound, keyboard.MainMenu(m))
		return
	}

	conv.SelectedTxID = tx.ID
	conv.Step = state.StepTxEditMenu
	h.state.Set(msg.From.ID, conv)

	detail := h.formatTxDetail(tx, m)
	var kb tgbotapi.ReplyKeyboardMarkup
	switch conv.ListTxType {
	case domain.TxDebt:
		kb = keyboard.TxEditMenuBorrow(m)
	case domain.TxLoan:
		kb = keyboard.TxEditMenuLoan(m)
	default:
		kb = keyboard.TxEditMenu(m)
	}
	h.sendWithKeyboard(msg.Chat.ID, fmt.Sprintf(m.TxEditMenuTitle, detail), kb)
}

func (h *Handler) handleTxEditMenu(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
	switch msg.Text {
	case m.BtnRecordPayment:
		h.startRecordPayment(ctx, msg, conv, m)
		return
	case m.BtnRecordRepay:
		h.startRecordPayment(ctx, msg, conv, m)
		return
	case m.BtnEditAmount:
		conv.Step = state.StepTxEditAmount
		h.state.Set(msg.From.ID, conv)
		// Show appropriate prompt based on tx type
		if conv.ListTxType == domain.TxSell || conv.ListTxType == domain.TxBuy {
			h.sendWithKeyboard(msg.Chat.ID, m.TxEditAskQty, keyboard.Cancel(m))
		} else {
			h.sendWithKeyboard(msg.Chat.ID, m.TxEditAskAmount, keyboard.Cancel(m))
		}
	case m.BtnEditNote:
		conv.Step = state.StepTxEditNote
		h.state.Set(msg.From.ID, conv)
		h.sendWithKeyboard(msg.Chat.ID, m.TxEditAskNote, keyboard.SkipCancel(m))
	case m.BtnDelete:
		tx, err := h.svc.GetTransaction(ctx, msg.From.ID, conv.SelectedTxID)
		if err != nil {
			h.state.Reset(msg.From.ID)
			h.sendWithKeyboard(msg.Chat.ID, m.TxNotFound, keyboard.MainMenu(m))
			return
		}
		conv.Step = state.StepTxDeleteConfirm
		h.state.Set(msg.From.ID, conv)
		detail := h.formatTxDetail(tx, m)
		h.sendWithKeyboard(msg.Chat.ID, fmt.Sprintf(m.TxDeleteConfirm, detail), keyboard.Confirm(m))
	default:
		h.send(msg.Chat.ID, m.InvalidChoice)
	}
}

func (h *Handler) handleTxEditAmount(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
	tx, err := h.svc.GetTransaction(ctx, msg.From.ID, conv.SelectedTxID)
	if err != nil {
		h.state.Reset(msg.From.ID)
		h.sendWithKeyboard(msg.Chat.ID, m.TxNotFound, keyboard.MainMenu(m))
		return
	}

	if tx.Type == domain.TxSell || tx.Type == domain.TxBuy {
		// Edit quantity — recalculate amount from unit price
		newQty, ok := h.parseQuantityInput(msg, m)
		if !ok {
			return
		}
		unitPrice := int64(0)
		if tx.Quantity > 0 {
			unitPrice = tx.Amount / tx.Quantity
		}
		newAmount := unitPrice * newQty

		if err := h.svc.UpdateTransactionAmount(ctx, msg.From.ID, tx.ID, tx, newAmount, newQty); err != nil {
			h.sendWithKeyboard(msg.Chat.ID, m.ErrorGeneric, keyboard.MainMenu(m))
			h.state.Reset(msg.From.ID)
			return
		}

		h.state.Reset(msg.From.ID)
		h.sendWithKeyboard(msg.Chat.ID, m.TxEditAmountDone, keyboard.MainMenu(m))
	} else {
		// Edit amount directly (borrow/loan)
		newAmount, ok := h.parseAmountInput(msg, m)
		if !ok {
			return
		}

		if err := h.svc.UpdateTransactionAmount(ctx, msg.From.ID, tx.ID, tx, newAmount, tx.Quantity); err != nil {
			h.sendWithKeyboard(msg.Chat.ID, m.ErrorGeneric, keyboard.MainMenu(m))
			h.state.Reset(msg.From.ID)
			return
		}

		h.state.Reset(msg.From.ID)
		h.sendWithKeyboard(msg.Chat.ID, m.TxEditAmountDone, keyboard.MainMenu(m))
	}
}

func (h *Handler) handleTxEditNote(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
	note := strings.TrimSpace(msg.Text)
	if note == m.BtnSkip {
		note = ""
	}
	if len(note) > 200 {
		note = note[:200]
	}

	if err := h.svc.UpdateTransactionNote(ctx, msg.From.ID, conv.SelectedTxID, note); err != nil {
		h.sendWithKeyboard(msg.Chat.ID, m.ErrorGeneric, keyboard.MainMenu(m))
		h.state.Reset(msg.From.ID)
		return
	}

	h.state.Reset(msg.From.ID)
	h.sendWithKeyboard(msg.Chat.ID, m.TxEditNoteDone, keyboard.MainMenu(m))
}

func (h *Handler) handleTxDeleteConfirm(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
	if msg.Text != m.BtnConfirm {
		h.send(msg.Chat.ID, m.InvalidChoice)
		return
	}

	tx, err := h.svc.GetTransaction(ctx, msg.From.ID, conv.SelectedTxID)
	if err != nil {
		h.state.Reset(msg.From.ID)
		h.sendWithKeyboard(msg.Chat.ID, m.TxNotFound, keyboard.MainMenu(m))
		return
	}

	if err := h.svc.DeleteTransaction(ctx, msg.From.ID, tx); err != nil {
		h.sendWithKeyboard(msg.Chat.ID, m.ErrorGeneric, keyboard.MainMenu(m))
		h.state.Reset(msg.From.ID)
		return
	}

	h.state.Reset(msg.From.ID)
	h.sendWithKeyboard(msg.Chat.ID, m.TxDeleteDone, keyboard.MainMenu(m))
}

// ─── Payment/Repayment ─────────────────────────────────────────────

func (h *Handler) startRecordPayment(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
	tx, err := h.svc.GetTransaction(ctx, msg.From.ID, conv.SelectedTxID)
	if err != nil || tx.CustomerID == nil {
		h.state.Reset(msg.From.ID)
		h.sendWithKeyboard(msg.Chat.ID, m.TxNotFound, keyboard.MainMenu(m))
		return
	}

	customer, err := h.svc.GetCustomer(ctx, msg.From.ID, *tx.CustomerID)
	if err != nil {
		h.state.Reset(msg.From.ID)
		h.sendWithKeyboard(msg.Chat.ID, m.ErrorGeneric, keyboard.MainMenu(m))
		return
	}

	// Store the customer info and outstanding balance for the payment step
	conv.CustomerID = customer.ID
	conv.Customer = customer.Name
	// For borrow (TxDebt): positive balance = they owe you
	// For loan (TxLoan): negative balance = you owe them
	if tx.Type == domain.TxDebt {
		conv.Amount = customer.Balance // positive: what they owe
	} else {
		conv.Amount = -customer.Balance // flip: negative balance → positive amount
	}
	conv.Step = state.StepTxRecordPayment
	h.state.Set(msg.From.ID, conv)

	outstanding := domain.FormatBirr(conv.Amount, m.Birr)
	if tx.Type == domain.TxDebt {
		h.sendWithKeyboard(msg.Chat.ID,
			fmt.Sprintf(m.AskPaymentAmount, escMD(customer.Name), outstanding),
			keyboard.PaymentAmount(m))
	} else {
		h.sendWithKeyboard(msg.Chat.ID,
			fmt.Sprintf(m.AskRepayAmount, escMD(customer.Name), outstanding),
			keyboard.PaymentAmount(m))
	}
}

func (h *Handler) handleRecordPayment(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
	text := strings.TrimSpace(msg.Text)

	var payAmount int64

	if text == m.BtnPayAll {
		payAmount = conv.Amount
	} else {
		amount, ok := h.parseAmountInput(msg, m)
		if !ok {
			return
		}
		payAmount = amount
	}

	if payAmount <= 0 {
		h.send(msg.Chat.ID, m.InvalidAmount)
		return
	}

	// Don't allow paying more than owed
	if payAmount > conv.Amount {
		h.sendWithKeyboard(msg.Chat.ID,
			fmt.Sprintf(m.AmountExceedsDebt, domain.FormatBirr(conv.Amount, m.Birr)),
			keyboard.PaymentAmount(m))
		return
	}

	custID := conv.CustomerID
	tx := &domain.Transaction{
		UserID:     msg.From.ID,
		CustomerID: &custID,
		Type:       domain.TxPayment,
		Amount:     payAmount,
		Note:       "",
	}

	isLoanRepay := conv.ListTxType == domain.TxLoan
	if err := h.svc.RecordPayment(ctx, tx, isLoanRepay); err != nil {
		h.sendWithKeyboard(msg.Chat.ID, m.ErrorGeneric, keyboard.MainMenu(m))
		h.state.Reset(msg.From.ID)
		return
	}

	h.state.Reset(msg.From.ID)

	remaining := conv.Amount - payAmount
	if remaining <= 0 {
		h.sendWithKeyboard(msg.Chat.ID, m.PaymentSettled, keyboard.MainMenu(m))
		return
	}

	payStr := domain.FormatBirr(payAmount, m.Birr)
	remStr := domain.FormatBirr(remaining, m.Birr)
	if conv.ListTxType == domain.TxDebt {
		h.sendWithKeyboard(msg.Chat.ID,
			fmt.Sprintf(m.PaymentDone, payStr, escMD(conv.Customer), remStr),
			keyboard.MainMenu(m))
	} else {
		h.sendWithKeyboard(msg.Chat.ID,
			fmt.Sprintf(m.RepayDone, payStr, escMD(conv.Customer), remStr),
			keyboard.MainMenu(m))
	}
}

func (h *Handler) formatTxDetail(tx *domain.Transaction, m *i18n.Messages) string {
	amount := domain.FormatBirr(tx.Amount, m.Birr)
	date := tx.CreatedAt.Format("02/01/2006")

	var lines []string

	switch tx.Type {
	case domain.TxSell, domain.TxBuy:
		var typeLabel string
		if tx.Type == domain.TxSell {
			typeLabel = m.BtnSell
		} else {
			typeLabel = m.BtnBuy
		}
		lines = append(lines, fmt.Sprintf(m.TxSummaryType, typeLabel))
		if tx.ProductName != "" {
			lines = append(lines, fmt.Sprintf(m.TxSummaryProduct, escMD(tx.ProductName)))
		}
		if tx.Quantity > 0 {
			lines = append(lines, fmt.Sprintf(m.TxSummaryQty, tx.Quantity))
		}
		lines = append(lines, fmt.Sprintf(m.TxSummaryTotal, amount))
	case domain.TxDebt, domain.TxLoan:
		var typeLabel string
		if tx.Type == domain.TxDebt {
			typeLabel = m.BtnBorrow
		} else {
			typeLabel = m.BtnLoan
		}
		lines = append(lines, fmt.Sprintf(m.TxSummaryType, typeLabel))
		if tx.CustomerName != "" {
			lines = append(lines, fmt.Sprintf(m.TxSummaryCustomer, escMD(tx.CustomerName)))
		}
		lines = append(lines, fmt.Sprintf(m.TxSummaryAmount, amount))
	default:
		lines = append(lines, fmt.Sprintf(m.TxSummaryAmount, amount))
	}

	if tx.Note != "" {
		lines = append(lines, fmt.Sprintf(m.TxSummaryNote, escMD(tx.Note)))
	}
	lines = append(lines, "📅 "+date)

	return strings.Join(lines, "\n")
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

// ─── BORROW FLOW ────────────────────────────────────────────────────
// Someone borrows from you → they owe you. Customer balance increases.

func (h *Handler) handleBorrowCustomer(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
	name := strings.TrimSpace(msg.Text)
	if name == "" {
		h.send(msg.Chat.ID, m.AskBorrowCustomer)
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
	conv.Step = state.StepBorrowAmount
	h.state.Set(msg.From.ID, conv)
	h.sendWithKeyboard(msg.Chat.ID, m.AskBorrowAmount, keyboard.Cancel(m))
}

func (h *Handler) handleBorrowAmount(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
	amount, ok := h.parseAmountInput(msg, m)
	if !ok {
		return
	}

	conv.Amount = amount
	conv.Step = state.StepBorrowProduct
	h.state.Set(msg.From.ID, conv)

	products, _ := h.svc.ListProducts(ctx, msg.From.ID)
	var names []string
	for _, p := range products {
		names = append(names, p.Name)
	}
	h.sendWithKeyboard(msg.Chat.ID, m.AskBorrowProduct, keyboard.ProductChoice(m, names))
}

func (h *Handler) handleBorrowProduct(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
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

	conv.Step = state.StepBorrowNote
	h.state.Set(msg.From.ID, conv)
	h.sendWithKeyboard(msg.Chat.ID, m.AskNote, keyboard.SkipCancel(m))
}

func (h *Handler) handleBorrowNote(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
	h.handleNoteStep(msg, conv, m, state.StepBorrowConfirm)
}

func (h *Handler) handleBorrowConfirm(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
	if msg.Text != m.BtnConfirm {
		h.send(msg.Chat.ID, m.InvalidChoice)
		return
	}

	custID := conv.CustomerID
	tx := &domain.Transaction{
		UserID:     msg.From.ID,
		CustomerID: &custID,
		ProductID:  conv.ProductID,
		Type:       domain.TxDebt,
		Amount:     conv.Amount,
		Note:       conv.Note,
	}

	if err := h.svc.AddTransaction(ctx, tx); err != nil {
		h.sendWithKeyboard(msg.Chat.ID, m.TxFailed, keyboard.MainMenu(m))
		h.state.Reset(msg.From.ID)
		return
	}

	h.state.Reset(msg.From.ID)
	h.sendWithKeyboard(msg.Chat.ID,
		fmt.Sprintf(m.BorrowConfirm, escMD(conv.Customer), domain.FormatBirr(conv.Amount, m.Birr)),
		keyboard.MainMenu(m))
}

// ─── LOAN FLOW ──────────────────────────────────────────────────────
// You borrow from someone → you owe them. Customer balance decreases (goes negative).

func (h *Handler) handleLoanPerson(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
	name := strings.TrimSpace(msg.Text)
	if name == "" {
		h.send(msg.Chat.ID, m.AskLoanPerson)
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
	conv.Step = state.StepLoanAmount
	h.state.Set(msg.From.ID, conv)
	h.sendWithKeyboard(msg.Chat.ID, m.AskLoanAmount, keyboard.Cancel(m))
}

func (h *Handler) handleLoanAmount(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
	amount, ok := h.parseAmountInput(msg, m)
	if !ok {
		return
	}

	conv.Amount = amount
	conv.Step = state.StepLoanNote
	h.state.Set(msg.From.ID, conv)
	h.sendWithKeyboard(msg.Chat.ID, m.AskNote, keyboard.SkipCancel(m))
}

func (h *Handler) handleLoanNote(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
	h.handleNoteStep(msg, conv, m, state.StepLoanConfirm)
}

func (h *Handler) handleLoanConfirm(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
	if msg.Text != m.BtnConfirm {
		h.send(msg.Chat.ID, m.InvalidChoice)
		return
	}

	custID := conv.CustomerID
	tx := &domain.Transaction{
		UserID:     msg.From.ID,
		CustomerID: &custID,
		Type:       domain.TxLoan,
		Amount:     conv.Amount,
		Note:       conv.Note,
	}

	if err := h.svc.AddTransaction(ctx, tx); err != nil {
		h.sendWithKeyboard(msg.Chat.ID, m.TxFailed, keyboard.MainMenu(m))
		h.state.Reset(msg.From.ID)
		return
	}

	h.state.Reset(msg.From.ID)
	h.sendWithKeyboard(msg.Chat.ID,
		fmt.Sprintf(m.LoanConfirm, domain.FormatBirr(conv.Amount, m.Birr), escMD(conv.Customer)),
		keyboard.MainMenu(m))
}

// ─── LEGACY FLOW (kept for compatibility) ───────────────────────────

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
	switch conv.TxType {
	case domain.TxSell, domain.TxBuy:
		summary = buildSellBuySummary(conv, m)
	case domain.TxDebt:
		summary = buildBorrowSummary(conv, m)
	case domain.TxLoan:
		summary = buildLoanSummary(conv, m)
	default:
		summary = buildLegacySummary(conv, m)
	}
	h.sendWithKeyboard(msg.Chat.ID, summary, keyboard.Confirm(m))
}

func buildBorrowSummary(conv *state.Conversation, m *i18n.Messages) string {
	s := fmt.Sprintf("%s\n\n%s\n%s\n%s",
		m.BorrowSummaryTitle,
		fmt.Sprintf(m.TxSummaryCustomer, escMD(conv.Customer)),
		fmt.Sprintf(m.TxSummaryType, m.BtnBorrow),
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

func buildLoanSummary(conv *state.Conversation, m *i18n.Messages) string {
	s := fmt.Sprintf("%s\n\n%s\n%s\n%s",
		m.LoanSummaryTitle,
		fmt.Sprintf(m.TxSummaryCustomer, escMD(conv.Customer)),
		fmt.Sprintf(m.TxSummaryType, m.BtnLoan),
		fmt.Sprintf(m.TxSummaryAmount, domain.FormatBirr(conv.Amount, m.Birr)),
	)
	if conv.Note != "" {
		s += "\n" + fmt.Sprintf(m.TxSummaryNote, escMD(conv.Note))
	}
	s += "\n\n" + m.TxSummaryConfirm
	return s
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

func (h *Handler) parseAmountInput(msg *tgbotapi.Message, m *i18n.Messages) (int64, bool) {
	text := strings.TrimSpace(msg.Text)
	text = strings.ReplaceAll(text, ",", "")
	amount, err := parseAmount(text)
	if err != nil || amount <= 0 {
		h.send(msg.Chat.ID, m.InvalidAmount)
		return 0, false
	}
	if amount > maxAmountCents {
		h.send(msg.Chat.ID, m.AmountTooLarge)
		return 0, false
	}
	return amount, true
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

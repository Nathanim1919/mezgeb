package handler

import (
	"context"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nathanim1919/mezgeb/internal/bot/keyboard"
	"github.com/nathanim1919/mezgeb/internal/bot/ratelimit"
	"github.com/nathanim1919/mezgeb/internal/bot/state"
	"github.com/nathanim1919/mezgeb/internal/domain"
	"github.com/nathanim1919/mezgeb/internal/i18n"
	"github.com/nathanim1919/mezgeb/internal/service"
)

type Handler struct {
	bot     *tgbotapi.BotAPI
	svc     *service.Service
	state   *state.Manager
	limiter *ratelimit.Limiter
}

func New(bot *tgbotapi.BotAPI, svc *service.Service, stateMgr *state.Manager, limiter *ratelimit.Limiter) *Handler {
	return &Handler{bot: bot, svc: svc, state: stateMgr, limiter: limiter}
}

func (h *Handler) HandleUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	msg := update.Message
	userID := msg.From.ID
	ctx := context.Background()

	// Rate limit: 30 messages per minute per user
	if !h.limiter.Allow(userID) {
		m := h.getMessages(ctx, userID)
		h.send(msg.Chat.ID, m.RateLimited)
		return
	}

	// Ensure user exists in DB
	if err := h.svc.EnsureUser(ctx, &domain.User{
		ID:           userID,
		FirstName:    msg.From.FirstName,
		Username:     msg.From.UserName,
		LanguageCode: msg.From.LanguageCode,
	}); err != nil {
		log.Printf("error upserting user: %v", err)
	}

	// Load user's language
	m := h.getMessages(ctx, userID)
	text := msg.Text

	// Global commands — check both languages for cancel
	if text == "/start" || text == "❌ Cancel" || text == "❌ ሰርዝ" {
		h.state.Reset(userID)
		h.sendWithKeyboard(msg.Chat.ID, m.Welcome, keyboard.MainMenu(m))
		return
	}

	// Check if user is in a conversation flow
	conv := h.state.Get(userID)
	if conv.Step != state.StepIdle {
		h.handleConversation(ctx, msg, conv, m)
		return
	}

	// Main menu routing — match against current language buttons
	switch text {
	case m.BtnAddTx:
		h.startTransaction(msg, m)
	case m.BtnReports:
		h.startReport(msg, m)
	case m.BtnToday, m.BtnThisWeek, m.BtnThisMonth:
		h.handleReportPeriod(ctx, msg, m)
	case m.BtnProducts:
		h.showProducts(ctx, msg, m)
	case m.BtnSettings:
		h.showSettings(msg, m)
	default:
		// Also try matching the OTHER language's buttons (user may have switched)
		h.sendWithKeyboard(msg.Chat.ID, m.NotUnderstood, keyboard.MainMenu(m))
	}
}

func (h *Handler) handleConversation(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
	switch conv.Step {
	// Transaction menu
	case state.StepTxMenu:
		h.handleTxMenu(ctx, msg, conv, m)

	// Transaction sub-menus
	case state.StepSellMenu:
		h.handleSellMenu(ctx, msg, conv, m)
	case state.StepBuyMenu:
		h.handleBuyMenu(ctx, msg, conv, m)
	case state.StepBorrowMenu:
		h.handleBorrowMenu(ctx, msg, conv, m)
	case state.StepLoanMenu:
		h.handleLoanMenu(ctx, msg, conv, m)

	// Transaction list → edit/delete
	case state.StepTxListSelect:
		h.handleTxListSelect(ctx, msg, conv, m)
	case state.StepTxEditMenu:
		h.handleTxEditMenu(ctx, msg, conv, m)
	case state.StepTxEditAmount:
		h.handleTxEditAmount(ctx, msg, conv, m)
	case state.StepTxEditNote:
		h.handleTxEditNote(ctx, msg, conv, m)
	case state.StepTxDeleteConfirm:
		h.handleTxDeleteConfirm(ctx, msg, conv, m)
	case state.StepTxRecordPayment:
		h.handleRecordPayment(ctx, msg, conv, m)

	// Sell flow
	case state.StepSellProduct:
		h.handleSellProduct(ctx, msg, conv, m)
	case state.StepSellNewName:
		h.handleSellNewName(ctx, msg, conv, m)
	case state.StepSellNewPrice:
		h.handleSellNewPrice(ctx, msg, conv, m)
	case state.StepSellNewStock:
		h.handleSellNewStock(ctx, msg, conv, m)
	case state.StepSellQuantity:
		h.handleSellQuantity(ctx, msg, conv, m)
	case state.StepSellNote:
		h.handleSellNote(ctx, msg, conv, m)
	case state.StepSellConfirm:
		h.handleSellConfirm(ctx, msg, conv, m)

	// Buy flow
	case state.StepBuyProduct:
		h.handleBuyProduct(ctx, msg, conv, m)
	case state.StepBuyNewName:
		h.handleBuyNewName(ctx, msg, conv, m)
	case state.StepBuyNewPrice:
		h.handleBuyNewPrice(ctx, msg, conv, m)
	case state.StepBuyPrice:
		h.handleBuyPrice(ctx, msg, conv, m)
	case state.StepBuyQuantity:
		h.handleBuyQuantity(ctx, msg, conv, m)
	case state.StepBuyNote:
		h.handleBuyNote(ctx, msg, conv, m)
	case state.StepBuyConfirm:
		h.handleBuyConfirm(ctx, msg, conv, m)

	// Borrow flow
	case state.StepBorrowCustomer:
		h.handleBorrowCustomer(ctx, msg, conv, m)
	case state.StepBorrowAmount:
		h.handleBorrowAmount(ctx, msg, conv, m)
	case state.StepBorrowProduct:
		h.handleBorrowProduct(ctx, msg, conv, m)
	case state.StepBorrowNote:
		h.handleBorrowNote(ctx, msg, conv, m)
	case state.StepBorrowConfirm:
		h.handleBorrowConfirm(ctx, msg, conv, m)

	// Loan flow
	case state.StepLoanPerson:
		h.handleLoanPerson(ctx, msg, conv, m)
	case state.StepLoanAmount:
		h.handleLoanAmount(ctx, msg, conv, m)
	case state.StepLoanNote:
		h.handleLoanNote(ctx, msg, conv, m)
	case state.StepLoanConfirm:
		h.handleLoanConfirm(ctx, msg, conv, m)

	// Legacy debt/payment flow
	case state.StepTxCustomerName:
		h.handleTxCustomerName(ctx, msg, conv, m)
	case state.StepTxType:
		h.handleTxType(ctx, msg, conv, m)
	case state.StepTxAmount:
		h.handleTxAmount(ctx, msg, conv, m)
	case state.StepTxProduct:
		h.handleTxProduct(ctx, msg, conv, m)
	case state.StepTxNote:
		h.handleTxNote(ctx, msg, conv, m)
	case state.StepTxConfirm:
		h.handleTxConfirm(ctx, msg, conv, m)

	// Product management
	case state.StepProductMenu:
		h.handleProductMenu(ctx, msg, conv, m)
	case state.StepProductName:
		h.handleProductName(ctx, msg, conv, m)
	case state.StepProductPrice:
		h.handleProductPrice(ctx, msg, conv, m)
	case state.StepProductStock:
		h.handleProductStock(ctx, msg, conv, m)

	// Product list → edit/delete
	case state.StepProductListSelect:
		h.handleProductListSelect(ctx, msg, conv, m)
	case state.StepProductEditMenu:
		h.handleProductEditMenu(ctx, msg, conv, m)
	case state.StepProductEditPrice:
		h.handleProductEditPrice(ctx, msg, conv, m)
	case state.StepProductEditStock:
		h.handleProductEditStock(ctx, msg, conv, m)
	case state.StepProductDeleteConfirm:
		h.handleProductDeleteConfirm(ctx, msg, conv, m)

	// Settings
	case state.StepSettingsMenu:
		h.handleSettingsMenu(ctx, msg, conv, m)
	case state.StepSettingsLang:
		h.handleSettingsLang(ctx, msg, conv, m)
	case state.StepClearDataConfirm:
		h.handleClearDataConfirm(ctx, msg, conv, m)
	}
}

func (h *Handler) getMessages(ctx context.Context, userID int64) *i18n.Messages {
	langStr, err := h.svc.GetLang(ctx, userID)
	if err != nil {
		langStr = "am"
	}
	return i18n.Get(i18n.Parse(langStr))
}

func (h *Handler) send(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	if _, err := h.bot.Send(msg); err != nil {
		log.Printf("error sending message: %v", err)
	}
}

func (h *Handler) sendWithKeyboard(chatID int64, text string, kb tgbotapi.ReplyKeyboardMarkup) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = kb
	if _, err := h.bot.Send(msg); err != nil {
		log.Printf("error sending message: %v", err)
	}
}

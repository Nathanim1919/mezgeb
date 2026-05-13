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
	case m.BtnCustomers:
		h.showCustomers(ctx, msg, m)
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
	case state.StepTxCustomerName:
		h.handleTxCustomerName(ctx, msg, conv, m)
	case state.StepTxType:
		h.handleTxType(ctx, msg, conv, m)
	case state.StepTxAmount:
		h.handleTxAmount(ctx, msg, conv, m)
	case state.StepTxProduct:
		h.handleTxProduct(ctx, msg, conv, m)
	case state.StepTxConfirm:
		h.handleTxConfirm(ctx, msg, conv, m)
	case state.StepProductName:
		h.handleProductName(ctx, msg, conv, m)
	case state.StepProductPrice:
		h.handleProductPrice(ctx, msg, conv, m)
	case state.StepSettingsMenu:
		h.handleSettingsMenu(ctx, msg, conv, m)
	case state.StepSettingsLang:
		h.handleSettingsLang(ctx, msg, conv, m)
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

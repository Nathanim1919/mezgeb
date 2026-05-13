package handler

import (
	"context"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nathanim1919/mezgeb/internal/bot/keyboard"
	"github.com/nathanim1919/mezgeb/internal/bot/state"
	"github.com/nathanim1919/mezgeb/internal/domain"
	"github.com/nathanim1919/mezgeb/internal/service"
)

type Handler struct {
	bot   *tgbotapi.BotAPI
	svc   *service.Service
	state *state.Manager
}

func New(bot *tgbotapi.BotAPI, svc *service.Service, stateMgr *state.Manager) *Handler {
	return &Handler{bot: bot, svc: svc, state: stateMgr}
}

func (h *Handler) HandleUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	msg := update.Message
	userID := msg.From.ID
	ctx := context.Background()

	// Ensure user exists in DB
	if err := h.svc.EnsureUser(ctx, &domain.User{
		ID:           userID,
		FirstName:    msg.From.FirstName,
		Username:     msg.From.UserName,
		LanguageCode: msg.From.LanguageCode,
	}); err != nil {
		log.Printf("error upserting user: %v", err)
	}

	text := msg.Text

	// Global commands
	if text == "/start" || text == "❌ Cancel" {
		h.state.Reset(userID)
		h.sendWithKeyboard(msg.Chat.ID, "Welcome to Mezgeb! 📒\n\nYour simple business assistant.\nWhat would you like to do?", keyboard.MainMenu())
		return
	}

	// Check if user is in a conversation flow
	conv := h.state.Get(userID)
	if conv.Step != state.StepIdle {
		h.handleConversation(ctx, msg, conv)
		return
	}

	// Main menu routing
	switch text {
	case "➕ Add Transaction":
		h.startTransaction(msg)
	case "📊 Reports":
		h.startReport(msg)
	case "📅 Today", "📆 This Week", "🗓 This Month":
		h.handleReportPeriod(ctx, msg)
	case "👥 Customers":
		h.showCustomers(ctx, msg)
	case "📦 Products":
		h.showProducts(ctx, msg)
	case "⚙️ Settings":
		h.send(msg.Chat.ID, "⚙️ Settings coming soon!\n\nFor now, just use the menu below.")
	default:
		h.sendWithKeyboard(msg.Chat.ID, "I didn't understand that. Use the menu below 👇", keyboard.MainMenu())
	}
}

func (h *Handler) handleConversation(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation) {
	switch conv.Step {
	case state.StepTxCustomerName:
		h.handleTxCustomerName(ctx, msg, conv)
	case state.StepTxType:
		h.handleTxType(ctx, msg, conv)
	case state.StepTxAmount:
		h.handleTxAmount(ctx, msg, conv)
	case state.StepTxProduct:
		h.handleTxProduct(ctx, msg, conv)
	case state.StepTxConfirm:
		h.handleTxConfirm(ctx, msg, conv)
	case state.StepProductName:
		h.handleProductName(ctx, msg, conv)
	case state.StepProductPrice:
		h.handleProductPrice(ctx, msg, conv)
	}
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

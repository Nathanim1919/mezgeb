package handler

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nathanim1919/mezgeb/internal/bot/keyboard"
	"github.com/nathanim1919/mezgeb/internal/bot/state"
	"github.com/nathanim1919/mezgeb/internal/i18n"
)

func (h *Handler) showSettings(msg *tgbotapi.Message, m *i18n.Messages) {
	conv := &state.Conversation{Step: state.StepSettingsMenu}
	h.state.Set(msg.From.ID, conv)
	h.sendWithKeyboard(msg.Chat.ID, m.SettingsTitle, keyboard.Settings(m))
}

func (h *Handler) handleSettingsMenu(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
	switch msg.Text {
	case m.BtnLanguage:
		conv.Step = state.StepSettingsLang
		h.state.Set(msg.From.ID, conv)
		h.sendWithKeyboard(msg.Chat.ID, m.ChooseLanguage, keyboard.LanguageChoice())
	default:
		h.send(msg.Chat.ID, m.InvalidChoice)
	}
}

func (h *Handler) handleSettingsLang(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
	var lang i18n.Lang

	switch msg.Text {
	case "🇪🇹 አማርኛ":
		lang = i18n.Amharic
	case "🇬🇧 English":
		lang = i18n.English
	default:
		h.send(msg.Chat.ID, m.InvalidChoice)
		return
	}

	if err := h.svc.SetLang(ctx, msg.From.ID, string(lang)); err != nil {
		h.send(msg.Chat.ID, m.ErrorGeneric)
		return
	}

	h.state.Reset(msg.From.ID)

	// Use the NEW language for the confirmation
	newM := i18n.Get(lang)
	h.sendWithKeyboard(msg.Chat.ID, newM.LanguageUpdated, keyboard.MainMenu(newM))
}

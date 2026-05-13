package handler

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nathanim1919/mezgeb/internal/bot/keyboard"
	"github.com/nathanim1919/mezgeb/internal/domain"
	"github.com/nathanim1919/mezgeb/internal/i18n"
)

func (h *Handler) showCustomers(ctx context.Context, msg *tgbotapi.Message, m *i18n.Messages) {
	customers, err := h.svc.ListCustomers(ctx, msg.From.ID)
	if err != nil {
		h.sendWithKeyboard(msg.Chat.ID, m.ErrorGeneric, keyboard.MainMenu(m))
		return
	}

	if len(customers) == 0 {
		h.sendWithKeyboard(msg.Chat.ID, m.CustomersEmpty, keyboard.MainMenu(m))
		return
	}

	text := m.CustomersTitle + "\n\n"
	for i, c := range customers {
		balanceStr := formatBalance(c.Balance, m)
		text += fmt.Sprintf("%d. *%s* — %s\n", i+1, escMD(c.Name), balanceStr)
	}

	text += "\n" + fmt.Sprintf(m.CustomersTotal, len(customers))
	h.sendWithKeyboard(msg.Chat.ID, text, keyboard.MainMenu(m))
}

func formatBalance(cents int64, m *i18n.Messages) string {
	if cents > 0 {
		return fmt.Sprintf(m.CustomerOwes, domain.FormatBirr(cents, m.Birr))
	} else if cents < 0 {
		return fmt.Sprintf(m.CustomerYouOwe, domain.FormatBirr(-cents, m.Birr))
	}
	return m.CustomerSettled
}

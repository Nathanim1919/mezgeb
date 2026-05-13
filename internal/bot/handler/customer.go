package handler

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nathanim1919/mezgeb/internal/bot/keyboard"
	"github.com/nathanim1919/mezgeb/internal/domain"
)

func (h *Handler) showCustomers(ctx context.Context, msg *tgbotapi.Message) {
	customers, err := h.svc.ListCustomers(ctx, msg.From.ID)
	if err != nil {
		h.sendWithKeyboard(msg.Chat.ID, "❌ Error loading customers.", keyboard.MainMenu())
		return
	}

	if len(customers) == 0 {
		h.sendWithKeyboard(msg.Chat.ID, "👥 No customers yet.\n\nAdd your first transaction to create a customer!", keyboard.MainMenu())
		return
	}

	text := "👥 *Your Customers*\n\n"
	for i, c := range customers {
		balanceStr := formatBalance(c.Balance)
		text += fmt.Sprintf("%d. *%s* — %s\n", i+1, c.Name, balanceStr)
	}

	text += fmt.Sprintf("\n_Total: %d customers_", len(customers))
	h.sendWithKeyboard(msg.Chat.ID, text, keyboard.MainMenu())
}

func formatBalance(cents int64) string {
	if cents > 0 {
		return fmt.Sprintf("owes you %s", domain.FormatBirr(cents))
	} else if cents < 0 {
		return fmt.Sprintf("you owe %s", domain.FormatBirr(-cents))
	}
	return "settled ✓"
}

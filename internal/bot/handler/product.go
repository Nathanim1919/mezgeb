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

func (h *Handler) showProducts(ctx context.Context, msg *tgbotapi.Message) {
	products, err := h.svc.ListProducts(ctx, msg.From.ID)
	if err != nil {
		h.sendWithKeyboard(msg.Chat.ID, "❌ Error loading products.", keyboard.MainMenu())
		return
	}

	text := "📦 *Your Products*\n\n"
	if len(products) == 0 {
		text += "_No products yet._\n"
	} else {
		for i, p := range products {
			text += fmt.Sprintf("%d. *%s* — %s\n", i+1, p.Name, domain.FormatBirr(p.Price))
		}
	}

	text += "\nTo add a product, type its name below or tap a menu button."

	conv := &state.Conversation{Step: state.StepProductName}
	h.state.Set(msg.From.ID, conv)

	h.sendWithKeyboard(msg.Chat.ID, text, keyboard.Cancel())
}

func (h *Handler) handleProductName(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation) {
	name := strings.TrimSpace(msg.Text)
	if name == "" {
		h.send(msg.Chat.ID, "Please enter a product name.")
		return
	}

	conv.Product = name
	conv.Step = state.StepProductPrice
	h.state.Set(msg.From.ID, conv)

	h.sendWithKeyboard(msg.Chat.ID, fmt.Sprintf("💰 Default price for *%s*? (in birr)", name), keyboard.Cancel())
}

func (h *Handler) handleProductPrice(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation) {
	text := strings.TrimSpace(msg.Text)
	text = strings.ReplaceAll(text, ",", "")

	price, err := strconv.ParseInt(text, 10, 64)
	if err != nil || price < 0 {
		h.send(msg.Chat.ID, "Please enter a valid price. Example: `500`")
		return
	}

	priceCents := price * 100
	product, err := h.svc.FindOrCreateProduct(ctx, msg.From.ID, conv.Product, priceCents)
	if err != nil {
		h.sendWithKeyboard(msg.Chat.ID, "❌ Error saving product.", keyboard.MainMenu())
		h.state.Reset(msg.From.ID)
		return
	}

	h.state.Reset(msg.From.ID)
	h.sendWithKeyboard(msg.Chat.ID,
		fmt.Sprintf("✅ Product *%s* saved at *%s*", product.Name, domain.FormatBirr(priceCents)),
		keyboard.MainMenu())
}

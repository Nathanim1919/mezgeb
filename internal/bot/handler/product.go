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

func (h *Handler) showProducts(ctx context.Context, msg *tgbotapi.Message, m *i18n.Messages) {
	products, err := h.svc.ListProducts(ctx, msg.From.ID)
	if err != nil {
		h.sendWithKeyboard(msg.Chat.ID, m.ErrorGeneric, keyboard.MainMenu(m))
		return
	}

	text := m.ProductsTitle + "\n\n"
	if len(products) == 0 {
		text += m.ProductsEmpty + "\n"
	} else {
		for i, p := range products {
			text += fmt.Sprintf("%d. *%s* — %s\n", i+1, escMD(p.Name), domain.FormatBirr(p.Price, m.Birr))
		}
	}

	text += "\n" + m.ProductAskName

	conv := &state.Conversation{Step: state.StepProductName}
	h.state.Set(msg.From.ID, conv)

	h.sendWithKeyboard(msg.Chat.ID, text, keyboard.Cancel(m))
}

func (h *Handler) handleProductName(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
	name := strings.TrimSpace(msg.Text)
	if name == "" {
		h.send(msg.Chat.ID, m.EnterCustomerName)
		return
	}

	conv.Product = name
	conv.Step = state.StepProductPrice
	h.state.Set(msg.From.ID, conv)

	h.sendWithKeyboard(msg.Chat.ID, fmt.Sprintf(m.ProductAskPrice, name), keyboard.Cancel(m))
}

func (h *Handler) handleProductPrice(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
	text := strings.TrimSpace(msg.Text)
	text = strings.ReplaceAll(text, ",", "")

	price, err := strconv.ParseInt(text, 10, 64)
	if err != nil || price < 0 {
		h.send(msg.Chat.ID, m.InvalidPrice)
		return
	}

	priceCents := price * 100
	product, err := h.svc.FindOrCreateProduct(ctx, msg.From.ID, conv.Product, priceCents)
	if err != nil {
		h.sendWithKeyboard(msg.Chat.ID, m.ProductError2, keyboard.MainMenu(m))
		h.state.Reset(msg.From.ID)
		return
	}

	h.state.Reset(msg.From.ID)
	h.sendWithKeyboard(msg.Chat.ID,
		fmt.Sprintf(m.ProductSaved, product.Name, domain.FormatBirr(priceCents, m.Birr)),
		keyboard.MainMenu(m))
}

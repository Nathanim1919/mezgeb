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

// showProducts shows the product sub-menu with Add/List buttons.
func (h *Handler) showProducts(ctx context.Context, msg *tgbotapi.Message, m *i18n.Messages) {
	conv := &state.Conversation{Step: state.StepProductMenu}
	h.state.Set(msg.From.ID, conv)
	h.sendWithKeyboard(msg.Chat.ID, m.ProductMenuTitle, keyboard.ProductMenu(m))
}

// handleProductMenu routes Add Product / List Products button presses.
func (h *Handler) handleProductMenu(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
	switch msg.Text {
	case m.BtnAddProduct:
		conv.Step = state.StepProductName
		h.state.Set(msg.From.ID, conv)
		h.sendWithKeyboard(msg.Chat.ID, m.ProductAskName, keyboard.Cancel(m))
	case m.BtnListProducts:
		h.listProducts(ctx, msg, m)
		h.state.Reset(msg.From.ID)
	default:
		h.send(msg.Chat.ID, m.InvalidChoice)
	}
}

// listProducts displays all products with price and stock info.
func (h *Handler) listProducts(ctx context.Context, msg *tgbotapi.Message, m *i18n.Messages) {
	products, err := h.svc.ListProducts(ctx, msg.From.ID)
	if err != nil {
		h.sendWithKeyboard(msg.Chat.ID, m.ErrorGeneric, keyboard.MainMenu(m))
		return
	}

	text := m.ProductsTitle + "\n\n"
	if len(products) == 0 {
		text += m.ProductsEmpty
	} else {
		for i, p := range products {
			stockLabel := fmt.Sprintf(m.ProductStock, p.Stock)
			if p.Stock <= 5 && p.Stock > 0 {
				stockLabel += " " + m.ProductLowStock
			} else if p.Stock == 0 {
				stockLabel = "⛔ 0"
			}
			text += fmt.Sprintf("%d. *%s*\n   💰 %s  |  📊 %s\n",
				i+1,
				escMD(p.Name),
				domain.FormatBirr(p.Price, m.Birr),
				stockLabel,
			)
		}
	}

	h.sendWithKeyboard(msg.Chat.ID, text, keyboard.MainMenu(m))
}

func (h *Handler) handleProductName(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
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
	conv.Step = state.StepProductPrice
	h.state.Set(msg.From.ID, conv)

	h.sendWithKeyboard(msg.Chat.ID, fmt.Sprintf(m.ProductAskPrice, escMD(name)), keyboard.Cancel(m))
}

func (h *Handler) handleProductPrice(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
	text := strings.TrimSpace(msg.Text)
	text = strings.ReplaceAll(text, ",", "")

	price, err := strconv.ParseInt(text, 10, 64)
	if err != nil || price <= 0 {
		h.send(msg.Chat.ID, m.InvalidPrice)
		return
	}

	priceCents := price * 100
	if priceCents > maxAmountCents {
		h.send(msg.Chat.ID, m.AmountTooLarge)
		return
	}

	conv.ProductPrice = priceCents
	conv.Step = state.StepProductStock
	h.state.Set(msg.From.ID, conv)

	h.sendWithKeyboard(msg.Chat.ID, fmt.Sprintf(m.ProductAskStock, escMD(conv.Product)), keyboard.Cancel(m))
}

func (h *Handler) handleProductStock(ctx context.Context, msg *tgbotapi.Message, conv *state.Conversation, m *i18n.Messages) {
	text := strings.TrimSpace(msg.Text)
	text = strings.ReplaceAll(text, ",", "")

	stock, err := strconv.ParseInt(text, 10, 64)
	if err != nil || stock < 0 {
		h.send(msg.Chat.ID, m.InvalidStock)
		return
	}

	product, err := h.svc.FindOrCreateProduct(ctx, msg.From.ID, conv.Product, conv.ProductPrice, stock)
	if err != nil {
		h.sendWithKeyboard(msg.Chat.ID, m.ProductError2, keyboard.MainMenu(m))
		h.state.Reset(msg.From.ID)
		return
	}

	h.state.Reset(msg.From.ID)
	h.sendWithKeyboard(msg.Chat.ID,
		fmt.Sprintf(m.ProductSaved, escMD(product.Name), domain.FormatBirr(conv.ProductPrice, m.Birr), stock),
		keyboard.MainMenu(m))
}

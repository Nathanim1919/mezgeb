package handler

import (
	"context"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nathanim1919/mezgeb/internal/bot/keyboard"
	"github.com/nathanim1919/mezgeb/internal/domain"
)

func (h *Handler) startReport(msg *tgbotapi.Message) {
	h.sendWithKeyboard(msg.Chat.ID, "📊 Choose a time period:", keyboard.ReportPeriod())
}

func (h *Handler) handleReportPeriod(ctx context.Context, msg *tgbotapi.Message) {
	now := time.Now()
	var from time.Time
	var label string

	switch msg.Text {
	case "📅 Today":
		from = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		label = "Today"
	case "📆 This Week":
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		from = time.Date(now.Year(), now.Month(), now.Day()-weekday+1, 0, 0, 0, 0, now.Location())
		label = "This Week"
	case "🗓 This Month":
		from = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		label = "This Month"
	default:
		return
	}

	report, err := h.svc.GetReport(ctx, msg.From.ID, from, now)
	if err != nil {
		h.sendWithKeyboard(msg.Chat.ID, "❌ Error generating report.", keyboard.MainMenu())
		return
	}

	text := formatReport(label, report)
	h.sendWithKeyboard(msg.Chat.ID, text, keyboard.MainMenu())
}

func formatReport(period string, r *domain.ReportData) string {
	s := fmt.Sprintf("📊 *Report — %s*\n\n", period)
	s += fmt.Sprintf("📝 Transactions: *%d*\n", r.TotalTransactions)
	s += fmt.Sprintf("💰 Revenue: *%s*\n", domain.FormatBirr(r.TotalRevenue))
	s += fmt.Sprintf("💸 New Debt: *%s*\n", domain.FormatBirr(r.TotalDebt))

	if len(r.TopProducts) > 0 {
		s += "\n📦 *Top Products:*\n"
		for i, p := range r.TopProducts {
			s += fmt.Sprintf("  %d. %s — %d sold (%s)\n", i+1, p.Name, p.Count, domain.FormatBirr(p.Total))
		}
	}

	if r.TotalTransactions == 0 {
		s += "\n_No transactions in this period._"
	}

	return s
}

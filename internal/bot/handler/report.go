package handler

import (
	"context"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nathanim1919/mezgeb/internal/bot/keyboard"
	"github.com/nathanim1919/mezgeb/internal/domain"
	"github.com/nathanim1919/mezgeb/internal/i18n"
)

func (h *Handler) startReport(msg *tgbotapi.Message, m *i18n.Messages) {
	h.sendWithKeyboard(msg.Chat.ID, m.ReportChoose, keyboard.ReportPeriod(m))
}

func (h *Handler) handleReportPeriod(ctx context.Context, msg *tgbotapi.Message, m *i18n.Messages) {
	now := time.Now()
	var from time.Time
	var label string

	switch msg.Text {
	case m.BtnToday:
		from = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		label = m.BtnToday
	case m.BtnThisWeek:
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		from = time.Date(now.Year(), now.Month(), now.Day()-weekday+1, 0, 0, 0, 0, now.Location())
		label = m.BtnThisWeek
	case m.BtnThisMonth:
		from = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		label = m.BtnThisMonth
	default:
		return
	}

	report, err := h.svc.GetReport(ctx, msg.From.ID, from, now)
	if err != nil {
		h.sendWithKeyboard(msg.Chat.ID, m.ReportError, keyboard.MainMenu(m))
		return
	}

	text := formatReport(label, report, m)
	h.sendWithKeyboard(msg.Chat.ID, text, keyboard.MainMenu(m))
}

func formatReport(period string, r *domain.ReportData, m *i18n.Messages) string {
	s := fmt.Sprintf(m.ReportTitle, period) + "\n\n"
	s += fmt.Sprintf(m.ReportTx, r.TotalTransactions) + "\n"
	s += fmt.Sprintf(m.ReportRevenue, domain.FormatBirr(r.TotalRevenue, m.Birr)) + "\n"
	s += fmt.Sprintf(m.ReportDebt, domain.FormatBirr(r.TotalDebt, m.Birr)) + "\n"

	if len(r.TopProducts) > 0 {
		s += "\n" + m.ReportTopProducts + "\n"
		for i, p := range r.TopProducts {
			s += fmt.Sprintf("  %d. %s — %d (%s)\n", i+1, p.Name, p.Count, domain.FormatBirr(p.Total, m.Birr))
		}
	}

	if r.TotalTransactions == 0 {
		s += "\n" + m.ReportNoTx
	}

	return s
}

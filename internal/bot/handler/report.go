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

	if r.TotalTransactions == 0 {
		s += m.ReportNoTx
		return s
	}

	s += fmt.Sprintf(m.ReportTx, r.TotalTransactions) + "\n\n"

	// Sell/Buy section (the main business metrics)
	hasSellBuy := r.TotalSales > 0 || r.TotalExpenses > 0
	if hasSellBuy {
		s += fmt.Sprintf(m.ReportSales, domain.FormatBirr(r.TotalSales, m.Birr)) + "\n"
		s += fmt.Sprintf(m.ReportExpenses, domain.FormatBirr(r.TotalExpenses, m.Birr)) + "\n"

		profit := r.TotalSales - r.TotalExpenses
		s += fmt.Sprintf(m.ReportProfit, domain.FormatBirr(profit, m.Birr)) + "\n"

		s += fmt.Sprintf(m.ReportItemsSold, r.ItemsSold) + "\n"
		s += fmt.Sprintf(m.ReportItemsBought, r.ItemsBought) + "\n"
	}

	// Legacy debt/payment section (only show if data exists)
	hasLegacy := r.TotalRevenue > 0 || r.TotalDebt > 0
	if hasLegacy {
		if hasSellBuy {
			s += "\n"
		}
		s += fmt.Sprintf(m.ReportRevenue, domain.FormatBirr(r.TotalRevenue, m.Birr)) + "\n"
		s += fmt.Sprintf(m.ReportDebt, domain.FormatBirr(r.TotalDebt, m.Birr)) + "\n"
	}

	// Top sold products
	if len(r.TopProducts) > 0 {
		s += "\n" + m.ReportTopProducts + "\n"
		for i, p := range r.TopProducts {
			s += fmt.Sprintf("  %d. %s — %d pcs (%s)\n", i+1, p.Name, p.Count, domain.FormatBirr(p.Total, m.Birr))
		}
	}

	return s
}

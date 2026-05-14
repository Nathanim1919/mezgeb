package keyboard

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nathanim1919/mezgeb/internal/i18n"
)

func MainMenu(m *i18n.Messages) tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(m.BtnAddTx),
			tgbotapi.NewKeyboardButton(m.BtnReports),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(m.BtnProducts),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(m.BtnSettings),
		),
	)
}

func SellMenu(m *i18n.Messages) tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(m.BtnNewSell),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(m.BtnListSells),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(m.BtnCancel),
		),
	)
}

func BuyMenu(m *i18n.Messages) tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(m.BtnNewBuy),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(m.BtnListBuys),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(m.BtnCancel),
		),
	)
}

func BorrowMenu(m *i18n.Messages) tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(m.BtnNewBorrow),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(m.BtnListBorrows),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(m.BtnCancel),
		),
	)
}

func LoanMenu(m *i18n.Messages) tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(m.BtnNewLoan),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(m.BtnListLoans),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(m.BtnCancel),
		),
	)
}

func TxEditMenu(m *i18n.Messages) tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(m.BtnEditAmount),
			tgbotapi.NewKeyboardButton(m.BtnEditNote),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(m.BtnDelete),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(m.BtnCancel),
		),
	)
}

func TransactionMenu(m *i18n.Messages) tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(m.BtnSell),
			tgbotapi.NewKeyboardButton(m.BtnBuy),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(m.BtnBorrow),
			tgbotapi.NewKeyboardButton(m.BtnLoan),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(m.BtnCancel),
		),
	)
}

func TransactionType(m *i18n.Messages) tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(m.BtnOwesMe),
			tgbotapi.NewKeyboardButton(m.BtnPaidMe),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(m.BtnBoughtProduct),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(m.BtnCancel),
		),
	)
}

func ProductChoice(m *i18n.Messages, products []string) tgbotapi.ReplyKeyboardMarkup {
	var rows [][]tgbotapi.KeyboardButton
	for _, p := range products {
		rows = append(rows, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(p)))
	}
	rows = append(rows, tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(m.BtnSkip),
		tgbotapi.NewKeyboardButton(m.BtnCancel),
	))
	return tgbotapi.NewReplyKeyboard(rows...)
}

func ProductChoiceWithNew(m *i18n.Messages, products []string) tgbotapi.ReplyKeyboardMarkup {
	var rows [][]tgbotapi.KeyboardButton
	for _, p := range products {
		rows = append(rows, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(p)))
	}
	rows = append(rows, tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(m.BtnNewProduct),
	))
	rows = append(rows, tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(m.BtnCancel),
	))
	return tgbotapi.NewReplyKeyboard(rows...)
}

func NotEnoughStock(m *i18n.Messages) tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(m.BtnSellAll),
			tgbotapi.NewKeyboardButton(m.BtnChangeProduct),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(m.BtnCancel),
		),
	)
}

func Confirm(m *i18n.Messages) tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(m.BtnConfirm),
			tgbotapi.NewKeyboardButton(m.BtnCancel),
		),
	)
}

func ReportPeriod(m *i18n.Messages) tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(m.BtnToday),
			tgbotapi.NewKeyboardButton(m.BtnThisWeek),
			tgbotapi.NewKeyboardButton(m.BtnThisMonth),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(m.BtnCancel),
		),
	)
}

func Cancel(m *i18n.Messages) tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(m.BtnCancel),
		),
	)
}

func SkipCancel(m *i18n.Messages) tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(m.BtnSkip),
			tgbotapi.NewKeyboardButton(m.BtnCancel),
		),
	)
}

func ProductEditMenu(m *i18n.Messages) tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(m.BtnEditPrice),
			tgbotapi.NewKeyboardButton(m.BtnEditStock),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(m.BtnDelete),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(m.BtnCancel),
		),
	)
}

func ProductMenu(m *i18n.Messages) tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(m.BtnAddProduct),
			tgbotapi.NewKeyboardButton(m.BtnListProducts),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(m.BtnCancel),
		),
	)
}

func Settings(m *i18n.Messages) tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(m.BtnLanguage),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(m.BtnClearData),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(m.BtnCancel),
		),
	)
}

func ClearDataConfirm(m *i18n.Messages) tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(m.BtnConfirmClear),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(m.BtnCancel),
		),
	)
}

func LanguageChoice() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🇪🇹 አማርኛ"),
			tgbotapi.NewKeyboardButton("🇬🇧 English"),
		),
	)
}

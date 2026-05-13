package i18n

type Lang string

const (
	Amharic Lang = "am"
	English Lang = "en"
)

func Default() Lang {
	return Amharic
}

func Parse(s string) Lang {
	switch s {
	case "en":
		return English
	case "am":
		return Amharic
	default:
		return Amharic
	}
}

type Messages struct {
	// Main menu
	Welcome          string
	MainMenuPrompt   string
	BtnAddTx         string
	BtnReports       string
	BtnCustomers     string
	BtnProducts      string
	BtnSettings      string

	// Transaction flow
	AskCustomerName  string
	AskTxType        string
	BtnOwesMe        string
	BtnPaidMe        string
	BtnBoughtProduct string
	AskAmount        string
	AskProduct       string
	BtnSkip          string
	BtnCancel        string
	BtnConfirm       string
	AskNote          string
	TxSummaryTitle   string
	TxSummaryNote    string
	TxSummaryCustomer string
	TxSummaryType    string
	TxSummaryAmount  string
	TxSummaryProduct string
	TxSummaryConfirm string
	TxConfirmDebt    string // "✅ %s now owes you %s"
	TxConfirmPayment string // "✅ Recorded %s payment from %s"
	TxConfirmPurchase string // "✅ %s bought %s for %s"
	TxConfirmGeneric string
	TxFailed         string
	InvalidAmount    string
	AmountTooLarge   string
	InvalidChoice    string
	EnterCustomerName string
	NameTooLong      string
	ProductError     string

	// Reports
	ReportChoose     string
	BtnToday         string
	BtnThisWeek      string
	BtnThisMonth     string
	ReportTitle      string
	ReportTx         string
	ReportRevenue    string
	ReportDebt       string
	ReportTopProducts string
	ReportNoTx       string
	ReportError      string

	// Customers
	CustomersTitle   string
	CustomersEmpty   string
	CustomerOwes     string // "owes you %s"
	CustomerYouOwe   string // "you owe %s"
	CustomerSettled  string
	CustomersTotal   string

	// Products
	ProductsTitle    string
	ProductsEmpty    string
	ProductAskName   string
	ProductAskPrice  string
	ProductSaved     string
	ProductError2    string
	InvalidPrice     string

	// Settings
	SettingsTitle    string
	BtnLanguage      string
	ChooseLanguage   string
	LanguageUpdated  string

	// General
	NotUnderstood    string
	ErrorGeneric     string
	RateLimited      string

	// Currency
	Birr             string
}

var translations = map[Lang]*Messages{
	Amharic: amharic(),
	English: english(),
}

func Get(lang Lang) *Messages {
	if m, ok := translations[lang]; ok {
		return m
	}
	return translations[Amharic]
}

func amharic() *Messages {
	return &Messages{
		Welcome:          "እንኳን ወደ መዝገብ በደህና መጡ! 📒\n\nቀላል የንግድ ረዳትዎ።\nምን ማድረግ ይፈልጋሉ?",
		MainMenuPrompt:   "ምን ማድረግ ይፈልጋሉ?",
		BtnAddTx:         "➕ ግብይት ጨምር",
		BtnReports:       "📊 ሪፖርት",
		BtnCustomers:     "👥 ደንበኞች",
		BtnProducts:      "📦 ምርቶች",
		BtnSettings:      "⚙️ ማስተካከያ",

		AskCustomerName:  "👤 የደንበኛ ስም?",
		AskTxType:        "ተቀብሏል! *%s*\n\nየግብይት አይነት?",
		BtnOwesMe:        "💸 ይበደረኛል",
		BtnPaidMe:        "💰 ከፍሎኛል",
		BtnBoughtProduct: "🛒 ምርት ገዝቷል",
		AskAmount:        "💰 ስንት ብር?",
		AskProduct:       "📦 የትኛው ምርት? (ወይም ዝለል)",
		BtnSkip:          "⏭ ዝለል",
		BtnCancel:        "❌ ሰርዝ",
		BtnConfirm:       "✅ አረጋግጥ",
		AskNote:          "📝 ማስታወሻ? (ወይም ዝለል)\nምሳሌ: ለ3 ኪሎ ስኳር",
		TxSummaryTitle:   "📋 *የግብይት ማጠቃለያ*",
		TxSummaryNote:    "📝 ማስታወሻ: *%s*",
		TxSummaryCustomer: "👤 ደንበኛ: *%s*",
		TxSummaryType:    "📝 አይነት: *%s*",
		TxSummaryAmount:  "💰 መጠን: *%s*",
		TxSummaryProduct: "📦 ምርት: *%s*",
		TxSummaryConfirm: "ያረጋግጡ?",
		TxConfirmDebt:    "✅ *%s* *%s* ይበደርዎታል",
		TxConfirmPayment: "✅ *%s* ከ *%s* ክፍያ ተመዝግቧል",
		TxConfirmPurchase: "✅ *%s* *%s* ገዝቷል በ *%s*",
		TxConfirmGeneric: "✅ ግብይት ተመዝግቧል!",
		TxFailed:         "❌ ግብይቱን ማስቀመጥ አልተቻለም። እባክዎ ደግመው ይሞክሩ።",
		InvalidAmount:    "እባክዎ ትክክለኛ መጠን ያስገቡ።\nምሳሌ: `1500`፣ `250.50`",
		AmountTooLarge:   "⚠️ ከፍተኛ መጠን 10,000,000 ብር ነው።",
		InvalidChoice:    "እባክዎ ከታች ያሉትን ቁልፎች ይጠቀሙ 👇",
		EnterCustomerName: "እባክዎ የደንበኛ ስም ያስገቡ።",
		NameTooLong:      "⚠️ ስሙ በጣም ረጅም ነው። እባክዎ ያሳጥሩት።",
		ProductError:     "❌ ከምርቱ ጋር ችግር ተፈጥሯል። እየዘለልን ነው።",

		ReportChoose:     "📊 ጊዜ ይምረጡ:",
		BtnToday:         "📅 ዛሬ",
		BtnThisWeek:      "📆 በዚህ ሳምንት",
		BtnThisMonth:     "🗓 በዚህ ወር",
		ReportTitle:      "📊 *ሪፖርት — %s*",
		ReportTx:         "📝 ግብይቶች: *%d*",
		ReportRevenue:    "💰 ገቢ: *%s*",
		ReportDebt:       "💸 አዲስ ብድር: *%s*",
		ReportTopProducts: "📦 *ምርጥ ምርቶች:*",
		ReportNoTx:       "_በዚህ ጊዜ ውስጥ ግብይት የለም።_",
		ReportError:      "❌ ሪፖርት ማዘጋጀት አልተቻለም።",

		CustomersTitle:   "👥 *ደንበኞችዎ*",
		CustomersEmpty:   "👥 እስካሁን ደንበኛ የለም።\n\nየመጀመሪያ ግብይትዎን ያስገቡ ደንበኛ ለመፍጠር!",
		CustomerOwes:     "ይበደርዎታል %s",
		CustomerYouOwe:   "እርስዎ ይበደራሉ %s",
		CustomerSettled:  "ተወራርዷል ✓",
		CustomersTotal:   "_ጠቅላላ: %d ደንበኞች_",

		ProductsTitle:    "📦 *ምርቶችዎ*",
		ProductsEmpty:    "_እስካሁን ምርት የለም።_",
		ProductAskName:   "ምርት ለመጨመር ስሙን ይጻፉ ወይም ከሜኒው ይምረጡ።",
		ProductAskPrice:  "💰 ለ *%s* ዋጋ? (በብር)",
		ProductSaved:     "✅ ምርት *%s* በ *%s* ተቀምጧል",
		ProductError2:    "❌ ምርቱን ማስቀመጥ አልተቻለም።",
		InvalidPrice:     "እባክዎ ትክክለኛ ዋጋ ያስገቡ። ምሳሌ: `500`",

		SettingsTitle:    "⚙️ *ማስተካከያ*",
		BtnLanguage:      "🌍 ቋንቋ",
		ChooseLanguage:   "🌍 ቋንቋ ይምረጡ:",
		LanguageUpdated:  "✅ ቋንቋ ተቀይሯል!",

		NotUnderstood:    "ይቅርታ፣ አልገባኝም። ከታች ያለውን ሜኒው ይጠቀሙ 👇",
		ErrorGeneric:     "❌ ችግር ተፈጥሯል። እባክዎ ደግመው ይሞክሩ።",
		RateLimited:      "⏳ እባክዎ ትንሽ ይጠብቁ። በጣም ብዙ መልዕክቶች።",
		Birr:             "ብር",
	}
}

func english() *Messages {
	return &Messages{
		Welcome:          "Welcome to Mezgeb! 📒\n\nYour simple business assistant.\nWhat would you like to do?",
		MainMenuPrompt:   "What would you like to do?",
		BtnAddTx:         "➕ Add Transaction",
		BtnReports:       "📊 Reports",
		BtnCustomers:     "👥 Customers",
		BtnProducts:      "📦 Products",
		BtnSettings:      "⚙️ Settings",

		AskCustomerName:  "👤 Customer name?",
		AskTxType:        "Got it! *%s*\n\nWhat type of transaction?",
		BtnOwesMe:        "💸 Owes Me",
		BtnPaidMe:        "💰 Paid Me",
		BtnBoughtProduct: "🛒 Bought Product",
		AskAmount:        "💰 How much? (in birr)",
		AskProduct:       "📦 Which product? (or skip)",
		BtnSkip:          "⏭ Skip",
		BtnCancel:        "❌ Cancel",
		BtnConfirm:       "✅ Confirm",
		AskNote:          "📝 Add a note? (or skip)\nExample: for 3kg sugar",
		TxSummaryTitle:   "📋 *Transaction Summary*",
		TxSummaryNote:    "📝 Note: *%s*",
		TxSummaryCustomer: "👤 Customer: *%s*",
		TxSummaryType:    "📝 Type: *%s*",
		TxSummaryAmount:  "💰 Amount: *%s*",
		TxSummaryProduct: "📦 Product: *%s*",
		TxSummaryConfirm: "Confirm?",
		TxConfirmDebt:    "✅ *%s* now owes you *%s*",
		TxConfirmPayment: "✅ Recorded *%s* payment from *%s*",
		TxConfirmPurchase: "✅ *%s* bought *%s* for *%s*",
		TxConfirmGeneric: "✅ Transaction recorded!",
		TxFailed:         "❌ Failed to save transaction. Please try again.",
		InvalidAmount:    "Please enter a valid amount.\nExamples: `1500`, `250.50`",
		AmountTooLarge:   "⚠️ Maximum amount is 10,000,000 birr.",
		InvalidChoice:    "Please choose from the buttons below 👇",
		EnterCustomerName: "Please enter a customer name.",
		NameTooLong:      "⚠️ Name is too long. Please shorten it.",
		ProductError:     "❌ Error with product. Skipping.",

		ReportChoose:     "📊 Choose a time period:",
		BtnToday:         "📅 Today",
		BtnThisWeek:      "📆 This Week",
		BtnThisMonth:     "🗓 This Month",
		ReportTitle:      "📊 *Report — %s*",
		ReportTx:         "📝 Transactions: *%d*",
		ReportRevenue:    "💰 Revenue: *%s*",
		ReportDebt:       "💸 New Debt: *%s*",
		ReportTopProducts: "📦 *Top Products:*",
		ReportNoTx:       "_No transactions in this period._",
		ReportError:      "❌ Error generating report.",

		CustomersTitle:   "👥 *Your Customers*",
		CustomersEmpty:   "👥 No customers yet.\n\nAdd your first transaction to create a customer!",
		CustomerOwes:     "owes you %s",
		CustomerYouOwe:   "you owe %s",
		CustomerSettled:  "settled ✓",
		CustomersTotal:   "_Total: %d customers_",

		ProductsTitle:    "📦 *Your Products*",
		ProductsEmpty:    "_No products yet._",
		ProductAskName:   "To add a product, type its name below or tap a menu button.",
		ProductAskPrice:  "💰 Default price for *%s*? (in birr)",
		ProductSaved:     "✅ Product *%s* saved at *%s*",
		ProductError2:    "❌ Error saving product.",
		InvalidPrice:     "Please enter a valid price. Example: `500`",

		SettingsTitle:    "⚙️ *Settings*",
		BtnLanguage:      "🌍 Language",
		ChooseLanguage:   "🌍 Choose language:",
		LanguageUpdated:  "✅ Language updated!",

		NotUnderstood:    "I didn't understand that. Use the menu below 👇",
		ErrorGeneric:     "❌ Something went wrong. Please try again.",
		RateLimited:      "⏳ Please wait a moment. Too many messages.",
		Birr:             "birr",
	}
}

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

	// Transaction menu
	TxMenuTitle      string
	BtnSell          string
	BtnBuy           string
	BtnBorrow        string
	BtnLoan          string
	ComingSoon       string

	// Sell/Buy flow
	AskSellProduct   string
	AskBuyProduct    string
	BtnNewProduct    string
	AskQuantity      string
	AskBuyPrice      string
	InvalidQuantity  string
	NotEnoughStock   string // "Not enough stock! You have %d"
	BtnChangeProduct string
	BtnSellAll       string
	TxSummaryTitle   string
	TxSummaryNote    string
	TxSummaryProduct string
	TxSummaryQty     string
	TxSummaryUnitPrice string
	TxSummaryTotal   string
	TxSummaryType    string
	TxSummaryConfirm string
	SellConfirm      string // "✅ Sold %d × %s for %s"
	BuyConfirm       string // "✅ Bought %d × %s for %s"
	TxConfirmGeneric string
	TxFailed         string

	// Legacy transaction flow (debt/payment — for borrow/loan later)
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
	TxSummaryCustomer string
	TxSummaryAmount  string
	TxConfirmDebt    string
	TxConfirmPayment string
	TxConfirmPurchase string
	InvalidAmount    string
	AmountTooLarge   string
	InvalidChoice    string
	EnterCustomerName string
	NameTooLong      string
	ProductError     string

	// Reports
	ReportChoose      string
	BtnToday          string
	BtnThisWeek       string
	BtnThisMonth      string
	ReportTitle       string
	ReportTx          string
	ReportSales       string
	ReportExpenses    string
	ReportProfit      string
	ReportItemsSold   string
	ReportItemsBought string
	ReportRevenue     string
	ReportDebt        string
	ReportTopProducts string
	ReportNoTx        string
	ReportError       string

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
	BtnAddProduct    string
	BtnListProducts  string
	ProductMenuTitle string
	ProductAskName   string
	ProductAskPrice  string
	ProductAskStock  string
	ProductSaved     string
	ProductError2    string
	InvalidPrice     string
	InvalidStock     string
	ProductStock     string // "Stock: %d"
	ProductLowStock  string // "⚠️ Low"

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
		BtnAddTx:         "💰 ግብይቶች",
		BtnReports:       "📊 ሪፖርት",
		BtnCustomers:     "👥 ደንበኞች",
		BtnProducts:      "📦 ምርቶች",
		BtnSettings:      "⚙️ ማስተካከያ",

		TxMenuTitle:      "💰 *ግብይቶች*\n\nምን ማድረግ ይፈልጋሉ?",
		BtnSell:          "🛒 ሽያጭ",
		BtnBuy:           "📥 ግዢ",
		BtnBorrow:        "🤝 ብድር",
		BtnLoan:          "💸 አበድር",
		ComingSoon:       "🔜 በቅርቡ ይመጣል!",

		AskSellProduct:   "🛒 የትኛውን ምርት ይሸጣሉ?",
		AskBuyProduct:    "📥 የትኛውን ምርት ይገዛሉ?",
		BtnNewProduct:    "➕ አዲስ ምርት",
		AskQuantity:      "🔢 ስንት? (ብዛት)",
		AskBuyPrice:      "💰 የግዢ ዋጋ ለአንዱ? (በብር)",
		InvalidQuantity:  "እባክዎ ትክክለኛ ቁጥር ያስገቡ። ምሳሌ: `5`",
		NotEnoughStock:   "⚠️ ክምችት በቂ አይደለም! ያለዎት: *%d*\n\nሌላ ቁጥር ያስገቡ፣ ሁሉንም ይሸጡ፣ ወይም ሌላ ምርት ይምረጡ 👇",
		BtnChangeProduct: "🔄 ሌላ ምርት",
		BtnSellAll:       "📦 ሁሉንም ሽጥ",
		TxSummaryTitle:   "📋 *የግብይት ማጠቃለያ*",
		TxSummaryNote:    "📝 ማስታወሻ: *%s*",
		TxSummaryProduct: "📦 ምርት: *%s*",
		TxSummaryQty:     "🔢 ብዛት: *%d*",
		TxSummaryUnitPrice: "💰 ዋጋ/አንድ: *%s*",
		TxSummaryTotal:   "💵 ጠቅላላ: *%s*",
		TxSummaryType:    "📝 አይነት: *%s*",
		TxSummaryConfirm: "ያረጋግጡ?",
		SellConfirm:      "✅ *%d* × *%s* ተሽጧል በ *%s*",
		BuyConfirm:       "✅ *%d* × *%s* ተገዝቷል በ *%s*",
		TxConfirmGeneric: "✅ ግብይት ተመዝግቧል!",
		TxFailed:         "❌ ግብይቱን ማስቀመጥ አልተቻለም። እባክዎ ደግመው ይሞክሩ።",

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
		TxSummaryCustomer: "👤 ደንበኛ: *%s*",
		TxSummaryAmount:  "💰 መጠን: *%s*",
		TxConfirmDebt:    "✅ *%s* *%s* ይበደርዎታል",
		TxConfirmPayment: "✅ *%s* ከ *%s* ክፍያ ተመዝግቧል",
		TxConfirmPurchase: "✅ *%s* *%s* ገዝቷል በ *%s*",
		InvalidAmount:    "እባክዎ ትክክለኛ መጠን ያስገቡ።\nምሳሌ: `1500`፣ `250.50`",
		AmountTooLarge:   "⚠️ ከፍተኛ መጠን 10,000,000 ብር ነው።",
		InvalidChoice:    "እባክዎ ከታች ያሉትን ቁልፎች ይጠቀሙ 👇",
		EnterCustomerName: "እባክዎ የደንበኛ ስም ያስገቡ።",
		NameTooLong:      "⚠️ ስሙ በጣም ረጅም ነው። እባክዎ ያሳጥሩት።",
		ProductError:     "❌ ከምርቱ ጋር ችግር ተፈጥሯል። እየዘለልን ነው።",

		ReportChoose:      "📊 ጊዜ ይምረጡ:",
		BtnToday:          "📅 ዛሬ",
		BtnThisWeek:       "📆 በዚህ ሳምንት",
		BtnThisMonth:      "🗓 በዚህ ወር",
		ReportTitle:       "📊 *ሪፖርት — %s*",
		ReportTx:          "📝 ግብይቶች: *%d*",
		ReportSales:       "🛒 ሽያጭ: *%s*",
		ReportExpenses:    "📥 ወጪ (ግዢ): *%s*",
		ReportProfit:      "📈 ትርፍ: *%s*",
		ReportItemsSold:   "📦 የተሸጡ ምርቶች: *%d*",
		ReportItemsBought: "📦 የተገዙ ምርቶች: *%d*",
		ReportRevenue:     "💰 ክፍያዎች: *%s*",
		ReportDebt:        "💸 አዲስ ብድር: *%s*",
		ReportTopProducts: "🏆 *ምርጥ የተሸጡ ምርቶች:*",
		ReportNoTx:        "_በዚህ ጊዜ ውስጥ ግብይት የለም።_",
		ReportError:       "❌ ሪፖርት ማዘጋጀት አልተቻለም።",

		CustomersTitle:   "👥 *ደንበኞችዎ*",
		CustomersEmpty:   "👥 እስካሁን ደንበኛ የለም።\n\nየመጀመሪያ ግብይትዎን ያስገቡ ደንበኛ ለመፍጠር!",
		CustomerOwes:     "ይበደርዎታል %s",
		CustomerYouOwe:   "እርስዎ ይበደራሉ %s",
		CustomerSettled:  "ተወራርዷል ✓",
		CustomersTotal:   "_ጠቅላላ: %d ደንበኞች_",

		ProductsTitle:    "📦 *ምርቶችዎ*",
		ProductsEmpty:    "_እስካሁን ምርት የለም።_",
		BtnAddProduct:    "➕ ምርት ጨምር",
		BtnListProducts:  "📋 ምርቶች ዝርዝር",
		ProductMenuTitle: "📦 *ምርቶች*\n\nምን ማድረግ ይፈልጋሉ?",
		ProductAskName:   "📦 የምርት ስም?",
		ProductAskPrice:  "💰 ለ *%s* ዋጋ? (በብር)",
		ProductAskStock:  "📊 የ *%s* ክምችት ብዛት?",
		ProductSaved:     "✅ ምርት *%s* በ *%s* ተቀምጧል\n📊 ክምችት: *%d*",
		ProductError2:    "❌ ምርቱን ማስቀመጥ አልተቻለም።",
		InvalidPrice:     "እባክዎ ትክክለኛ ዋጋ ያስገቡ። ምሳሌ: `500`",
		InvalidStock:     "እባክዎ ትክክለኛ ቁጥር ያስገቡ። ምሳሌ: `50`",
		ProductStock:     "ክምችት: %d",
		ProductLowStock:  "⚠️ ዝቅተኛ",

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
		BtnAddTx:         "💰 Transactions",
		BtnReports:       "📊 Reports",
		BtnCustomers:     "👥 Customers",
		BtnProducts:      "📦 Products",
		BtnSettings:      "⚙️ Settings",

		TxMenuTitle:      "💰 *Transactions*\n\nWhat would you like to do?",
		BtnSell:          "🛒 Sell",
		BtnBuy:           "📥 Buy",
		BtnBorrow:        "🤝 Borrow",
		BtnLoan:          "💸 Loan",
		ComingSoon:       "🔜 Coming soon!",

		AskSellProduct:   "🛒 Which product are you selling?",
		AskBuyProduct:    "📥 Which product are you buying?",
		BtnNewProduct:    "➕ New Product",
		AskQuantity:      "🔢 How many? (quantity)",
		AskBuyPrice:      "💰 Buy price per unit? (in birr)",
		InvalidQuantity:  "Please enter a valid number. Example: `5`",
		NotEnoughStock:   "⚠️ Not enough stock! You have: *%d*\n\nEnter a different number, sell all, or pick another product 👇",
		BtnChangeProduct: "🔄 Other Product",
		BtnSellAll:       "📦 Sell All",
		TxSummaryTitle:   "📋 *Transaction Summary*",
		TxSummaryNote:    "📝 Note: *%s*",
		TxSummaryProduct: "📦 Product: *%s*",
		TxSummaryQty:     "🔢 Quantity: *%d*",
		TxSummaryUnitPrice: "💰 Price/unit: *%s*",
		TxSummaryTotal:   "💵 Total: *%s*",
		TxSummaryType:    "📝 Type: *%s*",
		TxSummaryConfirm: "Confirm?",
		SellConfirm:      "✅ Sold *%d* × *%s* for *%s*",
		BuyConfirm:       "✅ Bought *%d* × *%s* for *%s*",
		TxConfirmGeneric: "✅ Transaction recorded!",
		TxFailed:         "❌ Failed to save transaction. Please try again.",

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
		TxSummaryCustomer: "👤 Customer: *%s*",
		TxSummaryAmount:  "💰 Amount: *%s*",
		TxConfirmDebt:    "✅ *%s* now owes you *%s*",
		TxConfirmPayment: "✅ Recorded *%s* payment from *%s*",
		TxConfirmPurchase: "✅ *%s* bought *%s* for *%s*",
		InvalidAmount:    "Please enter a valid amount.\nExamples: `1500`, `250.50`",
		AmountTooLarge:   "⚠️ Maximum amount is 10,000,000 birr.",
		InvalidChoice:    "Please choose from the buttons below 👇",
		EnterCustomerName: "Please enter a customer name.",
		NameTooLong:      "⚠️ Name is too long. Please shorten it.",
		ProductError:     "❌ Error with product. Skipping.",

		ReportChoose:      "📊 Choose a time period:",
		BtnToday:          "📅 Today",
		BtnThisWeek:       "📆 This Week",
		BtnThisMonth:      "🗓 This Month",
		ReportTitle:       "📊 *Report — %s*",
		ReportTx:          "📝 Transactions: *%d*",
		ReportSales:       "🛒 Sales: *%s*",
		ReportExpenses:    "📥 Expenses (Purchases): *%s*",
		ReportProfit:      "📈 Profit: *%s*",
		ReportItemsSold:   "📦 Items Sold: *%d*",
		ReportItemsBought: "📦 Items Bought: *%d*",
		ReportRevenue:     "💰 Payments: *%s*",
		ReportDebt:        "💸 New Debt: *%s*",
		ReportTopProducts: "🏆 *Top Sold Products:*",
		ReportNoTx:        "_No transactions in this period._",
		ReportError:       "❌ Error generating report.",

		CustomersTitle:   "👥 *Your Customers*",
		CustomersEmpty:   "👥 No customers yet.\n\nAdd your first transaction to create a customer!",
		CustomerOwes:     "owes you %s",
		CustomerYouOwe:   "you owe %s",
		CustomerSettled:  "settled ✓",
		CustomersTotal:   "_Total: %d customers_",

		ProductsTitle:    "📦 *Your Products*",
		ProductsEmpty:    "_No products yet._",
		BtnAddProduct:    "➕ Add Product",
		BtnListProducts:  "📋 List Products",
		ProductMenuTitle: "📦 *Products*\n\nWhat would you like to do?",
		ProductAskName:   "📦 Product name?",
		ProductAskPrice:  "💰 Price for *%s*? (in birr)",
		ProductAskStock:  "📊 Stock quantity for *%s*?",
		ProductSaved:     "✅ Product *%s* saved at *%s*\n📊 Stock: *%d*",
		ProductError2:    "❌ Error saving product.",
		InvalidPrice:     "Please enter a valid price. Example: `500`",
		InvalidStock:     "Please enter a valid number. Example: `50`",
		ProductStock:     "Stock: %d",
		ProductLowStock:  "⚠️ Low",

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

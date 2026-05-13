# Mezgeb Feature Roadmap

> The cashflow command center for Ethiopian small businesses.
> Know exactly where your money is. Save hours every day. Never lose a birr.
> Telegram-first. Amharic-native. Built for the way business actually works here.

---

## Phase 1 — Foundation (COMPLETED)

Core transaction engine and Telegram bot infrastructure.

- [x] **Transaction Recording** — Debt, Payment, Purchase with conversational flow
- [x] **Customer Management** — Auto-create on first mention, per-user isolation
- [x] **Product Catalog** — Create products with default prices
- [x] **Debt/Credit Tracking** — Automatic balance calculation per customer
- [x] **Basic Reports** — Today / This Week / This Month with revenue and debt summaries
- [x] **Amharic + English i18n** — Full bilingual support with language toggle
- [x] **Rate Limiting & Security** — 30 msg/min, parameterized queries, data isolation

---

## Phase 2 — Cashflow Visibility

> "Where is my money going?" — Every merchant's first question.
> Give merchants a real-time picture of money in vs. money out.

### 2.1 Expense Tracking & Categories
- [ ] **Business expense recording** — Rent, Inventory/Stock, Transport, Utilities, Salary, Tax, Packaging, Maintenance, Other
- [ ] Custom category creation ("ለኪራይ", "ለትራንስፖርት" — in their own words)
- [ ] Quick-expense button — one-tap common expenses (no typing)
- [ ] Daily expense log view
- [ ] Expense vs. income ratio per day/week/month

### 2.2 Real-Time Cashflow Dashboard
- [ ] **Today's cashflow snapshot** — Money In / Money Out / Net (one message, always up to date)
- [ ] Cash on hand tracker (opening balance + transactions = closing balance)
- [ ] "Where did my money go today?" — auto-breakdown by category
- [ ] Cash surplus/deficit alerts ("You spent 2,000 ETB more than you earned today")
- [ ] Weekly cashflow trend (Mon-Sun, which days you make/lose money)

### 2.3 Profit & Loss (Per Transaction)
- [ ] **Cost price vs. selling price** per product — know margin on every sale
- [ ] Per-transaction profit calculation (auto: selling price - cost price)
- [ ] Daily/weekly/monthly gross profit summary
- [ ] Most profitable products ranking
- [ ] Least profitable products (candidates for price adjustment or removal)
- [ ] Profit margin percentage per product

### 2.4 Advanced Reports
- [ ] Custom date range reports ("ከጥር 15 - የካቲት 20 አሳየኝ")
- [ ] Customer-level full transaction history
- [ ] Comparative reports — this month vs. last month with % change
- [ ] Busiest days of the week analysis
- [ ] Export as PDF or CSV (sent as Telegram file)

---

## Phase 3 — Debt & Receivables Control

> Ethiopian small business runs on trust and credit.
> Most merchants lose money not from bad sales, but from uncollected debt.

### 3.1 Smart Debt Tracking
- [ ] **Debt aging** — Classify debts by age (< 7 days, 7-30 days, 30-60 days, 60+ days)
- [ ] Overdue debt alerts ("አበበ 5,000 ብር ከ30 ቀን በላይ ይገባዎታል")
- [ ] Total outstanding debt dashboard (how much is owed to you right now)
- [ ] Debt collection priority list (sorted by amount and age)
- [ ] Partial payment tracking (customer pays 500 of 2,000 — remaining auto-updates)

### 3.2 Customer Credit Intelligence
- [ ] **Credit score per customer** — auto-calculated from payment history
  - Always pays on time = Green / Reliable
  - Sometimes late = Yellow / Caution
  - Frequently overdue = Red / Risky
- [ ] Credit limit per customer (alert when new debt would exceed limit)
- [ ] Payment pattern analysis ("Kebede always pays on the 15th")
- [ ] Customer ranking — top spenders, most reliable, most risky
- [ ] "Should I give credit to X?" — decision support based on history

### 3.3 Debt Reminders & Follow-ups
- [ ] **Merchant-triggered reminders** — "Remind me about Abebe's debt in 3 days"
- [ ] Scheduled reminder list (upcoming reminders dashboard)
- [ ] Suggested follow-up timing based on customer payment pattern
- [ ] Debt summary message you can forward to customer (polite, professional template)
- [ ] Batch reminder — "Show me everyone who owes more than 1,000 ETB"

### 3.4 Payables (What You Owe)
- [ ] Track debts TO suppliers separately from debts FROM customers
- [ ] Supplier payment due dates
- [ ] Payables vs. Receivables balance ("You're owed 50,000 but you owe 20,000 — net: +30,000")
- [ ] Payment scheduling — "Pay supplier X on Friday"

---

## Phase 4 — Inventory & Stock Control

> Stock is cash sitting on the shelf. Know what you have, what's moving, what's dead.

### 4.1 Stock Tracking
- [ ] **Stock quantity per product** — current units on hand
- [ ] Stock-in recording (purchased/restocked — quantity + cost)
- [ ] Stock-out on sale (auto-deduct when transaction recorded)
- [ ] Manual stock adjustment (breakage, loss, personal use)
- [ ] Current stock valuation (total birr sitting in inventory)

### 4.2 Stock Alerts & Intelligence
- [ ] **Low stock alerts** ("ስኳር 5 ብቻ ቀረ — ማዘዝ ይፈልጋሉ?")
- [ ] Reorder point per product (configurable threshold)
- [ ] Fast-moving vs. slow-moving product identification
- [ ] Dead stock detection (products not sold in 30+ days)
- [ ] Stock turnover rate (how fast each product sells)

### 4.3 Purchase & Restock Management
- [ ] Record purchase orders from suppliers
- [ ] Cost price history per product (track price changes from suppliers)
- [ ] Best supplier tracking (who gives better price for same product)
- [ ] Restock suggestions based on sales velocity ("You sell 10 Sugar/day, you have 15 left — restock in 1-2 days")

---

## Phase 5 — Time-Saving Automation

> Every minute a merchant spends on bookkeeping is a minute not spent with customers.

### 5.1 Quick Actions & Shortcuts
- [ ] **Favorite transactions** — one-tap repeat of common transactions
- [ ] Batch transaction entry ("5 coffees sold today at 50 ETB each")
- [ ] Smart defaults — pre-fill amount based on product price
- [ ] Transaction templates (merchant creates reusable patterns)
- [ ] Undo last transaction (within 5 minutes)

### 5.2 Daily Business Rituals
- [ ] **Morning briefing** (auto-sent at merchant's preferred time)
  - Yesterday's summary
  - Outstanding debts due today
  - Low stock alerts
  - Today's scheduled payments
- [ ] **End-of-day closing** — one-tap daily close with summary
  - Total sales, expenses, profit
  - Cash reconciliation (expected vs. actual cash on hand)
  - Unresolved items
- [ ] Weekly business health digest (every Sunday/Monday)

### 5.3 Recurring Transactions
- [ ] Define recurring expenses (rent, salary, utilities) with schedule
- [ ] Auto-reminder before due dates ("Rent due in 2 days — 15,000 ETB")
- [ ] One-tap confirmation to record
- [ ] Track missed/late recurring payments
- [ ] Monthly fixed cost overview

### 5.4 Smart Notifications
- [ ] Large transaction alerts (anything above merchant's normal range)
- [ ] Revenue milestone celebrations ("You passed 100,000 ETB this month!")
- [ ] Cashflow warnings ("At this rate, you'll be negative by Thursday")
- [ ] Quiet hours (no notifications between 10pm-7am)

---

## Phase 6 — Invoicing & Professionalism

> Help merchants look professional and get paid faster.

### 6.1 Digital Receipts
- [ ] **Auto-generate receipt** after each sale (image sent in chat)
- [ ] Business name, logo, and contact on receipt
- [ ] Unique receipt/invoice number
- [ ] Itemized breakdown (products, quantities, prices, total)
- [ ] Shareable — merchant can forward receipt to customer

### 6.2 Invoice Management
- [ ] Create invoices for credit sales (sent as image/PDF)
- [ ] Invoice status tracking (sent, viewed, paid, overdue)
- [ ] Payment recording against invoice
- [ ] Invoice history and search
- [ ] Monthly invoice summary

### 6.3 Business Profile
- [ ] Set business name, type, and location
- [ ] Business registration number (optional)
- [ ] TIN number for tax purposes (optional)
- [ ] This info appears on receipts and invoices

---

## Phase 7 — Growth & Insights

> From surviving to thriving. Data-driven decisions for merchants who want to grow.

### 7.1 Business Analytics
- [ ] **Monthly P&L statement** (revenue - COGS - expenses = net profit)
- [ ] Revenue trend graph (text-based chart in Telegram)
- [ ] Month-over-month growth rate
- [ ] Seasonal pattern recognition ("Sales peak in September, dip in March")
- [ ] Break-even analysis ("You need 150 ETB/day in sales to cover fixed costs")

### 7.2 Goal Setting & Tracking
- [ ] Set monthly revenue targets
- [ ] Daily progress toward goal ("You're 65% to your 50,000 ETB target")
- [ ] Savings goals (set aside profit for business expansion)
- [ ] Debt reduction goals ("Clear all debt over 30 days by end of month")

### 7.3 Multi-Account Management
- [ ] Track Cash, Bank, and Telebirr balances separately
- [ ] Record transfers between accounts
- [ ] Reconcile bank balance with recorded transactions
- [ ] Total business worth = all accounts + inventory + receivables - payables

### 7.4 Export & Tax Preparation
- [ ] Annual financial summary
- [ ] Tax-ready income/expense report
- [ ] Export all data (CSV) for accountant
- [ ] Print-ready reports (PDF via Telegram)

---

## Implementation Priority (What to Build Next)

Based on merchant pain points, this is the recommended build order:

| Priority | Feature | Why |
|----------|---------|-----|
| **P0** | Expense tracking (2.1) | Merchants have zero visibility on where money goes |
| **P0** | Cashflow dashboard (2.2) | The single most requested insight: "am I making money?" |
| **P1** | Cost price & profit per sale (2.3) | Merchants guess their margins — most guess wrong |
| **P1** | Debt aging & alerts (3.1) | Uncollected debt is the #1 silent killer of small business |
| **P1** | Customer credit scoring (3.2) | "Should I give this person credit?" needs a data-backed answer |
| **P2** | Morning briefing (5.2) | Saves 15 minutes every morning, builds daily habit |
| **P2** | Stock tracking (4.1) | Prevents stockouts and dead inventory |
| **P2** | Quick actions (5.1) | Reduces transaction time from 30 seconds to 5 seconds |
| **P3** | Low stock alerts (4.2) | Prevents lost sales from stockouts |
| **P3** | Debt reminders (3.3) | Automates the awkward "you owe me" conversation |
| **P3** | Digital receipts (6.1) | Professionalism + proof of transaction |
| **P4** | Advanced reports & export (2.4) | Power users and accountant handoff |
| **P4** | Recurring transactions (5.3) | Set-and-forget for predictable expenses |
| **P4** | Business analytics (7.1) | Growth-stage merchants making strategic decisions |

---

## Implementation Principles

1. **Telegram-first always** — Every feature must work in chat. Buttons over typing.
2. **Amharic-native** — Not translated, native. Numbers, dates, currency in Ethiopian format.
3. **5-second rule** — Any action should complete in under 5 seconds or it's too complex.
4. **No typing required** — Buttons, quick replies, and smart defaults wherever possible.
5. **Penny-accurate** — All money in cents. No floating point. Ever.
6. **Progressive disclosure** — New users see simple menu. Features unlock as they use more.
7. **Merchant language** — Say "ትርፍ" not "gross margin". Say "ዕዳ" not "accounts receivable".
8. **Offline-tolerant** — Assume unreliable connectivity. Queue and sync.
9. **Data is theirs** — Export everything, delete everything. Full data sovereignty.

---

*Build in priority order. Ship weekly. Validate with real merchants at every step.*
*The goal: a merchant opens Mezgeb once in the morning and knows exactly where their business stands.*

# Mezgeb

A Telegram-first micro-business operating system for small merchants in Ethiopia.

Mezgeb replaces notebooks, memory, and messy chat tracking with a fast, conversational Telegram bot.

## Features (Phase 1)

- **Add Transactions** — debt, payments, purchases in under 5 seconds
- **Debt/Credit Tracking** — automatic balance per customer
- **Product Tracking** — create and manage products
- **Reports** — daily, weekly, monthly summaries
- **Telegram Button Navigation** — no typing required for most actions

## Architecture

```
cmd/bot/              — entry point
internal/
  config/             — environment-based configuration
  domain/             — domain models and repository interfaces
  repository/postgres — PostgreSQL implementations
  service/            — business logic layer
  bot/
    handler/          — Telegram message handlers
    keyboard/         — reply keyboard builders
    state/            — in-memory conversation state
migrations/           — PostgreSQL migrations
```

## Prerequisites

- Go 1.21+
- PostgreSQL 14+
- [golang-migrate](https://github.com/golang-migrate/migrate) CLI
- A Telegram Bot Token (from [@BotFather](https://t.me/BotFather))

## Quick Start

### 1. Start PostgreSQL

```bash
make docker-up
```

### 2. Configure environment

```bash
cp .env.example .env
# Edit .env with your TELEGRAM_BOT_TOKEN
```

### 3. Run migrations

```bash
make migrate-up
```

### 4. Run the bot

```bash
make run
```

## Development

```bash
# Build binary
make build

# Create a new migration
make migrate-create name=add_something

# Roll back last migration
make migrate-down
```

## Docker (Full)

```bash
# Edit .env first, then:
docker compose up --build
```

## Database Schema

- **users** — Telegram users (ID = Telegram user ID)
- **customers** — merchant's customers with running balances
- **products** — product catalog with default prices
- **transactions** — every debt, payment, and purchase recorded

Amounts are stored in cents (birr * 100) for precision.

## Bot Flow

```
/start → Main Menu
  ├── ➕ Add Transaction → Customer → Type → Amount → Product → Confirm
  ├── 📊 Reports → Today / This Week / This Month
  ├── 👥 Customers → List with balances
  ├── 📦 Products → List / Add new
  └── ⚙️ Settings (coming soon)
```

## Design Principles

- Every interaction < 5 seconds
- Conversational, not dashboard-heavy
- Built for unstable internet and non-technical users
- Clear confirmations on every action
- Clean architecture for long-term maintainability

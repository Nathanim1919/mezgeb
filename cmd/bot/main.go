package main

import (
	"context"
	"log"

	_ "github.com/joho/godotenv/autoload"
	"github.com/nathanim1919/mezgeb/internal/bot"
	"github.com/nathanim1919/mezgeb/internal/config"
	"github.com/nathanim1919/mezgeb/internal/repository/postgres"
	"github.com/nathanim1919/mezgeb/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	ctx := context.Background()
	pool, err := postgres.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer pool.Close()

	svc := &service.Service{
		Users:        postgres.NewUserRepo(pool),
		Customers:    postgres.NewCustomerRepo(pool),
		Products:     postgres.NewProductRepo(pool),
		Transactions: postgres.NewTransactionRepo(pool),
		Reports:      postgres.NewReportRepo(pool),
	}

	b, err := bot.New(cfg.TelegramToken, svc)
	if err != nil {
		log.Fatalf("bot: %v", err)
	}

	b.Start()
}

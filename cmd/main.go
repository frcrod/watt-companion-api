package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/frcrod/watt-companion-api/internal/config"
	"github.com/frcrod/watt-companion-api/internal/database/out"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"
)

func run() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ctx := context.Background()

	pgxConfig, err := pgxpool.ParseConfig("user=postgres dbname=postgres password=password")
	if err != nil {
		slog.Error("Error can't connect to db")
		return nil
	}

	if err != nil {
		log.Fatal("Problem Connecting to DB")
	}

	pgxConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxUUID.Register(conn.TypeMap())
		return nil
	}

	pgxConnPool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		panic(err)
	}

	queries := out.New(pgxConnPool)

	e := config.CreateEchoInstance(queries)

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

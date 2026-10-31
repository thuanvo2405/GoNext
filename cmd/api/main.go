package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thuanvo2405/GoNext/internal/data"
	"github.com/thuanvo2405/GoNext/internal/token"
)

type config struct {
	port int
	db   struct {
		dsn string
	}
	jwt struct {
		secret string
		ttl    time.Duration
	}
}

type application struct {
	config config
	models data.Models
	tokens *token.Manager
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "HTTP network address")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("DATABASE_URL"), "PostgresSQL DSN")
	flag.DurationVar(&cfg.jwt.ttl, "jwt-ttl", 30*time.Minute, "JWT access token lifetime")
	flag.StringVar(&cfg.jwt.secret, "jwt-secret", os.Getenv("JWT_SECRET"), "JWT signing secret")

	flag.Parse()

	tokens, err := token.NewManager(cfg.jwt.secret, cfg.jwt.ttl)
	if err != nil {
		log.Fatal(err)
	}

	db, err := pgxpool.New(context.Background(), cfg.db.dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	app := &application{
		config: cfg,
		models: data.NewModels(db),
		tokens: tokens,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.Ping(ctx)
	if err != nil {
		log.Printf("database ping failed: %v", err)
		return
	}

	err = app.setupRouter().Run(fmt.Sprintf(":%v", app.config.port))
	if err != nil {
		fmt.Println(err)
		return
	}
}

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
)

type config struct {
	port int
	db   struct {
		dsn string
	}
}

type application struct {
	config config
	models data.Models
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "HTTP network address")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("DATABASE_URL"), "PostgresSQL DSN")

	flag.Parse()

	db, err := pgxpool.New(context.Background(), cfg.db.dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	app := &application{
		config: cfg,
		models: data.NewModels(db),
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

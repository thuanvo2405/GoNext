package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type config struct {
	port int
	db   struct {
		dsn string
	}
}

type application struct {
	config config
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 8080, "HTTP network address")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("DATABASE_URL"), "PostgresSQL DSN")

	flag.Parse()

	app := &application{
		config: cfg,
	}

	db, err := pgxpool.New(context.Background(), app.config.db.dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.Ping(ctx)
	if err != nil {
		log.Printf("database ping failed: %v", err)
		return
	}

	err = setupRouter().Run(fmt.Sprintf(":%v", app.config.port))
	if err != nil {
		fmt.Println(err)
		return
	}
}

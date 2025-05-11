package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/juancruzestevez/goSocial/internal/db"
	"github.com/juancruzestevez/goSocial/internal/env"
	"github.com/juancruzestevez/goSocial/internal/store"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error cargando .env")
	}

	cfg := config{
		addr: env.GetString("ADDR", ":3000"),
		db: dbConfig{
			addr: env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/gosocial?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime: env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}

	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		log.Panic("Error connecting to the database:", err)
	}

	defer db.Close()
	log.Println("Connected to the database")

	store := store.NewStore(db)

	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}

package main

import (
	"go-backend/config"
	"go-backend/database"
	"go-backend/internal/routers"
	"log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("%v", err)
	}
	db, err := database.Connect(cfg.DBPath)
	if err != nil {
		log.Fatalf("%v", err)
	}
	r := routers.New(db)
	if err := r.Run(cfg.ServerAddr); err != nil {
		log.Fatalf("%v", err)
	}
}

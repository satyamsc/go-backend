package main

import (
    "log"
    "go-backend-challenge/config"
    "go-backend-challenge/database"
    "go-backend-challenge/internal/routers"
)

func main() {
    cfg, err := config.Load()
    if err != nil { log.Fatalf("%v", err) }
    db, err := database.Connect(cfg.DBPath)
    if err != nil { log.Fatalf("%v", err) }
    r := routers.New(db)
    if err := r.Run(cfg.ServerAddr); err != nil { log.Fatalf("%v", err) }
}


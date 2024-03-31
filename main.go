package main

import (
	"fmt"
	"selfbotat-v2/bot"
	"time"

	"selfbotat-v2/bot/database"
	_ "selfbotat-v2/bot/commands"
	"selfbotat-v2/bot/config"
)

var cfg config.Config

func main() {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@localhost/%s?sslmode=disable", 
		cfg.Postgres.User, 
		cfg.Postgres.Password,
		cfg.Postgres.Database,
	)

	err := db.InitDatabase(connStr)
	if err != nil {
		panic(err)
	} else {
		bot.Log.Println("Postgres connected")
	}

	err, success := bot.InitClient()
	if err != nil {
			panic(err)
	} 

	if !success {
			panic("Failed to initialize client")
	}
}


func init() {
	var err error
	cfg, err = config.LoadConfig()
	if err != nil {
		panic(err)
	}
	
	bot.StartTime = time.Now()
}
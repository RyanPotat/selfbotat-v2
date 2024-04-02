package main

import (
	"fmt"
	"selfbotat-v2/bot"
	"time"

	"selfbotat-v2/bot/database"
	"selfbotat-v2/bot/config"
  "selfbotat-v2/bot/logger"

	_ "selfbotat-v2/bot/commands"
	client "selfbotat-v2/bot/client"
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
		Log.Debug.Print("Postgres connected")
	}

	client.Create()
}


func init() {
	var err error
	cfg, err = config.LoadConfig()
	if err != nil {
		panic(err)
	}
	
	bot.StartTime = time.Now()
}
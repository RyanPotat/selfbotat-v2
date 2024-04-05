package main

import (
	"fmt"
	"selfbotat-v2/bot"
	"time"

	"selfbotat-v2/bot/database"
	"selfbotat-v2/bot/config"
  "selfbotat-v2/bot/logger"
	"selfbotat-v2/bot/api"

	_ "selfbotat-v2/bot/commands"
	client "selfbotat-v2/bot/client"
)

var Config config.BotConfig

func main() {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@localhost/%s?sslmode=disable", 
		Config.Postgres.User, 
		Config.Postgres.Password,
		Config.Postgres.Database,
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
	Config, err = config.LoadConfig()
	if err != nil {
		panic(err)
	}
	
	bot.StartTime = time.Now()

	api.LoadQueries()
}
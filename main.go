package main

import (
	"fmt"
	"net/http"

	"selfbotat-v2/bot/api"
	"selfbotat-v2/bot/config"
	"selfbotat-v2/bot/client"
	"selfbotat-v2/bot/logger"
	"selfbotat-v2/bot/database"

	_ "net/http/pprof"
	_ "selfbotat-v2/bot/commands"
)

func main() {
	go func() {
		Log.Debug.Println(http.ListenAndServe("localhost:1337", nil))
	}()

	connStr := fmt.Sprintf(
		"postgres://%s:%s@localhost/%s?sslmode=disable", 
		config.Config.Postgres.User, 
		config.Config.Postgres.Password,
		config.Config.Postgres.Database,
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
	err = config.LoadConfig()
	if err != nil {
		panic(err)
	}

	err = api.LoadQueries()
	if err != nil {
		Log.Error.Panicln("Error loading TLA queries", err)
	}
}

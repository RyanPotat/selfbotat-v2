package commands

import (
	"selfbotat-v2/bot"
	"selfbotat-v2/bot/client"
	"selfbotat-v2/bot/logger"
	"strconv"
)

func init() {
	bot.AddCmd(bot.Command{
		Name: "clear",
		Aliases: []string{"c"},
		Execute: func(msg *bot.MessageData) {
			if len(msg.Args) == 0 {
				Log.Warn.Println("No amount specified")
				return
			}

			count, err := strconv.Atoi(msg.Args[0])
			if err != nil {	count = 1	}

			for i := 0; i < count; i++ {
				client.Say(msg.Channel.Login, "/clear")
			}
		},
	})
}


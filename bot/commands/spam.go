package commands

import (
	"strconv"
	"strings"

	"selfbotat-v2/bot"
	"selfbotat-v2/bot/client"
	"selfbotat-v2/bot/logger"
)

func init() {
	bot.AddCmd(bot.Command{
		Name:     "spam",
		Execute: func(msg *bot.MessageData) {
			count, err := strconv.Atoi(msg.Args[0])
			if err != nil {
				Log.Error.Print("Error parsing spam count", err)
				return
			}
			for i := 0; i < count; i++ {
				client.Say(msg.Channel.Login, strings.Join(msg.Args[1:], " "))
			}
		},
	})
}
package commands

import (
	"selfbotat-v2/bot"
	"selfbotat-v2/bot/database"
)

func init() {
	bot.AddCmd(bot.Command{
		Name:     "join",
		Execute: func(msg *bot.MessageData) {
			if len(msg.Args) == 0 {
				return
			}
			bot.Client.Join(msg.Args[0])
			db.NewChannel(msg.Channel.ID)
		},
	})
}
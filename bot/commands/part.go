package commands

import (
	"selfbotat-v2/bot"
	"selfbotat-v2/bot/database"
)

func init() {
	bot.AddCmd(bot.Command{
		Name:     "part",
		Execute: func(msg *bot.MessageData) {
			if len(msg.Args) == 0 {
				return
			}
			bot.Client.Depart(msg.Channel.Login)
			db.RemoveChannel(msg.Channel.ID)
		},
	})
}
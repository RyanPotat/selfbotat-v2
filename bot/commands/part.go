package commands

import (
	"selfbotat-v2/bot"
	"selfbotat-v2/bot/client"
	"selfbotat-v2/bot/database"
	"selfbotat-v2/bot/types"
)

func init() {
	bot.AddCmd(types.Command{
		Name:     "part",
		Execute: func(msg *types.MessageData) {
			if len(msg.Args) == 0 {
				return
			}

			client.Part(msg.Channel.Login)
			db.RemoveChannel(msg.Channel.ID)
		},
	})
}
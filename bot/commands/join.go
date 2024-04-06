package commands

import (
	"selfbotat-v2/bot"
	"selfbotat-v2/bot/client"
	"selfbotat-v2/bot/types"
	"selfbotat-v2/bot/database"
)

func init() {
	bot.AddCmd(types.Command{
		Name:     "join",
		Execute: func(msg *types.MessageData) {
			if len(msg.Args) == 0 {
				return
			}
			client.Join(msg.Args[0])
			db.NewChannel(msg.Channel.ID)
		},
	})
}
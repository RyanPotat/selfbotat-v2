package commands

import (
	"strconv"
	"strings"
	"selfbotat-v2/bot"
)

func init() {
	bot.AddCmd(bot.Command{
		Name:     "spam",
		Execute: func(msg *bot.MessageData) {
			count, err := strconv.Atoi(msg.Args[0])
			if err != nil {
				bot.Log.Println("Error parsing spam count", err)
				return
			}
			for i := 0; i < count; i++ {
				bot.Client.Say(msg.Channel.Login, strings.Join(msg.Args[1:], " "))
			}
		},
	})
}
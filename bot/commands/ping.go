package commands

import (
	"fmt"
	"time"
	"runtime"
	"strings"
	"selfbotat-v2/bot"
	"selfbotat-v2/bot/client"
)

func init() {
	bot.AddCmd(bot.Command{
		Name:     "ping",
		Aliases:  []string{"pong"},
		Execute: func(msg *bot.MessageData) {
				var used runtime.MemStats
		
				runtime.ReadMemStats(&used)
		
				response := []string{
					"GoldPLZ 🏓",
					fmt.Sprintf("Usage: %d MiB", used.TotalAlloc / 1024 / 1024),
					fmt.Sprintf("Uptime: %s", time.Since(bot.StartTime).Round(time.Second)),
				}
		
				client.Say(msg.Channel.Login, strings.Join(response, " ● "))
		},
	})
}
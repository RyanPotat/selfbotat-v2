package commands

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"selfbotat-v2/bot"
	"selfbotat-v2/bot/client"
	"selfbotat-v2/bot/types"
	"selfbotat-v2/bot/utils"
)

func init() {
	bot.AddCmd(types.Command{
		Name:     "ping",
		Aliases:  []string{"pong"},
		Execute: func(msg *types.MessageData) {
				var used runtime.MemStats
		
				runtime.ReadMemStats(&used)
		
				response := []string{
					"GoldPLZ 🏓",
					fmt.Sprintf("Usage: %d MiB", used.HeapInuse / 1024 / 1024),
					fmt.Sprintf("Allocated: %d MiB", (used.HeapIdle + used.HeapInuse) / 1024 / 1024),
					fmt.Sprintf("Uptime: %s", utils.Humanize(time.Since(bot.StartTime), 3)),
				}
		
				client.Say(msg.Channel.Login, strings.Join(response, " ● "))
		},
	})
}
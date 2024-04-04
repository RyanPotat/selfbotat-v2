package commands

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"selfbotat-v2/bot"
	"selfbotat-v2/bot/client"
	"selfbotat-v2/bot/logger"
)

var (
	stopSpam  chan struct{}
	spamming  bool
	spamMutex sync.Mutex
)

func init() {
	bot.AddCmd(bot.Command{
		Name:     "spam",
		Aliases: []string{"s"},
		Execute: func(msg *bot.MessageData) {
			// stop the spamming
			if stop, ok := msg.Params["stop"]; ok && stop.(bool) {
				if spamming {
					Log.Warn.Println("Stopping spam")
					close(stopSpam)
					return
				} 
				Log.Warn.Println("Not spamming")
				return
			}

			// prevent multiple spam loop
			if spamming {
				Log.Warn.Println("Already spamming")
				return
			}

			spamMutex.Lock()
      defer spamMutex.Unlock()

			var count, interval = 1, 1
			var err error

			// parse count
			if c, ok := msg.Params["c"]; ok {
				count, err = strconv.Atoi(c.(string))
				if err != nil {	count = 1	}
			}

			// parse interval
			if i, ok := msg.Params["i"]; ok {
				interval, err = strconv.Atoi(i.(string))
				if err != nil { interval = 1 }
			}

			// filter out params
			filteredArgs := bot.Filter(msg.Args, func(arg string) bool {
				return !(strings.HasPrefix(arg, "-") || strings.Contains(arg, ":"))
			})

			if len(filteredArgs) == 0 {
				Log.Warn.Print("No message to spam")
			}

			message := strings.Join(filteredArgs, " ")

			// announcements or actions
			prefix := ""
			if a, ok := msg.Params["a"]; ok && a.(bool) {
				prefix = "/announce "
			} else if m, ok := msg.Params["m"]; ok && m.(bool) {
				prefix = "/me "
			}

			// fill message
			if f, ok := msg.Params["f"]; ok && f.(bool) {
				for len(message) < 500 {
					message += " " + message
				}
				lastSpaceIndex := strings.LastIndex(message[:500], " ")

				if lastSpaceIndex != -1 {
						message = message[:lastSpaceIndex]
				} else { message = message[:500] }
			}

			stopSpam = make(chan struct{})
			go spam(msg.Channel.Login, count, interval, prefix + message)
		},
	})
}

func spam(channel string, count, interval int, message string) {
	spamming = true

	for i := 0; i < count; i++ {
		client.Say(channel, message)
		if interval > 0 {
			time.Sleep(time.Duration(interval) * time.Millisecond)
		}
		select {
		case <-stopSpam:
			spamming = false
			return
		default:
		}
	}
	spamming = false
}

package commands

import (
	"fmt"
	"strings"
	"os/exec"

	"selfbotat-v2/bot"
)

func init() {
	bot.AddCmd(bot.Command{
		Name:     "eval",
		Execute: func(msg *bot.MessageData) {
			code := fmt.Sprintf(`
			  package main
				
				import (
					"fmt"
					"time"
					"runtime"
					"strings"
				)	

				%s
			`, strings.Join(msg.Args, " "))

			fmt.Println(code)
			cmd := exec.Command("go", "run")
			cmd.Stdin = strings.NewReader(code)
		
			output, err := cmd.CombinedOutput()
			if err != nil {
				bot.Client.Say(msg.Channel.Login, string(err.Error()))
				return
			}
		
			bot.Client.Say(msg.Channel.Login, string(output))
		},
	})
}
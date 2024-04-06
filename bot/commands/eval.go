package commands

import (
	"fmt"
	"strings"
	"os/exec"

	"selfbotat-v2/bot"
	"selfbotat-v2/bot/types"
	"selfbotat-v2/bot/client"
)

func init() {
	bot.AddCmd(types.Command{
			Name: "eval",
			Execute: func(msg *types.MessageData) {
					code := strings.Join(msg.Args, " ")

					output, err := runGoCode(code)
					if err != nil {
						  client.Say(msg.Channel.Login, fmt.Sprintf("Error: %s", err.Error()))
							return
					}

					client.Say(msg.Channel.Login, output)
			},
	})
}

func runGoCode(code string) (string, error) {
	cmd := exec.Command("go", "run")

	cmd.Stdin = strings.NewReader(code)

	output, err := cmd.CombinedOutput()
	if err != nil {
			return "", fmt.Errorf("execution failed: %s", err.Error())
	}

	return string(output), nil
}

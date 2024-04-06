package bot

import (
	"time"

	"selfbotat-v2/bot/types"
)

var (
	Cmds = []types.Command{}
	StartTime time.Time
)

func init() {
	StartTime = time.Now()
}

func AddCmd(cmd types.Command) {
	Cmds = append(Cmds, cmd)
}

func FindCmd(command string) *types.Command {
	for _, cmd := range Cmds {
		if cmd.Name == command {
			return &cmd
		} else if len(cmd.Aliases) > 0 {
			for _, alias := range cmd.Aliases {
				if alias == command {
					return &cmd
				}
			}
		}
	}

	return nil
}

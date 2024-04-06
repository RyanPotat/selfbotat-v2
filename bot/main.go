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
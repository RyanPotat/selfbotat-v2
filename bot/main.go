package bot

import (
	"os"
	"log"
	"time"

	"github.com/gempir/go-twitch-irc/v4"
)

type Command struct {
	Name string
	Aliases []string
	Whitelist []string
	Execute func(msg *MessageData) 
	Params map[string]interface{}
	Requires string
}

type User struct {
	ID string
	Name string
	Login string
}

type Channel struct {
	ID string
	Login string
}

type MessageData struct {
	User User
	Channel Channel
	Text string
	Params map[string]interface{}
	Hashtags []string
	Command string
	Args []string
	Raw twitch.PrivateMessage
}

var (
	Client *twitch.Client
	StartTime time.Time
	Log = log.New(os.Stdout, "SelfBot: ", log.LstdFlags)
	cmds = []Command{}
)

func AddCmd(cmd Command) {
	cmds = append(cmds, cmd)
}




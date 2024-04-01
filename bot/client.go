package bot

import (
	"fmt"
	//"slices"
	"strings"

	"selfbotat-v2/bot/config"
	"selfbotat-v2/bot/database"

	twitch "github.com/gempir/go-twitch-irc/v4"
)

var (
	cfg config.Config
	joinedChannels map[string]bool
)

func InitClient() (error, bool) {
	Log.Println("Connecting to Twitch...")

	cfg = config.GetConfig()

	Client = twitch.NewClient(cfg.Login, fmt.Sprintf("oauth:%s", cfg.Password))

	Client.OnPrivateMessage(parseMessage)

	Client.OnSelfJoinMessage(func(msg twitch.UserJoinMessage) {
		Log.Printf("Joined #%s\n", msg.Channel)
	})

	Client.OnSelfPartMessage(func(msg twitch.UserPartMessage) {
		Log.Printf("Parted #%s\n", msg.Channel)
	})

	Client.OnConnect(func() {
		Log.Println("Connected to Twitch")
	})

	Client.SetJoinRateLimiter(twitch.CreateVerifiedRateLimiter())

	JoinChannels()

	err := Client.Connect()
	if err != nil {
		return err, false
	}

	return nil, true
}

func JoinChannels() {
	channels, err := db.GetChannels()
	if err != nil {
		Log.Println("Error getting channels", err)
		return
	}

	Client.Join("potatbotat")
	Client.Join("notohh")

	for _, channel := range channels {
		Client.Join(channel.Login)
	} 
}

func parseMessage(msg twitch.PrivateMessage) {
	_, found := db.GetUserByID(msg.User.ID)
	if !found {
		db.NewUser(
			msg.User.ID,
			msg.User.Name,
			msg.User.DisplayName,
		)
	}

	user := User{
		ID: msg.User.ID,
		Login: msg.User.Name,
		Name: msg.User.DisplayName,
	}

	channel := Channel{
		ID: msg.RoomID,
		Login: msg.Channel,
	}

	rawText := strings.TrimPrefix(msg.Message, cfg.Prefix)
	parts := strings.Split(strings.TrimSpace(rawText), " ")
	//params := slices.Contains(parts[1:], ":")

	cmd := parts[0]
	args := parts[1:]

	handleMessage(&MessageData{
		User: user,
		Channel: channel,
		Text: msg.Message,
		Command: cmd,
		Args: args,
		Params: make(map[string]interface{}),
		Hashtags: nil,
		Raw: msg,
	})
}

// temp handler while I figure some structure
func handleMessage(msg *MessageData) {
	if msg.User.ID != cfg.Id {
		return
	}

	if !strings.HasPrefix(msg.Text, cfg.Prefix) {
		return
	}

	for _, cmd := range cmds {
		if cmd.Name == msg.Command {
			cmd.Execute(msg)
		}
	}
}

package client

import (
	"fmt"
	//"slices"
	"strings"
	"sync"

	"selfbotat-v2/bot"
	"selfbotat-v2/bot/config"
	"selfbotat-v2/bot/database"
	"selfbotat-v2/bot/logger"
	"selfbotat-v2/bot/types"
	"selfbotat-v2/bot/utils"

	"github.com/gempir/go-twitch-irc/v4"
)

const (
	ensurePoolSize = 50
	maxClientSize  = 100
	maxMessageSize = 500
)

var (
	totalJoined     = make(map[string]bool)
	clientCount     = 1
	ClientPool      = make(map[int]*ChatClient)
	lastClientIndex int
)

type ChatClient struct {
	*twitch.Client
	joinedChannels map[string]bool
}

func Create() {
	ensurePool()
}

func ensurePool() {
	var wg sync.WaitGroup
	var poolMutex sync.Mutex

	for i := 0; i < ensurePoolSize; i++ {
		wg.Add(1)
		go func(clientID int) {
			defer wg.Done()

			poolMutex.Lock()
			client := createClient(clientID)
			poolMutex.Unlock()

			err := client.Connect()
			if err != nil {
				Log.Error.Println("Error connecting to Twitch", err)
			}
		}(i + 1)
	}

	wg.Wait()
}

func createClient(clientID int) *ChatClient {

	client := &ChatClient{
		twitch.NewClient(
			config.Config.Twitch.Login,
			fmt.Sprintf("oauth:%s", config.Config.Twitch.Password),
		),
		make(map[string]bool),
	}

	ClientPool[clientID] = client

	applyListeners(client, clientID)

	return client
}

func applyListeners(client *ChatClient, clientID int) {
	client.OnPrivateMessage(func(msg twitch.PrivateMessage) {
		parseMessage(msg)
	})

	client.OnSelfJoinMessage(func(msg twitch.UserJoinMessage) {
		Log.Info.Printf("Joined #%s", msg.Channel)
		client.joinedChannels[msg.Channel] = true
		totalJoined[msg.Channel] = true
	})

	client.OnSelfPartMessage(func(msg twitch.UserPartMessage) {
		Log.Info.Printf("Parted #%s", msg.Channel)
		client.joinedChannels[msg.Channel] = false
		totalJoined[msg.Channel] = false
	})

	client.OnConnect(func() {
		clientCount += 1
		if clientID == ensurePoolSize {
			Log.Debug.Printf("IRC connected")
			joinChannels()
		}
	})

	client.SetJoinRateLimiter(twitch.CreateVerifiedRateLimiter())
}

func joinChannels() {
	channels, err := db.GetChannels()
	if err != nil {
		Log.Error.Print("Error getting channels", err)
		return
	}

	Join("notohh")

	for _, channel := range channels {
		Join(channel.Login)
	}
}

func Join(channel string) bool {
	if totalJoined[channel] {
		return false
	}

	joined := false
	for _, client := range ClientPool {
		if len(client.joinedChannels) < maxClientSize {
			client.Join(channel)
			joined = true
			break
		}
	}

	if !joined {
		client := createClient(clientCount + 1)
		ClientPool[clientCount+1].Join(channel)

		err := client.Connect()
		if err != nil {
			Log.Error.Print("Error connecting to Twitch", err)
		}
	}

	return true
}

func Part(channel string) bool {
	parted := false
	for _, client := range ClientPool {
		if client.joinedChannels[channel] {
			client.Depart(channel)
			parted = true
		}
	}

	return parted
}

func Say(channel, message string) {
	lastClientIndex = (lastClientIndex + 1) % (len(ClientPool))
	client := ClientPool[lastClientIndex+1]

	if len(message) > maxMessageSize {
		message = message[:maxMessageSize]
	}

	client.Say(channel, message)
}

func parseMessage(msg twitch.PrivateMessage) {
	_, found := db.GetUser(msg.User.ID, false)
	if !found {
		db.NewUser(
			msg.User.ID,
			msg.User.Name,
			msg.User.DisplayName,
		)
	}

	user := types.User{
		ID:    msg.User.ID,
		Login: msg.User.Name,
		Name:  msg.User.DisplayName,
	}

	channel := types.Channel{
		ID:    msg.RoomID,
		Login: msg.Channel,
	}

	rawText := strings.TrimPrefix(msg.Message, config.Config.Prefix)
	parts := strings.Split(strings.TrimSpace(rawText), " ")
	cmd := parts[0]
	args := parts[1:]
	params := utils.CreateParams(args)
	hashtags := utils.CreateHashtags(msg.Message)

	handleMessage(&types.MessageData{
		User:     user,
		Channel:  channel,
		Text:     msg.Message,
		Command:  cmd,
		Args:     args,
		Params:   params,
		Hashtags: hashtags,
		Raw:      msg,
	})
}

// temp handler while I figure some structure
func handleMessage(msg *types.MessageData) {
	// check prefix
	if !strings.HasPrefix(msg.Text, config.Config.Prefix) {
		return
	}

	// find command
	cmd := bot.FindCmd(msg.Command)
	if cmd == nil {	return }

	// ignore self and non-whitelisted users
	isSelf := msg.User.ID == config.Config.Twitch.Id
	//isWhitelisted := slices.Contains(cmd.Whitelist, msg.User.ID)
	if !isSelf {
		return
	}

	// remove params from command arguments
	// for key, _ := range cmd.Params {
	// 	for _, arg := range msg.Args {
	// 		isStringParam := strings.HasPrefix(arg, key + ":")
	// 		isBoolParam := strings.HasPrefix(arg, "-" + key)
	// 		if isStringParam || isBoolParam {
	// 			msg.Args = slices.Remove(msg.Args, arg)
	// 		}
	// 	}
	// }

	// execute command
	cmd.Execute(msg)
}

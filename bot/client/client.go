package client

import (
	"fmt"
	"sync"
	//"slices"
	"strings"

	"selfbotat-v2/bot"
	"selfbotat-v2/bot/config"
	"selfbotat-v2/bot/database"

	Logger "selfbotat-v2/bot/logger"

	twitch "github.com/gempir/go-twitch-irc/v4"
)

const (
	ensurePoolSize = 50
	maxClientSize  = 100
)

var (
	cfg             config.Config
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
	cfg = config.GetConfig()

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
					Logger.Error("Error connecting to Twitch", err)
			}
		}(i + 1)
	}

	wg.Wait()
}

func createClient(clientID int) *ChatClient {

	client := &ChatClient{
			twitch.NewClient(cfg.Login, fmt.Sprintf("oauth:%s", cfg.Password)),
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
		Logger.Infof("Joined #%s\n", msg.Channel)
		client.joinedChannels[msg.Channel] = true
		totalJoined[msg.Channel] = true
	})

	client.OnSelfPartMessage(func(msg twitch.UserPartMessage) {
		Logger.Infof("Parted #%s\n", msg.Channel)
		client.joinedChannels[msg.Channel] = false
		totalJoined[msg.Channel] = false
	})

	client.OnConnect(func() {
		clientCount += 1
		if clientID == ensurePoolSize {
			Logger.Debug("IRC connected")
			joinChannels()
		}
	})

	client.SetJoinRateLimiter(twitch.CreateVerifiedRateLimiter())
}

func joinChannels() {
	channels, err := db.GetChannels()
	if err != nil {
		Logger.Error("Error getting channels", err)
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
		ClientPool[clientCount + 1].Join(channel)

		err := client.Connect()
		if err != nil {
				Logger.Error("Error connecting to Twitch", err)
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
	client := ClientPool[lastClientIndex + 1]

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

	user := bot.User{
		ID: msg.User.ID,
		Login: msg.User.Name,
		Name: msg.User.DisplayName,
	}

	channel := bot.Channel{
		ID: msg.RoomID,
		Login: msg.Channel,
	}

	rawText := strings.TrimPrefix(msg.Message, cfg.Prefix)
	parts := strings.Split(strings.TrimSpace(rawText), " ")
	//params := slices.Contains(parts[1:], ":")

	cmd := parts[0]
	args := parts[1:]

	handleMessage(&bot.MessageData{
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
func handleMessage(msg *bot.MessageData) {
	if msg.User.ID != cfg.Id {
		return
	}

	if !strings.HasPrefix(msg.Text, cfg.Prefix) {
		return
	}

	for _, cmd := range bot.GetCmds() {
		if cmd.Name == msg.Command {
			cmd.Execute(msg)
		}
	}
}

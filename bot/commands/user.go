package commands

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"selfbotat-v2/bot"
	"selfbotat-v2/bot/api"
	"selfbotat-v2/bot/client"
	"selfbotat-v2/bot/types"
	"selfbotat-v2/bot/utils"
)

func init() {
	bot.AddCmd(types.Command{
		Name: "user",
		Aliases: []string{"u"},
		Execute: func(msg *types.MessageData) {
			var input string
			if len(msg.Args) == 0 {
				input = msg.User.ID
			} else {
				input = msg.Args[0]
			}

			var wg sync.WaitGroup
			wg.Add(2)

			var user *types.TwitchUser
			go func() {
				defer wg.Done()
				var err error
				user, err = api.GetUserOrError(input)
				if err != nil {
					client.Say(msg.Channel.Login, "⚠️ " + err.Error())
					return
				}
	
				if user.ID == "" && user.Key == "" {
					userDoesntExist(msg)
					return
				}
	
				if user.Reason != "" {
					if user.Reason == "UNKNOWN" {
						userDoesntExist(msg)
						return
					}
	
					banType := user.Reason
	
					user, err = api.GetTwitchUser(input)
					user.Reason = banType
					if err != nil {
						client.Say(msg.Channel.Login, "⚠️ " + err.Error())
						return
					}
				}
			}()

			var ok bool
			var afkUser *types.AFKData
			go func() {
				defer wg.Done()
				afkUser, ok = api.GetAFK(input)
				if !ok { return	}
			}()

			wg.Wait()

			response := map[string]interface{}{
				"": nameBuilder(user),
				"ID": user.ID,
				"Roles": roleBuilder(user.Roles),
				"Followers": strconv.Itoa(user.Followers.TotalCount),
				"Following": strconv.Itoa(user.Follows.TotalCount),
				"Chatters": strconv.Itoa(user.Channel.Chatters.Count),
				"Prefix": user.EmotePrefix.Name,
				"Bio": user.Description,
				"Created": getAgo(user.CreatedAt, 3),
				"Last Live": offlineBuilder(user),
				"Currently AFK": afkBuilder(afkUser),
				"🔴 Live for": liveBuilder(user),
			}

			// Enforce order of response (idk how else to do this lol)
			keys := []string{
				"", "ID", "Roles", "Followers", "Following", "Chatters", "Prefix", "Bio", "Created", "Last Live", "Currently AFK", "🔴 Live for",
			}

			out := ""
			for _, key := range keys {
				if value, ok := response[key]; ok {
						if value != nil && value != "" && value != "0" {
								out += key + ": " + value.(string) + " ● "
						}
				}
			}

			client.Say(msg.Channel.Login, strings.TrimSuffix(out, " ● "))
		},
	})
}

func nameBuilder(user *types.TwitchUser) string {
	name := user.Login
	if strings.ToLower(user.DisplayName) == name {
		name = user.DisplayName
	}

	switch user.Reason {
	case INDEF:
			name = fmt.Sprintf("⛔ @%s is indefinitely banned ryanpo1Despair", name)
	case TEMP:
			name = fmt.Sprintf("⚠ @%s is temporarily banned", name)
	case DMCA:
			name = fmt.Sprintf("⚠ @%s is temporarily banned for DMCA", name)
	case DISABLED:
			name = fmt.Sprintf("⛔ @%s deactivated their account ryanpo1Despair", name)
	default:
			name = fmt.Sprintf("@%s", name)
	}
	return name
}

func roleBuilder(roles types.Roles) string {
	var out []string

	v := reflect.ValueOf(roles)

	for i := 0; i < v.NumField(); i++ {
		rerole := v.Field(i)

		// Check if the field is a pointer (nil)
		if rerole.Kind() == reflect.Ptr {
			if !rerole.IsNil() {
				b := rerole.Elem().Bool()
				if b {
					value := v.Type().Field(i).Name
					out = append(out, roleMap[value])
				}
			}
		// Handle bool directly
		} else {
			b := rerole.Bool()
			if b {
				value := v.Type().Field(i).Name
				out = append(out, roleMap[value])
			}
		}
	}
	
	return strings.Join(out, ", ")
}

func getAgo(t time.Time, l int) string {
	if t.IsZero() { return ""	}
	return utils.Humanize(time.Since(t), l) + " ago"
}

func offlineBuilder(user *types.TwitchUser) string {
	if user.Stream != nil { return "" }
	return getAgo(user.LastBroadcast.StartedAt, 2)
}

func liveBuilder(user *types.TwitchUser) string {
	if user.Stream == nil { return "" }

	liveSince := getAgo(user.Stream.CreatedAt, 2)
	game := ""
	if user.BroadcastSettings.Game.DisplayName != "" {
		name := user.BroadcastSettings.Game.DisplayName
		game = fmt.Sprintf(" streaming \"%s\"", name)
	}

  viewers := strconv.Itoa(user.Stream.ViewersCount)
	viewerString := ""
	if viewers != "0" {
		viewerString = fmt.Sprintf(" with %s viewers", viewers)
	}

	return fmt.Sprintf("%s%s%s", liveSince, game, viewerString)
}

func userDoesntExist(msg *types.MessageData) {
	// TODO: check name availability here
	client.Say(msg.Channel.Login, "⚠️ User not found")
}

func afkBuilder(afkUser *types.AFKData) string {
	if afkUser == nil {
		return ""
	}

	afkAgo := getAgo(afkUser.Started, 2)
	return fmt.Sprintf("%s (%s)", afkUser.Text, afkAgo)
}

var roleMap = map[string]string{
	"IsPartner":            "Partner",
	"IsAffiliate":          "Affiliate",
	"IsStaff":              "Staff",
	"IsAdmin":              "Admin",
	"IsGlobalMod":          "Global Mod",
	"IsExtensionDeveloper": "Extension Developer",
}

const (
	INDEF = "TOS_INDEFINITE"
	TEMP  = "TOS_TEMPORARY"
	DMCA = "DMCA"
	DISABLED = "DEACTIVATED"
)
package commands

import (
	"fmt"
	"reflect"
	"selfbotat-v2/bot"
	"selfbotat-v2/bot/api"
	"selfbotat-v2/bot/client"
	"selfbotat-v2/bot/types"
	"strconv"
	"strings"
	"time"
)

func init() {
	bot.AddCmd(bot.Command{
		Name: "user",
		Aliases: []string{"u"},
		Execute: func(msg *bot.MessageData) {
			if len(msg.Args) == 0 {
				client.Say(msg.Channel.Login, "FeelsDankMan")
				return
			}

			input := msg.Args[0]

			user, err := api.GetUserOrError(input)
			if err != nil {
				client.Say(msg.Channel.Login, "⚠️ " + err.Error())
				return
			}

			if user.ID == "" && user.Key == "" {
				userDoesntExist(msg)
				return
			}

			banType := ""
			if user.Reason != "" {
				if user.Reason == "UNKNOWN" {
					userDoesntExist(msg)
					return
				}

				banType = user.Reason

				user, err = api.GetTwitchUser(input)
				if err != nil {
					client.Say(msg.Channel.Login, "⚠️ " + err.Error())
					return
				}
			}

			name := user.Login
			if strings.ToLower(user.DisplayName) != name {
				name = user.DisplayName
			}

			switch banType {
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

			response := map[string]interface{}{
				"ID": user.ID,
				"Roles": buildRoles(user.Roles),
				"Followers": strconv.Itoa(user.Followers.TotalCount),
				"Following": strconv.Itoa(user.Follows.TotalCount),
				"Chatters": strconv.Itoa(user.Channel.Chatters.Count),
				"Prefix": user.EmotePrefix.Name,
				"Bio": user.Description,
				"Created": getAgo(user.CreatedAt, 3),
				"Last Live": offlineBuilder(user),
				"🔴 Live for": liveBuilder(user),
			}

			// Enforce order of response (idk how else to do this lol)
			keys := []string{
				"ID", "Roles", "Followers", "Following", "Chatters", "Prefix", "Bio", "Created", "Last Live", "🔴 Live for",
			}

			out := ""
			for _, key := range keys {
				if value, ok := response[key]; ok {
						if value != nil && value != "" && value != "0" {
								out += key + ": " + value.(string) + " ● "
						}
				}
			}

			out = name + " ● " + strings.TrimSuffix(out, " ● ")
			client.Say(msg.Channel.Login, out)
		},
	})
}

func buildRoles(roles types.Roles) string {
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
	return bot.Humanize(time.Since(t), l) + " ago"
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

func userDoesntExist(msg *bot.MessageData) {
	// TODO: check name availability here
	client.Say(msg.Channel.Login, "⚠️ User not found")
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
package api

import (
	"bytes"
	"encoding/json"
	"strings"

	"os"
	"regexp"
	"selfbotat-v2/bot/config"
	"selfbotat-v2/bot/logger"
	"selfbotat-v2/bot/types"
)

var (
	queries types.Queries
)

func LoadQueries() types.Queries {
	queriesJSON, err := os.ReadFile("queries.json")
	if err != nil {
		Log.Error.Panicln("error loading config: %w", err)
		return nil
	}

	err = json.Unmarshal(queriesJSON, &queries)
	if err != nil {
		Log.Error.Panicln("Error loading TLA queries", err)
		return nil
	}

	return queries
}

func GetTwitchUser(loginOrID string) (*types.TwitchUser, error) {
	re := regexp.MustCompile(`^[0-9]{1,25}`)

	input := make(map[string]string)
	if re.MatchString(loginOrID) {
		input["id"] = loginOrID
	} else {
		input["login"] = loginOrID
	}

	queryJSON, err := json.Marshal(types.TLAOP{
		OperationName: "User",
		Query: queries["User"],
		Variables: input,
	})
	if err != nil {
		Log.Error.Println("Error marshalling query", err)
		return nil, err
	}

	res, err := MakeRequest(
		"POST", 
		config.Config.Twitch.TLAURI, 
		config.Config.Twitch.TLAHeaders,
		bytes.NewBuffer(queryJSON),
	)
	if err != nil {
		Log.Error.Println("Error making request", err)
		return nil, err
	}

	defer res.Body.Close()

	var response types.TLAUserRes
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		Log.Error.Println("Error decoding TLA response ", err, response)
		return nil, err
	}

	if len(response.Errors) > 0 {
		Log.Error.Println("Error in TLA response", response.Errors)
		return nil, err
	}

	return &response.Data.User, nil
}

func GetUserOrError(loginOrID string) (*types.TwitchUser, error) {
	re := regexp.MustCompile(`^[0-9]{1,25}`)

	var inputType string
	input := make(map[string]string)
	if re.MatchString(loginOrID) {
		inputType = "Id"
	} else {
		inputType = "Login"
	}

	input[strings.ToLower(inputType)] = loginOrID

	queryJSON, err := json.Marshal(types.TLAOP{
		Query: queries["userResultBy" + inputType],
		Variables: input,
	})
	if err != nil {
		Log.Error.Println("Error marshalling query", err)
		return nil, err
	}

	res, err := MakeRequest(
		"POST", 
		config.Config.Twitch.TLAURI, 
		config.Config.Twitch.TLAHeaders,
		bytes.NewBuffer(queryJSON),
	)
	if err != nil {
		Log.Error.Println("Error making request", err)
		return nil, err
	}

	defer res.Body.Close()

	var response types.TLAUserOrErrorRes
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		Log.Error.Println("Error decoding TLA response ", err, response)
		return nil, err
	}

	if len(response.Errors) > 0 {
		Log.Error.Println("Error in TLA response", response.Errors)
		return nil, err
	}

	return &response.Data.User, nil
}



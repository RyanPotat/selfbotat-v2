package api

import (
	"encoding/json"
	"selfbotat-v2/bot/logger"
	"selfbotat-v2/bot/types"
)

const uri = "https://supinic.com/api/"

func GetAFK(login string) (*types.AFKData, bool) {
	res, err := Make.Get(uri + "bot/afk/check?username=" + login)
	if err != nil {
		Log.Error.Println("Supi/AFK - Error getting afk: ", err)
	}

	defer res.Body.Close()

	var afkData types.SupiAFKRes
	err = json.NewDecoder(res.Body).Decode(&afkData)
	if err != nil {
		Log.Error.Println("Supi/AFK - Error parsing afk data: ", err)
		return nil, false
	}

	if afkData.StatusCode != 200{
		if afkData.Error != nil {
			Log.Error.Println("Supi/AFK - Error getting afk data: ", afkData.Error.Message)
		} else {
			Log.Error.Println("Supi/AFK - Error getting afk data: ", afkData.StatusCode)
		}

		return nil, false
	}

	if afkData.Data == nil {
		return nil, false
	}

	return afkData.Data.Status, true
}
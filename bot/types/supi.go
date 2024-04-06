package types

import "time"

type GenericSupiRes struct {
	StatusCode int        `json:"statusCode"`
	Timestamp  int        `json:"timestamp"`
	Error      *SupiError `json:"error"`
}

type SupiAFKRes struct {
	GenericSupiRes
	Data *AFKStatus `json:"data"`
}

type AFKStatus struct {
	Status *AFKData `json:"status"`
}

type AFKData struct {
	ID        int       `json:"id"`
	UserAlias int       `json:"userAlias"`
	Started   time.Time `json:"started"`
	Text      string    `json:"text"`
	Status    string    `json:"status"`
	Name      string    `json:"name"`
	TwitchID  string    `json:"twitchId"`
}

type SupiError struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

package types

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type TwitchConfig struct {
	Login      string            `json:"login"`
	Id         string            `json:"id"`
	Password   string            `json:"password"`
	ClientID   string            `json:"client_id"`
	TLAHeaders map[string]string `json:"tlaHeaders"`
}

type URIs struct {
	PubSub string `json:"pubsub"`
	TLA    string `json:"tla"`
}

type BotConfig struct {
	Prefix   string         `json:"prefix"`
	Postgres PostgresConfig `json:"postgres"`
	Twitch   TwitchConfig   `json:"twitch"`
	URIs     URIs           `json:"uris"`
}
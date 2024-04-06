package types

import (
	"time"
)

type TwitchUsers struct {
	Users []TwitchUser `json:"users"`
}

type TwitchUser struct {
	ID                 string            `json:"id"`
	Login              string            `json:"login"`
	DisplayName        string            `json:"displayName"`
	ChatColor          string            `json:"chatColor"`
	Description        string            `json:"description"`
	ProfileImageURL    string            `json:"profileImageURL"`
	BannerImageURL     string            `json:"bannerImageURL"`
	CreatedAt          time.Time         `json:"createdAt"`
	UpdatedAt          time.Time         `json:"updatedAt"`
	HasPrime           bool              `json:"hasPrime"`
	HasTurbo           bool              `json:"hasTurbo"`
	Followers          Followers         `json:"followers"`
	Follows            Followers         `json:"follows"`
	EmoticonPrefix     EmoticonPrefix    `json:"emoticonPrefix"`
	BroadcastSettings  BroadcastSettings `json:"broadcastSettings"`
	Stream             *Stream           `json:"stream"`
	LastBroadcast      LastBroadcast     `json:"lastBroadcast"`
	Roles              Roles             `json:"roles"`
	Key                string            `json:"key,omitempty"`
	Reason             string            `json:"reason,omitempty"`
	Banned             bool              `json:"banned,omitempty"`
	PrimaryTeam        *PrimaryTeam      `json:"primaryTeam"`
	Channel            UserChannel       `json:"channel"`
	Team               *string           `json:"team"`
	TeamName           *string           `json:"teamName"`
	ChatterCount       int               `json:"chatterCount"`
	EmotePrefix        EmoticonPrefix    `json:"emotePrefix"`
}

type UserChannel struct {
	Chatters           TwitchChatters  `json:"chatters"`
	RecentChatMessages []TwitchMessage `json:"recentChatMessages"`
}

type TwitchMessage struct {
	Content   TwitchMessageContent `json:"content"`
	DeletedAt time.Time            `json:"deletedAt"`
	ID        string               `json:"id"`
	Sender    TwitchUser           `json:"sender"`
	ParentMessage interface{}      `json:"parentMessage"`
}

type TwitchMessageContent struct {
	Fragments []TwitchMessageFragment `json:"fragments"`
	Text      string                 `json:"text"`
}

type TwitchMessageFragment struct {
	Content TwitchMessageContent `json:"content"`
	Text    string              `json:"text"`
}

type TwitchChatters struct {
	Broadcasters []ChatterLogin `json:"broadcasters"`
	Moderators   []ChatterLogin `json:"moderators"`
	Viewers      []ChatterLogin `json:"viewers"`
	Staff        []ChatterLogin `json:"staff"`
	Vips         []ChatterLogin `json:"vips"`
	Count        int            `json:"count"`
}

type ChatterLogin struct {
	Login string `json:"login"`
}

type PrimaryTeam struct {
	DisplayName string `json:"displayName"`
	Name        string `json:"name"`
	ID          int    `json:"id"`
}

type BroadcastSettings struct {
	Title      string `json:"title"`
	Language   string `json:"language"`
	IsMature   bool   `json:"isMature"`
	Game       Game   `json:"game"`
}

type Game struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

type EmoticonPrefix struct {
	IsEditable bool   `json:"isEditable"`
	Name       string `json:"name"`
	State      string `json:"state"`
}

type Followers struct {
	TotalCount int `json:"totalCount,omitempty"`
}

type LastBroadcast struct {
	Game      Game      `json:"game"`
	StartedAt time.Time `json:"startedAt"`
	Title     string    `json:"title"`
}

type Roles struct {
	IsPartner           bool  `json:"isPartner"`
	IsStaff             *bool `json:"isStaff,omitempty"`
	IsAffiliate         bool  `json:"isAffiliate"`
	IsSiteAdmin         *bool `json:"isSiteAdmin,omitempty"`
	IsGlobalMod         *bool `json:"isGlobalMod,omitempty"`
	IsExtensionDeveloper bool `json:"isExtensionDeveloper"`
}

type Stream struct {
	Type                string        `json:"type"`
	Bitrate             int           `json:"bitrate"`
	BroadcastLanguage   string        `json:"broadcastLanguage"`
	BroadcasterSoftware string        `json:"broadcasterSoftware"`
	CreatedAt           time.Time     `json:"createdAt"`
	Language            string        `json:"language"`
	PreviewImageURL     string        `json:"previewImageURL"`
	ViewersCount        int           `json:"viewersCount"`
	Tags                []interface{} `json:"tags"`
	AverageFPS          int           `json:"averageFPS"`
	ClipCount           int           `json:"clipCount"`
	Height              int           `json:"height"`
	Width               int           `json:"width"`
}
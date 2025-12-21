package model

import "time"

type MusicState struct {
	Client        ClientInfo    `json:"client"`
	Metadata      MediaMetadata `json:"metadata"`
	PlaybackState PlaybackState `json:"playbackState"`
	UpdatedAt     time.Time     `json:"updatedAt"`
	Version       int64         `json:"version"`
}

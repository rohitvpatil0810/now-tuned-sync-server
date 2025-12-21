package model

type PlaybackState string

const (
	PlaybackPlaying PlaybackState = "playing"
	PlaybackPaused  PlaybackState = "paused"
	PlaybackStopped PlaybackState = "stopped"
)

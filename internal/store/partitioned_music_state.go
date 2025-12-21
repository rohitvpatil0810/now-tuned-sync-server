package store

import (
	"github.com/rohitvpatil0810/now-tuned-sync-server/internal/model"
)

type MusicState = model.MusicState
type PlaybackState = model.PlaybackState

type PartitionedMusicState struct {
	playing    map[string]MusicState // key = ClientID
	notPlaying map[string]MusicState
}

func NewPartitionedMusicState() *PartitionedMusicState {
	return &PartitionedMusicState{
		playing:    make(map[string]MusicState),
		notPlaying: make(map[string]MusicState),
	}
}

func (p *PartitionedMusicState) Update(state MusicState) {
	clientID := state.Client.ClientID

	if state.PlaybackState == model.PlaybackPlaying {
		delete(p.notPlaying, clientID)
		p.playing[clientID] = state
	} else {
		delete(p.playing, clientID)
		p.notPlaying[clientID] = state
	}
}

func (p *PartitionedMusicState) Remove(clientID string) {
	delete(p.playing, clientID)
	delete(p.notPlaying, clientID)
}

func (p *PartitionedMusicState) Winner() *MusicState {
	for _, state := range p.playing {
		return &state
	}
	for _, state := range p.notPlaying {
		return &state
	}
	return nil
}

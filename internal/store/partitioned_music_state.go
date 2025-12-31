package store

import (
	"sync"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/rohitvpatil0810/now-tuned-sync-server/internal/model"
)

type MusicState = model.MusicState
type PlaybackState = model.PlaybackState

type PartitionedMusicState struct {
	playing    map[string]MusicState // key = ClientID
	notPlaying map[string]MusicState
	clients    []chan *MusicState
	mu         sync.RWMutex
}

func NewPartitionedMusicState() *PartitionedMusicState {
	return &PartitionedMusicState{
		playing:    make(map[string]MusicState),
		notPlaying: make(map[string]MusicState),
	}
}

func (p *PartitionedMusicState) RegisterClient(ch chan *MusicState) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.clients = append(p.clients, ch)
}

func (p *PartitionedMusicState) UnregisterClient(ch chan *MusicState) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for i, c := range p.clients {
		if c == ch {
			p.clients = append(p.clients[:i], p.clients[i+1:]...)
			break
		}
	}
}

func (p *PartitionedMusicState) broadcast() {
	winner := p.Winner()
	for _, ch := range p.clients {
		select {
		case ch <- winner:
		default:
			// Skip if the channel is not ready to receive
		}
	}
}

func (p *PartitionedMusicState) Update(state MusicState) {
	p.mu.Lock()
	defer p.mu.Unlock()

	clientID := state.Client.ClientID
	oldWinner := p.Winner()

	if state.PlaybackState == model.PlaybackPlaying {
		delete(p.notPlaying, clientID)
		p.playing[clientID] = state
	} else {
		delete(p.playing, clientID)
		p.notPlaying[clientID] = state
	}

	// Broadcast only if the winner has changed
	newWinner := p.Winner()
	if !areStatesEqual(oldWinner, newWinner) {
		p.broadcast()
	}
}

func (p *PartitionedMusicState) Remove(clientID string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	oldWinner := p.Winner()

	delete(p.playing, clientID)
	delete(p.notPlaying, clientID)

	// Broadcast the updated winner
	newWinner := p.Winner()
	if !areStatesEqual(oldWinner, newWinner) {
		p.broadcast()
	}
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

func areStatesEqual(a, b *MusicState) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	return cmp.Equal(a, b,
		cmp.Comparer(func(t1, t2 time.Time) bool { return true }),
	)
}

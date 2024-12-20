package playlist

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type Playlist struct {
	Name   string   `json:"name"`
	Tracks []string `json:"tracks"`
	mu     sync.Mutex
}

// NewPlaylist creates a new, empty playlist
func NewPlaylist(name string) *Playlist {
	return &Playlist{Name: name, Tracks: []string{}}
}

// AddTrack adds a new track to the playlist after validation
func (p *Playlist) AddTrack(track string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, err := os.Stat(track); os.IsNotExist(err) {
		return fmt.Errorf("track does not exist: %s", track)
	}
	// Additional format checks can be added here if needed

	p.Tracks = append(p.Tracks, track)
	return nil
}

// RemoveTrack removes a track from the playlist by index
func (p *Playlist) RemoveTrack(index int) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if index < 0 || index >= len(p.Tracks) {
		return fmt.Errorf("invalid index")
	}
	p.Tracks = append(p.Tracks[:index], p.Tracks[index+1:]...)
	return nil
}

// Save saves the playlist to a file
func (p *Playlist) Save(filePath string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize playlist: %v", err)
	}
	return os.WriteFile(filePath, data, 0644)
}

// Load loads a playlist from a file
func Load(filePath string) (*Playlist, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read playlist file: %v", err)
	}
	var p Playlist
	err = json.Unmarshal(data, &p)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize playlist: %v", err)
	}
	return &p, nil
}

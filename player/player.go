package player

import (
    "fmt"
    "log"
    "os"
    "sync"
    "time"

    "github.com/rickcollette/megasound"
    "github.com/rickcollette/megasound/mp3"
    "github.com/rickcollette/megasound/speaker"
    "github.com/rickcollette/megasound/effects"
)
type MP3Player struct {
	streamer      megasound.StreamSeekCloser
	format        megasound.Format
	volumeCtrl    *effects.Volume
	volumeMutex   sync.Mutex
	playbackMutex sync.Mutex
	paused        bool
	done          chan bool
}
// NewMP3Player initializes a new MP3 player.
func NewMP3Player(filePath string) (*MP3Player, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open MP3 file: %v", err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		f.Close()
		return nil, fmt.Errorf("failed to decode MP3: %v", err)
	}

	// Initialize speaker once
	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		streamer.Close()
		f.Close()
		return nil, fmt.Errorf("failed to initialize speaker: %v", err)
	}

	volumeCtrl := &effects.Volume{
		Streamer: streamer,
		Base:     2,
		Volume:   0, // No volume change by default
		Silent:   false,
	}

	return &MP3Player{
		streamer:   streamer,
		format:     format,
		volumeCtrl: volumeCtrl,
		paused:     true,
		done:       make(chan bool),
	}, nil
}

// Play starts or resumes playback.
func (p *MP3Player) Play() {
	p.playbackMutex.Lock()
	defer p.playbackMutex.Unlock()

	if p.paused {
		speaker.Play(megasound.Seq(p.volumeCtrl, megasound.Callback(func() {
			p.done <- true
		})))
		p.paused = false
		log.Println("Playback started")
	}
}

// Pause stops playback.
func (p *MP3Player) Pause() {
	p.playbackMutex.Lock()
	defer p.playbackMutex.Unlock()

	if !p.paused {
		speaker.Lock()
		p.paused = true
		speaker.Unlock()
		speaker.Clear()
		log.Println("Playback paused")
	}
}

// Paused checks if the player is currently paused.
func (p *MP3Player) Paused() bool {
	p.playbackMutex.Lock()
	defer p.playbackMutex.Unlock()
	return p.paused
}

// SetVolume adjusts the volume of the player.
// volumeLevel: 0.0 (mute) to 1.0 (max).
func (p *MP3Player) SetVolume(volumeLevel float64) {
	p.volumeMutex.Lock()
	defer p.volumeMutex.Unlock()

	// Convert 0.0-1.0 range to natural volume control.
	p.volumeCtrl.Volume = (volumeLevel - 0.5) * 2 // Scale to range [-1, 1]
	p.volumeCtrl.Silent = volumeLevel <= 0
	log.Printf("Volume set to %.1f", p.volumeCtrl.Volume)
}

// Close releases resources.
func (p *MP3Player) Close() {
	p.streamer.Close()
	speaker.Clear()
	log.Println("MP3Player closed")
}
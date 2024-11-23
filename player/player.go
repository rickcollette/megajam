package player

import (
    "fmt"
    "log"
    "os"
    "sync"
    "time"

    "github.com/gopxl/beep/v2"
    "github.com/gopxl/beep/v2/mp3"
    "github.com/gopxl/beep/v2/speaker"
    "github.com/gopxl/beep/v2/volume"
)

type MP3Player struct {
    streamer      beep.StreamSeekCloser
    format        beep.Format
    volumeCtrl    *volume.Volume
    volumeMutex   sync.Mutex
    playbackMutex sync.Mutex
    paused        bool
    done          chan bool
}

// NewMP3Player initializes a new MP3 player
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

    // Create volume control streamer
    volumeCtrl := volume.New(volume.Caps{
        Min: -96, // Mute
        Max: 0,    // Max volume
    }, format.SampleRate, streamer)

    return &MP3Player{
        streamer:   streamer,
        format:     format,
        volumeCtrl: volumeCtrl,
        paused:     true,
        done:       make(chan bool),
    }, nil
}

// Play starts or resumes playback
func (p *MP3Player) Play() {
    p.playbackMutex.Lock()
    defer p.playbackMutex.Unlock()

    if p.paused {
        p.volumeMutex.Lock()
        p.volumeCtrl.Paused = false
        p.volumeMutex.Unlock()

        speaker.Play(beep.Seq(p.volumeCtrl, beep.Callback(func() {
            p.done <- true
        })))
        p.paused = false
        log.Println("Playback started")
    }
}

// Pause stops playback
func (p *MP3Player) Pause() {
    p.playbackMutex.Lock()
    defer p.playbackMutex.Unlock()

    if !p.paused {
        p.volumeMutex.Lock()
        p.volumeCtrl.Paused = true
        p.volumeMutex.Unlock()

        speaker.Clear()
        p.paused = true
        log.Println("Playback paused")
    }
}

// Paused checks if the player is currently paused
func (p *MP3Player) Paused() bool {
    p.playbackMutex.Lock()
    defer p.playbackMutex.Unlock()
    return p.paused
}

// SetVolume adjusts the volume of the player
// volumeLevel: 0.0 (mute) to 1.0 (max)
func (p *MP3Player) SetVolume(volumeLevel float64) {
    p.volumeMutex.Lock()
    defer p.volumeMutex.Unlock()

    // Convert 0.0-1.0 to -96 to 0 dB
    db := -96 + volumeLevel*96
    p.volumeCtrl.Gain = db
    log.Printf("Volume set to %.1f dB", db)
}

// Close releases resources
func (p *MP3Player) Close() {
    p.streamer.Close()
    speaker.Clear()
    log.Println("MP3Player closed")
}

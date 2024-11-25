package db

import (
	"gorm.io/gorm"
)

type Track struct {
	gorm.Model
	Title    string
	Artist   string
	Album    string
	Path     string
	Duration string
	BPM      int
	Key      string
}

type Playlist struct {
	gorm.Model
	Name   string
	Tracks []Track `gorm:"many2many:playlist_tracks;"`
}

type Crate struct {
	gorm.Model
	Name   string
	Tracks []Track `gorm:"many2many:crate_tracks;"`
}

type CuePoint struct {
    gorm.Model
    TrackID uint    `gorm:"index"` // Foreign key to Track
    Name    string  // Optional name for the cue point
    Time    float64 // Time in seconds
}

type Loop struct {
    gorm.Model
    TrackID uint    `gorm:"index"` // Foreign key to Track
    Name    string  // Optional name for the loop
    Start   float64 // Start time in seconds
    End     float64 // End time in seconds
}

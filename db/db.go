package db

import (
	"fmt"
	"megajam/logger"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDatabase initializes the database using the provided path.
func InitDatabase(databasePath string) {
	logger.Logger.Println("Initializing database")
	var err error
	DB, err = gorm.Open(sqlite.Open(databasePath), &gorm.Config{})
	if err != nil {
		logger.Logger.Printf("Failed to connect to the database: %v", err)
		return
	}

	logger.Logger.Println("Database connection established!")
	if err := migrateDB(); err != nil {
		logger.Logger.Printf("Database migration failed: %v", err)
	}
	logger.Logger.Println("DB Migration completed")
}

func migrateDB() error {
	logger.Logger.Println("DB Migration starting..")
	return DB.AutoMigrate(&Track{}, &Playlist{}, &Crate{}, &CuePoint{}, &Loop{})
}

func AddCuePoint(trackID uint, name string, time float64) error {
	cuePoint := CuePoint{TrackID: trackID, Name: name, Time: time}
	return DB.Create(&cuePoint).Error
}

func GetCuePoints(trackID uint) ([]CuePoint, error) {
	var cuePoints []CuePoint
	err := DB.Where("track_id = ?", trackID).Find(&cuePoints).Error
	return cuePoints, err
}

func AddLoop(trackID uint, name string, start, end float64) error {
	if start >= end {
		return fmt.Errorf("start time must be less than end time")
	}
	loop := Loop{TrackID: trackID, Name: name, Start: start, End: end}
	return DB.Create(&loop).Error
}

func GetLoops(trackID uint) ([]Loop, error) {
	var loops []Loop
	err := DB.Where("track_id = ?", trackID).Find(&loops).Error
	return loops, err
}

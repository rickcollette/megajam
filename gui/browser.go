package gui

import (
	"bytes"
	"fmt"
	"image"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"megajam/db"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/dhowden/tag"
)

// ExtractAlbumArt extracts album art from an MP3 file.
func ExtractAlbumArt(filePath string) *canvas.Image {
	f, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error opening file %s: %v", filePath, err)
		return nil
	}
	defer f.Close()

	meta, err := tag.ReadFrom(f)
	if err != nil {
		log.Printf("Error reading tags from file %s: %v", filePath, err)
		return nil
	}

	picture := meta.Picture()
	if picture == nil {
		log.Printf("No album art found in file %s", filePath)
		return nil
	}

	img, _, err := image.Decode(bytes.NewReader(picture.Data))
	if err != nil {
		log.Printf("Error decoding album art from file %s: %v", filePath, err)
		return nil
	}

	canvasImage := canvas.NewImageFromImage(img)
	canvasImage.SetMinSize(fyne.NewSize(100, 100)) // Fixed size for thumbnails
	return canvasImage
}

// createEnhancedBrowserSection creates the browser interface with search and track options.
func createEnhancedBrowserSection(myWindow fyne.Window) *fyne.Container {
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Search...")

	var tracks []db.Track
	if err := db.DB.Find(&tracks).Error; err != nil {
		log.Printf("Error loading tracks: %v", err)
	}

	filteredTracks := tracks
	var mutex sync.Mutex
	selectedTrackID := uint(0)

	trackList := widget.NewList(
		func() int { return len(filteredTracks) },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("Title"),
				widget.NewLabel("Artist"),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			mutex.Lock()
			defer mutex.Unlock()
			if id < 0 || id >= len(filteredTracks) {
				return
			}
			track := filteredTracks[id]
			labels := item.(*fyne.Container).Objects
			labels[0].(*widget.Label).SetText(track.Title)
			labels[1].(*widget.Label).SetText(track.Artist)
		},
	)

	trackList.OnSelected = func(id widget.ListItemID) {
		if id >= 0 && id < len(filteredTracks) {
			selectedTrackID = filteredTracks[id].ID
		}
	}

	searchEntry.OnChanged = func(query string) {
		mutex.Lock()
		defer mutex.Unlock()
		query = strings.ToLower(query)
		filteredTracks = []db.Track{}
		for _, track := range tracks {
			if strings.Contains(strings.ToLower(track.Title), query) || strings.Contains(strings.ToLower(track.Artist), query) {
				filteredTracks = append(filteredTracks, track)
			}
		}
		trackList.Refresh()
	}

	addCueButton := widget.NewButton("Add Cue", func() {
		if selectedTrackID == 0 {
			dialog.ShowInformation("No Track Selected", "Please select a track to add a cue point.", myWindow)
			return
		}

		// Create fields for the form.
		timeEntry := widget.NewEntry()
		nameEntry := widget.NewEntry()

		// Show the dialog form.
		dialog.ShowForm("Add Cue Point", "Add", "Cancel", []*widget.FormItem{
			widget.NewFormItem("Time (seconds)", timeEntry),
			widget.NewFormItem("Name", nameEntry),
		}, func(confirm bool) {
			if !confirm {
				return
			}

			// Parse time input.
			time, err := strconv.ParseFloat(timeEntry.Text, 64)
			if err != nil {
				dialog.ShowError(fmt.Errorf("invalid time format"), myWindow)
				return
			}

			// Add the cue point to the database.
			if err := db.AddCuePoint(selectedTrackID, nameEntry.Text, time); err != nil {
				dialog.ShowError(err, myWindow)
			} else {
				dialog.ShowInformation("Success", "Cue point added successfully.", myWindow)
			}
		}, myWindow)
	})

	addLoopButton := widget.NewButton("Add Loop", func() {
		if selectedTrackID == 0 {
			dialog.ShowInformation("No Track Selected", "Please select a track to add a loop.", myWindow)
			return
		}

		// Create fields for the form.
		startEntry := widget.NewEntry()
		endEntry := widget.NewEntry()
		nameEntry := widget.NewEntry()

		// Show the dialog form.
		dialog.ShowForm("Add Loop", "Add", "Cancel", []*widget.FormItem{
			widget.NewFormItem("Start Time (seconds)", startEntry),
			widget.NewFormItem("End Time (seconds)", endEntry),
			widget.NewFormItem("Name", nameEntry),
		}, func(confirm bool) {
			if !confirm {
				return
			}

			// Parse start and end times.
			start, err1 := strconv.ParseFloat(startEntry.Text, 64)
			end, err2 := strconv.ParseFloat(endEntry.Text, 64)
			if err1 != nil || err2 != nil || start >= end {
				dialog.ShowError(fmt.Errorf("invalid start or end time"), myWindow)
				return
			}

			// Add the loop to the database.
			if err := db.AddLoop(selectedTrackID, nameEntry.Text, start, end); err != nil {
				dialog.ShowError(err, myWindow)
			} else {
				dialog.ShowInformation("Success", "Loop added successfully.", myWindow)
			}
		}, myWindow)
	})

	return container.NewVBox(
		searchEntry,
		trackList,
		container.NewHBox(addCueButton, addLoopButton),
	)
}

package gui

import (
    "log"

    "megajam/config"
    "megajam/playlist"
    "megajam/player"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/dialog"
    "fyne.io/fyne/v2/widget"
)
// CreateToolbar creates the top toolbar.
func CreateToolbar() *fyne.Container {
	// Add buttons or actions as needed.
	openButton := widget.NewButton("Open", func() {
		log.Println("Open clicked")
		// Implement Open functionality
	})
	saveButton := widget.NewButton("Save", func() {
		log.Println("Save clicked")
		// Implement Save functionality
	})
	settingsButton := widget.NewButton("Settings", func() {
		log.Println("Settings clicked")
		// Open settings dialog
	})
	exitButton := widget.NewButton("Exit", func() {
		log.Println("Exit clicked")
		// Handle exit logic
	})

	// Toolbar layout
	return container.NewHBox(
		openButton,
		saveButton,
		settingsButton,
		exitButton,
	)
}

// CreateWaveformVisualizer creates a placeholder for the waveform visualizer.
func CreateWaveformVisualizer() *fyne.Container {
	// Replace this with the actual waveform visualizer implementation.
	return container.NewCenter(widget.NewLabel("Waveform Visualizer"))
}

func CreateGUI() {
	// Initialize Fyne app
	myApp := app.NewWithID("com.example.megajam")
	myWindow := myApp.NewWindow("Megajam DJ Controller")

	// Load configuration
	appConfig, err := config.LoadConfig("config/config.json")
	if err != nil {
		dialog.ShowError(err, myWindow)
		myWindow.ShowAndRun()
		return
	}

	// Validate configuration
	err = config.ValidateConfig(appConfig)
	if err != nil {
		dialog.ShowError(err, myWindow)
		myWindow.ShowAndRun()
		return
	}

	// Apply custom theme
	myApp.Settings().SetTheme(NewCustomTheme(&appConfig.Theme))

	// Resize window based on configuration
	myWindow.Resize(fyne.NewSize(
		float32(appConfig.Layout.WindowWidth),
		float32(appConfig.Layout.WindowHeight),
	))

	// Initialize player with a default track
	mp3Player, err := player.NewMP3Player("path/to/default.mp3") // Replace with actual default track path
	if err != nil {
		dialog.ShowError(err, myWindow)
	}
	defer mp3Player.Close()

	// Initialize sections with handlers
	leftDeck := CreateDeckSection(
		"Left Deck",
		"Track A",
		"3:45",
		"120",
		nil, // mp3Image can be set dynamically
		func() {
			if mp3Player.Paused() {
				mp3Player.Play()
				log.Println("Left Deck: Play")
			} else {
				mp3Player.Pause()
				log.Println("Left Deck: Pause")
			}
		},
		func() { log.Println("Left Deck: Sync button clicked") },
		func(value float64) {
			mp3Player.SetVolume(value / 100)
			log.Printf("Left Deck: Volume set to %.1f%%", value)
		},
	)

	rightDeck := CreateDeckSection(
		"Right Deck",
		"Track B",
		"4:10",
		"126",
		nil, // mp3Image can be set dynamically
		func() {
			if mp3Player.Paused() {
				mp3Player.Play()
				log.Println("Right Deck: Play")
			} else {
				mp3Player.Pause()
				log.Println("Right Deck: Pause")
			}
		},
		func() { log.Println("Right Deck: Sync button clicked") },
		func(value float64) {
			mp3Player.SetVolume(value / 100)
			log.Printf("Right Deck: Volume set to %.1f%%", value)
		},
	)

	mixer := CreateMixerSection()

	// Initialize playlist browser
	currentPlaylist := playlist.NewPlaylist("New Playlist")
	addTrackButton := widget.NewButton("Add Track", func() {
		dialog.ShowFileOpen(func(uri fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, myWindow)
				return
			}
			if uri == nil {
				return
			}
			err = currentPlaylist.AddTrack(uri.URI().Path())
			if err != nil {
				dialog.ShowError(err, myWindow)
				return
			}
			dialog.ShowInformation("Success", "Track added to playlist!", myWindow)
		}, myWindow)
	})
	removeTrackButton := widget.NewButton("Remove Track", nil)

	browserSection, selectedTrackChan := createEnhancedBrowserSection(currentPlaylist, addTrackButton, removeTrackButton, myWindow)

	// Handle Remove Track Button Click
	removeTrackButton.OnTapped = func() {
		select {
		case selectedIndex := <-selectedTrackChan:
			if err := currentPlaylist.RemoveTrack(selectedIndex); err != nil {
				dialog.ShowError(err, myWindow)
			} else {
				dialog.ShowInformation("Success", "Track removed from playlist!", myWindow)
			}
		default:
			dialog.ShowInformation("Info", "No track selected to remove.", myWindow)
		}
	}

	// Toolbar at the top
	toolbar := CreateToolbar()

	// Waveform visualizer below the toolbar
	waveformVisualizer := CreateWaveformVisualizer()

	// Decks and mixer in the middle
	decksAndMixer := container.NewGridWithColumns(
		3,
		leftDeck,
		mixer,
		rightDeck,
	)

	// Final layout
	mainLayout := container.NewBorder(
		container.NewVBox(toolbar, waveformVisualizer, decksAndMixer), // Top: Toolbar, Waveform visualizer, and decks/mixer
		browserSection,                                               // Bottom: Browser
		nil,                                                          // Left
		nil,                                                          // Right
		nil,                                                          // Center
	)

	// Set and show window content
	myWindow.SetContent(mainLayout)
	myWindow.ShowAndRun()
}
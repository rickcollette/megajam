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
                // Update button text to "Pause" (Requires access to the button)
                log.Println("Left Deck: Play")
            } else {
                mp3Player.Pause()
                // Update button text to "Play" (Requires access to the button)
                log.Println("Left Deck: Pause")
            }
        },
        func() {
            log.Println("Left Deck: Sync button clicked")
            // Implement sync functionality here
        },
        func(value float64) {
            mp3Player.SetVolume(value / 100) // Assuming volume range is 0-100
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
                // Update button text to "Pause" (Requires access to the button)
                log.Println("Right Deck: Play")
            } else {
                mp3Player.Pause()
                // Update button text to "Play" (Requires access to the button)
                log.Println("Right Deck: Pause")
            }
        },
        func() {
            log.Println("Right Deck: Sync button clicked")
            // Implement sync functionality here
        },
        func(value float64) {
            mp3Player.SetVolume(value / 100) // Assuming volume range is 0-100
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
            // Refresh the browser section if necessary
            // Implement a way to refresh thumbnails and track list
        }, myWindow)
    })
    removeTrackButton := widget.NewButton("Remove Track", func() {
        // Implement track removal logic using selectedTrackIndex
        // This requires accessing the selectedTrackIndex received from the channel
        // We'll set up a goroutine to listen to the channel
    })
    browserSection, selectedTrackChan := createEnhancedBrowserSection(currentPlaylist, addTrackButton, removeTrackButton, myWindow)

    // Handle Remove Track Button Click
    removeTrackButton.OnTapped = func() {
        selectedIndex := -1
        select {
        case idx := <-selectedTrackChan:
            selectedIndex = idx
        default:
            // No selection made
        }

        if selectedIndex >= 0 {
            err := currentPlaylist.RemoveTrack(selectedIndex)
            if err != nil {
                dialog.ShowError(err, myWindow)
                return
            }
            dialog.ShowInformation("Success", "Track removed from playlist!", myWindow)
            log.Printf("Track at index %d removed from playlist", selectedIndex)
            // Optionally, refresh the browser section
            browserSection.Refresh()
        } else {
            dialog.ShowInformation("Info", "No track selected to remove.", myWindow)
        }
    }

    // Layout components
    mainLayout := container.NewBorder(
        nil,                    // Top
        browserSection,         // Bottom
        nil,                    // Left
        nil,                    // Right
        container.NewHBox(leftDeck, mixer, rightDeck), // Center: Decks and mixer
    )

    // Set and show window content
    myWindow.SetContent(mainLayout)
    myWindow.ShowAndRun()
}

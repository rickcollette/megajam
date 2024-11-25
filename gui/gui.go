package gui

import (
	"encoding/json"
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"strings"

	"megajam/config"
	"megajam/db"
	"megajam/logger"
	"megajam/player"
	"megajam/playlist"
	"megajam/waveform"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// CreateWaveformVisualizer creates the waveform visualizer using the waveform package.
func CreateWaveformVisualizer(audioData []int32) *fyne.Container {
	logger.Logger.Println("Starting waveform visualizer")
	wave := waveform.NewWaveform(audioData)
	wave.SetMinSize(fyne.NewSize(400, 100))
	wave.OverrideForeground = true
	wave.OverrideForegroundColor = color.NRGBA{R: 255, G: 0, B: 0, A: 255} // Red waveform

	return container.NewVBox(
		wave,
	)
}

// CreateGUI initializes and runs the GUI application with the given AppConfig.
func CreateGUI(appConfig *config.AppConfig) {
	logger.Logger.Println("Starting UI")
	var err error

	// Initialize Fyne app.
	myApp := app.NewWithID("com.megalithiCode.megajam")
	myWindow := myApp.NewWindow("Megajam DJ Controller")

	// Initialize database with the configured path.
	db.InitDatabase(appConfig.DatabasePath)

	// Apply custom theme.
	logger.Logger.Println("Loading theme...")
	myApp.Settings().SetTheme(NewCustomTheme(&appConfig.Theme))
	myWindow.Resize(fyne.NewSize(
		float32(appConfig.Layout.WindowWidth),
		float32(appConfig.Layout.WindowHeight),
	))

	// Handle application mode
	switch appConfig.Mode {
	case "party":
		logger.Logger.Println("Running in Party mode.")
		// Implement party mode-specific configurations or UI changes
	case "hardcore":
		logger.Logger.Println("Running in Hardcore mode.")
		// Implement hardcore mode-specific configurations or UI changes
	default:
		logger.Logger.Printf("Unknown mode '%s', defaulting to Party mode.", appConfig.Mode)
	}

	// Check for tracks in the database.
	logger.Logger.Println("Checking for tracks in the database...")
	var trackCount int64
	if err := db.DB.Model(&db.Track{}).Count(&trackCount).Error; err != nil {
		dialog.ShowError(err, myWindow)
		logger.Logger.Fatalf("Failed to query database: %v", err)
		return
	}

	// Only initialize the MP3 player if there are tracks in the database.
	var mp3Player *player.MP3Player
	if trackCount > 0 {
		logger.Logger.Println("Tracks found in the database. Initializing player...")
		firstTrack := db.Track{}
		if err := db.DB.First(&firstTrack).Error; err != nil {
			dialog.ShowError(err, myWindow)
			logger.Logger.Fatalf("Failed to fetch first track: %v", err)
			return
		}

		mp3Player, err = player.NewMP3Player(firstTrack.Path)
		if err != nil {
			dialog.ShowError(err, myWindow)
			logger.Logger.Printf("Failed to load track '%s': %v", firstTrack.Path, err)
		} else {
			logger.Logger.Printf("Player initialized with track: %s", firstTrack.Title)
		}
	} else {
		logger.Logger.Println("No tracks found in the database. Skipping player initialization.")
	}

	// Ensure player is closed on exit.
	if mp3Player != nil {
		defer mp3Player.Close()
	}

	// Initialize playlist.
	logger.Logger.Println("Initializing Playlist...")
	currentPlaylist := playlist.NewPlaylist("My Playlist")
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
			} else {
				dialog.ShowInformation("Success", "Track added to playlist!", myWindow)
			}
		}, myWindow)
	})
	removeTrackButton := widget.NewButton("Remove Track", func() {
		if len(currentPlaylist.Tracks) == 0 {
			dialog.ShowInformation("Info", "No track to remove.", myWindow)
			return
		}
		// Add logic to remove track from playlist
		logger.Logger.Println("Track removed from playlist.")
	})

	// Create browser section.
	logger.Logger.Println("Initializing track browser...")
	browserSection := container.NewVBox(
		createEnhancedBrowserSection(myWindow),
		container.NewHBox(addTrackButton, removeTrackButton),
	)

	// Create decks with waveform visualization.
	logger.Logger.Println("Creating left Deck...")
	leftDeck := CreateDeckSection(
		"Left Deck", "Track A", "3:45", "120",
		nil,
		func() {
			if mp3Player != nil && mp3Player.Paused() {
				mp3Player.Play()
				logger.Logger.Println("Left Deck: Play")
			} else if mp3Player != nil {
				mp3Player.Pause()
				logger.Logger.Println("Left Deck: Pause")
			}
		},
		func() { logger.Logger.Println("Left Deck: Sync button clicked") },
		func(value float64) {
			if mp3Player != nil {
				mp3Player.SetVolume(value / 100)
				logger.Logger.Printf("Left Deck: Volume set to %.1f%%", value)
			}
		},
	)

	rightDeck := CreateDeckSection(
		"Right Deck", "Track B", "4:10", "126",
		nil,
		func() {
			if mp3Player != nil && mp3Player.Paused() {
				mp3Player.Play()
				logger.Logger.Println("Right Deck: Play")
			} else if mp3Player != nil {
				mp3Player.Pause()
				logger.Logger.Println("Right Deck: Pause")
			}
		},
		func() { logger.Logger.Println("Right Deck: Sync button clicked") },
		func(value float64) {
			if mp3Player != nil {
				mp3Player.SetVolume(value / 100)
				logger.Logger.Printf("Right Deck: Volume set to %.1f%%", value)
			}
		},
	)

	// Placeholder waveform data. Replace with real audio data from selected tracks.
	audioData := []int32{-500, 1000, -2000, 3000, -4000, 5000}
	waveformVisualizer := CreateWaveformVisualizer(audioData)

	// Determine initial background color based on mode.
	var backgroundColor color.Color
	switch appConfig.Mode {
	case "party":
		backgroundColor = color.NRGBA{R: 255, G: 0, B: 0, A: 255} // Red for party mode
	case "hardcore":
		backgroundColor = color.NRGBA{R: 0, G: 0, B: 255, A: 255} // Blue for hardcore mode
	default:
		backgroundColor = color.NRGBA{R: 255, G: 255, B: 255, A: 255} // Default to white
	}
	background := canvas.NewRectangle(backgroundColor)

	// Combine all sections into the main layout.
	mainLayout := container.NewVBox(
		CreateToolbar(myWindow, background),
		waveformVisualizer,
		container.NewGridWithColumns(2, leftDeck, rightDeck),
		browserSection,
	)
	content := container.NewMax(background, mainLayout)
	myWindow.SetContent(content)

	// Show the window
	myWindow.ShowAndRun()
}

// CreateToolbar creates the top toolbar with the Settings button.
func CreateToolbar(parent fyne.Window, background *canvas.Rectangle) *fyne.Container {
	openButton := widget.NewButton("Open", func() {
		logger.Logger.Println("Open clicked")
		// Add functionality to load saved playlists or configurations
	})
	saveButton := widget.NewButton("Save", func() {
		logger.Logger.Println("Save clicked")
		// Add functionality to save playlists or configurations
	})
	settingsButton := widget.NewButton("Settings", nil) // Handler will be set later
	exitButton := widget.NewButton("Exit", func() {
		logger.Logger.Println("Exit clicked")
		// Exit the application
		os.Exit(0)
	})

	toolbar := container.NewHBox(
		openButton,
		saveButton,
		settingsButton,
		exitButton,
	)

	// Assign handler to Settings button
	settingsButton.OnTapped = func() {
		openSettingsWindow(parent, background)
	}

	return toolbar
}

// openSettingsWindow creates and displays the settings window.
func openSettingsWindow(parent fyne.Window, background *canvas.Rectangle) {
	// Load current configuration
	appConfig, err := config.LoadConfig("config/config.json")
	if err != nil {
		dialog.ShowError(err, parent)
		return
	}

	// Retrieve available themes
	themes, err := getAvailableThemes()
	if err != nil {
		dialog.ShowError(err, parent)
		return
	}

	// Create theme selection dropdown
	themeOptions := []string{}
	for _, theme := range themes {
		themeOptions = append(themeOptions, theme.Name)
	}
	themeSelect := widget.NewSelect(themeOptions, nil)
	themeSelect.SetSelected(appConfig.ThemeName)

	// Create mode selection dropdown
	modeOptions := []string{}
	modeOptions = append(modeOptions, appConfig.Theme.AllowedModes...)
	modeSelect := widget.NewSelect(modeOptions, nil)
	modeSelect.SetSelected(appConfig.Mode)

	// Create a window for settings
	settingsWindow := fyne.CurrentApp().NewWindow("Settings")

	// Update mode options when theme changes
	themeSelect.OnChanged = func(selectedTheme string) {
		// Load the selected theme's allowed modes
		selectedThemeConfig, err := loadThemeConfig(selectedTheme)
		if err != nil {
			dialog.ShowError(err, parent)
			return
		}

		// Update mode options
		modeSelect.Options = selectedThemeConfig.AllowedModes
		if !contains(selectedThemeConfig.AllowedModes, appConfig.Mode) {
			if len(selectedThemeConfig.AllowedModes) > 0 {
				modeSelect.SetSelected(selectedThemeConfig.AllowedModes[0])
				appConfig.Mode = selectedThemeConfig.AllowedModes[0]
			} else {
				modeSelect.SetSelected("")
				appConfig.Mode = ""
			}
		} else {
			modeSelect.SetSelected(appConfig.Mode)
		}
		modeSelect.Refresh()
	}

	// Create Save and Cancel buttons
	saveButton := widget.NewButton("Save", func() {
		// Get selected theme and mode
		selectedTheme := themeSelect.Selected
		selectedMode := modeSelect.Selected

		// Load the selected theme's config
		selectedThemeConfig, err := loadThemeConfig(selectedTheme)
		if err != nil {
			dialog.ShowError(err, parent)
			return
		}

		// Validate the selected mode
		if !contains(selectedThemeConfig.AllowedModes, selectedMode) {
			dialog.ShowError(fmt.Errorf("selected mode '%s' is not allowed by the theme '%s'", selectedMode, selectedTheme), parent)
			return
		}

		// Update AppConfig
		appConfig.ThemeName = selectedTheme
		appConfig.Theme = selectedThemeConfig
		appConfig.Mode = selectedMode

		// Save the updated configuration
		if err := config.SaveConfig("config/config.json", appConfig); err != nil {
			dialog.ShowError(err, parent)
			return
		}

		// Apply the new theme
		newTheme := NewCustomTheme(&appConfig.Theme)
		fyne.CurrentApp().Settings().SetTheme(newTheme)

		// Update UI elements based on mode
		updateModeSpecificUI(background, appConfig.Mode)

		// Close the settings window
		settingsWindow.Close()
	})
	cancelButton := widget.NewButton("Cancel", func() {
		// Close the settings window without saving
		settingsWindow.Close()
	})
	buttons := container.NewHBox(saveButton, cancelButton)

	// Assemble the settings form
	settingsForm := container.NewVBox(
		widget.NewLabelWithStyle("Settings", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabel("Select Theme:"),
		themeSelect,
		widget.NewLabel("Select Mode:"),
		modeSelect,
		buttons,
	)

	settingsWindow.SetContent(container.NewCenter(settingsForm))
	settingsWindow.Resize(fyne.NewSize(400, 300))
	settingsWindow.Show()
}

// getAvailableThemes retrieves the list of available themes from the themes directory.
func getAvailableThemes() ([]ThemeInfo, error) {
	themesDir := "themes"
	entries, err := os.ReadDir(themesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read themes directory: %w", err)
	}

	themes := []ThemeInfo{}
	for _, entry := range entries {
		if entry.IsDir() {
			themeName := entry.Name()
			themeConfig, err := loadThemeConfig(themeName)
			if err != nil {
				logger.Logger.Printf("Skipping theme '%s' due to error: %v", themeName, err)
				continue
			}
			themes = append(themes, ThemeInfo{Name: themeName, Config: themeConfig})
		}
	}

	return themes, nil
}

// ThemeInfo holds theme name and its configuration.
type ThemeInfo struct {
	Name   string
	Config config.ThemeConfig
}

// loadThemeConfig loads the ThemeConfig for a given theme name.
func loadThemeConfig(themeName string) (config.ThemeConfig, error) {
	themePath := filepath.Join("themes", themeName, "theme.json")
	themeFile, err := os.Open(themePath)
	if err != nil {
		return config.ThemeConfig{}, fmt.Errorf("failed to open theme file '%s': %w", themePath, err)
	}
	defer themeFile.Close()

	var themeConfig config.ThemeConfig
	themeDecoder := json.NewDecoder(themeFile)
	if err := themeDecoder.Decode(&themeConfig); err != nil {
		return config.ThemeConfig{}, fmt.Errorf("failed to parse theme file '%s': %w", themePath, err)
	}

	return themeConfig, nil
}

// contains checks if a slice contains a specific string (case-insensitive).
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, item) {
			return true
		}
	}
	return false
}

// updateModeSpecificUI applies UI changes based on the selected mode.
func updateModeSpecificUI(background *canvas.Rectangle, mode string) {
	var newColor color.Color
	switch mode {
	case "party":
		newColor = color.NRGBA{R: 255, G: 0, B: 0, A: 255} // Red for party mode
	case "hardcore":
		newColor = color.NRGBA{R: 0, G: 0, B: 255, A: 255} // Blue for hardcore mode
	default:
		newColor = color.NRGBA{R: 255, G: 255, B: 255, A: 255} // Default to white
	}

	background.FillColor = newColor
	canvas.Refresh(background)
}

package config

import (
	"encoding/json"
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type ThemeConfig struct {
	BackgroundColor string   `json:"background_color"`
	TextColor       string   `json:"text_color"`
	ButtonColor     string   `json:"button_color"`
	HighlightColor  string   `json:"highlight_color"`
	WaveformColor   string   `json:"waveform_color"`
	KnobColor       string   `json:"knob_color"`
	AllowedModes    []string `json:"allowed_modes"` // New field for allowed modes
}

type LayoutConfig struct {
	WindowWidth  int `json:"window_width"`
	WindowHeight int `json:"window_height"`
}

type AppConfig struct {
	DatabasePath string       `json:"database_path"`
	ThemeName    string       `json:"theme_name"`
	Mode         string       `json:"mode"` // "party" or "hardcore"
	Theme        ThemeConfig  `json:"theme"`
	Layout       LayoutConfig `json:"layout"`
}

var configMutex sync.Mutex // Mutex for thread-safe operations

// LoadConfig reads the main configuration file and returns the AppConfig.
func LoadConfig(filePath string) (*AppConfig, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var config AppConfig
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	applyDefaults(&config)

	// Load theme from themes/$ThemeName/theme.json
	if config.ThemeName != "" {
		themePath := filepath.Join("themes", config.ThemeName, "theme.json")
		themeFile, err := os.Open(themePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open theme file '%s': %w", themePath, err)
		}
		defer themeFile.Close()

		var themeConfig ThemeConfig
		themeDecoder := json.NewDecoder(themeFile)
		if err := themeDecoder.Decode(&themeConfig); err != nil {
			return nil, fmt.Errorf("failed to parse theme file '%s': %w", themePath, err)
		}

		config.Theme = themeConfig
	}

	return &config, nil
}

// SaveConfig writes the current AppConfig back to the configuration file.
func SaveConfig(filePath string, config *AppConfig) error {
	configMutex.Lock()
	defer configMutex.Unlock()

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // For pretty-printing
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("failed to encode config to file: %w", err)
	}

	return nil
}

// applyDefaults sets default values for missing configuration fields.
func applyDefaults(config *AppConfig) {
	if config.Layout.WindowWidth <= 0 {
		config.Layout.WindowWidth = 1024
	}
	if config.Layout.WindowHeight <= 0 {
		config.Layout.WindowHeight = 768
	}
	if config.Theme.BackgroundColor == "" {
		config.Theme.BackgroundColor = "#000000"
	}
	if config.Theme.TextColor == "" {
		config.Theme.TextColor = "#FFFFFF"
	}
	if config.Theme.ButtonColor == "" {
		config.Theme.ButtonColor = "#444444"
	}
	if config.Theme.HighlightColor == "" {
		config.Theme.HighlightColor = "#00FF00"
	}
	if config.Theme.WaveformColor == "" {
		config.Theme.WaveformColor = "#FF0000"
	}
	if config.Theme.KnobColor == "" {
		config.Theme.KnobColor = "#0000FF"
	}
	if config.DatabasePath == "" {
		config.DatabasePath = "music_library.db" // Default database path
	}
	if config.Mode == "" {
		config.Mode = "party" // Default mode
	}
}

// ValidateConfig checks if the configuration is valid.
func ValidateConfig(config *AppConfig) error {
	if config.Layout.WindowWidth <= 0 || config.Layout.WindowHeight <= 0 {
		return fmt.Errorf("invalid window size in config")
	}
	if config.Mode != "party" && config.Mode != "hardcore" {
		return fmt.Errorf("invalid mode '%s' in config: must be 'party' or 'hardcore'", config.Mode)
	}

	// Check if the selected mode is allowed by the theme
	modeAllowed := false
	for _, allowedMode := range config.Theme.AllowedModes {
		if strings.EqualFold(allowedMode, config.Mode) {
			modeAllowed = true
			break
		}
	}
	if !modeAllowed {
		return fmt.Errorf("mode '%s' is not allowed by the selected theme", config.Mode)
	}

	colors := []string{
		config.Theme.BackgroundColor,
		config.Theme.TextColor,
		config.Theme.ButtonColor,
		config.Theme.HighlightColor,
		config.Theme.WaveformColor,
		config.Theme.KnobColor,
	}
	for _, hex := range colors {
		if _, err := ParseHexColor(hex); err != nil {
			return fmt.Errorf("invalid color '%s' in config: %w", hex, err)
		}
	}
	return nil
}

// ParseHexColor parses a hex color string into a color.Color.
func ParseHexColor(hex string) (color.Color, error) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return nil, fmt.Errorf("invalid hex color length")
	}
	var r, g, b uint8
	n, err := fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	if err != nil || n != 3 {
		return nil, fmt.Errorf("invalid hex color format")
	}
	return color.NRGBA{R: r, G: g, B: b, A: 255}, nil
}

package config

import (
	"encoding/json"
	"fmt"
	"image/color"
	"os"
	"strings"
)

type ThemeConfig struct {
	BackgroundColor string `json:"background_color"`
	TextColor       string `json:"text_color"`
	ButtonColor     string `json:"button_color"`
	HighlightColor  string `json:"highlight_color"`
	WaveformColor   string `json:"waveform_color"`
	KnobColor       string `json:"knob_color"`
}

type LayoutConfig struct {
	WindowWidth  int `json:"window_width"`
	WindowHeight int `json:"window_height"`
}

type AppConfig struct {
	Theme  ThemeConfig  `json:"theme"`
	Layout LayoutConfig `json:"layout"`
}

// LoadConfig reads the configuration file and returns the AppConfig.
func LoadConfig(filePath string) (*AppConfig, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var config AppConfig
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// ValidateConfig checks if the configuration is valid.
func ValidateConfig(config *AppConfig) error {
	if config.Layout.WindowWidth <= 0 || config.Layout.WindowHeight <= 0 {
		return fmt.Errorf("invalid window size in config")
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

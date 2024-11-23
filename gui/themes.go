package gui

import (
	"image/color"
	"log"

	"megajam/config"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// CustomTheme implements fyne.Theme based on user configuration.
type CustomTheme struct {
	config *config.ThemeConfig
}

// Ensure CustomTheme implements fyne.Theme.
var _ fyne.Theme = (*CustomTheme)(nil)

// NewCustomTheme creates a new CustomTheme with the given ThemeConfig.
func NewCustomTheme(config *config.ThemeConfig) *CustomTheme {
	return &CustomTheme{config: config}
}

// Color returns the color for the given name and variant.
func (t *CustomTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		c, err := config.ParseHexColor(t.config.BackgroundColor)
		if err != nil {
			log.Printf("Error parsing BackgroundColor: %v", err)
			return theme.DefaultTheme().Color(name, variant)
		}
		return c
	case theme.ColorNameButton:
		c, err := config.ParseHexColor(t.config.ButtonColor)
		if err != nil {
			log.Printf("Error parsing ButtonColor: %v", err)
			return theme.DefaultTheme().Color(name, variant)
		}
		return c
	case theme.ColorNameDisabledButton:
		return color.NRGBA{R: 51, G: 51, B: 51, A: 255} // #333333
	case theme.ColorNamePrimary:
		c, err := config.ParseHexColor(t.config.HighlightColor)
		if err != nil {
			log.Printf("Error parsing HighlightColor: %v", err)
			return theme.DefaultTheme().Color(name, variant)
		}
		return c
	case theme.ColorNameForeground:
		c, err := config.ParseHexColor(t.config.TextColor)
		if err != nil {
			log.Printf("Error parsing TextColor: %v", err)
			return theme.DefaultTheme().Color(name, variant)
		}
		return c
	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

// Icon returns the icon for the given name.
func (t *CustomTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// Font returns the font for the given text style.
func (t *CustomTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

// Size returns the size for the given size name.
func (t *CustomTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

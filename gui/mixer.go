package gui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"megajam/knobs"
)

// CreateMixerSection creates the mixer interface with sliders and knobs.
func CreateMixerSection() *fyne.Container {
	// Crossfader slider
	crossfader := widget.NewSlider(0, 100)
	crossfader.SetValue(50)
	crossfader.OnChanged = func(value float64) {
		log.Printf("Crossfader value changed to: %.1f", value)
		// Implement crossfader functionality here
	}

	// EQ Knobs with handlers
	lowKnob := knobs.CreateKnobWithLabel("Low", -10, 10, func(value float64) {
		log.Printf("Low EQ changed to: %.1f", value)
		// Implement low EQ adjustment here
	})
	midKnob := knobs.CreateKnobWithLabel("Mid", -10, 10, func(value float64) {
		log.Printf("Mid EQ changed to: %.1f", value)
		// Implement mid EQ adjustment here
	})
	highKnob := knobs.CreateKnobWithLabel("High", -10, 10, func(value float64) {
		log.Printf("High EQ changed to: %.1f", value)
		// Implement high EQ adjustment here
	})
	volumeKnob := knobs.CreateKnobWithLabel("Volume", 0, 100, func(value float64) {
		log.Printf("Volume changed to: %.1f", value)
		// Implement volume adjustment here
	})
	balanceKnob := knobs.CreateKnobWithLabel("Balance", -50, 50, func(value float64) {
		log.Printf("Balance changed to: %.1f", value)
		// Implement balance adjustment here
	})

	return container.NewVBox(
		widget.NewLabelWithStyle("Mixer", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		container.NewHBox(lowKnob, midKnob, highKnob),    // EQ knobs
		container.NewHBox(volumeKnob, balanceKnob),      // Volume and balance
		container.NewVBox(widget.NewLabel("Crossfader"), crossfader), // Crossfader
	)
}

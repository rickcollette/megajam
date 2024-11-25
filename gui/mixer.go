package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"megajam/knobs"
)

// CreateMixerSection creates the mixer interface with sliders, knobs, and crossfader.
func CreateMixerSection() *fyne.Container {
	// Volume sliders
	leftVolume := widget.NewSlider(0, 100)
	leftVolume.SetValue(50)
	leftVolume.OnChanged = func(value float64) {
		// Handle left volume change
	}

	rightVolume := widget.NewSlider(0, 100)
	rightVolume.SetValue(50)
	rightVolume.OnChanged = func(value float64) {
		// Handle right volume change
	}

	// EQ knobs
	leftEQ := container.NewVBox(
		widget.NewLabel("Hi"),
		knobs.CreateKnobWithLabel("", -10, 10, func(value float64) {
			// Handle Left Hi EQ
		}),
		widget.NewLabel("Mid"),
		knobs.CreateKnobWithLabel("", -10, 10, func(value float64) {
			// Handle Left Mid EQ
		}),
		widget.NewLabel("Low"),
		knobs.CreateKnobWithLabel("", -10, 10, func(value float64) {
			// Handle Left Low EQ
		}),
	)

	rightEQ := container.NewVBox(
		widget.NewLabel("Hi"),
		knobs.CreateKnobWithLabel("", -10, 10, func(value float64) {
			// Handle Right Hi EQ
		}),
		widget.NewLabel("Mid"),
		knobs.CreateKnobWithLabel("", -10, 10, func(value float64) {
			// Handle Right Mid EQ
		}),
		widget.NewLabel("Low"),
		knobs.CreateKnobWithLabel("", -10, 10, func(value float64) {
			// Handle Right Low EQ
		}),
	)

	// Color knob
	colorKnob := container.NewVBox(
		widget.NewLabel("COLOR"),
		knobs.CreateKnobWithLabel("", 0, 100, func(value float64) {
			// Handle color adjustment
		}),
	)

	// Cue buttons
	cueL := widget.NewButton("Cue L", func() {
		// Handle Cue L button
	})
	cueR := widget.NewButton("Cue R", func() {
		// Handle Cue R button
	})
	cueButtons := container.NewHBox(cueL, cueR)

	// Crossfader slider
	crossfader := widget.NewSlider(0, 100)
	crossfader.SetValue(50)
	crossfader.OnChanged = func(value float64) {
		// Handle crossfader movement
	}
	crossfaderSection := container.NewVBox(
		widget.NewLabel("Crossfader"),
		crossfader,
	)

	// Combine everything
	mixerLayout := container.NewVBox(
		container.NewHBox(leftVolume, leftEQ, colorKnob, rightEQ, rightVolume),
		container.NewVBox(cueButtons, crossfaderSection),
	)
	return mixerLayout
}

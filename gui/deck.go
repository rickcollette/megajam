package gui

import (
    "image/color"
    "log"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/canvas"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
)

// CreateDeckSection creates the deck interface with play/pause, sync, pitch control, and pads.
func CreateDeckSection(deckName, songTitle, timeLeft, bpm string, mp3Image *canvas.Image, playPauseHandler func(), syncHandler func(), pitchHandler func(float64)) *fyne.Container {
    // Sync Button
    syncButton := widget.NewButton("Sync", func() {
        if syncHandler != nil {
            syncHandler()
        } else {
            log.Println("Sync handler not implemented")
        }
    })

    // Circular Display
    titleLabel := widget.NewLabelWithStyle(songTitle, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
    timeLabel := widget.NewLabel(timeLeft)
    bpmLabel := widget.NewLabel("BPM: " + bpm)

    if mp3Image == nil {
        // Provide a default image if none is supplied
        // You need to provide a default.png resource or handle it accordingly
        mp3Image = canvas.NewImageFromResource(fyne.NewStaticResource("default.png", nil))
        mp3Image.FillMode = canvas.ImageFillContain
        mp3Image.SetMinSize(fyne.NewSize(100, 100))
    }

    mainDisplay := container.NewVBox(
        mp3Image,
        titleLabel,
        timeLabel,
        bpmLabel,
    )
    mainDisplayContainer := container.NewMax(
        canvas.NewCircle(color.Black), // Circular background
        mainDisplay,
    )

    // Pitch Control Slider
    pitchSlider := widget.NewSlider(-10, 10) // Adjust range as needed
    pitchSlider.Orientation = widget.Vertical
    pitchSlider.OnChanged = pitchHandler
    pitchControl := container.NewVBox(
        widget.NewLabelWithStyle("Pitch Control", fyne.TextAlignCenter, fyne.TextStyle{}),
        pitchSlider,
    )

    // Play/Pause Button
    playPauseButton := widget.NewButton("Play", func() {
        if playPauseHandler != nil {
            playPauseHandler()
        } else {
            log.Println("Play/Pause handler not implemented")
        }
    })

    // Pads
    pads := container.NewGridWithColumns(4,
        widget.NewButton("Pad 1", func() { log.Println("Pad 1 pressed") }),
        widget.NewButton("Pad 2", func() { log.Println("Pad 2 pressed") }),
        widget.NewButton("Pad 3", func() { log.Println("Pad 3 pressed") }),
        widget.NewButton("Pad 4", func() { log.Println("Pad 4 pressed") }),
        widget.NewButton("Pad 5", func() { log.Println("Pad 5 pressed") }),
        widget.NewButton("Pad 6", func() { log.Println("Pad 6 pressed") }),
        widget.NewButton("Pad 7", func() { log.Println("Pad 7 pressed") }),
        widget.NewButton("Pad 8", func() { log.Println("Pad 8 pressed") }),
    )

    // Assemble Deck Layout
    return container.NewVBox(
        widget.NewLabelWithStyle(deckName, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
        container.NewHBox(syncButton, widget.NewLabel("")),    // Sync button
        container.NewHBox(mainDisplayContainer, pitchControl), // Main Display and Pitch Control
        playPauseButton,                                       // Play/Pause button
        container.NewVBox(widget.NewLabel("PADS"), pads),      // Pads section
    )
}

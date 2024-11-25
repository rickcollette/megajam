# waveform fyne widget

Originally from <https://git.martyn.berlin/martyn>  
GPL Version 2  
Made modifications to the original code to update, optimize.  Mostly just for my own readability.

## Audio Waveform

A widget to display an audio waveform plot of `[]int32` data.

### Updates & Changes

- Moved `audioDataToImage` to the `Waveform` struct to fix scoping issues and align with Go best practices.
- Enhanced the `audioDataToImage` function to properly handle data mapping and bounds checking.
- Improved default values for colors (`Black` for `OverrideForegroundColor` and `White` for `OverrideBackgroundColor`).
- Simplified the `int32Map` function for readability.
- Updated `Refresh` to ensure proper theme-based or user-defined color handling.
- Improved comments and formatting for better code clarity and maintainability.
- Added bounds checking in `audioDataToImage` to prevent potential runtime errors.

### Widget Fields

| Field                      | Type          | Effect                                           | Default Value |
|----------------------------|---------------|-------------------------------------------------|---------------|
| `audioData`                | `[]int32`     | The data to display in the widget.              | `[]int32{}`   |
| `StretchSamples`           | `bool`        | Resample the samples to fit the widget size.    | `false`       |
| `TransparentBackground`    | `bool`        | Do not draw a background rectangle.             | `false`       |
| `OverrideForeground`       | `bool`        | Set the foreground color manually, or use theme.| `false`       |
| `OverrideForegroundColor`  | `color.Color` | Color to override the theme foreground.         | Black         |
| `OverrideBackground`       | `bool`        | Set the background color manually, or use theme.| `false`       |
| `OverrideBackgroundColor`  | `color.Color` | Color to override the theme background.         | White         |
| `SetMinSize(size)`         | `fyne.Size`   | Sets the widget's minimum size.                 | `200x64`      |

### Usage Example

```go
package main

import (
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "github.com/your-repo/waveform"
)

func main() {
    a := app.New()
    w := a.NewWindow("Waveform Example")

    data := []int32{-500, 1000, -2000, 3000, -4000, 5000}
    wave := waveform.NewWaveform(data)
    wave.SetMinSize(fyne.NewSize(400, 100))
    wave.OverrideForeground = true
    wave.OverrideForegroundColor = color.RGBA{255, 0, 0, 255} // Red color

    w.SetContent(container.NewVBox(
        wave,
    ))

    w.Resize(fyne.NewSize(500, 200))
    w.ShowAndRun()
}
```

### License

This widget is licensed under GPL Version 2.

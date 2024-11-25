package waveform

import (
	"image"
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/disintegration/imaging"
)

type Waveform struct {
	widget.BaseWidget
	audioData               []int32
	StretchSamples          bool
	TransparentBackground   bool
	OverrideForeground      bool
	OverrideForegroundColor color.Color
	OverrideBackground      bool
	OverrideBackgroundColor color.Color
	minSize                 fyne.Size
}

func NewWaveform(data []int32) *Waveform {
	w := &Waveform{
		audioData:               data,
		StretchSamples:          false,
		TransparentBackground:   false,
		OverrideForeground:      false,
		OverrideForegroundColor: color.Black,
		OverrideBackground:      false,
		OverrideBackgroundColor: color.White,
		minSize:                 fyne.NewSize(200, 64),
	}
	w.ExtendBaseWidget(w)
	return w
}

func (w *Waveform) CreateRenderer() fyne.WidgetRenderer {
	return newWaveformRenderer(w)
}

func (w *Waveform) SetMinSize(newSize fyne.Size) {
	w.minSize = newSize
}

func (w *Waveform) audioDataToImage(wd, ht int) image.Image {
	foregroundColor := theme.ForegroundColor()
	if w.OverrideForeground {
		foregroundColor = w.OverrideForegroundColor
	}

	upLeft := image.Point{0, 0}
	lowRight := image.Point{wd, ht}
	if w.StretchSamples {
		lowRight = image.Point{len(w.audioData), ht}
	}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	for i, sample := range w.audioData {
		sampleHeight := int32Map(sample, math.MinInt32, math.MaxInt32, 0, int32(ht/2))
		for y := ht/2 - int(sampleHeight); y <= ht/2+int(sampleHeight); y++ {
			if y >= 0 && y < ht {
				img.Set(i, y, foregroundColor)
			}
		}
	}
	if w.StretchSamples {
		return imaging.Resize(img, wd, ht, imaging.NearestNeighbor)
	}
	return img
}

type waveformRenderer struct {
	widget     *Waveform
	background *canvas.Rectangle
	canvas     *canvas.Raster
}

func newWaveformRenderer(widget *Waveform) *waveformRenderer {
	return &waveformRenderer{
		widget:     widget,
		background: canvas.NewRectangle(theme.BackgroundColor()),
		canvas:     canvas.NewRaster(widget.audioDataToImage),
	}
}

func int32Map(x, inMin, inMax, outMin, outMax int32) int32 {
	return int32((int64(x)-int64(inMin))*(int64(outMax)-int64(outMin))/(int64(inMax)-int64(inMin)) + int64(outMin))
}

func (r *waveformRenderer) Refresh() {
	bgColor := theme.BackgroundColor()
	if r.widget.OverrideBackground {
		bgColor = r.widget.OverrideBackgroundColor
	}
	r.background.FillColor = bgColor
	canvas.Refresh(r.background)
	canvas.Refresh(r.canvas)
}

func (r *waveformRenderer) Layout(size fyne.Size) {
	r.background.Resize(size)
	r.canvas.Resize(size)
}

func (r *waveformRenderer) MinSize() fyne.Size {
	return r.widget.minSize
}

func (r *waveformRenderer) Objects() []fyne.CanvasObject {
	if r.widget.TransparentBackground {
		return []fyne.CanvasObject{r.canvas}
	}
	return []fyne.CanvasObject{r.background, r.canvas}
}

func (r *waveformRenderer) Destroy() {}

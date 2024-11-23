package knobs

import (
	"fmt"
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Knob is a custom widget representing a rotary knob.
type Knob struct {
	widget.BaseWidget
	Value    float64
	Min      float64
	Max      float64
	OnChange func(float64)
}

// NewKnob creates a new Knob with specified min, max, and onChange handler.
func NewKnob(min, max float64, onChange func(float64)) *Knob {
	k := &Knob{Min: min, Max: max, Value: (max + min) / 2, OnChange: onChange}
	k.ExtendBaseWidget(k)
	return k
}

// CreateRenderer creates the renderer for the Knob.
func (k *Knob) CreateRenderer() fyne.WidgetRenderer {
	circle := canvas.NewCircle(color.NRGBA{R: 200, G: 200, B: 200, A: 255})
	indicator := canvas.NewLine(color.NRGBA{R: 0, G: 0, B: 0, A: 255})
	indicator.StrokeWidth = 2
	label := canvas.NewText("", color.NRGBA{R: 0, G: 0, B: 0, A: 255})

	return &knobRenderer{knob: k, circle: circle, indicator: indicator, label: label, objects: []fyne.CanvasObject{circle, indicator, label}}
}

type knobRenderer struct {
	knob      *Knob
	circle    *canvas.Circle
	indicator *canvas.Line
	label     *canvas.Text
	objects   []fyne.CanvasObject
}

func (r *knobRenderer) Layout(size fyne.Size) {
	r.circle.Resize(size)
	r.circle.Move(fyne.NewPos(0, 0))

	centerX, centerY := float64(size.Width/2), float64(size.Height/2)
	radius := float64(size.Width / 2 * 0.8)
	// Limit angle to 270 degrees (3/4 circle)
	minAngle := -135.0 * (math.Pi / 180.0)
	maxAngle := 135.0 * (math.Pi / 180.0)
	normalizedValue := (r.knob.Value - r.knob.Min) / (r.knob.Max - r.knob.Min)
	angle := minAngle + normalizedValue*(maxAngle-minAngle)

	endX := centerX + radius*math.Sin(angle)
	endY := centerY - radius*math.Cos(angle)

	r.indicator.Position1 = fyne.NewPos(float32(centerX), float32(centerY))
	r.indicator.Position2 = fyne.NewPos(float32(endX), float32(endY))

	r.label.Text = fmt.Sprintf("%.1f", r.knob.Value)
	r.label.Alignment = fyne.TextAlignCenter
	r.label.Resize(fyne.NewSize(size.Width, 20))
	r.label.Move(fyne.NewPos(0, size.Height-20))
}

func (r *knobRenderer) MinSize() fyne.Size {
	return fyne.NewSize(50, 70)
}

func (r *knobRenderer) Refresh() {
	canvas.Refresh(r.knob)
}

func (r *knobRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (r *knobRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *knobRenderer) Destroy() {}

// Dragged handles the drag event to adjust the knob's value.
func (k *Knob) Dragged(event *fyne.DragEvent) {
	centerX, centerY := float64(k.Size().Width)/2, float64(k.Size().Height)/2
	dx, dy := float64(event.Position.X)-centerX, float64(event.Position.Y)-centerY
	angle := math.Atan2(dy, dx)

	// Adjust angle to be within -135 to +135 degrees
	minAngle := -135.0 * (math.Pi / 180.0)
	maxAngle := 135.0 * (math.Pi / 180.0)
	if angle < minAngle {
		angle = minAngle
	} else if angle > maxAngle {
		angle = maxAngle
	}

	// Normalize angle to [0,1]
	normalized := (angle - minAngle) / (maxAngle - minAngle)
	value := k.Min + normalized*(k.Max - k.Min)
	value = math.Max(k.Min, math.Min(k.Max, value))

	if value != k.Value {
		k.Value = value
		if k.OnChange != nil {
			k.OnChange(k.Value)
		}
		k.Refresh()
	}
}

// DragEnd handles the end of a drag event.
func (k *Knob) DragEnd() {}

// CreateKnobWithLabel creates a knob with an accompanying label.
func CreateKnobWithLabel(label string, min, max float64, onChange func(float64)) *fyne.Container {
	knob := NewKnob(min, max, onChange)
	labelWidget := widget.NewLabel(label)
	return container.NewVBox(labelWidget, knob)
}

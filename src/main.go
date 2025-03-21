package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	title_gen "krox.editor/src/random_title_ending"
)

type C = layout.Context
type D = layout.Dimensions

type animation struct {
	start    time.Time
	duration time.Duration
	fps      int
}

// animate starts an animation at the current frame which will last for the provided duration.
func (a *animation) animate(gtx layout.Context, duration time.Duration) {
	a.start = gtx.Now
	a.duration = duration
	if a.fps == 0 {
		a.fps = 30
	}
	gtx.Execute(op.InvalidateCmd{At: gtx.Now.Add(duration / time.Duration(a.fps))})
}

// stop ends the animation immediately.
func (a *animation) stop() {
	a.duration = time.Duration(0)
}

// progress returns whether the animation is currently running and (if so) how far through the animation it is.
func (a animation) progress(gtx layout.Context) (animating bool, progress float32) {
	if gtx.Now.After(a.start.Add(a.duration)) {
		return false, 0
	}
	if a.fps == 0 {
		a.fps = 30
	}
	gtx.Execute(op.InvalidateCmd{At: gtx.Now.Add(a.duration / time.Duration(a.fps))})
	return true, float32(gtx.Now.Sub(a.start)) / float32(a.duration)
}

var progress float32
var progressIncrementer chan float32

func draw(w *app.Window, th *material.Theme) error {
	var ops op.Ops

	var startButton widget.Clickable
	var boilDurationInput widget.Editor

	boilDurationInput.SingleLine = true
	boilDurationInput.Alignment = text.Middle

	var boiling bool
	var boilDuration float32

	go func() {
		for p := range progressIncrementer {
			if boiling && progress < 1 {
				progress += p
				w.Invalidate()
			}
		}
	}()

	for {
		switch typ := w.Event().(type) {
		case app.FrameEvent:
			gtx := app.NewContext(&ops, typ)

			if startButton.Clicked(gtx) {
				boiling = !boiling
				if progress >= 1 {
					progress = 0
				}

				inputString := boilDurationInput.Text()
				inputString = strings.TrimSpace(inputString)
				inputFloat, err := strconv.ParseFloat(inputString, 32)
				if err != nil {
					boilDurationInput.SetText("Please enter a number")
					if boiling {
						boiling = false
						progress = 0
						// w.Invalidate()
					}
				}
				boilDuration = float32(inputFloat)
				boilDuration = boilDuration / (1 - progress)
			}

			layout.Flex{
				Axis:    layout.Vertical,
				Spacing: layout.SpaceStart,
			}.Layout(gtx,
				// egg
				layout.Rigid(
					func(gtx C) D {
						// Draw a custom path, shaped like an egg
						var eggPath clip.Path
						op.Offset(image.Pt(gtx.Dp(200), gtx.Dp(125))).Add(gtx.Ops)
						eggPath.Begin(gtx.Ops)
						// Rotate from 0 to 360 degrees
						for deg := 0.0; deg <= 360; deg++ {

							// Egg math (really) at this brilliant site. Thanks!
							// https://observablehq.com/@toja/egg-curve
							// Convert degrees to radians
							rad := deg / 360 * 2 * math.Pi
							// Trig gives the distance in X and Y direction
							cosT := math.Cos(rad)
							sinT := math.Sin(rad)
							// Constants to define the eggshape
							a := 110.0
							b := 150.0
							d := 20.0
							// The x/y coordinates
							x := a * cosT
							y := -(math.Sqrt(b*b-d*d*cosT*cosT) + d*sinT) * sinT
							// Finally the point on the outline
							p := f32.Pt(float32(x), float32(y))
							// Draw the line to this point
							eggPath.LineTo(p)
						}
						// Close the path
						eggPath.Close()

						// Get hold of the actual clip
						eggArea := clip.Outline{Path: eggPath.End()}.Op()

						// Fill the shape
						// color := color.NRGBA{R: 255, G: 239, B: 174, A: 255}
						color := color.NRGBA{R: 255, G: uint8(239 * (1 - progress)), B: uint8(174 * (1.3 - progress)), A: 255}
						paint.FillShape(gtx.Ops, color, eggArea)

						d := image.Point{Y: 375}
						return D{Size: d}
					},
				),
				// boil duration input
				layout.Rigid(
					func(gtx C) D {
						ed := material.Editor(th, &boilDurationInput, "Boil duration (seconds)")

						if boiling && progress < 1 {
							boilRemain := (1 - progress) * boilDuration
							// Format to 1 decimal.
							inputStr := fmt.Sprintf("%.1f", math.Round(float64(boilRemain)*10)/10)
							// Update the text in the inputbox
							boilDurationInput.SetText(inputStr)
						} else if progress >= 1 && boiling {
							boilDurationInput.SetText("Done!")
						}

						margins := layout.Inset{
							Top:    unit.Dp(0),
							Bottom: unit.Dp(30),
							Right:  unit.Dp(10),
							Left:   unit.Dp(10),
						}

						border := widget.Border{
							Color:        color.NRGBA{R: 204, G: 204, B: 204, A: 255},
							CornerRadius: unit.Dp(5),
							Width:        unit.Dp(2),
						}

						return margins.Layout(gtx,
							func(gtx C) D {
								return border.Layout(gtx, ed.Layout)
							},
						)
					},
				),
				// progress bar
				layout.Rigid(
					func(gtx C) D {
						bar := material.ProgressBar(th, progress)
						return bar.Layout(gtx)
					},
				),
				// start/stop button
				layout.Rigid(
					func(gtx C) D {
						margin := layout.UniformInset(unit.Dp(25))

						var text string
						if !boiling {
							text = "Start"
						} else {
							text = "Stop"
						}

						if progress >= 1 && boiling {
							text = "Finished"
						}

						return margin.Layout(gtx,
							func(gtx C) D {
								btn := material.Button(th, &startButton, text)
								if progress >= 1 && boiling {
									btn.Background = color.NRGBA{R: 61, G: 47, B: 76, A: 255}
								}
								return btn.Layout(gtx)
							},
						)
					},
				),
			)

			typ.Frame(gtx.Ops)
		case app.DestroyEvent:
			return typ.Err
		}
	}
}

func main() {
	go func() {
		w := new(app.Window)
		w.Option(app.Title(title_gen.GetRandomTitle()))
		w.Option(app.Size(unit.Dp(400), unit.Dp(600)))
		w.Option(app.MinSize(unit.Dp(300), unit.Dp(300)))

		th := material.NewTheme()
		if err := draw(w, th); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	progressIncrementer = make(chan float32)
	go func() {
		for {
			time.Sleep(time.Second / 25)
			progressIncrementer <- 0.004
		}
	}()

	app.Main()
}

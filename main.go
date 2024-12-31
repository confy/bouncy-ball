package main

import (
	"image"
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

var (
    WINDOW_WIDTH = unit.Dp(800)
    WINDOW_HEIGHT = unit.Dp(600)
    BALL_COLOR = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
    BACKGROUND_COLOR = color.NRGBA{R: 40, G: 40, B: 40, A: 255}
    CIRCLE_SIZE = 0.2 // Circle size as a fraction of window width
)

func main() {
    go func() {
        w := new(app.Window)
        w.Option(app.Title("Bouncy Ball"))
        w.Option(app.Size(WINDOW_WIDTH, WINDOW_HEIGHT))
        if err := draw(w); err != nil {
            log.Fatal(err)
        }
        os.Exit(0)
    }()

    app.Main()
}

type C = layout.Context
type D = layout.Dimensions

func draw(w *app.Window) error {
    // ops are the operations from the UI
    var ops op.Ops

    for {
        // listen for events
        switch e := w.Event().(type) {

        // this is sent when the application should re-render.
        case app.FrameEvent:
            gtx := app.NewContext(&ops, e)
            
            // Set the background color
            paint.Fill(gtx.Ops, BACKGROUND_COLOR)

            // Calculate the circle radius
            circleRadius := int(float64(gtx.Constraints.Max.X) * CIRCLE_SIZE / 2)

            // Draw the circle
            circle := clip.Ellipse{
                Min: image.Pt(gtx.Constraints.Max.X/2-circleRadius, gtx.Constraints.Max.Y/2-circleRadius),
                Max: image.Pt(gtx.Constraints.Max.X/2+circleRadius, gtx.Constraints.Max.Y/2+circleRadius),
            }.Op(gtx.Ops)
            paint.FillShape(gtx.Ops, BALL_COLOR, circle)

            // Draw the frame
            e.Frame(gtx.Ops)

        // on close
        case app.DestroyEvent:
            return e.Err
        }

    }
}
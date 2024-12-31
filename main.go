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
    WINDOW_WIDTH = unit.Dp(1000)
    WINDOW_HEIGHT = unit.Dp(1000)
    BALL_COLOR = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
    BACKGROUND_COLOR = color.NRGBA{R: 40, G: 40, B: 40, A: 255}
    CIRCLE_RADIUS = 80
    GRAVITY = 0.25
    DAMPING = 0.8
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


type Ball struct {
    X, Y, Radius int
    VelocityX, VelocityY float64
    AccelerationX, AccelerationY float64
}

func (b *Ball) Draw(gtx *layout.Context) {
    circle := clip.Ellipse{
        Min: image.Pt(b.X-b.Radius, b.Y-b.Radius),
        Max: image.Pt(b.X+b.Radius, b.Y+b.Radius),
    }.Op(gtx.Ops)
    paint.FillShape(gtx.Ops, BALL_COLOR, circle)
}

func (b *Ball) Update() {
    b.VelocityY += GRAVITY
    b.X += int(b.VelocityX)
    b.Y += int(b.VelocityY)

    // Check for collision with the window bounds and invert velocity and acceleration if necessary
    if b.X-b.Radius < 0 {
        b.X = b.Radius
        b.VelocityX = -b.VelocityX * DAMPING
        b.AccelerationX = -b.AccelerationX * DAMPING
    } else if b.X+b.Radius > int(WINDOW_WIDTH) {
        b.X = int(WINDOW_WIDTH) - b.Radius
        b.VelocityX = -b.VelocityX * DAMPING
        b.AccelerationX = -b.AccelerationX * DAMPING
    }

    if b.Y-b.Radius < 0 {
        b.Y = b.Radius
        b.VelocityY = -b.VelocityY * DAMPING
        b.AccelerationY = -b.AccelerationY * DAMPING
    } else if b.Y+b.Radius > int(WINDOW_HEIGHT) {
        b.Y = int(WINDOW_HEIGHT) - b.Radius
        b.VelocityY = -b.VelocityY * DAMPING
        b.AccelerationY = -b.AccelerationY * DAMPING
    }
}

func draw(w *app.Window) error {
    // ops are the operations from the UI
    var ops op.Ops

    ball := Ball{
        X: int(WINDOW_WIDTH) / 2,
        Y: int(WINDOW_HEIGHT) / 2,
        Radius: CIRCLE_RADIUS,
        VelocityX: 20,
        VelocityY: 20,
        AccelerationX: 2,
        AccelerationY: 1,
    }
    
    for {
        // listen for events
        switch e := w.Event().(type) {

        // this is sent when the application should re-render.
        case app.FrameEvent:
            gtx := app.NewContext(&ops, e)

            // Set the background color
            paint.Fill(gtx.Ops, BACKGROUND_COLOR)
            
            // update ball
            ball.Update()

            // draw ball
            ball.Draw(&gtx)

            // Draw the frame
            e.Frame(gtx.Ops)

            // Invalidate the window to request another frame, results in a FrameEvent
            w.Invalidate()
        // on close
        case app.DestroyEvent:
            return e.Err
        }

    }
}
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
    BACKGROUND_COLOR = color.NRGBA{R: 16, G: 16, B: 16, A: 255}
    
    BALL_RADIUS = 80
    BALL_COLOR = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
    
    TRAIL_START_RADIUS = 25
    TRAIL_COLOR = color.NRGBA{R: 0, G: 255, B: 255, A: 100}
    TRAIL_MAX_LENGTH = 100

    GRAVITY = 0.2
    DAMPING = 0.95
    ACCELERATION_DECAY = 0.99
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


type Circle struct {
    X, Y, Radius int
    Color color.NRGBA
}

func (c *Circle) Draw(gtx *layout.Context) {
    circle := clip.Ellipse{
        Min: image.Pt(c.X-c.Radius, c.Y-c.Radius),
        Max: image.Pt(c.X+c.Radius, c.Y+c.Radius),
    }.Op(gtx.Ops)
    paint.FillShape(gtx.Ops, c.Color, circle)
}


type Ball struct {
    Circle Circle
    VelocityX, VelocityY float64
    AccelerationX, AccelerationY float64
}

func (b *Ball) Update() {
    b.VelocityY += GRAVITY
    b.Circle.X += int(b.VelocityX)
    b.Circle.Y += int(b.VelocityY)

    b.VelocityX += b.AccelerationX
    b.VelocityY += b.AccelerationY

    b.AccelerationX *= ACCELERATION_DECAY
    b.AccelerationY *= ACCELERATION_DECAY

    // Check for collision with the window bounds and handle it
    b.handleCollision()
}

func (b *Ball) handleCollision() {
    // Check for collision with the left and right bounds
    if b.Circle.X-b.Circle.Radius < 0 {
        b.Circle.X = b.Circle.Radius
        b.VelocityX = -b.VelocityX * DAMPING
        b.AccelerationX = -b.AccelerationX * DAMPING
    } else if b.Circle.X+b.Circle.Radius > int(WINDOW_WIDTH) {
        b.Circle.X = int(WINDOW_WIDTH) - b.Circle.Radius
        b.VelocityX = -b.VelocityX * DAMPING
        b.AccelerationX = -b.AccelerationX * DAMPING
    }

    // Check for collision with the top and bottom bounds
    if b.Circle.Y-b.Circle.Radius < 0 {
        b.Circle.Y = b.Circle.Radius
        b.VelocityY = -b.VelocityY * DAMPING
        b.AccelerationY = -b.AccelerationY * DAMPING
    } else if b.Circle.Y+b.Circle.Radius > int(WINDOW_HEIGHT) {
        b.Circle.Y = int(WINDOW_HEIGHT) - b.Circle.Radius
        b.VelocityY = -b.VelocityY * DAMPING
        b.AccelerationY = -b.AccelerationY * DAMPING
    }
}


func drawTrail(gtx *layout.Context, trail []Circle) {
    // Draw the trail
    for i, c := range trail {
        realIndex := len(trail) - i - 1
        // Decrease the radius and alpha of the circle as it gets older
        var decay float64 = 1 - (float64(realIndex) / float64(len(trail)))
        c.Radius = int(float64(TRAIL_START_RADIUS) * decay)
        c.Color.A = uint8(255 * decay)
        c.Draw(gtx)
    }
}
func draw(w *app.Window) error {
    // ops are the operations from the UI
    var ops op.Ops
    
    ball := Ball{
        Circle: Circle{
            X: int(WINDOW_WIDTH) / 2,
            Y: int(WINDOW_HEIGHT) / 2,
            Radius: BALL_RADIUS,
            Color: BALL_COLOR,
        },
        VelocityX: 10,
        VelocityY: 5,
        AccelerationX: 0.1,
        AccelerationY: 0.1,
    }

    var trail []Circle
    for {
        switch e := w.Event().(type) {
        // on close
        case app.DestroyEvent:
            return e.Err

        // this is sent when the application should re-render.
        case app.FrameEvent:
            gtx := app.NewContext(&ops, e)

            // Set the background color
            paint.Fill(gtx.Ops, BACKGROUND_COLOR)
            
            // update ball
            ball.Update()

            // add ball to trail
            trail = append(trail, Circle{
                X: ball.Circle.X,
                Y: ball.Circle.Y,
                Radius: TRAIL_START_RADIUS,
                Color: TRAIL_COLOR,
            })
            if len(trail) > TRAIL_MAX_LENGTH {
                trail = trail[1:]
            }

            // draw trail
            drawTrail(&gtx, trail)

            // draw ball
            ball.Circle.Draw(&gtx)

            // Draw the frame
            e.Frame(gtx.Ops)

            // Invalidate the window to request another frame, results in a FrameEvent
            w.Invalidate()

        }
    }
}
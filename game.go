package main

import (
	"math"
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type Game struct {
	Ball   *Ball
	Paddle *Paddle

	Win *pixelgl.Window

	Player *Player

	Score int
}

func NewGame() (*Game, error) {
	var err error
	g := Game{}

	g.Ball, err = NewBall("assets/ball.png")
	if err != nil {
		return nil, err
	}

	g.Paddle, err = NewPaddle("assets/paddle.png")
	if err != nil {
		return nil, err
	}

	return &g, nil
}

// The main loop, must be called by pixel
func (g *Game) Run() {
	var err error
	ball := g.Ball
	paddle := g.Paddle

	// Create a new window
	g.Win, err = pixelgl.NewWindow(pixelgl.WindowConfig{
		Title:  "How to pong",
		Bounds: pixel.R(0, 0, 1024, 1024),
		VSync:  true,
	})
	if err != nil {
		panic(err)
	}

	// Set the ball and paddle positions relative to the window
	g.Ball.Position = g.Win.Bounds().Center()
	g.Paddle.Position = g.Win.Bounds().Center().Add(pixel.V(0, -400))

	// The main run loop
	win := g.Win
	for !win.Closed() {
		// Just show green while theres no player playing
		if g.Player == nil {
			win.Clear(colornames.Darkgreen)
			win.Update()
			continue
		}

		win.Clear(colornames.Skyblue)

		// Update
		g.Player.MovePaddle(g)
		g.Score += ball.update(win.Bounds(), paddle, 0.2)
		// Reset the paddle velocity now that we're done with update computations
		paddle.Velocity.X = 0

		// Draw
		ball.draw(win)
		paddle.draw(win)

		// Commit
		win.Update()
	}
}

func (g *Game) Reset() {
	g.Ball.Position = g.Win.Bounds().Center()
	g.Paddle.Position = g.Win.Bounds().Center().Add(pixel.V(0, -400))

	// Compute a new velocity that isnt too slow
	vx := rand.NormFloat64()*25 + 60
	vy := rand.NormFloat64()*25 + 60
	for math.Abs(vy) < 30 {
		vy = rand.NormFloat64()*25 + 60
	}
	if rand.Int()%2 == 0 {
		vx *= -1
	}
	if rand.Int()%2 == 0 {
		vy *= -1
	}
	g.Ball.Velocity = pixel.Vec{X: vx, Y: vy}

	g.Score = 0

	g.Player = nil
}

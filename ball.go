package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

// Ball represents the ball that players interact with.
// Ball is the one that is responsible for assigning points.
// (Game accumulates the points into the score).
type Ball struct {
	Position pixel.Vec
	Velocity pixel.Vec
	Sprite   *pixel.Sprite
}

// NewBall returns a ball with the sprite at spritePath
func NewBall(spritePath string) (*Ball, error) {
	pic, err := loadPicture(spritePath)
	if err != nil {
		return nil, err
	}

	sprite := pixel.NewSprite(pic, pic.Bounds())

	b := Ball{
		Velocity: pixel.V(25, 80),
		Sprite:   sprite,
	}

	return &b, nil
}

// update updates the balls position and velocity based on the location of its bounds and the paddle.
// update returns the number of points assigned during this update.
func (b *Ball) update(bounds pixel.Rect, paddle *Paddle, dt float64) int {
	hit := 0

	// If we hit the paddle, reverse direction vertically
	// and apply friction horizontally.
	// Hittinh the ball with the paddle is worth 1 point
	if b.hitPaddle(paddle) {
		b.Velocity.Y *= -1
		b.Velocity.X += paddle.Velocity.X
		hit = 1
	}

	w := b.Sprite.Frame().W()
	h := b.Sprite.Frame().H()

	// Project position into the future to check for collisions before they happen
	x := b.Position.X + b.Velocity.X*dt
	y := b.Position.Y + b.Velocity.Y*dt

	// If we hit the vertical walls of our bounds, bounce back horizontally
	if x < bounds.Min.X+w/2 || x > bounds.Max.X-w/2 {
		b.Velocity.X *= -1
		x = b.Position.X + b.Velocity.X*dt
	}

	// If we hit the horizontal walls of our bounds, bounce back vertically
	if y < bounds.Min.Y+h/2 || y > bounds.Max.Y-h/2 {
		// Hitting the top is worth 2 points
		// Hitting the bottom is worth -3 points
		if y > bounds.Max.Y-h/2 {
			hit = 2
		} else {
			hit = -3
		}
		b.Velocity.Y *= -1
		y = b.Position.Y + b.Velocity.Y*dt
	}

	// Update our position
	b.Position.X = x
	b.Position.Y = y

	return hit
}

// draw draws the ball at the correct position in w
func (b *Ball) draw(w *pixelgl.Window) {
	mat := pixel.IM.Moved(b.Position)
	b.Sprite.Draw(w, mat)
}

// hitPaddle returns true if the ball and paddle are touching
func (b *Ball) hitPaddle(p *Paddle) bool {
	bbox := b.Sprite.Frame().Moved(b.Position)
	pbox := p.Sprite.Frame().Moved(p.Position)

	return bbox.Intersects(pbox)
}

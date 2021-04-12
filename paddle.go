package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

// Paddle represents the paddle that the player can move.
type Paddle struct {
	Position pixel.Vec
	Velocity pixel.Vec // reset at the end of every update
	Sprite   *pixel.Sprite
}

func NewPaddle(spritePath string) (*Paddle, error) {
	pic, err := loadPicture(spritePath)
	if err != nil {
		return nil, err
	}

	sprite := pixel.NewSprite(pic, pic.Bounds())

	p := Paddle{
		Sprite: sprite,
	}

	return &p, nil
}

func (p *Paddle) draw(w *pixelgl.Window) {
	mat := pixel.IM.Moved(p.Position)
	p.Sprite.Draw(w, mat)
}

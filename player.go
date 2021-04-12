package main

import (
	"math"
	"math/rand"
	"sort"
)

const maxPaddleVelocity float64 = 20
const minPaddleVelocity float64 = 3

func sigmoid(x float64) float64 {
	e := math.Exp(x)
	return (e / (1.0 + e)) - 0.5
}

type Player struct {
	Genes [8]float64
	Score int
}

func NewPlayer() Player {
	player := Player{
		Score: -math.MaxInt32,
	}

	for i := 0; i < 8; i++ {
		player.Genes[i] = rand.NormFloat64() * 15
	}

	return player
}

func (p Player) MovePaddle(g *Game) {
	x := float64(0)
	ball := g.Ball
	paddle := g.Paddle
	width := g.Win.Bounds().W()
	height := g.Win.Bounds().H()

	// Normalize input
	bpx := ball.Position.X / width
	bpy := ball.Position.Y / height

	bvx := ball.Velocity.X / width
	bvy := ball.Velocity.Y / height

	ppx := paddle.Position.X / width
	ppy := paddle.Position.Y / height

	// Linear combination over the genes
	x += bpx * p.Genes[0]
	x += bpy * p.Genes[1]

	x += bvx * p.Genes[2]
	x += bvy * p.Genes[3]

	x += ppx * p.Genes[4]
	x += ppy * p.Genes[5]

	dx := bpx - ppx
	dy := bpy - ppy
	x += math.Sqrt(dx*dx+dy*dy) * p.Genes[6]

	x += (1 - ppx) * p.Genes[7]

	// Normalize and scale
	x = maxPaddleVelocity * sigmoid(x)

	// Move paddle
	g.Paddle.Position.X += x
	g.Paddle.Velocity.X = x
}

func Breed(a, b Player) Player {
	c := Player{}

	for i := 0; i < 8; i++ {
		r := rand.Int() % 10
		switch {
		case r < 3:
			c.Genes[i] = a.Genes[i]
		case r >= 3 && r < 6:
			c.Genes[i] = b.Genes[i]
		case r >= 6 && r < 8:
			c.Genes[i] = (a.Genes[i] + b.Genes[i]) / 2.0
		default:
			c.Genes[i] = rand.NormFloat64() * 50
		}
	}

	return c
}

type Generation []Player

func (g Generation) Len() int {
	return len(g)
}

func (g Generation) Less(i, j int) bool {
	// Reverse so we sort high -> low
	return g[i].Score > g[j].Score
}

func (g Generation) Swap(i, j int) {
	g[i], g[j] = g[j], g[i]
}

func (g Generation) Next(cull float64) Generation {
	if cull >= 1.0 || cull <= 0 {
		cull = 0.5
	}

	nextGen := make(Generation, len(g))
	sort.Sort(g)

	// Keep the top cull% of players
	n := int32(math.Floor(float64(len(g))*cull)) + 1

	for idx := range nextGen {
		aIDX := rand.Int31n(n)
		bIDX := rand.Int31n(n)

		// Keep sampling until we have two different parents
		for aIDX == bIDX {
			bIDX = rand.Int31n(n)
		}

		nextGen[idx] = Breed(g[aIDX], g[bIDX])
	}

	return nextGen
}

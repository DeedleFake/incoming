package main

import (
	"github.com/DeedleFake/incoming/engine"
	"image"
)

type Asteroid struct {
	anim *engine.Anim
	rect image.Rectangle
}

func NewAsteroid(anim *engine.Anim, start image.Point) *Asteroid {
	return &Asteroid{
		anim: anim,
		rect: image.Rectangle{
			Min: start,
			Max: start.Add(anim.Size()),
		},
	}
}

func (a Asteroid) Anim() *engine.Anim {
	return a.anim
}

func (a Asteroid) Bounds() image.Rectangle {
	return a.rect
}

func (a *Asteroid) Move(x, y int) {
	a.rect.Min.X += x
	a.rect.Min.Y += y

	a.rect.Max.X += x
	a.rect.Max.Y += y
}

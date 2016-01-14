package main

import (
	"github.com/DeedleFake/incoming/engine"
	"image"
)

type Player struct {
	anim *engine.Anim
	loc  image.Point
}

func NewPlayer(anim *engine.Anim, loc image.Point) *Player {
	return &Player{
		anim: anim,
		loc:  loc,
	}
}

func (p Player) Anim() *engine.Anim {
	return p.anim
}

func (p Player) Loc() image.Point {
	return p.loc
}

func (p *Player) Move(x, y int) {
	p.loc.X += x
	p.loc.Y += y
}

func (p Player) Bounds() image.Rectangle {
	// TODO: This is inefficient.

	size := p.anim.Size()

	return image.Rect(
		p.loc.X+(size.X/4),
		p.loc.Y+(size.Y/8),
		size.X/2,
		size.Y,
	)
}

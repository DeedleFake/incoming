package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Anim struct {
	image  *sdl.Texture
	frames []sdl.Rect
	cur    int
}

func NewAnim(data []byte, frameW int) *Anim {
	panic("Not implemented.")
}

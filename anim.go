package main

import (
	"fmt"
	"golang.org/x/exp/shiny/screen"
	"image"
)

type Anim struct {
	image screen.Texture
	len   int
	cur   image.Rectangle
}

func newAnim(tex screen.Texture, frameW int) (*Anim, error) {
	bnds := tex.Bounds()
	if bnds.Dx()%frameW != 0 {
		return nil, &InvalidFrameWidthError{bnds.Dx(), frameW}
	}

	return &Anim{
		image: tex,
		len:   bnds.Dx() / frameW,
		cur:   image.Rect(bnds.Min.X, bnds.Min.Y, frameW, bnds.Dy()),
	}, nil
}

func (anim *Anim) Advance() {
	anim.cur = anim.cur.Add(image.Pt(anim.cur.Dx(), 0))
	if anim.cur.Min.X >= anim.image.Size().X {
		anim.cur.Sub(image.Pt(anim.image.Size().X, 0))
	}
}

func (anim Anim) Frames() int {
	return anim.len
}

type InvalidFrameWidthError struct {
	ImageW int
	FrameW int
}

func (err InvalidFrameWidthError) Error() string {
	return fmt.Sprintf("Image width (%v) is not divisible by frame width (%v)", err.ImageW, err.FrameW)
}

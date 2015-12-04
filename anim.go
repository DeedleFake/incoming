package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
	"unsafe"
)

type Anim struct {
	ren *sdl.Renderer

	image  *sdl.Texture
	frames []sdl.Rect
	cur    int
}

func loadTexture(ren *sdl.Renderer, data []byte) (*sdl.Texture, error) {
	rw := sdl.RWFromMem(unsafe.Pointer(&data[0]), len(data))
	defer rw.FreeRW()

	s, err := img.Load_RW(rw, false)
	if err != nil {
		return nil, err
	}

	return ren.CreateTextureFromSurface(s)
}

func NewAnim(ren *sdl.Renderer, data []byte, frameW int) (*Anim, error) {
	img, err := loadTexture(ren, data)
	if err != nil {
		return nil, err
	}

	_, _, w, h, err := img.Query()
	if err != nil {
		return nil, err
	}
	if int(w)%frameW != 0 {
		return nil, &InvalidFrameWidthError{int(w), frameW}
	}

	frames := make([]sdl.Rect, 0, int(w)/frameW)
	for i := int32(0); i < w; i += int32(frameW) {
		frames = append(frames, sdl.Rect{
			X: i,
			Y: 0,
			W: int32(frameW),
			H: h,
		})
	}

	return &Anim{
		ren: ren,

		image:  img,
		frames: frames,
	}, nil
}

func (anim Anim) Draw(dst *sdl.Rect) error {
	return anim.ren.Copy(anim.image, &anim.frames[anim.cur], dst)
}

func (anim *Anim) Advance() {
	anim.cur++
	if anim.cur >= len(anim.frames) {
		anim.cur = 0
	}
}

type InvalidFrameWidthError struct {
	ImageW int
	FrameW int
}

func (err InvalidFrameWidthError) Error() string {
	return fmt.Sprintf("Image width (%v) is not divisible by frame width (%v)", err.ImageW, err.FrameW)
}

package main

import (
	"fmt"
	"golang.org/x/exp/shiny/screen"
	"image"
	"sync"
	"time"
)

type Anim struct {
	image screen.Texture
	cur   image.Rectangle

	// TODO: Find a cleaner way to do this.
	m     sync.Mutex
	delay chan time.Duration
	done  chan struct{}
}

func newAnim(tex screen.Texture, frameW int) (*Anim, error) {
	bnds := tex.Bounds()
	if bnds.Dx()%frameW != 0 {
		return nil, &InvalidFrameWidthError{bnds.Dx(), frameW}
	}

	return &Anim{
		image: tex,
		cur:   image.Rect(bnds.Min.X, bnds.Min.Y, frameW, bnds.Dy()),
	}, nil
}

func (anim *Anim) Copy() *Anim {
	anim.m.Lock()
	defer anim.m.Unlock()

	bnds := anim.image.Bounds()

	return &Anim{
		image: anim.image,
		cur: image.Rectangle{
			Min: bnds.Min,
			Max: bnds.Min.Add(anim.cur.Size()),
		},
	}
}

func (anim *Anim) advance() {
	anim.m.Lock()
	defer anim.m.Unlock()

	s := anim.image.Size()
	w := anim.cur.Dx()

	anim.cur.Min.X += w
	anim.cur.Max.X += w
	if anim.cur.Min.X >= s.X {
		anim.cur.Min.X, anim.cur.Max.X = 0, w
	}
}

func (anim *Anim) animate(done <-chan struct{}, delay time.Duration) {
	t := time.NewTicker(delay)
	last := delay

	for {
		select {
		case <-t.C:
			anim.advance()

		case delay := <-anim.delay:
			if delay == last {
				continue
			}

			t.Stop()
			t = time.NewTicker(delay)
			last = delay

		case <-done:
			t.Stop()
			return
		}
	}
}

func (anim *Anim) Start(delay time.Duration) {
	if anim.done != nil {
		anim.delay <- delay
		return
	}

	anim.delay = make(chan time.Duration)
	anim.done = make(chan struct{})

	go anim.animate(anim.done, delay)
}

func (anim *Anim) Stop() {
	if anim.done == nil {
		return
	}

	select {
	case <-anim.done:
	default:
		close(anim.done)
	}

	anim.done = nil
}

func (anim *Anim) Frames() int {
	anim.m.Lock()
	defer anim.m.Unlock()

	return anim.image.Size().X / anim.cur.Dx()
}

func (anim *Anim) Size() image.Point {
	anim.m.Lock()
	defer anim.m.Unlock()

	return anim.cur.Size()
}

func (anim *Anim) Image() (screen.Texture, image.Rectangle) {
	anim.m.Lock()
	defer anim.m.Unlock()

	return anim.image, anim.cur
}

type InvalidFrameWidthError struct {
	ImageW int
	FrameW int
}

func (err InvalidFrameWidthError) Error() string {
	return fmt.Sprintf("Image width (%v) is not divisible by frame width (%v)", err.ImageW, err.FrameW)
}

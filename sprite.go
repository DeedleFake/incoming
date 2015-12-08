package main

import (
	"image"
	"time"
)

type Sprite struct {
	Anim   Anim
	Bounds image.Rectangle

	done chan struct{}
}

func NewSprite(anim Anim, delay time.Duration) (s *Sprite) {
	s = &Sprite{
		Anim: anim,
		Bounds: image.Rectangle{
			Min: image.ZP,
			Max: anim.Size(),
		},

		done: make(chan struct{}),
	}
	runtime.SetFinalizer(s, (*Sprite).release)

	go s.animate(delay)

	return
}

// TODO: This will never get run. The loop in animate maintains a
// reference to s.
func (s *Sprite) release() {
	close(s.done)
}

func (s *Sprite) animate(delay time.Duration) {
	t := time.NewTicker(delay)

	for {
		select {
		case <-t.C:
			// TODO: Fix data race.
			s.Anim.Advance()
		case <-s.done:
			t.Stop()
			return
		}
	}
}

func (s *Sprite) Move(x, y int) {
	s.Bounds.Min.X += x
	s.Bounds.Min.Y += y

	s.Bounds.Max.X += x
	s.Bounds.Max.Y += y
}
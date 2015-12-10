package main

import (
	"image"
	"sync"
	"time"
)

// TODO: Delete Sprite and work the animation handling into Anim.
type Sprite struct {
	Loc image.Point

	anim  *Anim
	animM sync.RWMutex
	delay chan time.Duration

	done chan struct{}
}

func NewSprite(anim *Anim, delay time.Duration) (s *Sprite) {
	s = &Sprite{
		anim:  anim,
		delay: make(chan time.Duration),
	}

	s.StartAnim(delay)

	return
}

func (s *Sprite) StartAnim(delay time.Duration) {
	s.StopAnim()

	s.done = make(chan struct{})
	go s.animate(delay)
}

func (s *Sprite) AdjustDelay(delay time.Duration) {
	s.delay <- delay
}

func (s *Sprite) StopAnim() {
	if s.done == nil {
		return
	}

	select {
	case <-s.done:
	default:
		close(s.done)
	}
}

func (s *Sprite) animate(delay time.Duration) {
	t := time.NewTicker(delay)

	// Prevent potential data race.
	done := s.done

	for {
		select {
		case <-t.C:
			s.animM.Lock()
			s.anim.Advance()
			s.animM.Unlock()

		case delay := <-s.delay:
			t.Stop()
			t = time.NewTicker(delay)

		case <-done:
			t.Stop()
			return
		}
	}
}

func (s *Sprite) Move(x, y int) {
	s.Loc.X += x
	s.Loc.Y += y
}

func (s *Sprite) Draw(state *State) {
	s.animM.RLock()
	defer s.animM.RUnlock()

	state.Draw(s.anim, s.Loc)
}

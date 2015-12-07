package main

import (
	"bytes"
	"image"
	"log"
	"time"
)

//go:generate go build ./cmd/bintogo

//go:generate ./bintogo ./images/title.png

type Title struct {
	s *State

	bg      *Anim
	bgDelay <-chan time.Time
}

func NewTitle(s *State) *Title {
	bg, err := s.LoadAnim(bytes.NewReader(titleData[:]), 160)
	if err != nil {
		log.Fatalf("Failed to load title BG: %v", err)
	}

	return &Title{
		s: s,

		bg:      bg,
		bgDelay: time.Tick(time.Second / 20),
	}
}

func (t *Title) Enter() {
}

func (t *Title) Update() {
	select {
	case <-t.bgDelay:
		t.bg.Advance()
	default:
	}

	t.s.Draw(t.bg, image.ZP)
}

//go:generate ./bintogo ./images/player.png
//go:generate ./bintogo ./images/a1.png

type Game struct {
}

//go:generate ./bintogo ./images/lose.png

type Lose struct {
}

//go:generate ./bintogo ./images/win.png

type Win struct {
}

//go:generate rm -f ./bintogo

func main() {
	s := NewState()
	opts := StateOptions{
		Width:  240,
		Height: 160,
		FPS:    60,
	}

	err := s.Run(&opts, func() bool {
		s.AddRoom("title", NewTitle(s))
		s.EnterRoom("title")

		return true
	})
	if err != nil {
		log.Fatalf("Failed to run game: %v", err)
	}
}

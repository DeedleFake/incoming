package main

import (
	"log"
	"time"
)

//go:generate go build ./cmd/bintogo

//go:generate ./bintogo ./images/title.png

type Title struct {
	s *State

	bg *Anim
}

func NewTitle(s *State) *Title {
	bg, err := s.NewAnim(titleData[:], 160)
	if err != nil {
		log.Fatalf("Failed to load title BG: %v", err)
	}

	return &Title{
		s: s,

		bg: bg,
	}
}

func (t *Title) Enter() {
}

func (t *Title) Update() {
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
	s, err := NewState("Incoming!", 240, 160)
	if err != nil {
		log.Fatalf("Failed to initialize game state: %v", err)
	}
	s.AddRoom("title", NewTitle(s))
	s.EnterRoom("title")

	fps := time.Tick(time.Second / 60)
	for s.Update() {
		<-fps
	}
}

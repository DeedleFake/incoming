package main

import (
	"log"
	"time"
)

type Title struct {
	s *State

	//bg Anim
}

func (t *Title) Enter() {
}

func (t *Title) Update() {
}

//go:generate go build ./cmd/bintogo
//
//go:generate ./bintogo ./images/title.png
//go:generate ./bintogo ./images/player.png
//go:generate ./bintogo ./images/a1.png
//go:generate ./bintogo ./images/lose.png
//go:generate ./bintogo ./images/win.png
//
//go:generate rm -f ./bintogo
func main() {
	s, err := NewState("Incoming!", 240, 160)
	if err != nil {
		log.Fatalf("Failed to initialize game state: %v", err)
	}

	s.AddRoom("title", &Title{
		s: s,

		//bg: NewAnim(titleData),
	})

	s.EnterRoom("title")

	fps := time.Tick(time.Second / 60)
	for s.Update() {
		<-fps
	}
}

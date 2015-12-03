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

func main() {
	s, err := NewState("Incoming!", 240, 160)
	if err != nil {
		log.Fatalf("Failed to initialize game state: %v", err)
	}

	s.AddRoom("title", &Title{
		s: s,

		//bg: NewAnim(titleAnim),
	})

	s.EnterRoom("title")

	fps := time.Tick(time.Second / 60)
	for s.Update() {
		<-fps
	}
}

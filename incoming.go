package main

import (
	"bytes"
	"golang.org/x/mobile/event/key"
	"image"
	"image/color"
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
	bg, err := s.LoadAnim(bytes.NewReader(titleData[:]), s.Bounds().Dx())
	if err != nil {
		log.Fatalf("Failed to load title BG: %v", err)
	}

	return &Title{
		s: s,

		bg:      bg,
		bgDelay: time.Tick(time.Second / 3),
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

	if t.s.KeyPress(key.CodeReturnEnter) {
		t.s.EnterRoom("game")
		return
	}

	t.s.Draw(t.bg, image.ZP)

	t.s.Publish()
}

//go:generate ./bintogo ./images/player.png
//go:generate ./bintogo ./images/a1.png

type Game struct {
	s *State

	player *Sprite
}

func NewGame(s *State) *Game {
	player, err := s.LoadAnim(bytes.NewReader(playerData[:]), 20)
	if err != nil {
		log.Fatalf("Failed to load player: %v", err)
	}

	return &Game{
		s: s,

		player: NewSprite(player, time.Second/6),
	}
}

func (g *Game) Enter() {
}

func (g *Game) Update() {
	g.s.Fill(g.s.Bounds(), color.Black)

	g.player.Draw(g.s)

	g.s.Publish()
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
		s.AddRoom("game", NewGame(s))
		s.EnterRoom("title")

		return true
	})
	if err != nil {
		log.Fatalf("Failed to run game: %v", err)
	}
}

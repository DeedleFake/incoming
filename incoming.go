package main

import (
	"bytes"
	"golang.org/x/mobile/event/key"
	"image"
	"image/color"
	"log"
	"math/rand"
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

		bg: bg,
	}
}

func (t *Title) Enter() {
	t.bg.Start(time.Second / 3)
}

func (t *Title) Update() {
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

	player    *Anim
	playerLoc image.Point

	asteroids []*Anim
	asteroidb []image.Rectangle

	startFrame int
	won        bool
}

func NewGame(s *State) *Game {
	player, err := s.LoadAnim(bytes.NewReader(playerData[:]), 20)
	if err != nil {
		log.Fatalf("Failed to load player: %v", err)
	}

	return &Game{
		s: s,

		player:    player,
		playerLoc: image.Pt(10, 10),
	}
}

func (g *Game) Enter() {
	g.startFrame = g.s.Frame()
}

func (g *Game) Update() {
	const (
		Length = 1000

		AsteroidNum    = 3
		AsteroidChance = 25
	)

	if g.s.Frame() > g.startFrame+Length {
		g.won = true
	}

	delay := time.Second / 6
	if g.s.KeyDown(key.CodeUpArrow) {
		g.playerLoc.Y--
	}
	if g.s.KeyDown(key.CodeDownArrow) {
		g.playerLoc.Y++
	}
	if g.s.KeyDown(key.CodeLeftArrow) {
		delay = time.Second / 3
		g.playerLoc.X--
	}
	if g.s.KeyDown(key.CodeRightArrow) {
		delay = time.Second / 12
		g.playerLoc.X++
	}
	g.player.Start(delay)

	if !g.won && (len(g.asteroids) < AsteroidNum) {
		ready := true
		for _, a := range g.asteroidb {
			if a.Min.Y < a.Dy() {
				ready = false
				break
			}
		}

		if ready && (rand.Int()%AsteroidChance != 0) {
			g.asteroidAdd()
		}
	}

	g.s.Fill(g.s.Bounds(), color.Black)

	g.s.Draw(g.player, g.playerLoc)

	g.s.Publish()
}

func (g *Game) asteroidAdd() {
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

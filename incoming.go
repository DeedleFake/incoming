package main

import (
	"bytes"
	"github.com/DeedleFake/incoming/engine"
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
	s *engine.State

	bg      *engine.Anim
	bgDelay <-chan time.Time
}

func NewTitle(s *engine.State) *Title {
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
	s *engine.State

	player    *engine.Anim
	playerLoc image.Point

	asteroid *engine.Anim
	// TODO: Make an asteroid struct so that only one list is necessary.
	asteroids []*engine.Anim
	asteroidb []image.Rectangle

	startFrame int
	won        bool
}

func NewGame(s *engine.State) *Game {
	player, err := s.LoadAnim(bytes.NewReader(playerData[:]), 40)
	if err != nil {
		log.Fatalf("Failed to load player: %v", err)
	}

	asteroid, err := s.LoadAnim(bytes.NewReader(a1Data[:]), 32)
	if err != nil {
		log.Fatalf("Failed to load asteroid: %v", err)
	}

	return &Game{
		s: s,

		player:    player,
		playerLoc: image.Pt(10, 10),

		asteroid: asteroid,
	}
}

func (g *Game) Enter() {
	g.startFrame = g.s.Frame()
}

func (g *Game) Update() {
	const (
		Length = 1000

		PlayerSpeed = 2

		AsteroidNum    = 3
		AsteroidChance = 25
		AsteroidSpeed  = 2
	)

	if g.s.Frame() > g.startFrame+Length {
		g.won = true
	}

	delay := time.Second / 6
	if g.s.KeyDown(key.CodeUpArrow) {
		g.playerLoc.Y -= PlayerSpeed
	}
	if g.s.KeyDown(key.CodeDownArrow) {
		g.playerLoc.Y += PlayerSpeed
	}
	if g.s.KeyDown(key.CodeLeftArrow) {
		delay = time.Second / 3
		g.playerLoc.X -= PlayerSpeed
	}
	if g.s.KeyDown(key.CodeRightArrow) {
		delay = time.Second / 12
		g.playerLoc.X += PlayerSpeed
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

	for i := 0; i < len(g.asteroidb); i++ {
		g.asteroidb[i].Min.Y += AsteroidSpeed
		g.asteroidb[i].Max.Y += AsteroidSpeed
		if g.asteroidb[i].Min.Y >= g.s.Bounds().Max.Y {
			g.asteroidRemove(i)
			i--
		}
	}

	g.s.Fill(g.s.Bounds(), color.Black)

	g.s.Draw(g.player, g.playerLoc)

	for i, a := range g.asteroids {
		g.s.Draw(a, g.asteroidb[i].Min)
	}

	g.s.Publish()
}

func (g *Game) asteroidAdd() {
	a := g.asteroid.Copy()
	a.Start(time.Second / time.Duration(5+rand.Intn(2)))
	g.asteroids = append(g.asteroids, a)

	w := a.Size().X

	x := rand.Intn(g.s.Bounds().Dx() - w)
	y := -a.Size().Y
	g.asteroidb = append(g.asteroidb, image.Rectangle{
		Min: image.Pt(x, y),
		Max: image.Pt(x+w, 0),
	})
}

func (g *Game) asteroidRemove(i int) {
	g.asteroids[i].Stop()
	g.asteroids = append(g.asteroids[:i], g.asteroids[i+1:]...)

	g.asteroidb = append(g.asteroidb[:i], g.asteroidb[i+1:]...)
}

//go:generate ./bintogo ./images/lose.png

type Lose struct {
}

//go:generate ./bintogo ./images/win.png

type Win struct {
}

//go:generate rm -f ./bintogo

func main() {
	s := engine.NewState()
	opts := engine.StateOptions{
		Width:  480,
		Height: 320,
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

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
	t.bg.Start(time.Second/3, true)
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

	player *Player

	asteroid  *engine.Anim
	asteroids []*Asteroid

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

		player: NewPlayer(player,
			image.Pt(
				s.Bounds().Dx()/2-player.Size().X/2,
				s.Bounds().Dy()/2-player.Size().Y/2,
			),
		),

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
		g.player.Move(0, -PlayerSpeed)
	}
	if g.s.KeyDown(key.CodeDownArrow) {
		g.player.Move(0, PlayerSpeed)
	}
	if g.s.KeyDown(key.CodeLeftArrow) {
		delay = time.Second / 3
		g.player.Move(-PlayerSpeed, 0)
	}
	if g.s.KeyDown(key.CodeRightArrow) {
		delay = time.Second / 12
		g.player.Move(PlayerSpeed, 0)
	}
	g.player.Anim().Start(delay, true)

	playerB := g.player.Bounds()

	if !g.won && (len(g.asteroids) < AsteroidNum) {
		ready := true
		for _, a := range g.asteroids {
			if a.Bounds().Min.Y < a.Bounds().Dy() {
				ready = false
				break
			}
		}

		if ready && (rand.Int()%AsteroidChance != 0) {
			g.asteroidAdd()
		}
	}

	for i := 0; i < len(g.asteroids); i++ {
		asteroid := g.asteroids[i]

		asteroid.Move(0, AsteroidSpeed)
		if asteroid.Bounds().Min.Y >= g.s.Bounds().Max.Y {
			g.asteroidRemove(i)

			i--
			continue
		}

		if asteroid.Bounds().Overlaps(playerB) {
			g.s.EnterRoom("lose")
			return
		}
	}

	g.s.Fill(g.s.Bounds(), color.Black)

	g.s.Draw(g.player.Anim(), g.player.Loc())

	for _, a := range g.asteroids {
		g.s.Draw(a.Anim(), a.Bounds().Min)
	}

	g.s.Publish()
}

func (g *Game) asteroidAdd() {
	a := g.asteroid.Copy()
	a.Start(time.Second/time.Duration(5+rand.Intn(2)), true)

	w := a.Size().X
	s := image.Pt(
		rand.Intn(g.s.Bounds().Dx()-w),
		-a.Size().Y,
	)

	g.asteroids = append(g.asteroids, NewAsteroid(a, s))
}

func (g *Game) asteroidRemove(i int) {
	g.asteroids[i].Anim().Stop()
	g.asteroids = append(g.asteroids[:i], g.asteroids[i+1:]...)
}

//go:generate ./bintogo ./images/lose.png

type Lose struct {
	s *engine.State

	bg *engine.Anim
}

func NewLose(s *engine.State) *Lose {
	bg, err := s.LoadAnim(bytes.NewReader(loseData[:]), s.Bounds().Dx())
	if err != nil {
		log.Fatal("Failed to load lose background: %v", err)
	}

	return &Lose{
		s: s,

		bg: bg,
	}
}

func (l *Lose) Enter() {
	l.bg.Start(5*time.Second/6, false)
}

func (l *Lose) Update() {
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
		s.AddRoom("lose", NewLose(s))
		s.EnterRoom("title")

		return true
	})
	if err != nil {
		log.Fatalf("Failed to run game: %v", err)
	}
}

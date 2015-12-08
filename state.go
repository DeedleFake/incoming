package main

import (
	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"io"
	"log"
	"time"

	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

type State struct {
	frame int

	rooms map[string]Room
	room  Room

	s    screen.Screen
	win  screen.Window
	bnds image.Rectangle

	eventsDone chan struct{}
	fps        *time.Ticker
}

func NewState() *State {
	return &State{
		rooms: make(map[string]Room),
	}
}

func (s State) LoadAnim(r io.Reader, frameW int) (*Anim, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}

	buf, err := s.s.NewBuffer(img.Bounds().Size())
	if err != nil {
		return nil, err
	}
	defer buf.Release()

	draw.Draw(buf.RGBA(), buf.Bounds(), img, img.Bounds().Min, draw.Src)

	tex, err := s.s.NewTexture(buf.Size())
	if err != nil {
		return nil, err
	}

	tex.Upload(image.ZP, buf, buf.Bounds(), s.win)

	return newAnim(tex, frameW)
}

func (s *State) Draw(anim *Anim, dst image.Point) {
	screen.Copy(s.win, dst, anim.image, anim.cur, draw.Over, nil)
}

func (s *State) Fill(r image.Rectangle, c color.Color) {
	s.win.Fill(r, c, draw.Over)
}

func (s *State) Publish() {
	s.win.Publish()
}

func (s *State) Bounds() image.Rectangle {
	return s.bnds
}

func (s *State) eventsStart() {
	ev := s.win.Events()

	keys := make(map[key.Code]bool)

	s.eventsDone = make(chan struct{})
	for {
		select {
		case ev := <-ev:
			switch ev := ev.(type) {
			case key.Event:
				keys[ev.Code] = ev.Direction != key.DirRelease
			case error:
				log.Printf("Event error: %v", ev)
			}

		case <-s.eventsDone:
			return
		}
	}
}

func (s *State) eventsStop() {
	close(s.eventsDone)
}

func (s *State) Run(opts *StateOptions, init func() bool) (reterr error) {
	if opts == nil {
		opts = &DefaultStateOptions
	}

	driver.Main(func(scrn screen.Screen) {
		s.s = scrn

		win, err := scrn.NewWindow(&screen.NewWindowOptions{
			Width:  opts.Width,
			Height: opts.Height,
		})
		if err != nil {
			reterr = err
			return
		}
		s.win = win
		s.bnds = image.Rect(0, 0, opts.Width, opts.Height)

		if !init() {
			return
		}

		go s.eventsStart()
		defer s.eventsStop()

		s.fps = time.NewTicker(time.Second / time.Duration(opts.FPS))
		defer s.fps.Stop()

		for {
			s.room.Update()
			s.frame++

			<-s.fps.C
		}
	})

	return
}

func (s State) Frame() int {
	return s.frame
}

func (s *State) AddRoom(name string, room Room) {
	s.rooms[name] = room
}

func (s *State) EnterRoom(name string) {
	s.room = s.rooms[name]
	s.room.Enter()
}

type StateOptions struct {
	Width  int
	Height int
	FPS    int
}

var DefaultStateOptions = StateOptions{
	Width:  640,
	Height: 480,
	FPS:    60,
}

type Room interface {
	Enter()
	Update()
}

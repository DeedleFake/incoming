package main

import (
	"golang.org/x/exp/shiny/driver/gldriver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"io"

	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

type State struct {
	frame int

	rooms map[string]Room
	room  Room

	s   screen.Screen
	win screen.Window

	eventsDone chan struct{}
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

	sender := make(dummySender)
	tex.Upload(image.ZP, buf, buf.Bounds(), sender)
	ev := <-sender
	if err, ok := ev.(error); ok {
		return nil, err
	}

	return newAnim(tex, frameW)
}

func (s *State) Draw(anim *Anim, dst image.Point) {
	screen.Copy(s.win, dst, anim.image, anim.cur, draw.Over, nil)
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
			}

		case <-s.eventsDone:
			return
		}
	}
}

func (s *State) eventsStop() {
	close(s.eventsDone)
}

func (s *State) Run(opts *StateOptions, init func(), tick func()) (reterr error) {
	if opts == nil {
		opts = &DefaultStateOptions
	}

	gldriver.Main(func(scrn screen.Screen) {
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

		init()

		go s.eventsStart()
		defer s.eventsStop()

		for {
			tick()
			s.frame++
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
}

var DefaultStateOptions = StateOptions{
	Width:  640,
	Height: 480,
}

type Room interface {
	Enter()
	Update()
}

type dummySender chan interface{}

func (ds dummySender) Send(ev interface{}) {
	ds <- ev
}

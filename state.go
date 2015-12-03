package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type State struct {
	frame int

	rooms map[string]Room
	room  Room

	win *sdl.Window
	ren *sdl.Renderer
}

func NewState(title string, w, h int) (*State, error) {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		return nil, err
	}

	win, ren, err := sdl.CreateWindowAndRenderer(w, h, 0)
	if err != nil {
		sdl.Quit()
		return nil, err
	}
	win.SetTitle(title)

	return &State{
		rooms: make(map[string]Room),

		win: win,
		ren: ren,
	}, nil
}

func (s *State) Update() bool {
	defer func() {
		s.frame++
	}()

	for {
		ev := sdl.PollEvent()
		if ev == nil {
			break
		}

		switch ev.(type) {
		case *sdl.QuitEvent:
			return false
		}
	}

	s.room.Update()

	return true
}

func (s *State) Frame() int {
	return s.frame
}

func (s *State) AddRoom(name string, room Room) {
	s.rooms[name] = room
}

func (s *State) EnterRoom(name string) {
	s.room = s.rooms[name]
	s.room.Enter()
}

type Room interface {
	Enter()
	Update()
}

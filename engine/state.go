package engine

import (
	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"io"
	"log"
	"sync"
	"time"

	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

// State is the main struct of the engine. It keeps track of pretty
// much everything, as well as running the actual game.
type State struct {
	frame int

	rooms map[string]Room
	room  Room

	s    screen.Screen
	win  screen.Window
	bnds image.Rectangle

	eventsDone chan struct{}
	fps        *time.Ticker

	kqpool sync.Pool
	kqc    chan *keyQuery
}

// NewState initializes a new state.
func NewState() *State {
	kqc := make(chan *keyQuery)

	return &State{
		rooms: make(map[string]Room),

		kqpool: sync.Pool{
			New: func() interface{} {
				return &keyQuery{
					c: kqc,
					r: make(chan bool),
				}
			},
		},
		kqc: kqc,
	}
}

// LoadAnim loads an animation from r with the given frame width. The
// width of the image loaded must be evenly divisible by the frame
// width.
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

	tex.Upload(image.ZP, buf, buf.Bounds())

	return newAnim(tex, frameW)
}

// Draw draws an image onto the screen at the specified point.
func (s State) Draw(img Imager, dst image.Point) {
	image, clip := img.Image()
	// TODO: According to the race detector, this is causing a data
	// race. I'm not sure how that's possible, since everything being
	// dealt with here should only ever be accessed from one thread.
	// Maybe it's a Shiny bug?
	screen.Copy(s.win, dst, image, clip, draw.Over, nil)
}

// Fill fills r with c on the screen.
func (s State) Fill(r image.Rectangle, c color.Color) {
	s.win.Fill(r, c, draw.Over)
}

// Publish updates the screen. Changes to the screen are not
// guarunteed to actually appear until Publish() has been called.
func (s State) Publish() {
	s.win.Publish()
}

// Bounds returns the screen's bounds.
func (s State) Bounds() image.Rectangle {
	return s.bnds
}

func (s *State) eventsStart() {
	ev := s.win.Events()

	keys := make(map[key.Code]bool)
	keysCheck := make(map[key.Code]bool)

	// TODO: This causes a potential data race.
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

		case kq := <-s.kqc:
			down := keys[kq.code]
			if !kq.press {
				kq.r <- keys[kq.code]
				continue
			}

			if down {
				kq.r <- !keysCheck[kq.code]
				keysCheck[kq.code] = true
			} else {
				kq.r <- false
				keysCheck[kq.code] = false
			}

		case <-s.eventsDone:
			return
		}
	}
}

func (s *State) eventsStop() {
	close(s.eventsDone)
}

// KeyDown returns whether or not a specific key is currently being
// held down.
func (s *State) KeyDown(code key.Code) bool {
	// Can you say 'overkill'?
	kq := s.kqpool.Get().(*keyQuery)
	defer s.kqpool.Put(kq)

	return kq.Q(code, false)
}

// KeyPress returns true the first time it is called after the key
// represented by code is pressed, but returns false for that key
// until the key is released and then pressed again.
func (s *State) KeyPress(code key.Code) bool {
	kq := s.kqpool.Get().(*keyQuery)
	defer s.kqpool.Put(kq)

	return kq.Q(code, true)
}

// Run starts the game. It blocks until an unrecoverable error occurs
// or until the game otherwise exits.
//
// Before the game enters the mainloop, init is called. If init
// returns false, the game exits. init may be nil.
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

		if (init != nil) && !init() {
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

// Frame returns the current frame of the game, with 0 being the
// initial frame.
func (s State) Frame() int {
	return s.frame
}

// AddRoom adds a room to the game, or overwrites an existing room if
// one named name already exists.
func (s *State) AddRoom(name string, room Room) {
	s.rooms[name] = room
}

// EnterRoom switches to the room named name, calling its Enter()
// method.
func (s *State) EnterRoom(name string) {
	s.room = s.rooms[name]
	s.room.Enter()
}

// StateOptions are options for initializing a State.
type StateOptions struct {
	// The width and height of the screen.
	Width  int
	Height int

	// The target FPS of the game.
	//
	// TODO: Make it possible to remove the FPS cap.
	FPS int
}

// DefaultStateOptions are the options used if NewState is passed nil.
var DefaultStateOptions = StateOptions{
	Width:  640,
	Height: 480,
	FPS:    60,
}

// Room represents a room of the game. This can be almost anything,
// from a menu, to a level, to credits.
type Room interface {
	Enter()
	Update()
}

// Imager is a type that can represent itself as a piece of a texture.
type Imager interface {
	// Image returns a texture and a clipping rectangle.
	Image() (screen.Texture, image.Rectangle)
}

type keyQuery struct {
	c chan *keyQuery
	r chan bool

	code  key.Code
	press bool
}

func (kq *keyQuery) Q(code key.Code, press bool) bool {
	kq.code = code
	kq.press = press

	kq.c <- kq
	return <-kq.r
}

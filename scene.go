package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type scene struct {
	time  int
	bg    *sdl.Texture
	bird  *bird
	birds []*sdl.Texture
	pipes *pipes
}

func newScene(r *sdl.Renderer) (*scene, error) {
	bg, err := img.LoadTexture(r, "res/imgs/background.png")
	if err != nil {
		return nil, fmt.Errorf("could not load background image: %v", err)
	}
	b, err := newBird(r)
	if err != nil {
		return nil, err
	}

	ps, err := newPipes(r)
	if err != nil {
		return nil, err
	}

	return &scene{bg: bg, bird: b, pipes: ps}, nil
}

func (s *scene) run(events <-chan sdl.Event, r *sdl.Renderer) chan error {
	errc := make(chan error)
	go func() {
		defer close(errc)
		done := false
		tick := time.Tick(10 * time.Millisecond)
		for !done {
			select {
			case e := <-events:
				done = s.handleEvent(e)
				log.Printf("event: %T", e)
			case <-tick:
				s.update()
				if s.bird.isDead() {
					drawTitle(r, "Game Over", 30)
					time.Sleep(time.Second)
					score := fmt.Sprintf("Score: %d", passed)
					drawTitle(r, score, 30)
					time.Sleep(time.Second)
					s.restart()
				}
				if err := s.paint(r); err != nil {
					errc <- err
				}
			}
		}
	}()
	return errc
}

func (s *scene) update() {
	s.bird.update()
	s.pipes.update()
	s.pipes.touch(s.bird)
}

func (s *scene) restart() {
	passed = 0
	s.bird.restart()
	s.pipes.restart()
}
func (s *scene) paint(r *sdl.Renderer) error {
	s.time++
	r.Clear()
	if err := r.Copy(s.bg, nil, nil); err != nil {
		return fmt.Errorf("could not copy background: %v", err)
	}

	if err := s.bird.paint(r); err != nil {
		return err
	}

	if err := s.pipes.paint(r); err != nil {
		return err
	}

	r.Present()
	return nil
}

func (s *scene) destroy() {
	s.bg.Destroy()
	s.bird.destroy()
	s.pipes.destroy()
}

func (s *scene) handleEvent(event sdl.Event) bool {
	switch event.(type) {
	case *sdl.QuitEvent:
		return true
	case *sdl.MouseButtonEvent:
		s.bird.jump()
		return false
	case *sdl.MouseMotionEvent, *sdl.WindowEvent, *sdl.TouchFingerEvent, *sdl.CommonEvent:
		return false
	case *sdl.KeyUpEvent, *sdl.KeyDownEvent:
		// for s.bird.x < 100 {
		// 	s.bird.x++
		// }
		// count := 0
		// for s.bird.x > 20 {
		// 	count++
		// 	if count%2 == 0 {
		// 		s.bird.x--
		// 	}
		// }
		return false
	default:
		log.Printf("unknown event %T", event)
	}
	return false
}

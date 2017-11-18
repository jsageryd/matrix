package main

import (
	"container/list"
	"io"
	"math/rand"
	"time"

	"github.com/gdamore/tcell"
)

type matrix struct {
	screen    tcell.Screen
	now       time.Duration
	seedFeed  *rand.Rand
	feed      io.Reader
	color     string
	xdensity  float32 // the probability (0-1) that a stream will spawn at a column
	frequency int     // the number of stream rows per second
	segments  *list.List
}

func newMatrix(seed int64, screen tcell.Screen, color string) *matrix {
	return &matrix{
		screen:    screen,
		seedFeed:  rand.New(rand.NewSource(seed)),
		feed:      randomRuneFeed{runes: []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!\"#€%&/()=?<>,.-;:_'^*$|[]\\{}")},
		color:     color,
		frequency: 10,
		xdensity:  0.03,
		segments:  list.New(),
	}
}

func (m *matrix) enter() error {
	if err := m.screen.Init(); err != nil {
		return err
	}

	stop := make(chan struct{})

	go func() {
		for {
			ev, ok := m.screen.PollEvent().(*tcell.EventKey)
			if !ok {
				continue
			}

			if ev.Key() == tcell.KeyEsc || ev.Key() == tcell.KeyCtrlC || (ev.Key() == tcell.KeyRune && ev.Rune() == 'q') {
				close(stop)
			}

			if ev.Key() == tcell.KeyRune {
				switch ev.Rune() {
				case 'w':
					m.color = "white"
				case 'b':
					m.color = "blue"
				case 'g':
					m.color = "green"
				case 'r':
					m.color = "red"
				case 'y':
					m.color = "yellow"
				case 'o':
					m.color = "orange"
				case 'm':
					m.color = "magenta"
				case 'c':
					m.color = "cyan"
				}
			}
		}
	}()

	go func() {
		for range time.Tick(time.Second / 20) {
			select {
			case <-stop:
				return
			default:
				m.step(time.Second / 20)
				m.draw()
				m.screen.Show()
			}
		}
	}()

	<-stop

	return nil
}

func (m *matrix) step(d time.Duration) {
	width, height := m.screen.Size()

	// Kill old segments
	var next *list.Element
	for e := m.segments.Front(); e != nil; e = next {
		next = e.Next()
		s := e.Value.(*segment)
		if s.position(m.now)-s.length() > height {
			m.segments.Remove(e)
		}
	}

	// Round to closest previous multiple of the frequency
	now := int64(m.now)
	step := int64(time.Second) / int64(m.frequency)
	start := now - now%step

	// Move forward in case we are in the past
	if start < now {
		start += step
	}

	for n := start; n < now+int64(d); n += step {
		rng := rand.New(rand.NewSource(m.seedFeed.Int63()))
		for x := 0; x < width; x++ {
			if rng.Float32() <= m.xdensity {
				len := rng.Intn(height/2) + 3
				speed := rng.Intn(15) + 5
				shiny := rng.Float32() > 0.8
				s := newSegment(m.screen, m.feed, x, len, m.now, m.color, speed, shiny)
				m.segments.PushBack(s)
			}
		}
	}

	m.now += d
}

func (m *matrix) draw() {
	for e := m.segments.Front(); e != nil; e = e.Next() {
		e.Value.(*segment).draw(m.now)
	}
}

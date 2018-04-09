package main

import (
	"container/list"
	"io"
	"math/rand"
	"sync"
	"time"

	"github.com/gdamore/tcell"
)

var (
	feedAlpha = randomRuneFeed{runes: []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!\"#€%&/()=?<>,.-;:_'^*$|[]\\{}")}
	feedCyril = randomRuneFeed{runes: []rune("АБВГДЕЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯабвгдежзийклмнопрстуфхцчшщъыьэюя0123456789!\"#€%&/()=?<>,.-;:_'^*$|[]\\{}")}
	feedDot   = randomRuneFeed{runes: []rune(".")}
	feedKata  = randomRuneFeed{runes: []rune("アイウエオカキクケコガギグゲゴサシスセソザジズゼゾタチツテトダヂヅデドナニヌネノハヒフヘホバビブベボパピプペポマミムメモヤユヨラリルレロワン")}
)

type matrix struct {
	mu        sync.RWMutex
	screen    tcell.Screen
	now       time.Duration
	seedFeed  *rand.Rand
	feed      io.Reader
	color     string
	xdensity  float32 // the probability (0-1) that a stream will spawn at a column
	frequency int     // the number of stream rows per second
	segments  *list.List
}

func newMatrix(seed int64, screen tcell.Screen, color string, feed io.Reader) *matrix {
	return &matrix{
		screen:    screen,
		seedFeed:  rand.New(rand.NewSource(seed)),
		feed:      feed,
		color:     color,
		frequency: 10,
		xdensity:  0.03,
		segments:  list.New(),
	}
}

func (m *matrix) getColor() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.color
}

func (m *matrix) setColor(c string) {
	m.mu.Lock()
	m.color = c
	m.mu.Unlock()
}

func (m *matrix) getFeed() io.Reader {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.feed
}

func (m *matrix) setFeed(f io.Reader) {
	m.mu.Lock()
	m.feed = f
	m.mu.Unlock()
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
				return
			}

			if ev.Key() == tcell.KeyRune {
				switch ev.Rune() {
				case 'w':
					m.setColor("white")
				case 'b':
					m.setColor("blue")
				case 'g':
					m.setColor("green")
				case 'r':
					m.setColor("red")
				case 'y':
					m.setColor("yellow")
				case 'o':
					m.setColor("orange")
				case 'm':
					m.setColor("magenta")
				case 'c':
					m.setColor("cyan")
				case 'a':
					m.setFeed(feedAlpha)
				case 'C':
					m.setFeed(feedCyril)
				case 'd':
					m.setFeed(feedDot)
				case 'k':
					m.setFeed(feedKata)
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

	color := m.getColor()
	feed := m.getFeed()

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
				s := newSegment(m.screen, feed, x, len, m.now, color, speed, shiny)
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

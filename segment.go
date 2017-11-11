package main

import (
	"bufio"
	"container/ring"
	"io"
	"math/rand"

	termbox "github.com/nsf/termbox-go"
)

type segment struct {
	feed  *bufio.Reader
	r     *ring.Ring
	white bool
}

func newSegment(feed io.Reader, length int) *segment {
	return &segment{
		feed:  bufio.NewReader(feed),
		r:     ring.New(length),
		white: rand.Float32() > 0.8,
	}
}

func (s *segment) step() {
	s.r = s.r.Next()
	s.r.Value, _, _ = s.feed.ReadRune()
}

func (s *segment) rune(n int) rune {
	if r, ok := s.r.Move(n).Value.(rune); ok {
		return r
	}
	return 0
}

func (s *segment) draw(x, y int) {
	if y >= s.r.Len() {
		termbox.SetCell(x, y-s.r.Len(), ' ', termbox.ColorDefault, termbox.ColorDefault)
	}
	if s.white {
		for offset := 0; offset < min(5, s.r.Len()); offset++ {
			termbox.SetCell(x, y-offset, s.rune(-offset), termbox.Attribute(255-offset*3), termbox.ColorDefault)
		}
	} else {
		termbox.SetCell(x, y, s.rune(0), termbox.Attribute(240), termbox.ColorDefault)
	}
}

func min(i, j int) int {
	if i < j {
		return i
	}
	return j
}

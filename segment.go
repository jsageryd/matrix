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
	color string
	shiny bool
}

func newSegment(feed io.Reader, length int, color string) *segment {
	return &segment{
		feed:  bufio.NewReader(feed),
		r:     ring.New(length),
		color: color,
		shiny: rand.Float32() > 0.8,
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
	if s.shiny {
		for offset := 0; offset < min(5, s.r.Len()); offset++ {
			termbox.SetCell(x, y-offset, s.rune(-offset), termbox.Attribute(s.colorShade(4-offset)), termbox.ColorDefault)
		}
	} else {
		termbox.SetCell(x, y, s.rune(0), termbox.Attribute(s.colorShade(0)), termbox.ColorDefault)
	}
}

func (s *segment) colorShade(n int) int {
	n = max(0, n)
	n = min(4, n)

	switch s.color {
	case "white":
		return []int{240, 244, 248, 252, 255}[n]
	default:
		return []int{240, 244, 248, 252, 255}[n]
	}
}

func max(i, j int) int {
	if i > j {
		return i
	}
	return j
}

func min(i, j int) int {
	if i < j {
		return i
	}
	return j
}

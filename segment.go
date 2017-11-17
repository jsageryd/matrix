package main

import (
	"bufio"
	"container/ring"
	"io"
	"time"

	termbox "github.com/nsf/termbox-go"
)

type segment struct {
	feed      *bufio.Reader
	r         *ring.Ring
	column    int
	startTime time.Duration
	color     string
	shiny     bool
	speed     int // blocks per second
	runesRead int
}

func newSegment(feed io.Reader, column, length int, startTime time.Duration, color string, speed int, shiny bool) *segment {
	return &segment{
		feed:      bufio.NewReader(feed),
		r:         ring.New(length),
		color:     color,
		startTime: startTime,
		shiny:     shiny,
		column:    column,
		speed:     speed,
	}
}

func (s *segment) rune(n int) rune {
	if r, ok := s.r.Move(n).Value.(rune); ok {
		return r
	}
	return 0
}

func (s *segment) position(now time.Duration) int {
	return int(((now - s.startTime) * time.Duration(s.speed)) / time.Second)
}

func (s *segment) length() int {
	return s.r.Len()
}

func (s *segment) draw(now time.Duration) {
	y := s.position(now)

	for s.runesRead <= y {
		s.r = s.r.Next()
		s.r.Value, _, _ = s.feed.ReadRune()
		s.runesRead++
	}

	// This takes care of cleaning the path behind the segment as it progresses.
	// It will break if the speed of the segments is higher than the refresh rate,
	// but this is typically not the case so it should be ok to ignore for now.
	if y >= s.r.Len() {
		termbox.SetCell(s.column, y-s.r.Len(), ' ', termbox.ColorDefault, termbox.ColorDefault)
	}

	if s.shiny {
		for offset := 0; offset < min(5, s.r.Len()); offset++ {
			termbox.SetCell(s.column, y-offset, s.rune(-offset), termbox.Attribute(s.colorShade(4-offset)), termbox.ColorDefault)
		}
	} else {
		termbox.SetCell(s.column, y, s.rune(0), termbox.Attribute(s.colorShade(0)), termbox.ColorDefault)
	}
}

func (s *segment) colorShade(n int) int {
	n = max(0, n)
	n = min(4, n)

	switch s.color {
	case "white":
		return []int{240, 244, 248, 252, 255}[n] + 1
	case "blue":
		return []int{18, 19, 20, 27, 33}[n] + 1
	case "green":
		return []int{22, 28, 34, 40, 46}[n] + 1
	case "red":
		return []int{52, 88, 124, 160, 196}[n] + 1
	case "yellow":
		return []int{94, 94, 136, 220, 226}[n] + 1
	case "orange":
		return []int{166, 202, 208, 214, 220}[n] + 1
	case "pink":
		return []int{89, 126, 162, 198, 199}[n] + 1
	default:
		return []int{240, 244, 248, 252, 255}[n] + 1
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

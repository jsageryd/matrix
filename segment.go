package main

import (
	"bufio"
	"container/ring"
	"io"
	"time"

	"github.com/gdamore/tcell"
)

type segment struct {
	screen    tcell.Screen
	feed      *bufio.Reader
	r         *ring.Ring
	column    int
	startTime time.Duration
	color     string
	shiny     bool
	speed     int // blocks per second
	runesRead int
}

func newSegment(screen tcell.Screen, feed io.Reader, column, length int, startTime time.Duration, color string, speed int, shiny bool) *segment {
	return &segment{
		screen:    screen,
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
		s.screen.SetContent(s.column, y-s.r.Len(), ' ', nil, tcell.StyleDefault)
	}

	if s.shiny {
		for offset := 0; offset < min(5, s.r.Len()); offset++ {
			s.screen.SetContent(s.column, y-offset, s.rune(-offset), nil, tcell.StyleDefault.Foreground(tcell.Color(s.colorShade(4-offset))))
		}
	} else {
		s.screen.SetContent(s.column, y, s.rune(0), nil, tcell.StyleDefault.Foreground(tcell.Color(s.colorShade(0))))
	}
}

func (s *segment) colorShade(n int) int {
	n = max(0, n)
	n = min(4, n)

	switch s.color {
	case "white":
		return []int{240, 244, 248, 252, 255}[n]
	case "blue":
		return []int{18, 19, 20, 27, 33}[n]
	case "green":
		return []int{22, 28, 34, 40, 46}[n]
	case "red":
		return []int{52, 88, 124, 160, 196}[n]
	case "yellow":
		return []int{94, 94, 136, 220, 226}[n]
	case "orange":
		return []int{166, 202, 208, 214, 220}[n]
	case "magenta":
		return []int{89, 126, 162, 198, 199}[n]
	case "cyan":
		return []int{31, 45, 51, 87, 195}[n]
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

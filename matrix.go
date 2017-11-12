package main

import (
	"io"
	"math/rand"
	"time"

	termbox "github.com/nsf/termbox-go"
)

type matrix struct {
	feed             io.Reader
	streamsPerSecond int
}

func newMatrix() *matrix {
	return &matrix{
		feed:             randomRuneFeed{runes: []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!\"#€%&/()=?<>,.-;:_'^*$|[]\\{}")},
		streamsPerSecond: 50,
	}
}

func (m *matrix) enter() {
	stop := make(chan struct{})

	go func() {
		for {
			ev := termbox.PollEvent()
			if ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyCtrlC || ev.Ch == 'q' {
				close(stop)
			}
		}
	}()

	go func() {
		for range time.Tick(time.Second / 60) {
			select {
			case <-stop:
				return
			default:
				termbox.Flush()
			}
		}
	}()

ticker:
	for range time.Tick(time.Second / time.Duration(m.streamsPerSecond)) {
		select {
		case <-stop:
			break ticker
		default:
			go m.stream(stop)
		}
	}
}

func (m *matrix) stream(stop <-chan struct{}) {
	width, height := termbox.Size()
	x := rand.Intn(width)
	segmentLength := rand.Intn(height/2) + 1
	s := newSegment(m.feed, segmentLength, "white")

	speed := 50 + (rand.Intn(100))

	for y := 0; y < height+segmentLength; y++ {
		select {
		case <-stop:
			break
		default:
			time.Sleep(time.Duration(speed) * time.Millisecond)
			s.step()
			s.draw(x, y)
		}
	}
}

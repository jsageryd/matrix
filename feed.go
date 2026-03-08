package main

import (
	"math/rand"
	"strings"
)

type randomRuneFeed struct {
	runes []rune
}

func (r randomRuneFeed) Read(p []byte) (int, error) {
	return strings.NewReader(string(r.runes[rand.Intn(len(r.runes))])).Read(p)
}

type randomRuneRangeFeed struct {
	from, to  rune
	maxRunes  int
	rangeSize int
	offset    int
}

func (r *randomRuneRangeFeed) Read(p []byte) (int, error) {
	if r.rangeSize == 0 {
		r.rangeSize = int(r.to) - int(r.from) + 1
		if r.maxRunes == 0 {
			r.maxRunes = r.rangeSize
		}
		r.maxRunes = min(r.maxRunes, r.rangeSize)
		r.offset = rand.Intn(r.rangeSize)
	}
	return strings.NewReader(string(r.from + rune((rand.Intn(r.maxRunes)*r.rangeSize/r.maxRunes+r.offset)%r.rangeSize))).Read(p)
}

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
	from, to rune
}

func (r randomRuneRangeFeed) Read(p []byte) (int, error) {
	return strings.NewReader(string(rune(int(r.from) + rand.Intn(int(r.to)-int(r.from)+1)))).Read(p)
}

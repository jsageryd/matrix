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

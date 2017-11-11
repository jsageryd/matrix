package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	termbox "github.com/nsf/termbox-go"
)

func main() {
	if err := termbox.Init(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer termbox.Close()

	termbox.SetOutputMode(termbox.Output256)

	rand.Seed(time.Now().UnixNano())

	newMatrix().enter()
}

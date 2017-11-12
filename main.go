package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	termbox "github.com/nsf/termbox-go"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	color := flag.String("c", "white", "colour (white, blue, green, red, yellow, orange)")
	flag.Parse()

	if err := termbox.Init(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer termbox.Close()

	termbox.SetOutputMode(termbox.Output256)

	newMatrix(*color).enter()
}

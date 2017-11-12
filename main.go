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
	if err := termbox.Init(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer termbox.Close()

	termbox.SetOutputMode(termbox.Output256)

	rand.Seed(time.Now().UnixNano())

	color := flag.String("c", "white", "colour")
	flag.Parse()

	newMatrix(*color).enter()
}

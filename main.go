package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	termbox "github.com/nsf/termbox-go"
)

func main() {
	color := flag.String("c", "white", "colour (white, blue, green, red, yellow, orange, pink)")
	flag.Parse()

	if err := termbox.Init(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer termbox.Close()

	termbox.SetOutputMode(termbox.Output256)

	newMatrix(time.Now().UnixNano(), *color).enter()
}

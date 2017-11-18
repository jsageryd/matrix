package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Matrix\n\nFlags:")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Valid colours are:\n  white, blue, green, red, yellow, orange, magenta, cyan")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Colour can also be changed by typing the initial\n(e.g. 'c' for 'cyan') while running.")
	}

	color := flag.String("c", "white", "colour")
	flag.Parse()

	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer screen.Fini()

	m := newMatrix(time.Now().UnixNano(), screen, *color)
	if err := m.enter(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

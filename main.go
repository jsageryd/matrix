package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gdamore/tcell"
)

func main() {
	exitCode := 0
	defer func() { os.Exit(exitCode) }()

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Matrix\n\nFlags:")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "| Valid colours:\n|   white, blue, green, red, yellow, orange, magenta, cyan")
		fmt.Fprintln(os.Stderr, "|")
		fmt.Fprintln(os.Stderr, "| Colour can also be changed by typing the initial\n| (e.g. 'c' for 'cyan') while running.")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "| Valid feeds: alpha, cyril, dot, kata")
		fmt.Fprintln(os.Stderr, "|")
		fmt.Fprintln(os.Stderr, "| Feed can also be changed by typing the initial\n| (e.g. 'a' for 'alpha', special case 'C' for 'cyril') while running.")
	}

	color := flag.String("c", "white", "colour")
	feedStr := flag.String("f", "alpha", "feed")
	flag.Parse()

	var feed io.Reader

	switch *feedStr {
	case "alpha":
		feed = feedAlpha
	case "cyril":
		feed = feedCyril
	case "dot":
		feed = feedDot
	case "kata":
		feed = feedKata
	default:
		feed = feedAlpha
	}

	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitCode = 1
		return
	}
	defer screen.Fini()

	m := newMatrix(time.Now().UnixNano(), screen, *color, feed)
	if err := m.enter(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitCode = 1
		return
	}
}

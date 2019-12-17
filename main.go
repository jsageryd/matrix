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
		fmt.Fprintln(os.Stderr, "| Valid feeds: alpha, binary, block, cyril, dot, greek, hangeul, hex, hira, kata, line, num, stdin, zh")
		fmt.Fprintln(os.Stderr, "|")
		fmt.Fprintln(os.Stderr, "| Feed can also be changed by typing the upper-case initial\n| (e.g. 'A' for 'alpha') while running.")
		fmt.Fprintln(os.Stderr, "| Hex feed uses X. Binary feed uses 0.")
	}

	color := flag.String("c", "white", "colour")
	feedStr := flag.String("f", "alpha", "feed")
	flag.Parse()

	var feed io.Reader

	switch *feedStr {
	case "alpha":
		feed = feedAlpha
	case "binary":
		feed = feedBinary
	case "block":
		feed = feedBlock
	case "cyril":
		feed = feedCyril
	case "dot":
		feed = feedDot
	case "greek":
		feed = feedGreek
	case "hangeul":
		feed = feedHangeul
	case "hex":
		feed = feedHex
	case "hira":
		feed = feedHira
	case "kata":
		feed = feedKata
	case "line":
		feed = feedLine
	case "num":
		feed = feedNum
	case "stdin":
		feed = feedStdin
	case "zh":
		feed = feedZh
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

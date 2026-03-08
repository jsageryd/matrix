# Matrix

[![License MIT](https://img.shields.io/badge/license-MIT-lightgrey.svg?style=flat)](#license)
[![Follow rabbit](https://img.shields.io/badge/follow-rabbit-lightgrey.svg?style=flat)](#)

Unfortunately, no one can be told what Matrix is. You will have to see it for
yourself.

![Screenshot](./matrix.png)

## Installation
```
$ go install github.com/jsageryd/matrix@latest
```

## Usage
```
$ matrix -h
```
```
Matrix

Flags:
  -c string
      colour (default "white")
  -f string
      feed (default "alpha")

| Valid colours:
|   white, blue, green, red, yellow, orange, magenta, cyan
|
| Colour can also be changed by typing the initial
| (e.g. 'c' for 'cyan') while running.

| Valid feeds: alpha, binary, block, cyril, dot, greek, hangeul, jamo, hex, hira, kata, line, num, stdin, zh
|
| Feed can also be changed by typing the upper-case initial
| (e.g. 'A' for 'alpha') while running.
| Hex feed uses X. Binary feed uses 0.
```

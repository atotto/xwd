package main

import (
	"context"
	"flag"
	"image/png"
	"log"
	"os"

	"github.com/atotto/xwd"
)

var (
	filename = flag.String("out", "screenshot.png", "output filename")
)

func main() {
	flag.Parse()

	ctx := context.Background() // TODO: signal

	f, err := os.OpenFile(*filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	m, err := xwd.Capture(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if err := png.Encode(f, m); err != nil {
		log.Fatal(err)
	}
}

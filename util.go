package xwd

import (
	"bytes"
	"context"
	"image"
	"io"
	"log"
	"os/exec"
)

func DumpXWindowImage(ctx context.Context, w io.Writer) error {
	cmd := exec.CommandContext(ctx, "xwd", "-root", "-display", ":0")

	cmd.Stdout = w
	var b bytes.Buffer
	cmd.Stderr = &b

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func Capture(ctx context.Context) (image.Image, error) {
	r, w := io.Pipe()

	go func() {
		if err := DumpXWindowImage(ctx, w); err != nil {
			log.Fatal(err)
		}
		w.Close()
	}()

	return Decode(r)
}

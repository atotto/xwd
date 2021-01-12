package xwd

import (
	"bytes"
	"context"
	"image"
	"io"
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

	var err2 error
	go func() {
		if err2 = DumpXWindowImage(ctx, w); err2 != nil {
			r.CloseWithError(io.EOF)
		}
		w.Close()
	}()

	m, err := Decode(r)
	if err2 != nil {
		return nil, err2
	}
	if err != nil {
		return nil, err
	}
	return m, nil
}

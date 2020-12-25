# xwd image library for golang

example:

```go
func main() {
	f, err := os.Create("screenshot.png")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	ctx := context.Background()
	m, err := xwd.Capture(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if err := png.Encode(f, m); err != nil {
		log.Fatal(err)
	}
}
```


```go
// $ xwd -root -display :0 -out screenshot.xwd
func xwd2png(w io.Writer, r io.Reader) error {
	m, err := xwd.Decode(r)
	if err != nil {
		return err
	}
	if err := png.Encode(w, m); err != nil {
		return err
	}
	return nil
}
```
// Package xwd implements the xwd (X Window dump) image library.
//
// https://en.wikipedia.org/wiki/Xwd
package xwd

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"io"
)

const (
	xyBitmap = 0
	xyPixmap = 1
	zPixmap  = 2
)

// xwdFileHeader
// see /usr/include/X11/XWDFile.h
type xwdFileHeader struct {
	HeaderSize        uint32 /* header_size = SIZEOF(XWDheader) + length of null-terminated window name. */
	FileVersion       uint32 /* = XWD_FILE_VERSION above */
	PixmapFormat      uint32 /* zPixmap or XYPixmap */
	PixmapDepth       uint32 /* Pixmap depth */
	PixmapWidth       uint32 /* Pixmap width */
	PixmapHeight      uint32 /* Pixmap height */
	XOffset           uint32 /* Bitmap x offset, normally 0 */
	ByteOrder         uint32 /* of image data: MSBFirst, LSBFirst */
	BitmapUnit        uint32 /* bitmap_unit applies to bitmaps (depth 1 format XY) only. It is the number of bits that each scanline is padded to. */
	BitmapBitOrder    uint32 /* bitmaps only: MSBFirst, LSBFirst */
	BitmapPad         uint32 /* bitmap_pad applies to pixmaps (non-bitmaps) only. It is the number of bits that each scanline is padded to. */
	BitsPerPixel      uint32 /* Bits per pixel */
	BytesPerLine      uint32 /* bytes_per_line is pixmap_width padded to bitmap_unit (bitmaps) or bitmap_pad (pixmaps).  It is the delta (in bytes) to get to the same x position on an adjacent row. */
	VisualClass       uint32 /* Class of colormap */
	RedMask           uint32 /* Z red mask */
	GreenMask         uint32 /* Z green mask */
	BlueMask          uint32 /* Z blue mask */
	BitsPerRgb        uint32 /* Log2 of distinct color values */
	ColormapEntries   uint32 /* Number of entries in colormap; not used? */
	NColors           uint32 /* Number of xwdColor structures */
	WindowWidth       uint32 /* Window width */
	WindowHeight      uint32 /* Window height */
	WindowX           uint32 /* Window upper left X coordinate */
	WindowY           uint32 /* Window upper left Y coordinate */
	WindowBorderWidth uint32 /* Window border width */
}

// xwdColor
// see /usr/include/X11/XWDFile.h
type xwdColor struct {
	Pixel uint32
	Red   uint16
	Green uint16
	Blue  uint16
	Flags uint8
	Pad   uint8
}

func Decode(r io.Reader) (img image.Image, err error) {
	header := xwdFileHeader{}
	if err := binary.Read(r, binary.BigEndian, &header); err != nil {
		return nil, err
	}

	if header.FileVersion != 7 {
		return nil, fmt.Errorf("not suppoted file version: %d", header.FileVersion)
	}
	if header.PixmapFormat != zPixmap {
		return nil, fmt.Errorf("not suppoted yet")
	}

	// null-terminated window name
	windowNameSize := header.HeaderSize - 100
	windowName := make([]byte, windowNameSize)
	n, err := r.Read(windowName)
	if err != nil {
		return nil, err
	}
	if n != int(windowNameSize) {
		return nil, fmt.Errorf("cannot read window name")
	}

	colorMaps := make([]xwdColor, header.NColors)
	if err := binary.Read(r, binary.BigEndian, colorMaps); err != nil {
		return nil, err
	}

	m := image.NewNRGBA(image.Rect(0, 0, int(header.PixmapWidth), int(header.PixmapHeight)))

	readByteN := int(header.BitsPerPixel) / 8
	pad := make([]byte, header.PixmapWidth*(header.BitmapPad)/8)
	buf := make([]byte, readByteN)
	for y := 0; y < int(header.PixmapHeight); y++ {
		for x := 0; x < int(header.PixmapWidth); x++ {
			n, err := io.ReadFull(r, buf)
			if err != nil {
				if err == io.EOF {
					return m, nil
				} else {
					return nil, err
				}
			}
			if n != readByteN {
				return nil, fmt.Errorf("invalid read size")
			}
			m.Set(x, y, color.NRGBA{
				// TODO: use mask
				R: buf[2],
				G: buf[1],
				B: buf[0],
				A: 0xFF,
			})
		}
		// skip padding
		n, err = io.ReadFull(r, pad)
		if err != nil {
			if err == io.EOF {
				return m, nil
			} else {
				return nil, err
			}
		}
		if n != int(header.PixmapWidth) {
			return nil, fmt.Errorf("invalid read size")
		}
	}
	return m, nil
}

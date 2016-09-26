package screenshot

// #cgo CFLAGS:   -pipe -O2 -mmacosx-version-min=10.7 -fPIC
// #cgo CXXFLAGS: -pipe -O2 -mmacosx-version-min=10.7 -fPIC -std=c++11 -stdlib=libc++
// #cgo LDFLAGS: -framework CoreGraphics -framework CoreFoundation
// #include <stdint.h>
// #include "screenshot.h"
import "C"

import (
	"image"
	"errors"
	"unsafe"
	"image/color"
)

// Capture screen.
// x and y represent distance from the upper-left corner of main display.
// Y-axis is downward direction. This means coordinates system is similar to Windows OS.
func Capture(x, y, width, height int) (*image.RGBA, error) {
	var ptr *C.uint32_t
	ptr = C.Capture(C.int(x), C.int(y), C.int(width), C.int(height))
	if ptr == nil {
		return nil, errors.New("failed capturing display")
	}
	defer C.Dispose(ptr)

	length := width * height
	buffer := (*[1 << 30]C.uint32_t)(unsafe.Pointer(ptr))[:length:length]

	rect := image.Rect(0, 0, width, height)
	img := image.NewRGBA(rect)

	var col color.RGBA
	i := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := buffer[i]

			col.B = (uint8)(0xff & (c >> 24))
			col.G = (uint8)(0xff & (c >> 16))
			col.R = (uint8)(0xff & (c >> 8))
			col.A = 255

			img.Set(x, y, col)
			i++
		}
	}

	return img, nil
}

// Get the number of active displays.
func NumActiveDisplays() int {
	return int(C.NumActiveDisplays())
}

// Get the bounds of display_index'th display.
// The main display is display_index = 0.
func GetDisplayBounds(display_index int) image.Rectangle {
	var x, y, w, h C.int
	x = 0
	y = 0
	w = 0
	h = 0
	C.GetDisplayBounds(C.int(display_index), &x, &y, &w, &h)
	return image.Rect(int(x), int(y), int(x + w), int(y + h))
}

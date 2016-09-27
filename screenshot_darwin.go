package screenshot

// #cgo CFLAGS:   -pipe -O2 -mmacosx-version-min=10.7 -fPIC
// #cgo CXXFLAGS: -pipe -O2 -mmacosx-version-min=10.7 -fPIC -std=c++11 -stdlib=libc++
// #cgo LDFLAGS: -framework CoreGraphics -framework CoreFoundation
// #include <stdint.h>
// #include "screenshot_darwin.h"
import "C"

import (
	"image"
	"errors"
	"unsafe"
)

func Capture(x, y, width, height int) (*image.RGBA, error) {
	rect := image.Rect(0, 0, width, height)
	img := image.NewRGBA(rect)

	ret := C.Capture(C.int(x), C.int(y),
		C.int(width), C.int(height),
		(*C.uint32_t)(unsafe.Pointer(&img.Pix[0])), C.int(img.Stride))
	if ret == 0 {
		return img, nil
	} else {
		return nil, errors.New("")
	}
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

//go:build !cgo

package screenshot

import (
	"errors"
	"image"
)

var ErrNoCGO = errors.New("requires CGO on OSX")

// Capture returns screen capture of specified desktop region.
// x and y represent distance from the upper-left corner of main display.
// Y-axis is downward direction. This means coordinates system is similar to Windows OS.
func Capture(x, y, width, height int) (*image.RGBA, error) {
	return nil, ErrNoCGO
}

// NumActiveDisplays returns the number of active displays.
func NumActiveDisplays() int {
	return 0
}

// GetDisplayBounds returns the bounds of displayIndex'th display.
// The main display is displayIndex = 0.
func GetDisplayBounds(displayIndex int) image.Rectangle {
	return image.Rectangle{}
}

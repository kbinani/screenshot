//go:build !s390x && !ppc64le && !darwin && !windows && (linux || freebsd || openbsd || netbsd)

package screenshot

import (
	"github.com/kbinani/screenshot/internal"
	"image"
)

// Capture returns screen capture of specified desktop region.
// x and y represent distance from the upper-left corner of primary display.
// Y-axis is downward direction. This means coordinates system is similar to Windows OS.
func Capture(x, y, width, height int) (*image.RGBA, error) {
	return internal.Capture(x, y, width, height)
}

// NumActiveDisplays returns the number of active displays.
func NumActiveDisplays() int {
	return internal.NumActiveDisplays()
}

// GetDisplayBounds returns the bounds of displayIndex'th display.
// The main display is displayIndex = 0.
func GetDisplayBounds(displayIndex int) image.Rectangle {
	return internal.GetDisplayBounds(displayIndex)
}

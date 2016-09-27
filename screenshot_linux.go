package screenshot

import (
	"image"
	"errors"
)

// Capture screen.
// x and y represent distance from the upper-left corner of main display.
// Y-axis is downward direction. This means coordinates system is similar to Windows OS.
func Capture(x, y, width, height int) (*image.RGBA, error) {
	return nil, errors.New("screenshot.Capture not supported for linux")
}

// Get the number of active displays.
func NumActiveDisplays() int {
	return 0
}

// Get the bounds of displayIndex'th display.
// The main display is displayIndex = 0.
func GetDisplayBounds(displayIndex int) image.Rectangle {
	return image.ZR
}

// Captures screen-shot image as image.RGBA.
// Mac(with cgo), and Windows(without cgo) are supported.
package screenshot

import (
	"image"
)

func CaptureDisplay(display_index int) (*image.RGBA, error) {
	rect := GetDisplayBounds(display_index)
	return Capture(rect.Min.X, rect.Min.Y, rect.Dx(), rect.Dy())
}

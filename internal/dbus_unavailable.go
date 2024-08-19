//go:build freebsd

package internal

import (
	"image"
)

func Capture(x, y, width, height int) (img *image.RGBA, e error) {
	return captureXinerama(x, y, width, height)
}

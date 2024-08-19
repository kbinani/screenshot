//go:build !s390x && !ppc64le && !darwin && !windows && (linux || openbsd || netbsd)

package internal

import (
	"image"
	"os"
)

func Capture(x, y, width, height int) (img *image.RGBA, e error) {
	sessionType := os.Getenv("XDG_SESSION_TYPE")
	if sessionType == "wayland" {
		return captureDbus(x, y, width, height)
	} else {
		return captureXinerama(x, y, width, height)
	}
}

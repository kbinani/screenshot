// Package screenshot captures screen-shot image as image.RGBA.
// Mac, Windows, Linux, FreeBSD, OpenBSD, NetBSD, and Solaris are supported.
package screenshot

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"strconv"
)

// CaptureDisplay captures whole region of displayIndex'th display.
func CaptureDisplay(displayIndex int) (*image.RGBA, error) {
	rect := GetDisplayBounds(displayIndex)
	return CaptureRect(rect)
}

// CaptureRect captures specified region of desktop.
func CaptureRect(rect image.Rectangle) (*image.RGBA, error) {
	return Capture(rect.Min.X, rect.Min.Y, rect.Dx(), rect.Dy())
}

//Encodes and Saves the image in PNG format 
func SavePng(img *image.RGBA, path string) (error) {
	file, err := os.Create(path+".png")
	if err != nil {
		file.Close()
		return err
	}
	png.Encode(file, img)
        file.Close()
	return nil
}

//Encodes and Saves the image in JPEG format
func SaveJpeg(img *image.RGBA, path string, imgQuality int) (error) {
	if imgQuality > 100 {
		return errors.New("ImageQuality must be smaller than 100. But provided "+strconv.Itoa(imgQuality)+".Which is greater than 100");
	}

	file, err := os.Create(path+".jpg")
	if err != nil {
		file.Close()
		return err
	}
	jpeg.Encode(file, img, &jpeg.Options{imgQuality})
	file.Close()
        return nil	
}

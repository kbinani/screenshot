screenshot
==========

[![GoDoc](https://godoc.org/github.com/kbinani/screenshot?status.svg)](https://godoc.org/github.com/kbinani/screenshot)

* Go library to capture desktop screen.
* Support Windows and Mac environment.
* `cgo` free for Windows.
* Multiple display supported.

example
=======
```go
package main

import (
	"github.com/kbinani/screenshot"
	"image/png"
	"os"
	"fmt"
)

func main() {
	n := screenshot.NumActiveDisplays()

	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)
		fmt.Printf("#%d : %v\n", i, bounds)

		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			panic(err)
		}
		fileName := fmt.Sprintf("%d_%dx%d.png", i, bounds.Dx(), bounds.Dy())
		file, _ := os.Create(fileName)
		defer file.Close()
		png.Encode(file, img)
	}
}
```

screen coordinate
=================
Y-axis is downward direction in this library. This means coordinate system is similar to Windows OS. The origin of coordinate is upper-left corner of main display.

license
=======

MIT Licence

author
======

kbinani

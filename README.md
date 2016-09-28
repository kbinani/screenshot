screenshot
==========

[![](https://img.shields.io/badge/godoc-reference-5566aa.svg)](https://godoc.org/github.com/kbinani/screenshot)
![](https://img.shields.io/badge/license-MIT-green.svg?style=flat)

* Go library to capture desktop screen.
* Support Windows, Mac, Linux, and FreeBSD environment.
* `cgo` free for Windows, Linux, and FreeBSD.
* Multiple display supported.

example
=======

* sample program `main.go`

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

			img, err := screenshot.CaptureRect(bounds)
			if err != nil {
				panic(err)
			}
			fileName := fmt.Sprintf("%d_%dx%d.png", i, bounds.Dx(), bounds.Dy())
			file, _ := os.Create(fileName)
			defer file.Close()
			png.Encode(file, img)

			fmt.Printf("#%d : %v \"%s\"\n", i, bounds, fileName)
		}
	}
	```

* output example
	
	```bash
	$ go run main.go
	#0 : (0,0)-(1280,800) "0_1280x800.png"
	#1 : (-293,-1440)-(2267,0) "1_2560x1440.png"
	#2 : (-1373,-1812)-(-293,108) "2_1080x1920.png"
	$ ls -1
	0_1280x800.png
	1_2560x1440.png
	2_1080x1920.png
	main.go
	```

coordinate
=================
Y-axis is downward direction in this library. The origin of coordinate is upper-left corner of main display. This means coordinate system is similar to Windows OS

license
=======

MIT Licence

author
======

kbinani

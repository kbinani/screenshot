//go:build cgo

package screenshot

/*
#cgo LDFLAGS: -framework CoreGraphics -framework CoreFoundation -framework ScreenCaptureKit
#cgo CXXFLAGS: -std=c++17
#include <CoreGraphics/CoreGraphics.h>
#include <ScreenCaptureKit/ScreenCaptureKit.h>

CGImageRef CaptureScreen() {
    NSError *error = nil;
    
    // Create the capture session
    SCStreamConfiguration *config = [[SCStreamConfiguration alloc] init];
    SCStream *stream = [[SCStream alloc] initWithConfiguration:config error:&error];
    
    if (error) {
        NSLog(@"Failed to create stream: %@", error);
        return NULL;
    }
    
    // Start the capture session
    [stream startCapture:&error];
    
    if (error) {
        NSLog(@"Failed to start capture: %@", error);
        return NULL;
    }
    
    // Get the screen image
    CGImageRef image = [stream copyImageForDisplayID:kCGDirectMainDisplay];
    
    if (image == NULL) {
        NSLog(@"Failed to capture image");
    }
    
    [stream stopCapture];
    
    return image;
}

void CompatCGImageRelease(void* image) {
    CGImageRelease((CGImageRef)image);
}
*/
import "C"

import (
    "errors"
    "image"
    "unsafe"

    "github.com/kbinani/screenshot/internal/util"
)

func Capture(x, y, width, height int) (*image.RGBA, error) {
    if width <= 0 || height <= 0 {
        return nil, errors.New("width or height should be > 0")
    }

    rect := image.Rect(0, 0, width, height)
    img, err := util.CreateImage(rect)
    if err != nil {
        return nil, err
    }

    imgRef := C.CaptureScreen()
    defer C.CompatCGImageRelease(imgRef)
    if imgRef == nil {
        return nil, errors.New("cannot capture display")
    }

    widthInt := int(C.CGImageGetWidth(imgRef))
    heightInt := int(C.CGImageGetHeight(imgRef))
    if widthInt != width || heightInt != height {
        return nil, errors.New("captured image size mismatch")
    }

    // Copy image data
    colorSpace := createColorspace()
    ctx := createBitmapContext(width, height, (*C.uint32_t)(unsafe.Pointer(&img.Pix[0])), img.Stride)
    if ctx == 0 {
        return nil, errors.New("cannot create bitmap context")
    }
    defer C.CGColorSpaceRelease(colorSpace)

    rect := C.CGRectMake(0, 0, C.CGFloat(width), C.CGFloat(height))
    C.CGContextDrawImage(ctx, rect, imgRef)

    i := 0
    for iy := 0; iy < height; iy++ {
        j := i
        for ix := 0; ix < width; ix++ {
            img.Pix[j], img.Pix[j+1], img.Pix[j+2], img.Pix[j+3] = img.Pix[j+1], img.Pix[j+2], img.Pix[j+3], 255
            j += 4
        }
        i += img.Stride
    }

    return img, nil
}

func NumActiveDisplays() int {
    var count C.uint32_t = 0
    if C.CGGetActiveDisplayList(0, nil, &count) == C.kCGErrorSuccess {
        return int(count)
    } else {
        return 0
    }
}

func GetDisplayBounds(displayIndex int) image.Rectangle {
    id := getDisplayId(displayIndex)
    main := C.CGMainDisplayID()

    var rect image.Rectangle

    bounds := getCoreGraphicsCoordinateOfDisplay(id)
    rect.Min.X = int(bounds.origin.x)
    if main == id {
        rect.Min.Y = 0
    } else {
        mainBounds := getCoreGraphicsCoordinateOfDisplay(main)
        mainHeight := mainBounds.size.height
        rect.Min.Y = int(mainHeight - (bounds.origin.y + bounds.size.height))
    }
    rect.Max.X = rect.Min.X + int(bounds.size.width)
    rect.Max.Y = rect.Min.Y + int(bounds.size.height)

    return rect
}

func getDisplayId(displayIndex int) C.CGDirectDisplayID {
    main := C.CGMainDisplayID()
    if displayIndex == 0 {
        return main
    } else {
        n := NumActiveDisplays()
        ids := make([]C.CGDirectDisplayID, n)
        if C.CGGetActiveDisplayList(C.uint32_t(n), (*C.CGDirectDisplayID)(unsafe.Pointer(&ids[0])), nil) != C.kCGErrorSuccess {
            return 0
        }
        index := 0
        for i := 0; i < n; i++ {
            if ids[i] == main {
                continue
            }
            index++
            if index == displayIndex {
                return ids[i]
            }
        }
    }

    return 0
}

func getCoreGraphicsCoordinateOfDisplay(id C.CGDirectDisplayID) C.CGRect {
    main := C.CGDisplayBounds(C.CGMainDisplayID())
    r := C.CGDisplayBounds(id)
    return C.CGRectMake(r.origin.x, -r.origin.y-r.size.height+main.size.height,
        r.size.width, r.size.height)
}

func createBitmapContext(width int, height int, data *C.uint32_t, bytesPerRow int) C.CGContextRef {
    colorSpace := createColorspace()
    if colorSpace == 0 {
        return 0
    }
    defer C.CGColorSpaceRelease(colorSpace)

    return C.CGBitmapContextCreate(unsafe.Pointer(data),
        C.size_t(width),
        C.size_t(height),
        8, // bits per component
        C.size_t(bytesPerRow),
        colorSpace,
        C.kCGImageAlphaNoneSkipFirst)
}

func createColorspace() C.CGColorSpaceRef {
    return C.CGColorSpaceCreateWithName(C.kCGColorSpaceSRGB)
}

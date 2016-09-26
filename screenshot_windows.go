package screenshot

import (
	"image"
	"image/color"
	"unsafe"
	"syscall"
	win "github.com/lxn/win"
)

var (
	libUser32, _ = syscall.LoadLibrary("user32.dll")
	funcGetDesktopWindow, _ = syscall.GetProcAddress(syscall.Handle(libUser32), "GetDesktopWindow")
	funcEnumDisplayMonitors, _ = syscall.GetProcAddress(syscall.Handle(libUser32), "EnumDisplayMonitors")
)

func Capture(x, y, width, height int) (*image.RGBA, error) {
	hwnd := getDesktopWindow()
	hdc := win.GetDC(hwnd)
	data := make([]uint32, width * height)
	memory_device := win.CreateCompatibleDC(hdc)
	bitmap := win.CreateCompatibleBitmap(hdc, int32(width), int32(height))
	var header win.BITMAPINFOHEADER
	header.BiSize = uint32(unsafe.Sizeof(header))
	header.BiPlanes = 1
	header.BiBitCount = 32
	header.BiWidth = int32(width)
	header.BiHeight = int32(-height)
	header.BiCompression = win.BI_RGB
	header.BiSizeImage = 0
	win.SelectObject(memory_device, win.HGDIOBJ(bitmap))
	win.BitBlt(memory_device, 0, 0, int32(width), int32(height), hdc, int32(x), int32(y), win.SRCCOPY)
	win.GetDIBits(hdc, bitmap, 0, uint32(height), (*byte)(unsafe.Pointer(&data[0])), (*win.BITMAPINFO)(unsafe.Pointer(&header)), win.DIB_RGB_COLORS)
	win.ReleaseDC(hwnd, hdc)
	win.DeleteDC(memory_device)

	rect := image.Rect(0, 0, width, height)
	img := image.NewRGBA(rect)

	var col color.RGBA
	i := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := data[i]

			col.A = 255
			col.R = (uint8)(0xff & (c >> 16))
			col.G = (uint8)(0xff & (c >> 8))
			col.B = (uint8)(0xff & (c))

			img.Set(x, y, col)
			i++
		}
	}

	return img, nil
}

func NumActiveDisplays() int {
	var count int = 0
	enumDisplayMonitors(win.HDC(0), nil, syscall.NewCallback(countupMonitorCallback), uintptr(unsafe.Pointer(&count)))
	return count
}

func GetDisplayBounds(display_index int) image.Rectangle {
	var ctx getMonitorBoundsContext
	ctx.Index = display_index
	ctx.Count = 0
	enumDisplayMonitors(win.HDC(0), nil, syscall.NewCallback(getMonitorBoundsCallback), uintptr(unsafe.Pointer(&ctx)))
	return image.Rect(
		int(ctx.Rect.Left), int(ctx.Rect.Top),
		int(ctx.Rect.Right - ctx.Rect.Left), int(ctx.Rect.Bottom - ctx.Rect.Top))
}

func getDesktopWindow() win.HWND {
	ret, _, _ := syscall.Syscall(funcGetDesktopWindow, 0, 0, 0, 0)
	return win.HWND(ret)
}

func enumDisplayMonitors(hdc win.HDC, lprcClip *win.RECT, lpfnEnum uintptr, dwData uintptr) bool {
	ret, _, _ := syscall.Syscall6(funcEnumDisplayMonitors, 4,
		uintptr(hdc),
		uintptr(unsafe.Pointer(lprcClip)),
		lpfnEnum,
		dwData,
		0,
		0)
	return int(ret) != 0
}

func countupMonitorCallback(hMonitor win.HMONITOR, hdcMonitor win.HDC, lprcMonitor *win.RECT, dwData uintptr) uintptr {
	var count *int
	count = (*int)(unsafe.Pointer(dwData))
	*count = *count + 1
	return uintptr(1)
}

type getMonitorBoundsContext struct {
	Index int
	Rect win.RECT
	Count int
}

func getMonitorBoundsCallback(hMonitor win.HMONITOR, hdcMonitor win.HDC, lprcMonitor *win.RECT, dwData uintptr) uintptr {
	var ctx *getMonitorBoundsContext
	ctx = (*getMonitorBoundsContext)(unsafe.Pointer(dwData))
	if ctx.Count == ctx.Index {
		ctx.Rect = *lprcMonitor
		return uintptr(0)
	} else {
		ctx.Count = ctx.Count + 1
		return uintptr(1)
	}
}

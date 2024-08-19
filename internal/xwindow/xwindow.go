package xwindow

import (
	"fmt"
	"github.com/gen2brain/shm"
	"github.com/godbus/dbus/v5"
	"github.com/jezek/xgb"
	mshm "github.com/jezek/xgb/shm"
	"github.com/jezek/xgb/xinerama"
	"github.com/jezek/xgb/xproto"
	"github.com/kbinani/screenshot/internal/util"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"net/url"
	"os"
	"sync/atomic"
)

var gCounter uint64 = 0

func Capture(x, y, width, height int) (img *image.RGBA, e error) {
	sessionType := os.Getenv("XDG_SESSION_TYPE")
	if sessionType == "wayland" {
		return captureDbus(x, y, width, height)
	} else {
		return captureXinerama(x, y, width, height)
	}
}

func captureDbus(x, y, width, height int) (img *image.RGBA, e error) {
	c, err := dbus.ConnectSessionBus()
	if err != nil {
		return nil, fmt.Errorf("dbus.SessionBus() failed: %v", err)
	}
	defer func(c *dbus.Conn) {
		err := c.Close()
		if err != nil {
			e = err
		}
	}(c)
	token := atomic.AddUint64(&gCounter, 1)
	options := map[string]dbus.Variant{
		"modal":        dbus.MakeVariant(false),
		"interactive":  dbus.MakeVariant(false),
		"handle_token": dbus.MakeVariant(token),
	}
	obj := c.Object("org.freedesktop.portal.Desktop", dbus.ObjectPath("/org/freedesktop/portal/desktop"))
	call := obj.Call("org.freedesktop.portal.Screenshot.Screenshot", 0, "", options)
	var path dbus.ObjectPath
	err = call.Store(&path)
	if err != nil {
		return nil, fmt.Errorf("dbus.Store() failed: %v", err)
	}
	ch := make(chan *dbus.Message)
	c.Eavesdrop(ch)
	for msg := range ch {
		o, ok := msg.Headers[dbus.FieldPath]
		if !ok {
			continue
		}
		s, ok := o.Value().(dbus.ObjectPath)
		if !ok {
			return nil, fmt.Errorf("dbus.FieldPath value does't have ObjectPath type")
		}
		if s != path {
			continue
		}
		for _, body := range msg.Body {
			v, ok := body.(map[string]dbus.Variant)
			if !ok {
				continue
			}
			uri, ok := v["uri"]
			if !ok {
				continue
			}
			path, ok := uri.Value().(string)
			if !ok {
				return nil, fmt.Errorf("uri is not a string")
			}
			fpath, err := url.Parse(path)
			if err != nil {
				return nil, fmt.Errorf("url.Parse(%v) failed: %v", path, err)
			}
			if fpath.Scheme != "file" {
				return nil, fmt.Errorf("uri is not a file path")
			}
			file, err := os.Open(fpath.Path)
			if err != nil {
				return nil, fmt.Errorf("os.Open(%s) failed: %v", path, err)
			}
			defer func(file *os.File) {
				_ = file.Close()
				_ = os.Remove(fpath.Path)
			}(file)
			img, err := png.Decode(file)
			if err != nil {
				return nil, fmt.Errorf("png.Decode(%s) failed: %v", path, err)
			}
			canvas, err := util.CreateImage(image.Rect(0, 0, width, height))
			if err != nil {
				return nil, fmt.Errorf("util.CreateImage(%v) failed: %v", path, err)
			}
			draw.Draw(canvas, image.Rect(0, 0, width, height), img, image.Point{x, y}, draw.Src)
			return canvas, e
		}
	}
	return nil, fmt.Errorf("dbus.Message doesn't contain uri")
}

func captureXinerama(x, y, width, height int) (img *image.RGBA, e error) {
	defer func() {
		err := recover()
		if err != nil {
			img = nil
			e = fmt.Errorf("%v", err)
		}
	}()
	c, err := xgb.NewConn()
	if err != nil {
		return nil, err
	}
	defer c.Close()

	err = xinerama.Init(c)
	if err != nil {
		return nil, err
	}

	reply, err := xinerama.QueryScreens(c).Reply()
	if err != nil {
		return nil, err
	}

	primary := reply.ScreenInfo[0]
	x0 := int(primary.XOrg)
	y0 := int(primary.YOrg)

	useShm := true
	err = mshm.Init(c)
	if err != nil {
		useShm = false
	}

	screen := xproto.Setup(c).DefaultScreen(c)
	wholeScreenBounds := image.Rect(0, 0, int(screen.WidthInPixels), int(screen.HeightInPixels))
	targetBounds := image.Rect(x+x0, y+y0, x+x0+width, y+y0+height)
	intersect := wholeScreenBounds.Intersect(targetBounds)

	rect := image.Rect(0, 0, width, height)
	img, err = util.CreateImage(rect)
	if err != nil {
		return nil, err
	}

	// Paint with opaque black
	index := 0
	for iy := 0; iy < height; iy++ {
		j := index
		for ix := 0; ix < width; ix++ {
			img.Pix[j+3] = 255
			j += 4
		}
		index += img.Stride
	}

	if !intersect.Empty() {
		var data []byte

		if useShm {
			shmSize := intersect.Dx() * intersect.Dy() * 4
			shmId, err := shm.Get(shm.IPC_PRIVATE, shmSize, shm.IPC_CREAT|0777)
			if err != nil {
				return nil, err
			}

			seg, err := mshm.NewSegId(c)
			if err != nil {
				return nil, err
			}

			data, err = shm.At(shmId, 0, 0)
			if err != nil {
				return nil, err
			}

			mshm.Attach(c, seg, uint32(shmId), false)

			defer mshm.Detach(c, seg)
			defer func() {
				_ = shm.Rm(shmId)
			}()
			defer func() {
				_ = shm.Dt(data)
			}()

			_, err = mshm.GetImage(c, xproto.Drawable(screen.Root),
				int16(intersect.Min.X), int16(intersect.Min.Y),
				uint16(intersect.Dx()), uint16(intersect.Dy()), 0xffffffff,
				byte(xproto.ImageFormatZPixmap), seg, 0).Reply()
			if err != nil {
				return nil, err
			}
		} else {
			xImg, err := xproto.GetImage(c, xproto.ImageFormatZPixmap, xproto.Drawable(screen.Root),
				int16(intersect.Min.X), int16(intersect.Min.Y),
				uint16(intersect.Dx()), uint16(intersect.Dy()), 0xffffffff).Reply()
			if err != nil {
				return nil, err
			}

			data = xImg.Data
		}

		// BitBlt by hand
		offset := 0
		for iy := intersect.Min.Y; iy < intersect.Max.Y; iy++ {
			for ix := intersect.Min.X; ix < intersect.Max.X; ix++ {
				r := data[offset+2]
				g := data[offset+1]
				b := data[offset]
				img.SetRGBA(ix-(x+x0), iy-(y+y0), color.RGBA{r, g, b, 255})
				offset += 4
			}
		}
	}

	return img, e
}

func NumActiveDisplays() (num int) {
	defer func() {
		e := recover()
		if e != nil {
			num = 0
		}
	}()

	c, err := xgb.NewConn()
	if err != nil {
		return 0
	}
	defer c.Close()

	err = xinerama.Init(c)
	if err != nil {
		return 0
	}

	reply, err := xinerama.QueryScreens(c).Reply()
	if err != nil {
		return 0
	}

	num = int(reply.Number)
	return num
}

func GetDisplayBounds(displayIndex int) (rect image.Rectangle) {
	defer func() {
		e := recover()
		if e != nil {
			rect = image.Rectangle{}
		}
	}()

	c, err := xgb.NewConn()
	if err != nil {
		return image.Rectangle{}
	}
	defer c.Close()

	err = xinerama.Init(c)
	if err != nil {
		return image.Rectangle{}
	}

	reply, err := xinerama.QueryScreens(c).Reply()
	if err != nil {
		return image.Rectangle{}
	}

	if displayIndex >= int(reply.Number) {
		return image.Rectangle{}
	}

	primary := reply.ScreenInfo[0]
	x0 := int(primary.XOrg)
	y0 := int(primary.YOrg)

	screen := reply.ScreenInfo[displayIndex]
	x := int(screen.XOrg) - x0
	y := int(screen.YOrg) - y0
	w := int(screen.Width)
	h := int(screen.Height)
	rect = image.Rect(x, y, x+w, y+h)
	return rect
}

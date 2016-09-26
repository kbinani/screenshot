#include "screenshot.h"
#include <CoreGraphics/CoreGraphics.h>
#include <utility> // std::swap

template <class T>
class scoped_cfref
{
public:
    scoped_cfref(T o)
        : o_(o)
    {}

    ~scoped_cfref()
    {
        if (o_) {
            release();
        }
    }

    scoped_cfref(scoped_cfref const& other) = delete;
    scoped_cfref& operator = (scoped_cfref const& other) = delete;

    scoped_cfref(scoped_cfref&& other)
    {
        std::swap(o_, other.o_);
    }

    scoped_cfref& operator = (scoped_cfref&& other)
    {
        std::swap(o_, other.o_);
        return *this;
    }

    operator T() const { return o_; }

private:
    void release();

private:
    T o_ = NULL;
};

template <> inline void scoped_cfref<CGColorSpaceRef>::release() { CGColorSpaceRelease(o_); }
template <> inline void scoped_cfref<CGImageRef>::release() { CGImageRelease(o_); }
template <> inline void scoped_cfref<CGContextRef>::release() { CGContextRelease(o_); }

#ifdef __cplusplus
extern "C" {
#endif

static uint32_t* createImage(int width, int height)
{
    uint32_t* data = (uint32_t*)malloc(width * height * sizeof(uint32_t));
    memset(data, 0, sizeof(uint32_t) * width * height);
    return data;
}

static void disposeImage(uint32_t* data)
{
    if (data) {
        free(data);
    }
}

static CGColorSpaceRef macCreateColorspace()
{
    return CGColorSpaceCreateWithName(kCGColorSpaceSRGB);
}

static CGContextRef macCreateBitmapContext(size_t width, size_t height)
{
    scoped_cfref<CGColorSpaceRef> colorSpace = macCreateColorspace();
    if (colorSpace == NULL) {
        return NULL;
    }

    int const bits_per_component = 8;
    CGContextRef context = CGBitmapContextCreate(NULL,
                                                 width,
                                                 height,
                                                 bits_per_component,
                                                 0,
                                                 colorSpace,
                                                 kCGImageAlphaNoneSkipFirst);

    return context;
}

static CGDirectDisplayID macGetDisplayId(int display_index)
{
    CGDirectDisplayID main = CGMainDisplayID();
    if (display_index == 0) {
        return main;
    } else {
        CGDirectDisplayID ids[128];
        uint32_t count = 0;
        CGGetActiveDisplayList(sizeof(ids) / sizeof(CGDirectDisplayID), ids, &count);
        int index = 0;
        for (uint32_t i = 0; i < count; ++i) {
            if (ids[i] == main) {
                continue;
            }
            ++index;
            if (index == display_index) {
                return ids[i];
            }
        }
    }

    return 0;
}

static CGPoint macGetCoreGraphicsCoordinateFromWindowsCoordinate(CGPoint p, CGRect mainDisplayBounds)
{
    return CGPointMake(p.x, mainDisplayBounds.size.height - p.y);
}

static CGRect macGetCoreGraphicsCoordinateOfDisplay(CGDirectDisplayID id)
{
    CGRect main = CGDisplayBounds(CGMainDisplayID());
    CGRect r = CGDisplayBounds(id);
    return CGRectMake(r.origin.x, -r.origin.y - r.size.height + main.size.height,
                      r.size.width, r.size.height);
}

uint32_t* Capture(int x, int y, int width, int height)
{
    if (width <= 0 || height <= 0) {
        return NULL;
    }

    CGRect mainDisplayBounds = macGetCoreGraphicsCoordinateOfDisplay(CGMainDisplayID());

    CGPoint bottomLeftWin = CGPointMake(x, y + height);
    CGPoint bottomLeft = macGetCoreGraphicsCoordinateFromWindowsCoordinate(bottomLeftWin, mainDisplayBounds);
    CGRect captureBounds = CGRectMake(bottomLeft.x, bottomLeft.y, width, height);

    CGDirectDisplayID ids[128];
    uint32_t count = 0;
    CGGetActiveDisplayList(sizeof(ids) / sizeof(CGDirectDisplayID), ids, &count);

    uint32_t* data = createImage(width, height);
    if (!data) {
        return NULL;
    }
    scoped_cfref<CGContextRef> cgctx = macCreateBitmapContext(width, height);
    if (!cgctx) {
        disposeImage(data);
        return NULL;
    }

    scoped_cfref<CGColorSpaceRef> colorSpace = macCreateColorspace();
    if (!colorSpace) {
        CFRelease(cgctx);
        disposeImage(data);
        return NULL;
    }

    for (uint32_t i = 0; i < count; ++i) {
        CGDirectDisplayID id = ids[i];
        CGRect bounds = macGetCoreGraphicsCoordinateOfDisplay(id);
        CGRect intersect = CGRectIntersection(bounds, captureBounds);
        if (CGRectIsNull(intersect)) {
            continue;
        }
        if (intersect.size.width <= 0 || intersect.size.height <= 0) {
            continue;
        }

        // CGDisplayCreateImageForRect potentially fail in case width/height is odd number.
        if ((int)intersect.size.width % 2 != 0) {
            intersect.size.width = (int)intersect.size.width + 1;
        }
        if ((int)intersect.size.height % 2 != 0) {
            intersect.size.height = (int)intersect.size.height + 1;
        }

        CGRect intersectDisplayLocal = CGRectMake(intersect.origin.x - bounds.origin.x,
                                                  bounds.origin.y + bounds.size.height - (intersect.origin.y + intersect.size.height),
                                                  intersect.size.width, intersect.size.height);
        scoped_cfref<CGImageRef> captured = CGDisplayCreateImageForRect(id, intersectDisplayLocal);
        if (!captured) {
            continue;
        }

        scoped_cfref<CGImageRef> image = CGImageCreateCopyWithColorSpace(captured, colorSpace);
        if (!image) {
            continue;
        }
        CGRect drawRect = CGRectMake(intersect.origin.x - captureBounds.origin.x, intersect.origin.y - captureBounds.origin.y,
                                     intersect.size.width, intersect.size.height);
        CGContextDrawImage(cgctx, drawRect, image);
    }

    uint8_t* ptr = (uint8_t*)CGBitmapContextGetData(cgctx);
    int const bytesPerRow = CGBitmapContextGetBytesPerRow(cgctx);
    int const bytesPerPixel = CGBitmapContextGetBitsPerPixel(cgctx) / 8;
    for (int iy = 0; iy < height; ++iy) {
        for (int ix = 0; ix < width; ++ix) {
            uint8_t* p = ptr + (iy * bytesPerRow + ix * bytesPerPixel);
            data[iy * width + ix] = *(uint32_t*)p;
        }
    }

    return data;
}

uint32_t NumActiveDisplays()
{
    CGDirectDisplayID ids[128];
    uint32_t count = 0;
    CGGetActiveDisplayList(sizeof(ids) / sizeof(CGDirectDisplayID), ids, &count);
    return count;
}

void GetDisplayBounds(int display_index, int* x, int* y, int* width, int* height)
{
    CGDirectDisplayID id = macGetDisplayId(display_index);
    CGRect bounds = macGetCoreGraphicsCoordinateOfDisplay(id);
    if (x) {
        *x = (int)bounds.origin.x;
    }
    if (y) {
        CGDirectDisplayID main = CGMainDisplayID();
        if (main == id) {
            *y = 0;
        } else {
            CGRect mainBounds = macGetCoreGraphicsCoordinateOfDisplay(main);
            CGFloat mainHeight = mainBounds.size.height;
            *y = (int)(mainHeight - (bounds.origin.y + bounds.size.height));
        }
    }
    if (width) {
        *width = (int)bounds.size.width;
    }
    if (height) {
        *height = (int)bounds.size.height;
    }
}

void Dispose(uint32_t* data)
{
    disposeImage(data);
}

#ifdef __cplusplus
} // extern "C"
#endif

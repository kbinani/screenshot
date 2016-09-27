#pragma once
#include <stdint.h>
#include <stdlib.h>

#ifdef __cplusplus
extern "C" {
#endif

int Capture(int x, int y, int width, int height, uint32_t* dest, int bytesPerRow);
uint32_t NumActiveDisplays();
void GetDisplayBounds(int displayIndex, int* x, int* y, int* width, int* height);

#ifdef __cplusplus
} // extern "C"
#endif

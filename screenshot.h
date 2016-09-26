#pragma once
#include <stdint.h>
#include <stdlib.h>

#ifdef __cplusplus
extern "C" {
#endif

uint32_t* Capture(int x, int y, int width, int height);
void Dispose(uint32_t* data);
uint32_t NumActiveDisplays();
void GetDisplayBounds(int display_index, int* x, int* y, int* width, int* height);

#ifdef __cplusplus
} // extern "C"
#endif

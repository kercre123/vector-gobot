#ifndef __LIBJPEG_INTERFACE_H
#define __LIBJPEG_INTERFACE_H

#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>

#ifdef __cplusplus
extern "C" {
#endif

void encodeToJPEG(unsigned char* yuvData, int width, int height, int quality, unsigned char** jpegBuf, unsigned long* jpegSize);
void GetFrameAsJPEGDownSampled(const uint8_t* rawData, int width, int height, int quality, uint8_t** jpegBuf, unsigned long* jpegSize);


#ifdef __cplusplus
}
#endif

#endif // __LIBJPEG_INTERFACE_H

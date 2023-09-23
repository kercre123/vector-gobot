#include <cstdio>
#include <stdlib.h>
#include <turbojpeg.h>
#include "libjpeg_interface.h"
#include <vector>
#include <iostream>
#include <chrono>
#include <omp.h>

void encodeToJPEG(unsigned char* yuvData, int width, int height, int quality, unsigned char** jpegBuf, unsigned long* jpegSize) {
    tjhandle tjInstance = tjInitCompress();
    if (!tjInstance) {
        fprintf(stderr, "Error initializing TurboJPEG: %s\n", tjGetErrorStr());
        exit(1);
    }

    if (width % 2 != 0 || height % 2 != 0) {
        fprintf(stderr, "Width and height must be divisible by 2 for TJSAMP_420\n");
        return;
    }

    if (tjCompressFromYUV(tjInstance, yuvData, width, 1, height, TJSAMP_420, jpegBuf, jpegSize, quality, TJFLAG_FASTDCT) != 0) {
        fprintf(stderr, "Error with TJ compression: %s\n", tjGetErrorStr());
    }
    tjDestroy(tjInstance);
}

/*
this is now here because there seemed to be an issue with how Go was handling data types. putting it all into C++ allowed it to just work
*/
// works, so im not getting rid of it
// void GetFrameAsJPEGDownSampled(const uint8_t* rawData, int width, int height, int quality, uint8_t** jpegBuf, unsigned long* jpegSize) {

//     // unpack rgb10
//     std::vector<uint16_t> unpackedData(width * height);
//     for (int i = 0; i < width * height / 4 * 5; i += 5) {
//     int idx = (i / 5 * 4);
//         unpackedData[idx + 0] = (rawData[i + 0] << 2) | ((rawData[i + 4] >> 6) & 0x03);
//         unpackedData[idx + 1] = (rawData[i + 1] << 2) | ((rawData[i + 4] >> 4) & 0x03);
//         unpackedData[idx + 2] = (rawData[i + 2] << 2) | ((rawData[i + 4] >> 2) & 0x03);
//         unpackedData[idx + 3] = (rawData[i + 3] << 2) | ((rawData[i + 4] >> 0) & 0x03);
//     }

//     // debayer rggb
//     std::vector<uint8_t> rgbImage(width * height * 3 / 4);
//     for (int y = 0; y < height; y += 2) {
//       for (int x = 0; x < width; x += 2) {
//         uint16_t r = unpackedData[y * width + x];
//         uint16_t g = (unpackedData[y * width + x + 1] + unpackedData[(y + 1) * width + x]) / 2;
//         uint16_t b = unpackedData[(y + 1) * width + x + 1];
//         int idx = (y / 2 * width + x) / 2 * 3;
//         rgbImage[idx + 0] = r >> 2;
//         rgbImage[idx + 1] = g >> 2;
//         rgbImage[idx + 2] = b >> 2;
//       }
//     }

//     // compress
//     tjhandle tjInstance = tjInitCompress();
//     if (!tjInstance) {
//       std::cerr << "Error initializing TurboJPEG: " << tjGetErrorStr() << std::endl;
//       return;
//     }

//     // 4:4:4
//     if (tjCompress2(tjInstance, rgbImage.data(), width / 2, 0, height / 2, TJPF_RGB, jpegBuf, jpegSize, TJSAMP_444, quality, TJFLAG_FASTDCT) != 0) {
//       std::cerr << "Error with TJ compression: " << tjGetErrorStr() << std::endl;
//     }
//     tjDestroy(tjInstance);
//   }


// // success with NEON
// void GetFrameAsJPEGDownSampled(const uint8_t* rawData, int width, int height, int quality, uint8_t** jpegBuf, unsigned long* jpegSize) {
//     // Allocate a buffer for the debayered image
//     uint8_t* debayeredData = new uint8_t[width * height / 2 * 3];

//     for (int i = 0; i < height; i += 2) {
//         for (int j = 0; j < width; j += 8) {
//             // Load 5 bytes containing 4 packed 10-bit pixels
//             uint8x8_t packed_pixels1 = vld1_u8(rawData + ((i * width + j) / 4 * 5));
//             uint8x8_t packed_pixels2 = vld1_u8(rawData + ((i * width + j + 4) / 4 * 5));

//             // Unpack the 10-bit values
//             uint16x8_t unpacked_pixels1 = vshll_n_u8(packed_pixels1, 2);
//             uint16x8_t unpacked_pixels2 = vshll_n_u8(packed_pixels2, 2);

//             // Extracting R, G, B values
//             uint16_t r1 = vgetq_lane_u16(unpacked_pixels1, 0);
//             uint16_t g1 = (vgetq_lane_u16(unpacked_pixels1, 1) + vgetq_lane_u16(unpacked_pixels1, 2)) / 2;
//             uint16_t b1 = vgetq_lane_u16(unpacked_pixels1, 3);

//             uint16_t r2 = vgetq_lane_u16(unpacked_pixels2, 0);
//             uint16_t g2 = (vgetq_lane_u16(unpacked_pixels2, 1) + vgetq_lane_u16(unpacked_pixels2, 2)) / 2;
//             uint16_t b2 = vgetq_lane_u16(unpacked_pixels2, 3);

//             // Set the RGB values in the debayered image (downsampled to half width and height)
//             int idx = (i / 2 * (width / 2) + j / 2) * 3;
//             debayeredData[idx + 0] = uint8_t(r1 >> 2);
//             debayeredData[idx + 1] = uint8_t(g1 >> 2);
//             debayeredData[idx + 2] = uint8_t(b1 >> 2);

//             debayeredData[idx + 3] = uint8_t(r2 >> 2);
//             debayeredData[idx + 4] = uint8_t(g2 >> 2);
//             debayeredData[idx + 5] = uint8_t(b2 >> 2);
//         }
//     }

//     // Initialize TurboJPEG
//     tjhandle tjInstance = tjInitCompress();
//     if (!tjInstance) {
//         delete[] debayeredData;
//         return;
//     }

//     // Compress to JPEG
//     if (tjCompress2(tjInstance, debayeredData, width / 2, 0, height / 2, TJPF_RGB, jpegBuf, jpegSize, TJSAMP_444, quality, TJFLAG_FASTDCT) != 0) {
//         delete[] debayeredData;
//         tjDestroy(tjInstance);
//         return;
//     }

//     // Cleanup
//     delete[] debayeredData;
//     tjDestroy(tjInstance);
// }

// faster, but not fast enough!
// void GetFrameAsJPEGDownSampled(const uint8_t* rawData, int width, int height, int quality, uint8_t** jpegBuf, unsigned long* jpegSize) {
//     uint16_t* unpackedData = new uint16_t[width * height];
// #pragma omp parallel for
//     for (int i = 0; i < width * height / 4 * 5; i += 5) {
//         int idx = (i / 5 * 4);
//         unpackedData[idx + 0] = (rawData[i + 0] << 2) | ((rawData[i + 4] >> 6) & 0x03);
//         unpackedData[idx + 1] = (rawData[i + 1] << 2) | ((rawData[i + 4] >> 4) & 0x03);
//         unpackedData[idx + 2] = (rawData[i + 2] << 2) | ((rawData[i + 4] >> 2) & 0x03);
//         unpackedData[idx + 3] = (rawData[i + 3] << 2) | ((rawData[i + 4] >> 0) & 0x03);
//     }
//     uint8_t* rgbImage = new uint8_t[width * height * 3 / 4];
// #pragma omp parallel for
// for (int y = 0; y < height; y += 2) {
//     for (int x = 0; x < width; x += 2) { // Process 2x2 pixel block at a time
//         uint16_t r = unpackedData[y * width + x];
//         uint16_t g = (unpackedData[y * width + x + 1] + unpackedData[(y + 1) * width + x]) / 2;
//         uint16_t b = unpackedData[(y + 1) * width + x + 1];
//         int idx = (y / 2 * width + x) / 2 * 3;
//         rgbImage[idx + 0] = r >> 2;
//         rgbImage[idx + 1] = g >> 2;
//         rgbImage[idx + 2] = b >> 2;
//     }
// }
//     // compress
//     tjhandle tjInstance = tjInitCompress();
//     if (!tjInstance) {
//         std::cerr << "Error initializing TurboJPEG: " << tjGetErrorStr() << std::endl;
//         delete[] unpackedData;
//         delete[] rgbImage;
//         return;
//     }

//     // 4:4:4
//     if (tjCompress2(tjInstance, rgbImage, width / 2, 0, height / 2, TJPF_RGB, jpegBuf, jpegSize, TJSAMP_444, quality, TJFLAG_FASTDCT) != 0) {
//         std::cerr << "Error with TJ compression: " << tjGetErrorStr() << std::endl;
//     }

//     tjDestroy(tjInstance);

//     delete[] unpackedData;
//     delete[] rgbImage;
// }


// btw, this is here so that i have a specific place for testing compiler optimizations without affecting the other code
void GetFrameAsJPEGDownSampled(const uint8_t* rawData, int width, int height, int quality, uint8_t** jpegBuf, unsigned long* jpegSize) {
    uint8_t* rgbImage = new uint8_t[width * height * 3 / 4];

#pragma omp parallel for
    for (int y = 0; y < height; y += 2) {
        for (int x = 0; x < width; x += 2) { // Process 2x2 pixel block at a time
            int idx_raw = (y * width + x) / 4 * 5;
            int idx_rgb = (y / 2 * width + x) / 2 * 3;

            uint16_t r = (rawData[idx_raw + 0] << 2) | ((rawData[idx_raw + 4] >> 6) & 0x03);
            uint16_t g1 = (rawData[idx_raw + 1] << 2) | ((rawData[idx_raw + 4] >> 4) & 0x03);
            uint16_t g2 = (rawData[idx_raw + 2] << 2) | ((rawData[idx_raw + 4] >> 2) & 0x03);
            uint16_t b = (rawData[idx_raw + 3] << 2) | ((rawData[idx_raw + 4] >> 0) & 0x03);

            uint16_t g = (g1 + g2) >> 1;

            rgbImage[idx_rgb + 0] = r >> 2;
            rgbImage[idx_rgb + 1] = g >> 2;
            rgbImage[idx_rgb + 2] = b >> 2;
        }
    }

    // compress
    tjhandle tjInstance = tjInitCompress();
    if (!tjInstance) {
        std::cerr << "Error initializing TurboJPEG: " << tjGetErrorStr() << std::endl;
        delete[] rgbImage;
        return;
    }

    // 4:2:0
    if (tjCompress2(tjInstance, rgbImage, width / 2, 0, height / 2, TJPF_RGB, jpegBuf, jpegSize, TJSAMP_420, quality, TJFLAG_FASTDCT) != 0) {
        std::cerr << "Error with TJ compression: " << tjGetErrorStr() << std::endl;
    }

    tjDestroy(tjInstance);
    delete[] rgbImage;
}




// auto start = std::chrono::high_resolution_clock::now();
// auto stop = std::chrono::high_resolution_clock::now();
// auto duration = std::chrono::duration_cast<std::chrono::microseconds>(stop - start);

// std::cout << "Time taken: " << duration.count() << " microseconds" << std::endl;





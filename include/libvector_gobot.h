#ifndef __LIBVECTOR_GOBOT_H
#define __LIBVECTOR_GOBOT_H

#include <stdint.h>
#include "spine.h"

#ifdef __cplusplus
extern "C" {
#endif

int spine_full_init();
void close_spine();
void spine_full_update(uint32_t seq, int16_t* motors_data, uint32_t* leds_data);
spine_dataframe_t iterate();
void init_lcd();
void set_pixels(uint16_t *pixels);

typedef struct IMUData {
    float gx;
    float gy;
    float gz;
    float ax;
    float ay;
    float az;
} IMUData;

int imu_init();
IMUData getIMUData();


#ifdef __cplusplus
}
#endif

#endif // __LIBVECTOR_GOBOT_H

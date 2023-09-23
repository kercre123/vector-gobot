#include "utils.h"
#include "lcd.h"
#include "libvector_gobot.h"

#define BUF_SZ 184*96
uint16_t buff[BUF_SZ];

int lcd_spi;

#define LINE_SIZE 184
#define LINES 96

void init_lcd() {
    lcd_spi = lcd_init();
}

    void set_pixels(uint16_t *pixels) {
        uint16_t buff[BUF_SZ];
        for(int i=0; i<BUF_SZ; i++) {
            buff[i] = pixels[i];
        }
        lcd_write_data(lcd_spi, (char *)buff, BUF_SZ*2);
    }

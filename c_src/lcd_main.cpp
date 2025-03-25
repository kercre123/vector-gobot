#include "utils.h"
#include "lcd.h"
#include "libvector_gobot.h"

#define BUF_SZ  (LCD_FRAME_WIDTH_SANTEK * LCD_FRAME_HEIGHT_SANTEK)
uint16_t buff[BUF_SZ];

#define BUF_SZ_MIDAS  (LCD_FRAME_WIDTH_MIDAS * LCD_FRAME_HEIGHT_MIDAS)
uint16_t buff_midas[BUF_SZ_MIDAS];
int lcd_spi;

#define LINE_SIZE  LCD_FRAME_WIDTH_SANTEK
#define LINES      LCD_FRAME_HEIGHT_SANTEK

#define LINE_SIZE_MIDAS  LCD_FRAME_WIDTH_MIDAS
#define LINES_MIDAS      LCD_FRAME_HEIGHT_MIDAS

bool midas;

extern "C" void core_common_on_exit(void)
{
  lcd_shutdown();
}

void init_lcd() {
    lcd_spi = lcd_init();
    midas = IsXray();
}

void set_pixels(uint16_t *pixels) {
    if (midas) {
        lcd_draw_frame2(pixels, BUF_SZ_MIDAS * 2);
    } else {
        lcd_draw_frame2(pixels, BUF_SZ * 2);
    }
}

// void set_pixels_midas(uint16_t *pixels) {
//     // for midas, swap byte order for each pixel
//     uint16_t conv[BUF_SZ_MIDAS];
//     for (int i = 0; i < BUF_SZ_MIDAS; i++) {
//         conv[i] = __builtin_bswap16(pixels[i]);
//     }
//     lcd_write_data(lcd_spi, (char *)conv, BUF_SZ_MIDAS * 2);
// }

bool is_midas() {
    return IsXray();
}

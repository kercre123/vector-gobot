package main

import (
	"fmt"
	"os"
	"time"

	"github.com/kercre123/vector-gobot/pkg/vcam"
	"github.com/kercre123/vector-gobot/pkg/vscreen"
)

func GetFrameAsUint16Array(rawData []byte, width int, height int) []uint16 {
	downWidth := 184
	downHeight := 96
	rgb565Image := make([]uint16, downWidth*downHeight)

	for y := 0; y < downHeight; y++ {
		for x := 0; x < downWidth; x++ {
			srcX := x * width / downWidth
			srcY := y * height / downHeight
			idxRaw := (srcY*width + srcX) / 4 * 5

			r := (uint16(rawData[idxRaw+0]) << 2) | ((uint16(rawData[idxRaw+4]) >> 6) & 0x03)
			g1 := (uint16(rawData[idxRaw+1]) << 2) | ((uint16(rawData[idxRaw+4]) >> 4) & 0x03)
			g2 := (uint16(rawData[idxRaw+2]) << 2) | ((uint16(rawData[idxRaw+4]) >> 2) & 0x03)
			b := (uint16(rawData[idxRaw+3]) << 2) | ((uint16(rawData[idxRaw+4]) >> 0) & 0x03)
			g := (g1 + g2) >> 1

			pixel := (r&0xF8)<<8 | (g&0xFC)<<3 | b>>3
			rgb565Image[y*downWidth+x] = pixel
		}
	}

	return rgb565Image
}

func main() {
	vcam.InitCam(true)
	vscreen.InitLCD()
	vscreen.BlackOut()
	fmt.Println("Show camera data on screen...")
	for {
		frame, err := vcam.GetFrame()
		if err != nil {
			fmt.Println("error getting frame: ", err)
			os.Exit(1)
		}
		scrnData := GetFrameAsUint16Array(frame, 1280, 720)
		vscreen.SetScreen(scrnData)
		time.Sleep(time.Millisecond * 20)
	}
}

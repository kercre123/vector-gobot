package vscreen

// #cgo LDFLAGS: -lvector-gobot -ldl
// #cgo CFLAGS: -I${SRCDIR}/../../include -w
// #include "libvector_gobot.h"
// #include "lcd.h"
import "C"
import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"os/exec"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

var screenInitted bool

/*
Init the LCD. This must be run before any screen functions are run.
*/
func InitLCD() {
	exec.Command("/bin/bash", "-c", "chmod 666 /sys/module/spidev/parameters/bufsiz").Run()
	exec.Command("/bin/bash", "-c", "echo 35328 > /sys/module/spidev/parameters/bufsiz").Run()
	exec.Command("/bin/bash", "-c", "chmod 444 /sys/module/spidev/parameters/bufsiz").Run()
	C.init_lcd()
	screenInitted = true
	BlackOut()
}

/*
Check if LCD is initiated.
*/
func IsInited() bool {
	return screenInitted
}

func wrapText(text string, lineWidth int) []string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return nil
	}
	lines := words[:1]
	j := 0
	for _, word := range words[1:] {
		if len(lines[j]+" "+word) <= lineWidth {
			lines[j] += " " + word
		} else {
			lines = append(lines, word)
			j++
		}
	}
	return lines
}

/*
Make every pixel on the screen black
*/
func BlackOut() error {
	if !screenInitted {
		return errors.New("init screen first")
	}
	pixels := make([]uint16, 184*96)
	for i := range pixels {
		pixels[i] = 0x000000
	}
	SetScreen(pixels)
	return nil
}

/*
Create screen data from text. It will automatically wrap
*/
func CreateTextImage(text string) []uint16 {
	const W, H = 184, 96
	img := image.NewRGBA(image.Rect(0, 0, W, H))
	black := color.RGBA{0, 0, 0, 255}
	white := color.RGBA{255, 255, 255, 255}

	draw.Draw(img, img.Bounds(), &image.Uniform{black}, image.Point{}, draw.Src)

	d := &font.Drawer{
		Dst:  img,
		Src:  &image.Uniform{white},
		Face: basicfont.Face7x13,
		Dot:  fixed.P(0, 13),
	}

	//fmt.Println(13 * fixed.I(13))

	// Wrap text
	lines := wrapText(text, W/7) // assume each character is ~7px wide
	for _, line := range lines {
		d.Dot.X = 0
		d.DrawString(line)
		d.Dot.Y += fixed.I(13) // move down for the next line
	}

	pixels := make([]uint16, W*H)
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			// Convert the color format from RGBA to RGB565
			pixel := (r>>8&0xF8)<<8 | (g>>8&0xFC)<<3 | b>>8>>3
			pixels[y*W+x] = uint16(pixel)
		}
	}

	return pixels
}

/*
Create screen data from a slice of text
*/
func CreateTextImageFromSlice(lines []string) []uint16 {
	const W, H = 184, 96
	img := image.NewRGBA(image.Rect(0, 0, W, H))
	black := color.RGBA{0, 0, 0, 255}
	white := color.RGBA{255, 255, 255, 255}

	draw.Draw(img, img.Bounds(), &image.Uniform{black}, image.Point{}, draw.Src)

	d := &font.Drawer{
		Dst:  img,
		Src:  &image.Uniform{white},
		Face: basicfont.Face7x13,
		Dot:  fixed.P(0, 13),
	}

	// Wrap text
	for _, line := range lines {
		d.Dot.X = 0
		d.DrawString(line)
		d.Dot.Y += fixed.I(13) // move down for the next line
	}

	pixels := make([]uint16, W*H)
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			// Convert the color format from RGBA to RGB565
			pixel := (r>>8&0xF8)<<8 | (g>>8&0xFC)<<3 | b>>8>>3
			pixels[y*W+x] = uint16(pixel)
		}
	}

	return pixels
}

type Line struct {
	Text  string
	Color color.Color
}

/*
A line is defined as:

	type Line struct {
		Text  string
		Color color.Color
	}
*/
func CreateTextImageFromLines(lines []Line) []uint16 {
	const W, H = 184, 96
	img := image.NewRGBA(image.Rect(0, 0, W, H))
	black := color.RGBA{0, 0, 0, 255}
	white := color.RGBA{255, 255, 255, 255}

	draw.Draw(img, img.Bounds(), &image.Uniform{black}, image.Point{}, draw.Src)

	d := &font.Drawer{
		Dst:  img,
		Src:  &image.Uniform{white},
		Face: basicfont.Face7x13,
		Dot:  fixed.P(0, 13),
	}

	// Wrap text
	for _, line := range lines {
		d.Src = &image.Uniform{line.Color}
		d.Dot.X = 0
		d.DrawString(line.Text)
		d.Dot.Y += fixed.I(13) // move down for the next line
	}

	pixels := make([]uint16, W*H)
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			// Convert the color format from RGBA to RGB565
			pixel := (r>>8&0xF8)<<8 | (g>>8&0xFC)<<3 | b>>8>>3
			pixels[y*W+x] = uint16(pixel)
		}
	}

	return pixels
}

/*
Applies data to the screen
*/
func SetScreen(pixels []uint16) error {
	if !screenInitted {
		return errors.New("screen is not inited")
	}
	C.set_pixels((*C.uint16_t)(&pixels[0]))
	return nil
}

func StopLCD() {
	// the program does not setup a constant communication channel with the screen
	// this function just makes sure we don't send any data to the screen after it is run
	screenInitted = false
}

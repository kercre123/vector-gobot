package main

import (
	"fmt"
	"os"
	"time"

	"github.com/kercre123/vector-gobot/pkg/vbody"
	"github.com/kercre123/vector-gobot/pkg/vscreen"
)

var isMidas bool

func main() {
	fmt.Println("Initing body...")
	err := vbody.InitSpine()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Initing screen...")
	vscreen.InitLCD()
	isMidas, err = vscreen.IsMidas()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(isMidas)
	vbody.SetLEDs(vbody.LED_BLUE, vbody.LED_BLUE, vbody.LED_BLUE)
	fmt.Println("Show readout of sensor values on screen for 10 seconds...")
	exit := false
	go func() {
		time.Sleep(time.Second * 12)
		exit = true
	}()
	frameChan := vbody.GetFrameChan()
	for frame := range frameChan {
		if exit {
			fmt.Println("Done")
			os.Exit(0)
		}
		scrnLines := []string{
			"Touch: " + fmt.Sprint(frame.Touch),
			"Cliffs: " + fmt.Sprint(frame.Cliffs[0]) + " " + fmt.Sprint(frame.Cliffs[1]) + " " + fmt.Sprint(frame.Cliffs[2]) + " " + fmt.Sprint(frame.Cliffs[3]),
		}
		scrnData := vscreen.CreateTextImageFromSlice(scrnLines)
		vscreen.SetScreen(scrnData)
	}
}

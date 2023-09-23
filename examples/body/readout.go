package main

import (
	"fmt"
	"os"
	"time"

	"github.com/kercre123/vector-gobot/pkg/vbody"
	"github.com/kercre123/vector-gobot/pkg/vscreen"
)

func main() {
	fmt.Println("Initing body...")
	err := vbody.InitSpine()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Initing screen...")
	vscreen.InitLCD()
	vbody.SetMotors(0, 0, -100, -200)
	time.Sleep(time.Second * 2)
	vbody.SetMotors(0, 0, 0, 300)
	time.Sleep(time.Second * 1)
	vbody.SetMotors(0, 0, 0, 0)
	fmt.Println("Show readout of sensor values on screen for 10 seconds...")
	exit := false
	go func() {
		time.Sleep(time.Second * 12)
		exit = true
	}()
	for {
		if exit {
			fmt.Println("Done")
			os.Exit(0)
		}
		frame, err := vbody.GetFrame()
		if err != nil {
			fmt.Println("error getting frame: ", err)
			os.Exit(1)
		}
		scrnLines := []string{
			"Touch: " + fmt.Sprint(frame.Touch),
			"Cliffs: " + fmt.Sprint(frame.Cliffs[0]) + " " + fmt.Sprint(frame.Cliffs[1]) + " " + fmt.Sprint(frame.Cliffs[2]) + " " + fmt.Sprint(frame.Cliffs[3]),
		}
		scrnData := vscreen.CreateTextImageFromSlice(scrnLines)
		vscreen.SetScreen(scrnData)
		time.Sleep(time.Millisecond * 5)
	}
}

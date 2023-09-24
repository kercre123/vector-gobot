package main

import (
	"fmt"
	"math"
	"time"

	"github.com/kercre123/vector-gobot/pkg/vimu"
	"github.com/kercre123/vector-gobot/pkg/vscreen"
)

func main() {
	vscreen.InitLCD()
	vscreen.BlackOut()
	err := vimu.InitIMU()
	if err != nil {
		panic(err)
	}
	ticker := time.NewTicker(time.Millisecond * 10)
	for range ticker.C {
		frame, err := vimu.GetFrame()
		if err != nil {
			panic(err)
		}
		scrnLines := []string{
			"IMU:",
			"gX: " + fmt.Sprint(math.Round(float64(frame.Gyro.X))),
			"gY: " + fmt.Sprint(math.Round(float64(frame.Gyro.Y))),
			"gZ: " + fmt.Sprint(math.Round(float64(frame.Gyro.Z))),
			"aX: " + fmt.Sprint(math.Round(float64(frame.Accel.X))),
			"aY: " + fmt.Sprint(math.Round(float64(frame.Accel.Y))),
			"aZ: " + fmt.Sprint(math.Round(float64(frame.Accel.Z))),
		}
		scrnData := vscreen.CreateTextImageFromSlice(scrnLines)
		vscreen.SetScreen(scrnData)
	}
}

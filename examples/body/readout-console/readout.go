package main

import (
	"fmt"
	"time"

	"github.com/inancgumus/screen"
	"github.com/kercre123/vector-gobot/pkg/vbody"
)

func main() {
	screen.Clear()
	vbody.ReadOnly = true
	vbody.InitSpine()
	fchan := vbody.GetFrameChan()
	fps := 0
	fpsFinal := 0
	timeBefore := time.Now()
	for frame := range fchan {
		fps++
		screen.MoveTopLeft()
		fmt.Println("Frames recieved per second: ", fpsFinal)
		fmt.Println("Button: ", frame.ButtonState)
		fmt.Println("Cliff sensors: ", frame.Cliffs)
		fmt.Println("Charger voltage: ", frame.ChargerVoltage)
		if time.Since(timeBefore) >= time.Second {
			fpsFinal = fps
			fps = 0
			timeBefore = time.Now()
		}
	}
}

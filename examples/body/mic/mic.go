package main

import (
	"encoding/binary"
	"os"
	"time"

	"github.com/kercre123/vector-gobot/pkg/vbody"
)

func main() {
	file, err := os.Create("/tmp/output.pcm")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	err = vbody.InitSpine()
	if err != nil {
		panic(err)
	}
	shouldStop := false
	go func() {
		time.Sleep(time.Second * 10)
		shouldStop = true
	}()
	go func() {
		vbody.SetLEDs(vbody.LED_GREEN, vbody.LED_BLUE, vbody.LED_RED)
		vbody.SetMotors(0, 0, 0, 60)
		time.Sleep(time.Second * 3)
		vbody.SetMotors(0, 0, 0, 0)
		vbody.SetMotors(0, 0, 0, -60)
		time.Sleep(time.Second * 3)
	}()
	frameChan := vbody.GetFrameChan()
	for frame := range frameChan {
		if shouldStop {
			return
		}
		err := binary.Write(file, binary.LittleEndian, frame.MicData)
		if err != nil {
			panic(err)
		}
	}
	vbody.StopSpine()
}

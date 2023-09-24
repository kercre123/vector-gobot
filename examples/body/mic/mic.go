package main

import (
	"encoding/binary"
	"fmt"
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
	vbody.InitSpine()
	shouldStop := false
	go func() {
		time.Sleep(time.Second * 12)
		shouldStop = true
	}()
	ticker := time.NewTicker(time.Millisecond * 10)
	for range ticker.C {
		if shouldStop {
			return
		}
		frame, _ := vbody.GetFrame()
		fmt.Println(len(frame.MicData))
		for _, sample := range frame.MicData {
			err := binary.Write(file, binary.BigEndian, sample)
			if err != nil {
				panic(err)
			}
		}
	}
	vbody.StopSpine()
}

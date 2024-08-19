package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"sync"

	"github.com/kercre123/vector-gobot/pkg/vbody"
	vosk "github.com/kercre123/vosk-api/go"
	"github.com/maxhawkins/go-webrtcvad"
)

var (
	rec         *vosk.VoskRecognizer
	vadthing    *webrtcvad.VAD
	chunkBuffer [][]int16
	bufferMutex sync.Mutex
	micDataBuf  []int16
)

func main() {

	vosk.SetLogLevel(-1)

	var err error
	vadthing, err = webrtcvad.New()
	if err != nil {
		panic(err)
	}

	model, err := vosk.NewModel("./vosk-model")
	if err != nil {
		panic(err)
	}

	rec, err = vosk.NewRecognizer(model, 16000)
	if err != nil {
		panic(err)
	}

	vadthing.SetMode(3)
	vbody.ReadOnly = true
	err = vbody.InitSpine()
	if err != nil {
		panic(err)
	}
	defer vbody.StopSpine()

	go frameGetter()

	fmt.Println("Say something to Vector! This program will transcribe all mic data and print it to console.")
	for {
		chunk := getNextChunkFromBuffer()
		if chunk == nil {
			continue
		}

		iVoled := increaseVolume(chunk, 15)
		var bufBytes []byte
		binchunk := bytes.NewBuffer(bufBytes)
		binary.Write(binchunk, binary.LittleEndian, iVoled)

		rec.AcceptWaveform(binchunk.Bytes())

		if IsDoneSpeaking(binchunk.Bytes()) {
			var jsonMap map[string]string
			json.Unmarshal([]byte(rec.FinalResult()), &jsonMap)
			fmt.Println(jsonMap["text"])
			rec.Reset()
		}
	}
}

func frameGetter() {
	frameChan := vbody.GetFrameChan()
	for frame := range frameChan {
		smashed := smashPCM(frame.MicData)
		fullBuf, _, isFull := fillBuf(smashed)
		if isFull {
			fillBuffer(fullBuf)
		}
	}
}

func fillBuffer(data []int16) {
	bufferMutex.Lock()
	defer bufferMutex.Unlock()

	chunkBuffer = append(chunkBuffer, data)
}

func fillBuf(in []int16) (full []int16, leftover []int16, filled bool) {
	for i, inny := range in {
		micDataBuf = append(micDataBuf, inny)
		if len(micDataBuf) == 320 {
			mbuf := micDataBuf
			micDataBuf = []int16{}
			return mbuf, in[i+1:], true
		}
	}
	return nil, micDataBuf, false
}

// Retrieve the next chunk from the buffer
func getNextChunkFromBuffer() []int16 {
	bufferMutex.Lock()
	defer bufferMutex.Unlock()

	if len(chunkBuffer) > 0 {
		chunk := chunkBuffer[0]
		chunkBuffer = chunkBuffer[1:]
		return chunk
	}
	return nil
}

var activeCount int
var inactiveCount int

func IsDoneSpeaking(chunk320 []byte) bool {
	// technically lower than 16000 but whatevs
	active, err := vadthing.Process(16000, chunk320)
	if err != nil {
		panic(err)
	}
	if active {
		inactiveCount = 0
		activeCount++
	} else {
		inactiveCount++
		if inactiveCount == 15 {
			if activeCount >= 15 {
				activeCount = 0
				inactiveCount = 0
				return true
			} else {
				activeCount = 0
			}
		}
	}
	return false
}

func increaseVolume(input []int16, factor int16) []int16 {
	output := make([]int16, len(input))

	for i, sample := range input {
		newSample := int32(sample) * int32(factor)
		if newSample > math.MaxInt16 {
			newSample = math.MaxInt16
		} else if newSample < math.MinInt16 {
			newSample = math.MinInt16
		}

		output[i] = int16(newSample)
	}

	return output
}

func smashPCM(input []int16) []int16 {
	if len(input) != 320 {
		panic("gotta be 320 m8")
	}

	output := make([]int16, 80)

	// Extract only the 2nd channel
	for i := 0; i < 80; i++ {
		output[i] = input[i*4]
	}

	return output
}

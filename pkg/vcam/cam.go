package vcam

/*
 #cgo LDFLAGS: -L${SRCDIR}/.
 #cgo CFLAGS: -I${SRCDIR}/. -I${SRCDIR}/inc
 #include "camera_client.h"
*/
import "C"
import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"math"
	"os"
	"os/exec"
	"sort"
	"strings"
	"sync"
	"time"
	"unsafe"
)

var camera *C.struct_anki_camera_handle
var frameBuffer []byte
var frameBufferMutex sync.Mutex
var stopLooping bool
var readyForFrames bool

func sleep(ms int) {
	time.Sleep(time.Millisecond * time.Duration(ms))
}

/*
Init camera. This initiates the communication with mm-anki-camera and starts getting frames.
It also modifies the service file to get full framerate from the camera.
If you want a horrible autoexposure implementation, give it `true`.
*/
func InitCam(autoExposure bool) error {
	file, err := os.ReadFile("/lib/systemd/system/mm-anki-camera.service")
	if err != nil {
		panic("mm-anki-camera service doesn't exist. is this an Anki Vector?")
	}
	if !strings.Contains(string(file), "mm-anki-camera -r 1") {
		fmt.Println("vcam: Adding -r 1 to mm-anki-camera service file for faster framerate (only happens on first init)...")
		original := `/usr/bin/mm-anki-camera $MM_ANKI_CAMERA_OPTS`
		replacement := `/usr/bin/mm-anki-camera -r 1 $MM_ANKI_CAMERA_OPTS`
		err := exec.Command("sudo", "sed", "-i", fmt.Sprintf("s|%s|%s|g", original, replacement), "/lib/systemd/system/mm-anki-camera.service").Run()
		if err != nil {
			fmt.Println("failed to execute sed command:", err)
		}
		exec.Command("/bin/bash", "-c", "systemctl daemon-reload").Run()
		exec.Command("/bin/bash", "-c", "systemctl restart mm-anki-camera").Run()
		sleep(100)
		fmt.Println("success")
	}
	readyForFrames = false

	rc := C.camera_init(&camera)
	if rc != 0 {
		return fmt.Errorf("failed to initialize camera camera_init()")
	}

	sleep(1000)
	rc = C.camera_start(camera)
	if rc != 0 {
		return fmt.Errorf("failed to start camera camera_start()")
	}

	stopLooping = false
	go func() {
		// wait for camera to be ready
		for C.camera_status(camera) != C.ANKI_CAMERA_STATUS_RUNNING && !stopLooping {
			sleep(30)
		}
		var r C.int

		// frame-buffer-fill loop
		for C.camera_status(camera) == C.ANKI_CAMERA_STATUS_RUNNING && !stopLooping {
			sleep(30)
			var frame *C.anki_camera_frame_t

			r = C.camera_frame_acquire(camera, 0, &frame)
			if r != 0 {
				continue
			}

			frameSize := int(frame.height) * int(frame.bytes_per_row)
			frameData := C.GoBytes(unsafe.Pointer(&frame.data), C.int(frameSize))

			//fmt.Println(int(frame.width), int(frame.height), int(frame.bits_per_pixel), int(frame.bytes_per_row))

			frameBufferMutex.Lock()
			frameBuffer = make([]byte, frameSize)
			copy(frameBuffer, frameData)
			frameBufferMutex.Unlock()
			C.camera_frame_release(camera, frame.frame_id)
			if !readyForFrames && len(frameData) > 0 {
				readyForFrames = true
			}
		}
	}()
	if autoExposure {
		// setup auto exposure
		go func() {
			for {
				if !readyForFrames {
					time.Sleep(time.Millisecond * 50)
				} else {
					break
				}
			}
			time.Sleep(time.Second * 1)
			for !stopLooping {
				sleep(200)
				frame, _ := GetFrame()
				runAutoExposure(frame)
			}
		}()
	}

	// only return when frames are ready
	for {
		if !readyForFrames {
			time.Sleep(time.Millisecond * 50)
		} else {
			break
		}
	}
	return nil
}

/*
Gets a frame from the camera when the next one is ready.
Read about output data here: https://github.com/digital-dream-labs/vector/blob/main/docs/vision/Debayering.md
DebayerRGGBBilinear() can convert this to a nice image.Image.
vjpeg includes a function for converting this to JPEG.
*/
func GetFrame() ([]byte, error) {
	if !readyForFrames {
		return nil, errors.New("camera not inited")
	}
	frameBufferMutex.Lock()
	defer frameBufferMutex.Unlock()
	return frameBuffer, nil
}

/*
Stop the camera. This should be run when you are done using the camera.
*/
func StopCam() error {
	if !readyForFrames {
		return fmt.Errorf("camera already stopped")
	}
	readyForFrames = false
	stopLooping = true
	fmt.Println("Stopping Camera...")
	rc := C.camera_stop(camera)
	if rc != 0 {
		return fmt.Errorf("failed to stop camera")
	}

	rc = C.camera_release(camera)
	if rc != 0 {
		return fmt.Errorf("failed to release camera")
	}

	return nil
}

/*
Check if vcam is recieving frames from the camera
*/
func IsInited() bool {
	return readyForFrames
}

/*
Set exposure.
Accepts milliseconds and gain.
I still need to experiment with limits and combinations.
*/
func SetExposure(ms uint16, gain float64) {
	if !readyForFrames {
		fmt.Println("must init camera before setting exposure")
		return
	}
	C.camera_set_exposure(camera, C.uint16_t(ms), C.float(gain))
}

func runAutoExposure(rawData []byte) (uint16, float64) {
	width := 1280
	height := 720
	// debayering, go style
	rgbImage := make([][][]uint8, height/2)
	for i := range rgbImage {
		rgbImage[i] = make([][]uint8, width/2)
		for j := range rgbImage[i] {
			rgbImage[i][j] = make([]uint8, 3)
		}
	}

	for y := 0; y < height; y += 2 {
		for x := 0; x < width; x += 2 {
			idxRaw := (y*width + x) / 4 * 5
			idxY := y / 2
			idxX := x / 2

			r := (uint16(rawData[idxRaw+0]) << 2) | ((uint16(rawData[idxRaw+4]) >> 6) & 0x03)
			g1 := (uint16(rawData[idxRaw+1]) << 2) | ((uint16(rawData[idxRaw+4]) >> 4) & 0x03)
			g2 := (uint16(rawData[idxRaw+2]) << 2) | ((uint16(rawData[idxRaw+4]) >> 2) & 0x03)
			b := (uint16(rawData[idxRaw+3]) << 2) | ((uint16(rawData[idxRaw+4]) >> 0) & 0x03)

			g := (g1 + g2) >> 1

			rgbImage[idxY][idxX][0] = uint8(r >> 2)
			rgbImage[idxY][idxX][1] = uint8(g >> 2)
			rgbImage[idxY][idxX][2] = uint8(b >> 2)
		}
	}

	// brightest spot alg
	brightnessValues := make([]float64, 0, (height/2)*(width/2))
	sumBrightness := 0.0
	for _, row := range rgbImage {
		for _, pixel := range row {
			brightness := float64(pixel[0]+pixel[1]+pixel[2]) / 3.0
			brightnessValues = append(brightnessValues, brightness)
			sumBrightness += brightness
		}
	}

	sort.Float64s(brightnessValues)
	percentile95 := brightnessValues[int(0.95*float64(len(brightnessValues)))]

	meanBrightness := sumBrightness / float64(len(brightnessValues))

	targetBrightness := 130.0

	referenceBrightness := 0.7*percentile95 + 0.3*meanBrightness

	exposureMsFloat := 100 * (targetBrightness - referenceBrightness) / targetBrightness
	exposureMs := uint16(math.Max(1, math.Min(100, exposureMsFloat))) // 100 seems to be effective max

	gain := 4 * (targetBrightness - referenceBrightness) / targetBrightness
	gain = math.Max(0, math.Min(5, gain))
	SetExposure(exposureMs, gain)

	return exposureMs, gain
}

func unpackRaw10Bilinear(rawData []byte) []uint16 {
	unpackedData := make([]uint16, len(rawData)*8/10)
	for i := 0; i < len(rawData)/5*5; i += 5 {
		unpackedData[i/5*4+0] = uint16(rawData[i+0])<<2 | uint16(rawData[i+4]>>6)&0x03
		unpackedData[i/5*4+1] = uint16(rawData[i+1])<<2 | uint16(rawData[i+4]>>4)&0x03
		unpackedData[i/5*4+2] = uint16(rawData[i+2])<<2 | uint16(rawData[i+4]>>2)&0x03
		unpackedData[i/5*4+3] = uint16(rawData[i+3])<<2 | uint16(rawData[i+4]>>0)&0x03
	}
	return unpackedData
}

/*
Takes camera data and converts it to an image.NRGBA.
Does full 1280x720, though it's a little slow. Meant for single frames, not streams.
*/
func DebayerRGGBBilinear(camData []uint16, width, height int) *image.NRGBA {
	rgbImage := image.NewNRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			idx := y*width + x
			r, g, b := uint16(0), uint16(0), uint16(0)
			if y%2 == 0 && x%2 == 0 {
				r = camData[idx]
			} else if y%2 == 1 && x%2 == 1 {
				b = camData[idx]
			} else {
				g = camData[idx]
			}
			rgbImage.SetNRGBA(x, y, color.NRGBA{R: uint8(r >> 2), G: uint8(g >> 2), B: uint8(b >> 2), A: 255})
		}
	}

	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			c := rgbImage.NRGBAAt(x, y)
			if c.R == 0 {
				c.R = uint8((uint16(rgbImage.NRGBAAt(x-1, y).R) + uint16(rgbImage.NRGBAAt(x+1, y).R) + uint16(rgbImage.NRGBAAt(x, y-1).R) + uint16(rgbImage.NRGBAAt(x, y+1).R)) / 4)
			}
			if c.G == 0 {
				c.G = uint8((uint16(rgbImage.NRGBAAt(x-1, y).G) + uint16(rgbImage.NRGBAAt(x+1, y).G) + uint16(rgbImage.NRGBAAt(x, y-1).G) + uint16(rgbImage.NRGBAAt(x, y+1).G)) / 4)
			}
			if c.B == 0 {
				c.B = uint8((uint16(rgbImage.NRGBAAt(x-1, y).B) + uint16(rgbImage.NRGBAAt(x+1, y).B) + uint16(rgbImage.NRGBAAt(x, y-1).B) + uint16(rgbImage.NRGBAAt(x, y+1).B)) / 4)
			}
			rgbImage.SetNRGBA(x, y, c)
		}
	}

	return rgbImage
}

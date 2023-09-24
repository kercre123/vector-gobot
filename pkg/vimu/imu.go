package vimu

// #cgo LDFLAGS: -lvector-gobot -ldl
// #cgo CFLAGS: -I${SRCDIR}/../../include -w
// #include "libvector_gobot.h"
import "C"

import (
	"fmt"
	"sync"
	"time"
)

type IMUFrame struct {
	Gyro struct {
		X float32
		Y float32
		Z float32
	}
	Accel struct {
		X float32
		Y float32
		Z float32
	}
}

var CurrentIMUFrame struct {
	mu sync.Mutex
	IMUFrame
}

var imuSPI int
var IMUInited bool

// Init the IMU, must be run before you get a frame
func InitIMU() error {
	spi := C.imu_init()
	if int(spi) != 0 {
		return fmt.Errorf("error initializing imu: " + fmt.Sprint(int(spi)))
	}
	imuSPI = int(spi)
	IMUInited = true
	go commsLoop()
	time.Sleep(time.Millisecond * 200)
	return nil
}

func commsLoop() {
	ticker := time.NewTicker(time.Millisecond * 10)
	defer ticker.Stop()
	for range ticker.C {
		if !IMUInited {
			break
		}
		data := C.getIMUData()
		CurrentIMUFrame.mu.Lock()
		CurrentIMUFrame.IMUFrame.Gyro.X = float32(data.gx)
		CurrentIMUFrame.IMUFrame.Gyro.Y = float32(data.gy)
		CurrentIMUFrame.IMUFrame.Gyro.Z = float32(data.gz)
		CurrentIMUFrame.IMUFrame.Accel.X = float32(data.ax)
		CurrentIMUFrame.IMUFrame.Accel.Y = float32(data.ay)
		CurrentIMUFrame.IMUFrame.Accel.Z = float32(data.az)
		CurrentIMUFrame.mu.Unlock()
	}
}

// Stop the IMU. Stops comms loop
func StopIMU() {
	IMUInited = false
	time.Sleep(time.Millisecond * 100)
}

// Get a frame from the IMU
func GetFrame() (IMUFrame, error) {
	if !IMUInited {
		return IMUFrame{}, fmt.Errorf("imu not inited")
	}
	CurrentIMUFrame.mu.Lock()
	defer CurrentIMUFrame.mu.Unlock()
	return CurrentIMUFrame.IMUFrame, nil
}

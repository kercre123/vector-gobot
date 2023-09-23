# vector-gobot

This is a Go wrapper for [poc.vic-hack](https://github.com/torimos/poc.vic-hack). It allows you to directly communicate with the Anki Vector robot's hardware with easy-to-use Go functions and definitions.

## Modules

-   vcam
    -   Communicates with the built-in camera daemon to get frames from the camera.
    -   Can set exposure and gain as well.
-   vbody
    -   Fully functional bodyboard communication.
-   vscreen
    -   Fully functional LCD communication.
-   vjpeg
    -   Meant to be used in conjunction with vcam.
    -   Takes a camera frame, unpacks it, debayers it, and converts it to JPEG as fast as possible with turbojpeg.
    -   Needs an extra lib to be built.

## Building

1. Clone the repo (with submodules):

```
git clone --recurse-submodules https://github.com/kercre123/vector-gobot
```

2. Build the libs:

```
make
```

-   At this point, you can use vcam, vscreen, and vbody. To use vjpeg, you must compile libjpeg-turbo.so:

```
make libjpeg-turbo
```

3. Build your program (in this case, the body example program):

```
CC="$(pwd)/vic-toolchain/arm-linux-gnueabi/bin/arm-linux-gnueabi-gcc" \
CGO_LDFLAGS="-Lbuild" \
GOARM=7 \
GOARCH=arm \
CGO_ENABLED=1 \
go build -o build/main examples/body/readout.go
```

-   You can define your own toolchain. Just make sure the same one used to build the libs is also the one you are using to build the Go program. Example: `TOOLCHAIN=$HOME/vchain/arm-linux-gnueabi/bin/arm-linux-gnueabi- make`

## Installing a program on a bot

1. The anki-robot.target must be stopped, so SSH in and run:

```
systemctl stop anki-robot.target
```

2. The built libs (in ./build) need to be in a place where programs can find them. You can do this by copying them over to the bot's /lib folder or by putting them in any directory and running a program like this:

```
LD_LIBRARY_PATH=/data/gobot_libs ./main
```

-   Note: by default, /data is mounted as noexec, meaning you can't run any programs in it. To change this:

```
mount -o rw,remount,exec /data
```

## Features

1. Spine (vspine)
    -   [x] LEDs
    -   [x] Motors
    -   [x] Encoders
    -   [x] Mics
    -   [x] Touch
    -   [x] Battery/Charger Voltage
    -   [x] Body Temperature
    -   [ ] ToF Sensor
        -   It is able to get data, but I haven't figured out how to calculate it all into a nice mm value
    -   [x] Cliff Sensors
2. Camera (vcam)
    -   Communicates with mm-anki-camera to get data
    -   Changes service file to get full framerate from camera
    -   Gets RGGB10 data from camera. Functions to convert that are included
3. JPEG (vjpeg)
    -   Meant to be used in conjunction with camera 
    -   Communicates with libturbojpeg to quickly compress frames, like for an MJPEG stream
    -   Includes a direct RGGB10-to-JPEG function
4. Screen (vscreen)
    -   Fully working

## TODO

-   PCM (speaker)
-   IMU
-   ToF sensor calculations
-   Maybe read from calibration files
-   Implement a way to normalize a camera image and remove distortion from lens
-   Nice mic interface functions
-   Remove all memory leaks
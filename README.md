# vector-gobot

A Go wrapper for [poc.vic-hack](https://github.com/torimos/poc.vic-hack).

## Features

1. Spine (vspine)
-   [x] LEDs
-   [x] Motors
-   [x] Encoders
-   [x] Mics
-   [x] Touch
-   [x] Battery/Charger Voltage
-   [x] Body Temperature
-   [/] ToF Sensor
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
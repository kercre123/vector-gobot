package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/kercre123/vector-gobot/pkg/vcam"
	"github.com/kercre123/vector-gobot/pkg/vjpeg"
)

const quality = 60

func mjpegStream(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary=frame")

	for {
		frame, _ := vcam.GetFrame()
		jpegData := vjpeg.RGGB10ToJPEGDownSample(frame, quality)
		_, err := fmt.Fprintf(w, "--frame\r\nContent-Type: image/jpeg\r\nContent-Length: %d\r\n\r\n", len(jpegData))
		if err != nil {
			fmt.Println("stopping mjpeg stream: " + err.Error())
			break
		}
		_, err = w.Write(jpegData)
		if err != nil {
			fmt.Println("stopping mjpeg stream: " + err.Error())
			break
		}
		_, err = w.Write([]byte("\r\n"))
		if err != nil {
			fmt.Println("stopping mjpeg stream: " + err.Error())
			break
		}
		time.Sleep(time.Second / 30)
	}
}

func BeginServer() {
	vcam.InitCam(true)
	http.HandleFunc("/stream", mjpegStream)
	fmt.Println("listening at port 8888")
	http.ListenAndServe(":8888", nil)
}

func main() {
	BeginServer()
}

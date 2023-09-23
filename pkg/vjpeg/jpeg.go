package vjpeg

/*
#cgo CFLAGS: -I${SRCDIR}/../../include -I${SRCDIR}/../../include/libjpeg_turbo
#cgo LDFLAGS: -lturbojpeg -ljpeg_interface
#include "libjpeg_interface.h"
*/
import "C"
import "unsafe"

/*
Encode YUV data to JPEG. This is not for raw camera data, more just for testing.
*/
func EncodeToJPEG(yuvData []byte, quality int, width int, height int) []byte {
	var jpegSize C.ulong
	var jpegBuf *C.uchar
	C.encodeToJPEG((*C.uchar)(&yuvData[0]), C.int(width), C.int(height), C.int(quality), &jpegBuf, &jpegSize)
	goSlice := C.GoBytes(unsafe.Pointer(jpegBuf), C.int(jpegSize))
	C.free(unsafe.Pointer(jpegBuf))
	return goSlice
}

/*
RGGB10-debayer-downsample-to-JPEG. Ends up with a 640x480 resolution.
rawData should come directly from getFrame(), quality should be between 1-100
*/
func RGGB10ToJPEGDownSample(rawData []byte, quality int) []byte {
	width := 1280
	height := 720
	var jpegSize C.ulong
	var jpegBuf *C.uchar
	C.GetFrameAsJPEGDownSampled((*C.uint8_t)(unsafe.Pointer(&rawData[0])), C.int(width), C.int(height), C.int(quality), &jpegBuf, &jpegSize)
	jpegData := C.GoBytes(unsafe.Pointer(jpegBuf), C.int(jpegSize))
	C.free(unsafe.Pointer(jpegBuf))
	return jpegData
}

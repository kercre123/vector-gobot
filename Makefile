COMPILEFILE := main.go

# Get the absolute path of the current directory
ABSPATH := $(shell pwd)

# Check if TOOLCHAIN is not set
ifeq ($(TOOLCHAIN),)
  # Set TOOLCHAIN to a default value based on ABSPATH
  TOOLCHAIN_DIR := $(ABSPATH)/vic-toolchain/arm-linux-gnuebai/bin
  TOOLCHAIN := $(TOOLCHAIN_DIR)/arm-linux-gnueabi-

  ifeq ($(shell test -d $(dir $(TOOLCHAIN_DIR)) && echo yes),)
    $(error The directory $(dir $(TOOLCHAIN_DIR)) does not exist. You must define a $$TOOLCHAIN or follow the README instructions to get a toolchain.)
  endif
endif

COMMON_FLAGS := -O3 -mfpu=neon-vfpv4 -mfloat-abi=softfp -mcpu=cortex-a7 -ffast-math
GPP_FLAGS := -w -shared -Iinclude -fPIC -std=c++11

all: vector-gobot jpeg_interface

vector-gobot:
	$(TOOLCHAIN)g++ \
	$(GPP_FLAGS) $(COMMON_FLAGS) \
	-o build/libvector-gobot.so \
	c_src/*.cpp \
	c_src/libs/*.cpp

jpeg_interface:
	$(TOOLCHAIN)g++ $(GPP_FLAGS) $(COMMON_FLAGS) -o build/libjpeg_interface.so hacksrc/jpeg.cpp -Ilibjpeg-turbo/include -fopenmp

go_build:
	CC="$(TOOLCHAIN)gcc -w -Lbuild" \
	CGO_CFLAGS="-Iinclude $(COMMON_FLAGS) -Ilibjpeg-turbo/include" \
	CGO_LDFLAGS="-ldl" \
	GOARM=7 \
	GOARCH=arm \
	CGO_ENABLED=1 \
	go build -ldflags '-w -s' -o build/main $(COMPILEFILE)

clean:
	rm -f build/librobot.so build/libjpeg_interface.so build/libanki-camera.so build/main

.PHONY: all librobot.so vector-gobot jpeg_interface go_build clean

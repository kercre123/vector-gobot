COMPILEFILE := ./examples/body/readout.go

ABSPATH := $(shell pwd)

ifeq ($(TOOLCHAIN),)
  TOOLCHAIN_DIR := $(ABSPATH)/vic-toolchain/arm-linux-gnueabi/bin
  TOOLCHAIN := $(TOOLCHAIN_DIR)/arm-linux-gnueabi-
  ifeq ($(shell test -d $(dir $(TOOLCHAIN_DIR)) && echo yes),)
    $(error The directory $(dir $(TOOLCHAIN_DIR)) does not exist. You must define a $$TOOLCHAIN or follow the README instructions to get a toolchain.)
  endif
endif

COMMON_FLAGS := -O3 -mfpu=neon-vfpv4 -mfloat-abi=softfp -mcpu=cortex-a7 -ffast-math
GPP_FLAGS := -w -shared -Iinclude -fPIC -std=c++11

all: vector-gobot jpeg_interface
	echo "Successfully compiled libvector-gobot.so and libjpeg_interface.so to ./build."

vector-gobot:
	mkdir -p build
	$(TOOLCHAIN)g++ \
	$(GPP_FLAGS) $(COMMON_FLAGS) \
	-o build/libvector-gobot.so \
	c_src/*.cpp \
	c_src/libs/*.cpp

libjpeg-turbo:
	mkdir -p build
	./make-turbojpeg.sh

jpeg_interface:
	mkdir -p build
	$(TOOLCHAIN)g++ $(GPP_FLAGS) $(COMMON_FLAGS) -o build/libjpeg_interface.so c_src/jpeg/jpeg.cpp -Ilibjpeg-turbo -fopenmp

example:
	CC="$(TOOLCHAIN)gcc" \
	CGO_LDFLAGS="-Lbuild" \
	GOARM=7 \
	GOARCH=arm \
	CGO_ENABLED=1 \
	go build -o build/main $(COMPILEFILE)

clean:
	mkdir -p build
	rm -f build/libvector-gobot.so build/libjpeg_interface.so build/libjpeg_turbo build/main

.PHONY: all librobot.so vector-gobot jpeg_interface example clean libjpeg-turbo

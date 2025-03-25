COMPILEFILE := ./examples/body/readout.go

ABSPATH := $(shell pwd)

ifeq ($(GCC),)
  TOOLCHAIN_DIR := $(ABSPATH)/vic-toolchain/arm-linux-gnueabi/bin
  TOOLCHAIN := $(TOOLCHAIN_DIR)/arm-linux-gnueabi-
  GCC := ${TOOLCHAIN}gcc
  GPP := ${TOOLCHAIN}g++
  ifeq ($(shell test -d $(dir $(TOOLCHAIN_DIR)) && echo yes),)
    $(error The directory $(dir $(TOOLCHAIN_DIR)) does not exist. You must define a $$TOOLCHAIN or follow the README instructions to get a toolchain.)
  endif
endif

ifneq (,$(findstring gnueabihf,$(TOOLCHAIN)))
    COMMON_FLAGS := -O3 -mfpu=neon-vfpv4 -mfloat-abi=hard -mcpu=cortex-a7 -ffast-math -fpermissive
else ifneq (,$(findstring gnueabi,$(TOOLCHAIN)))
    COMMON_FLAGS := -O3 -mfpu=neon-vfpv4 -mfloat-abi=softfp -mcpu=cortex-a7 -ffast-math
else ifneq (,$(findstring aarch64,$(TOOLCHAIN)))
    COMMON_FLAGS := -O3 -ffast-math
else
    COMMON_FLAGS := -O3 -ffast-math -fpermissive
endif

GPP_FLAGS := -w -shared -Iinclude -fPIC -std=c++11 -Wno-c++11-narrowing
GCC_FLAGS := -w -shared -Iinclude -fPIC

all: vector-gobot
	echo "Successfully compiled libvector-gobot.so and libjpeg_interface.so to ./build."

vector-gobot:
	mkdir -p build
	$(GPP) $(GPP_FLAGS) $(COMMON_FLAGS) -fPIC -c c_src/*.cpp c_src/libs/*.cpp
	$(GCC) $(GCC_FLAGS) $(COMMON_FLAGS) -fPIC -c c_src/libs/*.c
	$(GPP) $(GPP_FLAGS) $(COMMON_FLAGS) -shared -latomic -o build/libvector-gobot.so *.o
	rm -f *.o


libjpeg-turbo:
	mkdir -p build
	./make-turbojpeg.sh

jpeg_interface:
	mkdir -p build
	$(GPP) $(GPP_FLAGS) $(COMMON_FLAGS) -o build/libjpeg_interface.so c_src/jpeg/jpeg.cpp -Ilibjpeg-turbo -fopenmp -static-libstdc++

example:
	CC="$(GCC)" \
	CGO_LDFLAGS="-Lbuild" \
	GOARM=7 \
	GOARCH=arm \
	CGO_ENABLED=1 \
	go build -o build/main $(COMPILEFILE)

clean:
	mkdir -p build
	rm -f build/libvector-gobot.so build/libjpeg_interface.so build/libjpeg_turbo build/main

.PHONY: all librobot.so vector-gobot jpeg_interface example clean libjpeg-turbo

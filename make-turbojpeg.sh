#!/bin/bash

ABSPATH="$(pwd)"


if [[ $TOOLCHAIN ]]; then
	if [[ ! -f ${TOOLCHAIN}-gcc ]]; then
		echo "The toolchain you have provided is invalid. ${TOOLCHAIN}-gcc does not exist."
		exit 1
	fi
else
	if [[ ! -f ${ABSPATH}/vic-toolchain/arm-linux-gnueabi/bin/arm-linux-gnueabi-gcc ]]; then
		echo "Toolchain not found. Either define one or follow the README instructions to clone the vic-toolchain submodule."
		exit 1
	else
		TOOLCHAIN=${ABSPATH}/vic-toolchain/arm-linux-gnueabi/bin/arm-linux-gnueabi-
	fi
fi

if [[ ! -f libjpeg-turbo/CMakeLists.txt ]]; then
	echo "libjpeg-turbo has not been cloned. Follow the README instructions to clone the submodule."
else
	echo "Everything is in place! Building libjpeg-turbo..."
	cd libjpeg-turbo
	rm -rf build
	mkdir -p build
	cd build
	ARMCC_FLAGS="-mfloat-abi=softfp -mfpu=neon-vfpv4 -mcpu=cortex-a7 -O3 -ffast-math -fopenmp"
	ARMCC_PREFIX=${TOOLCHAIN}
	cmake -DCMAKE_C_COMPILER=${ARMCC_PREFIX}gcc \
  	-DCMAKE_CXX_COMPILER=${ARMCC_PREFIX}g++ \
  	-DCMAKE_C_FLAGS="${ARMCC_FLAGS}" \
  	-DCMAKE_CXX_FLAGS="${ARMCC_FLAGS}" \
  	-DCMAKE_VERBOSE_MAKEFILE:BOOL=ON \
  	-DCMAKE_INSTALL_PREFIX=${ABSPATH}/build/libjpeg-turbo \
  	-DCMAKE_SYSTEM_NAME=Linux \
  	-DCMAKE_SYSTEM_PROCESSOR=armv7 \
  	..
	cd build
	make -j
	make install
	echo
	echo "libjpeg-turbo has been built! ${ABSPATH}/build/libjpeg-turbo/lib"
fi

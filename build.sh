#!/bin/bash

CC="$(pwd)/vic-toolchain/arm-linux-gnueabihf/bin/arm-linux-gnueabihf-gcc" \
CGO_LDFLAGS="-L$(pwd)/build -L$(pwd)/build/libjpeg-turbo/lib" \
GOARM=7 \
GOARCH=arm \
CGO_ENABLED=1 \
go build -o main $@

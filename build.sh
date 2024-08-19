#!/bin/bash

CC="$(pwd)/vic-toolchain/arm-linux-gnueabi/bin/arm-linux-gnueabi-gcc" \
CGO_LDFLAGS="-L$(pwd)/build -L$(pwd)/build/libjpeg-turbo/lib" \
GOARM=7 \
GOARCH=arm \
CGO_ENABLED=1 \
go build -o main $@

#!/bin/bash

#CC="$(pwd)/vic-toolchain/arm-linux-gnueabihf/bin/arm-linux-gnueabihf-gcc" \
CGO_LDFLAGS="-L$(pwd)/build -L$(pwd)/build/libjpeg-turbo/lib" \
CGO_ENABLED=1 \
go build -o main $@

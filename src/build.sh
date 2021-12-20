#!/bin/bash

# MedHash Tools
# Copyright (c) 2021 GHIFARI160
# MIT License

# buildArgs="-trimpath"
# ldArgs="-s -w -X github.com/ghifari160/medhash-tools/src/common.VERSION=$ver"

# Build a given tool for a given platform and architecture
# $1: GOOS
# $2: GOARCH
# $3: Tool name
# $4: Suffix
build_helper () {
    echo "Building $3 for $1 / $2 (Suffix: $4)"
    GOOS=$1 GOARCH=$2 go build $buildArgs -ldflags="$ldArgs" -o $buildDir/$3$4 ./$3
}

# Build all tools
# $1: GOOS
# $2: GOARCH
# $3: Suffix
build_all () {
    build_helper $1 $2 medhash-gen $3
    build_helper $1 $2 medhash-chk $3
    build_helper $1 $2 medhash-upgrade $3
}

# Merge binaries into a universal binary
# $1: Tool name
# $2: Intel binary
# $3: M1 binary
# $4: Platform
# $5: Suffix
merge () {
    echo "Building universal binary for $1"
    lipo -create -output $buildDir/$1-$4-$ver$5 $buildDir/$2 $buildDir/$3
}

# Move binary to the release directory
# $1: Tool Name
# $2: Source
# $3: Platform_Architecture
# $4: Suffix
move_helper () {
    echo "Moving $1 to $releaseDir/$ver/$1-$3-$ver$4"
    mv $buildDir/$2 $buildDir/$1-$3-$ver$4
}

# Linux
build_linux () {
    build_all linux 386 -386

    move_helper medhash-gen medhash-gen-386 linux_386
    move_helper medhash-chk medhash-chk-386 linux_386
    move_helper medhash-upgrade medhash-upgrade-386 linux_386
}

# macOS
build_macos () {
    if [[ "$OSTYPE" == "darwin"* ]]; then
        build_all darwin amd64 -intel
        build_all darwin arm64 -m1

        merge medhash-gen medhash-gen-intel medhash-gen-m1 macos
        merge medhash-chk medhash-chk-intel medhash-chk-m1 macos
        merge medhash-upgrade medhash-upgrade-intel medhash-upgrade-m1 macos
    else
        echo "Cannot build macOS universal binary on non-macOS host machine"
    fi
}

# Windows
build_windows () {
    build_all windows 386 -386

    move_helper medhash-gen medhash-gen-386 windows_x86 .exe
    move_helper medhash-chk medhash-chk-386 windows_x86 .exe
    move_helper medhash-upgrade medhash-upgrade-386 windows_x86 .exe
}

source config

if [ ! -d $buildDir ]; then
    echo "Creating build directory"
    mkdir -p $buildDir
fi

case $1 in
    "linux")
        build_linux
        ;;

    "macos" | "darwin")
        build_macos
        ;;

    "windows")
        build_windows
        ;;

    *)
        build_linux
        build_macos
        build_windows
        ;;
esac

echo "Done"

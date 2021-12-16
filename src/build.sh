#!/bin/bash

# MedHash Tools
# Copyright (c) 2021 GHIFARI160
# MIT License

ver="0.3.0"
buildDir="../out"
releaseDir="../dist"

buildArgs="-trimpath"
ldArgs="-s -w -X github.com/ghifari160/medhash-tools/src/common.VERSION=$ver"

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
    lipo -create -output $releaseDir/$ver/$1-$4-$ver$5 $buildDir/$2 $buildDir/$3
}

# Move binary to the release directory
# $1: Tool name
# $2: Platform
# $3: Architecture
# $4: Suffix
move_helper () {
    echo "Moving $1 to $releaseDir/$ver/$1-$2_$3-$ver$4"
    mv $buildDir/$1 $releaseDir/$ver/$1-$2_$3-$ver$4
}

# Linux
build_linux () {
    build_all linux 386 -386

    move_helper medhash-gen-386 linux 386
    move_helper medhash-chk-386 linux 386
    move_helper medhash-upgrade-386 linux 386
}

# macOS
build_macos () {
    if [ $(uname -s) == Darwin* ]; then
        build_all darwin amd64 -intel
        build_all darwin arm64 -m1

        merge medhash-gen medhash-gen-intel medhash-gen-m1 macos
        merge medhash-chk medhash-chk-intel medhash-chk-m1 macos
        merge medhash-upgrade medhash-upgrade-intel medhash-upgrade-m1 macos

        rm -f $buildDir/medhash*
    else
        echo "Cannot build for macOS on non-macOS host machine"
    fi
}

# Windows
build_windows () {
    build_all windows 386 -386

    move_helper medhash-gen-386 windows x86 .exe
    move_helper medhash-chk-386 windows x86 .exe
    move_helper medhash-upgrade-386 windows x86 .exe
}

if [ ! -d $buildDir ]; then
    echo "Creating build directory"
    mkdir -p $buildDir
fi

if [ ! -d $releaseDir/$ver ]; then
    echo "Creating release directory for v$ver"
    mkdir -p $releaseDir/$ver
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

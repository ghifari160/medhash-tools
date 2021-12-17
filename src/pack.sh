#!/bin/bash

# MedHash Tools
# Copyright (c) 2021 GHIFARI160
# MIT License

ver="0.3.0"

releaseDir="../dist"

# Pack tools into an archive
# $1: Archive format
# $2: Platform_Arch
# $3: Suffix
pack_helper () {
    archiver=""
    arg=""
    ext=""

    if [ "$1" == "gzip" ]; then
        archiver="tar"
        arg="-zcf"
        ext="tar.gz"
    elif [ "$1" == "bzip2" ]; then
        archiver="tar"
        arg="-jcf"
        ext="tar.bz2"
    elif [ "$1" == "zip" ]; then
        archiver="zip"
        ext="zip"
        arg="-rq"
    fi

    if [ "$archiver" == "tar" ]; then
        tar $arg $releaseDir/medhash-tools-$2-$ver.$ext -C $releaseDir/$ver/ medhash-{gen,chk,upgrade}-$2-$ver$3
    elif [ "$archiver" == "zip" ]; then
        (cd $releaseDir/$ver && zip $arg - medhash-{gen,chk,upgrade}-$2-$ver$3) > $releaseDir/medhash-tools-$2-$ver.$ext
    fi
}

# Pack for Linux
pack_linux () {
    if [ -f $releaseDir/$ver/medhash-gen-linux_386-$ver ]; then
        echo -n "Packing for Linux_386... "
        pack_helper gzip linux_386
        echo "done"
    else
        echo "Skipping Linux_386"
    fi
}

# Pack for macOS
pack_macos () {
    if [ -f $releaseDir/$ver/medhash-gen-macos-$ver ]; then
        echo -n "Packing for macOS... "
        pack_helper gzip macos
        echo "done"
    else
        echo "Skipping macOS"
    fi
}

# Pack for Windows
pack_windows () {
    if [ -f $releaseDir/$ver/medhash-gen-windows_x86-$ver.exe ]; then
        echo -n "Packing for Windows_x86... "
        pack_helper zip windows_x86 .exe
        echo "done"
    else
        echo "Skipping Windows_x86"
    fi
}

case $1 in
    "linux")
        pack_linux
        ;;

    "macos" | "darwin")
        pack_macos
        ;;

    "windows")
        pack_windows
        ;;

    *)
        pack_linux
        pack_macos
        pack_windows
        ;;
esac

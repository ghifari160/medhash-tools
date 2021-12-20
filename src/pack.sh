#!/bin/bash

# MedHash Tools
# Copyright (c) 2021 GHIFARI160
# MIT License

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
        tar $arg $releaseDir/$ver/medhash-tools-$2-$ver.$ext -C $pkgDir/$2 .
    elif [ "$archiver" == "zip" ]; then
        (cd $pkgDir/$2 && zip $arg - .) > $releaseDir/$ver/medhash-tools-$2-$ver.$ext
    fi
}

# Prepare binary directory
# $1: Platform_Arch
prep_bin_dir () {
    if [ ! -d $pkgDir/$1/bin ]; then
        mkdir -p $pkgDir/$1/bin
    fi
}

# Copy binary
# $1: Tool name
# $2: Platform_Architecture
# $3: Suffix
copy_helper () {
    prep_bin_dir $2

    echo "Copying $1"
    cp $buildDir/$1-$2-$ver$3 $pkgDir/$2/bin/$1$3
}

# Build license RTF
# $1: Path to RTF
# $2: Path to source
license_builder () {
    touch $1

    licenseRTF=$(cat $2)
    licenseRTF=$(echo "$licenseRTF" | sed 's/MIT License$/\\f0\\fs48MIT License\n\\f0\\fs28/g')
    licenseRTF=$(echo "$licenseRTF" | sed 's/^$/\\line\\line/g')
    licenseRTF=$(echo "$licenseRTF" | tr '\n' ' ')

    echo "{\\rtf1\\ansi\\def0 {\\fonttbl \\f0 HelveticaNeue-Light;}" > $1
    echo "$licenseRTF" >> $1
    echo "}" >> $1
}

# Create macOS component package
# $1: Platform_Arch
# $2: Suffix
pkg_helper () {
    echo "Building macOS component package"
    pkgbuild --root $pkgDir/macos \
             --identifier $identifier \
             --version $ver \
             --install-location $defaultInstallLocation_macOS \
             --scripts $pkgScriptsDir \
             --filter "(.*).DS_Store(.*)" \
             --filter "(.*).git(.*)" \
             --filter "(.*).pkg" \
             --filter Resources \
             $pkgDir/macos/medhash-tools$2

    echo "Generating license text"
    mkdir -p $pkgDir/macos/Resources/en.lproj
    license_builder $pkgDir/macos/Resources/en.lproj/LICENSE.rtf ../LICENSE

    echo "Building macOS installation package"
    productbuild --distribution Distribution.xml \
                 --package-path $pkgDir/macos \
                 --resources $pkgDir/macos/Resources \
                 $releaseDir/$ver/medhash-tools-$1-$ver$2
}

# Pack for Linux
pack_linux () {
    if [ -f $buildDir/medhash-gen-linux_386-$ver ]; then
        copy_helper medhash-gen linux_386
        copy_helper medhash-chk linux_386
        copy_helper medhash-upgrade linux_386

        cp ../LICENSE $pkgDir/linux_386/LICENSE

        echo "Packing for Linux_386... "
        pack_helper gzip linux_386
        echo "Done"
    else
        echo "Skipping Linux_386"
    fi
}

# Pack for macOS: Tarball
pack_macos_tarball () {
    pack_helper gzip macos
}

# Pack for macOS: Package
pack_macos_pkg () {
    mkdir -p $pkgDir/macos/medhash-tools

    # Move bin dir into medhash-tools dir
    mv $pkgDir/macos/bin $pkgDir/macos/medhash-tools/

    mv $pkgDir/macos/LICENSE $pkgDir/macos/medhash-tools/
    cp $pkgScriptsDir/uninstall.sh "$pkgDir/macos/medhash-tools/Uninstall MedHash Tools"

    pkg_helper macos .pkg

    # Move bin dir out of medhash-tools dir
    mv $pkgDir/macos/medhash-tools/bin $pkgDir/macos/
}

# Pack for macOS
pack_macos () {
    if [ -f $buildDir/medhash-gen-macos-$ver ]; then
        copy_helper medhash-gen macos
        copy_helper medhash-chk macos
        copy_helper medhash-upgrade macos

        cp ../LICENSE $pkgDir/macos/LICENSE

        echo "Packing for macOS... "
        pack_macos_tarball
        pack_macos_pkg
        echo "Done"
    else
        echo "Skipping macOS"
    fi
}

# Pack for Windows
pack_windows () {
    if [ -f $buildDir/medhash-gen-windows_x86-$ver.exe ]; then
        copy_helper medhash-gen windows_x86 .exe
        copy_helper medhash-chk windows_x86 .exe
        copy_helper medhash-upgrade windows_x86 .exe

        cp ../LICENSE $pkgDir/windows_x86/LICENSE

        echo "Packing for Windows_x86... "
        pack_helper zip windows_x86 .exe
        echo "Done"
    else
        echo "Skipping Windows_x86"
    fi
}

source config

if [ ! -d $releaseDir ]; then
    mkdir -p $releaseDir/$ver
fi

if [ ! -d $pkgDir ]; then
    mkdir -p $pkgDir
fi

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

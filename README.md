# MedHash Tools

Simple tools for verifying media file integrity.

## Usage

Generating medhash

``` shell
medhash-gen [target dir]
```

Verifying medhash

``` shell
medhash-chk [target dir]
```

Upgrading medhash from previous versions

``` shell
medhash-upgrade [target dir]
```

**Note:** specifying a target directory will run the appropriate tool in scoped mode.
In scoped mode, both `medhash-gen` and `medhash-chk` will generate and verify hashes for media in the target directory.
`medhash-upgrade` will _enter_ the directory and attempt to upgrade medhash file to the current format.

## Building

Building requires a working Go 1.15+ installation.

Clone and enter the repository

``` shell
git clone https://github.com/ghifari160/medhash-tools
cd medhash-tools/src
```

Build the binaries (Linux and macOS)

``` shell
./build.sh
```

Build the binaries (Windows)

``` shell
build.bat
```

The binaries are stored in `dist/{VERSION}`.

You can specificy a build target if necessary

``` shell
./build.sh linux
```

**Note:** building the macOS target requires a macOS host machine.
The macOS target builds a universal binary.
This is done by building an Intel binary and an M1 binary, then merging the two with `lipo`.

To clean the build environment, run `clean.sh` in Linux and macOS, and `clean.bat` in Windows.

## Packing for release

Packing requires a working `tar` and `zip` installationg.

Enter the source directory and run packing script

``` shell
./pack.sh
```

The packing script will skip a platform if its archives are missing.

You can specify a packing target if necessary

``` shell
./pack.sh linux
```

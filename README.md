# MedHash Tools

A collection of tools for dealing with media hash (_medhash_) files.

## Building

Building requires a working Go 1.15+ installation and Make.

Clone and enter the repository

``` shell
git clone https://github.com/ghifari160/medhash-tools
cd medhash-tools
```

Build the binaries

``` shell
make
```

The binaries are stored in `out/build` relative to the current directory. You can override the
target directory by setting the `BUILD_DIR` environment variable.

### Cross Compilation

You can compile the binaries for Linux, macOS, or Windows by replacing the `<target platform>` with
`linux` for Linux, `macos` or `darwin` for macOS, and `windows` for Windows.

``` shell
make <target platform>
```

You can set the `GOOS` and `GOARCH` environment variables to cross compile to other platforms.

## Installing

Installation requires a working Go 1.15+ installation.

You can install individual tools with

``` shell
go get -u github.com/ghifari160/src/<tool name>
```

Alternatively, you can install _all_ tools with

``` shell
go get -u github.com/ghifari160/src/...
```

## Packaging

Setup the directory structure

``` shell
mkdir -p out/packaging out/release
```

Prepare the files for packaging

``` shell
cd out/packaging && \
cp ../build/* . && \
cp ../../LICENSE . && \
cp ../../README.md . && \
cd ../../
```

Run ReMan's packaging tool (or, alternatively, [package without ReMan](#packaging-without-reman))

``` shell
cd out/packaging && \
package ../release medhash-tools-macos 0.2.0 && \
cd ../../
```

**Note:** Replace `macos` with the target platform `linux` for Linux, `macos` for macOS, and `win`
for Windows) and `0.2.0` with medhash-tools' version.

The release packages are stored in `out/release` relative to the current directory.

### Packaging without ReMan

Create the release packages

``` shell
cd out/packaging && \
zip -r -X ../release/medhash-tools-macos-v0.2.0.zip . -x "*.git*" && \
tar --exclude="*.git*" -zcvf ../release/medhash-tools-v0.2.0.tar.gz . && \
tar --exclude="*.git*" -jcvf ../release/medhash-tools-v0.2.0.tar.bz2 . && \
cd ../../
```

**Note:** Replace `macos` with the target platform `linux` for Linux, `macos` for macOS, and `win`
for Windows) and `0.2.0` with medhash-tools' version.

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

**Note:** specifying a target directory will run the appropriate tool in scoped mode. In scoped
mode, both `medhash-gen` and `medhash-chk` will generate and verify hashes for media in the target
directory. `medhash-upgrade` will _enter_ the directory and attempt to upgrade medhash file to the
current format.

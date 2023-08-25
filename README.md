# MedHash Tools

Simple tools for verifying media file integrity.

## Rewrite

This project has been rewritten.

This project was repurposed from a prototype.
Rewriting this project meant that it could be designed in a more logical way.

Code from before the rewrite are available under the
[legacy](https://github.com/Ghifari160/medhash-tools/tree/legacy) branch.

The legacy branch will be removed after the release of v1.0.0.

### Performance comparison

| Version | Duration (s) | Rate (GiB/s) | Payload Size (GiB) |
|---------|--------------|--------------|--------------------|
| v0.5.0  | 227.5550     | 0.0440       | 10.0000            |
| v0.6.0  | 171.0000     | 0.0585       | 10.0000            |

**Note:** v0.6.0 is the rewrite.

These figures are collected using the program in [`_bench`](_bench),
built using the following parameters

``` text
go build .
```

and executed using the follow parameters.

``` text
./_bench -q 5 -s 10G -r <path to store report> <-c command [args...]...>
```

The program ran on a 2019 MacBook Pro (MacBookPro15,1) with Intel Core i7-9750H,
16 GB 2400 MHz DDR4, Radeon Pro 560X 4 GB, macOS Ventura 13.4.1, Go 1.20.6.

## Usage

Generating medhash

``` shell
medhash gen [target dir]
```

Verifying medhash

``` shell
medhash chk [target dir]
```

Upgrading medhash from previous versions

``` shell
medhash upgrade [target dir]
```

## Building

Building requires a working Go 1.20+ installation.

Clone and enter the repository

``` shell
git clone https://github.com/ghifari160/medhash-tools
cd medhash-tools
```

Build the binaries

``` shell
go build . -o out/bin/medhash
```

You can specificy a build target if necessary

``` shell
GOOS=linux GOARCH=386 go build . -o out/bin/medhash
```

In the past, building the macOS target requires a macOS host machine.
This is because the macOS target builds a universal binary using a macOS-specific tool (`lipo`).
This is no longer the case.
The universal binary for macOS is now generated with
[randall77/makefat](https://github.com/randall77/makefat).

## Release

Release binaries are automatically built with GitHub Actions.

Artifacts are automatically uploaded to GHIFARI160's download server.
They are available for download from
`https://projects.gassets.space/medhash-tools/{version}/medhash-{os_arch}.tar.gz`.

## License

MedHash Tools is distributed under the terms of the [MIT License](LICENSE).

<!-- markdownlint-disable MD024 -->

# Changelog

All notable changes in MedHash Tools will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).
This project attempts to adhere to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.6.1] - 2023-08-30

### Changed

- Upgraded `gopkg.in/yaml.v3` to v3.0.0.

## [0.6.0] - 2023-08-24

### Added

- Added support for multiple target directories.
  Pass each target directory as positional arguments to each tool.
- Added `--ignore` parameter to `gen` and `upgrade`.
- Added `--force` parameter to `upgrade`.
- Added preset parameters to `gen`, `chk`, and `upgrade`.
  `gen` and `upgrade` will generate hashes using the preset.
  `chk` will only attempt to verify hashes in the preset.
  Non-existent hashes are ignored.
  - Added `--default` preset parameter.
  - Added `--all` preset parameter.
- Added hash algorithm parameters to `gen`, `chk`, and `upgrade`.
  `gen` and `upgrade` will only generate hashes using the specified algorithm.
  `chk` will only verify hashes of the specified algorithm.
  Non-existent hashes are ignored.
  - Added `--xxh3` parameter.
  - Added `--sha3` parameter.
  - Added `--sha256` parameter.
  - Added `--sha1` parameter.
  - Added `--md5` parameter.
- Added `version` command.
- Added colored status labels.
  It works on Windows too!

### Changed

- `medhash-*` commands are now subcommands.
  Pass them as the first parameter to `medhash`.
  For example, `medhash-gen` is now `medhash gen`.
- Updated MedHash Manifest Specification to v0.5.0.
  `upgrade` now upgrades Manifest v0.4.0 without `--force`.
  Manifest v0.5.0 can be regenerated with `--force`.
- Default preset now only generates [XXH3_64](https://xxhash.com) hash.
  For cryptographic use, specify the appropriate algorithm or use `--all`.

### Deprecated

- Manifest field `media.hash.sha3-256` is now deprecated.
  Use `media.hash.sha3` instead.
  For compatibility, MedHash Tools populate SHA3 hashes into both fields.
  `chk` prefers `media.hash.sha3` to `media.hash.sha3-256`, meaning that it will attempt to verify
  the former before attempting to verify the latter.

### Removed

- Removed `-v` parameter.
  Previously, this parameter enables verbose mode.

### Fixed

- Fixed [#1](https://github.com/Ghifari160/medhash-tools/issues/1): invalid memory address when
  files are not found.

### Security

- Upgraded `golang.org/x/crypto` to v0.11.0.

## [0.5.0] - 2021-12-29

### Added

- `medhash-gen` and `medhash-chk` will now generate and verify SHA3-256 hash for every Media.

### Changed

- Manifest Specification has been updated to v0.4.0.
  v0.3.0 Manifests are compatible and will continue to work with this version of MedHash Tools.
  However, you will be prompted to rerun `medhash-gen` to upgrade the Manifest.
- `medhash-gen` now credits itself for Manifest generation in the generator field.
- `medhash-chk` now prints the Manifest generator in verbose mode (`-v` flag).
- `medhash-chk` now validates the Manifest version.
  Compatible Manifests will prompt the user to rerun `medhash-gen` for upgrade.
  Non-compatible Manifests will prevent the tool from working.

### Deprecated

- SHA1 support is now deprecated.
- MD5 support is now deprecated.

### Fixed

- Fixed a bug where `medhash-gen` and `medhash-check` would display paths with `/` separator on Windows.
- Optimized `medhash-gen` and `medhash-chk`.

## [0.4.0] - 2021-12-17

Now in color!

### Added

- Added verbose flag to `medhash-gen` and `medhash-chk`.
  Both tools will no longer print target and working directory paths without the verbose flag.
- Added version flag to `medhash-gen` and `medhash-chk`.
  Both tools will print the toolset version and exit with the flag set.
- Added manifest flag to `medhash-gen` and `medhash-chk`.
  In `medhash-gen`, the Manifest will be stored in the specified path.
  In `medhash-chk`, the Manifest will be read from the specified path.
- Added file flag to `medhash-chk`.
  Verify only the specified file.
- Added status color to `medhash-gen` and `medhash-chk`.
  `medhash-gen` will color the sanity check status for each Media appropriately.
  `medhash-chk` will color the verification status for each Media appropriately.
- macOS binaries are now universal binaries.
  They will work in both Intel and M1 Macs.

### Changed

- Changed target directory behavior in `medhash-gen`.
  Previously, the Manifest will always be stored in the current working directory, with Media paths relative to the current working directory.
  Now, the Manifest will be stored in either the target directory or the path specified in `-manifest`. The Media paths will always be relative to the target directory.
- Rewrote internal library.
  The `data` internal library has been rewritten to `medhash` internal library.
  The original `data` library is deprecated and will be removed in the future.
- Optimized `medhash-gen` and `medhash-chk`.
  Both tools now are now more efficient and behave more similarly to each other.
  Errors will be printed and the tools will now attempt to recover before exiting.

### Deprecated

- Deprecated `data` internal library.

## [0.3.0] - 2020-11-12

The internal logic of medhash-tools did not use buffers.
As a result, hashing large media can consume large amounts of memory (up to triple the media size).
This update should fix that issue.

### Changed

- Optimized hash generation with buffered I/O.

## [0.2.0] - 2020-11-11

Manage your media file hashes better with medhash-tools!
medhash-tools handles the creation (`medhash-gen`) and verification (`medhash-chk`) of media hashes.
If you are a user of the legacy medhash-tools from [here](https://github.com/Ghifari160/infra/tree/main/renderer),
you can upgrade your legacy medhash files with `medhash-upgrade`.
See more usage details [here](https://github.com/Ghifari160/medhash-tools/tree/718609ec57f7366b99c718fe91d5959345c9fbfa#usage).

### Added

- Added medhash upgrade tool (`medhash-upgrade`).

### Changed

- Optimized medhash generation with `medhash-gen`.
- Optimized medhash verification with `medhash-chk`.
- Switched medhash format to JSON format.

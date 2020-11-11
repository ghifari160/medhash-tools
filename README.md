# MedHash Tools

A collection of tools for dealing with media hash (_medhash_) files.

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

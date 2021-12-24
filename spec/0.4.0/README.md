# MedHash Manifest Specification v0.4.0

MedHash Tools stores media hashes in a _Manifest_.
This documents the specification of the _Manifest_ format.

## File Format

The _Manifest_ is a JSON text document.
It must be encoded as UTF-8.

## `version` Field

The `version` field denotes the _Manifest Specification_ format.
This version number may be different than the toolset version as some tool updates do not require any changes to the Manifest Specification.

## `generator` Field

The `generator` field denotes the _Generator_ of the _Manifest_.
This field is optional.

## `media` Field

The `media` field is an array of [Media](#media-object) objects.
Each entry describes a media and its hashes.
Hashes are provided in multiple algorithms for compatibility reasons.

### Media Object

The media object describes the _Media_ and its hashes through the use of a _Hash Container_.

| Field  | Type                        | Description                                           |
|--------|-----------------------------|-------------------------------------------------------|
| `path` | string                      | Path to media file, relative to the Manifest location |
| `hash` | [Hash](#hash-object) object | Hash Container                                        |

**Note:** For the purposes of the _Manifest_, `path` must use `/` as the path separator.

### Hash Object

The hash object describes a _Hash Container_.
A single Hash Container contains hashes of the same Media, each calculated from a different algorithm.
This is done for compatibility reasons.

| Field      | Type       | Description                |
|------------|------------|----------------------------|
| `sha256`   | string     | SHA256 hash                |
| `sha3-256` | string     | SHA3-256 hash              |
| ~~`sha1`~~ | ~~string~~ | ~~SHA1 hash~~ _Deprecated_ |
| ~~`md5`~~  | ~~string~~ | ~~MD5 hash~~ _Deprecated_  |

**Note:** SHA1 and MD5 are both _deprecated_.
These algorithms were included and supported in the _Specification_ due to compatibility reasons.
For that same reason, tools must maintain support for these hashes until they are _removed_ from the specification.
SHA3-256 is now supported.
For the most part, SHA3-256 is a drop-in replacement of SHA256.

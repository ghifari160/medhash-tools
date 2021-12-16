# MedHash Manifest Specification v0.3.0

MedHash Tools stores media hashes in a Manifest.
This documents the specification of the Manifest format.

## File Format

The Manifest is a JSON text document.
It must be encoded as UTF-8.

## `version` Field

The `version` field denotes the _Manifest Specification_ format.
The version number must follow the Semantic Versioning format.
This version number may be different than the toolset version as some tool updates do not require any changes to the Manifest Specification.

## `media` Field

The `media` field is an array of [Media](#media-object) objects.
Each entry describes a media and its hashes.
Hashes are provided in multiple algorithms for compatibility reasons.

### Media Object

The media object describes the Media and its hashes through the use of a Hash Container.

| Field  | Type                        | Description                                           |
|--------|-----------------------------|-------------------------------------------------------|
| `path` | string                      | Path to media file, relative to the Manifest location |
| `hash` | [Hash](#hash-object) object | Hash Container                                        |

### Hash Object

The hash object describes a Hash Container.
A single Hash Container contains hashes of the same Media, each calculated from a different algorithm.
This is done for compatibility reasons.

| Field    | Type   | Description |
|----------|--------|-------------|
| `sha256` | string | SHA256 hash |
| `sha1`   | string | SHA1 hash   |
| `md5`    | string | MD5 hash    |

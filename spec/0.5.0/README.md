# MedHash Manifest Specification v0.5.0

MedHash Tools stores media hashes in a _Manifest_.
This documents the specification of the _Manifest_ format.

## File Format

The _Manifest_ is a JSON text document.
It must be encoded as UTF-8.

## Fields

| Field       | Type        | Notes                                    |
|-------------|-------------|------------------------------------------|
| `version`   | string      | Required. Manifest Specification format. |
| `generator` | string      | Optional. Generator of the Manifest.     |
| `media`     | \[\][Media] | Required. Array of Media objects.        |

### `version` field

The `version` field denotes the _Manifest Specification_ format.
This version number may be different than the toolset version as some tool updates do not require
any changes to the Manifest Specification.

### `generator` field

The `generator` field denotes the _Generator_ of the _Manifest_.
This field is optional.

### `media` field

The `media` field is an array of [Media](#media-object) objects.
Each entry describes a media and its hashes.
Hashes are provided in multiple algorithms for compatibility reasons.

### Media object

The media object describes the _Media_ and its hashes through the use of a _Hash Container_.

| Field  | Type   | Notes                                                           |
|--------|--------|-----------------------------------------------------------------|
| `path` | string | Required. Path to media file, relative to the Manifest location |
| `hash` | [Hash] | Required. Hash Container                                        |

**Note:** For the purposes of the _Manifest_, `path` must use `/` as the path separator.

### Hash object

The hash object describes a _Hash Container_.
A single Hash Container contains one or more hashes of the same Media.

| Field          | Type   | Notes                        |
|----------------|--------|------------------------------|
| `xxh3`         | string | Preferred. xxHash (XXH3_64). |
| `sha256`       | string | SHA256 hash.                 |
| `sha3`         | string | SHA3-256 hash.               |
| ~~`sha3-256`~~ | string | Deprecated: use `sha3`.      |
| `sha1`         | string | SHA1 hash.                   |
| `md5`          | string | MD5 hash.                    |

**Notes:**

- [xxHash] (XXH3_64) is now the preferred hash.
- [MedHash Manifest Specification v0.4.0] introduced SHA3-256 support under the `sha3-256` field.
  SHA3-256 hash has been moved to `sha3`.
  `sha3-256` is deprecated.
- SHA1 and MD5 were _previously deprecated_ in MedHash Manifest Specification v0.4.0.
  This is no longer the case.
  While the use of both hashes should be discouraged, many tools (notably Git) still depend on
  these hashes.

## Presets

A Preset is a previously determined set of hash algorithms.
When generating and upgrading Manifests, tools should generate hashes using _only_ the algorithms
contained by the preset.
When checking Manifests, tools should _attempt_ to verify hashes using _only_ the algorithms
contained by the preset.
If no compatible hashes are present in the Manifest, the Media passes the check.

A number of presets are defined in this specification.
Tools may implement additional presets.

### Default preset

This preset must be the default behavior of tools.

This preset _only_ contains `xxh3`.

### All preset

This preset contains _all_ supported hash algorithms.

### Legacy preset

This preset contains the hash algorithms supported by the legacy MedHash Tools:

- `sha3-256`
- `sha256`
- `sha1`
- `md5`

[Media]: #media-object
[Hash]: #hash-object
[xxHash]: https://xxhash.com/
[MedHash Manifest Specification v0.4.0]: https://github.com/Ghifari160/medhash-tools/tree/0b85f13fbabd6e724efe4ea872e08b60ef48da89/spec/0.4.0

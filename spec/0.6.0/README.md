# MedHash Manifest Specification v0.6.0

MedHash Tools stores media hashes in a _Manifest_.
This documents the specification of the _Manifest_ format.

## File Format

The _Manifest_ is a JSON text document.
It must be encoded as UTF-8.

## Reproducibility

In most context, the order of JSON fields do not matter.
However, they do matter when considering reproducibility.
To maintain the reproducibility of the Manifest, all fields must in the same order as they appear
in this Specification.
Additional care must be taken for some fields, as noted in their section.

## Fields

| Field       | Type        | Required? | Notes                          |
|-------------|-------------|-----------|--------------------------------|
| `version`   | string      | Yes       | Manifest Specification format. |
| `generator` | string      | No        | Generator of the Manifest.     |
| `media`     | \[\][Media] | Yes       | Array of Media objects.        |
| `signature` | [Signature] | No        | Signature of the Manifest.     |

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

To maintain reproducibility, the contents of the array must be sorted by their path in ascending
order.

### Media object

The media object describes the _Media_ and its hashes through the use of a _Hash Container_.

| Field  | Type   | Required? | Notes                                                           |
|--------|--------|-----------|-----------------------------------------------------------------|
| `path` | string | Yes       | Required. Path to media file, relative to the Manifest location |
| `hash` | [Hash] | Yes       | Required. Hash Container                                        |

**Note:** For the purposes of the _Manifest_, `path` must use `/` as the path separator.

### Signature object

The signature object describes the signature of the Manifest.
All fields are optional, but only specified fields are verified (see [Signature and verification]).

| Field      | Type   | Required? | Notes                         |
|------------|--------|-----------|-------------------------------|
| `ed25519`  | string | No        | Preferred. Ed25519 signature. |
| `minisign` | string | No        | Minisign signature.           |
| `pgp`      | string | No        | PGP signature.                |

### Hash object

The hash object describes a _Hash Container_.
A single Hash Container contains one or more hashes of the same Media.

| Field          | Type   | Required? | Notes                        |
|----------------|--------|-----------|------------------------------|
| `xxh3`         | string | No        | Preferred. xxHash (XXH3_64). |
| `sha512`       | string | No        | SHA512 hash.                 |
| `sha256`       | string | No        | SHA256 hash.                 |
| `sha3`         | string | No        | SHA3-256 hash.               |
| ~~`sha3-256`~~ | string | No        | Deprecated: use `sha3`.      |
| `sha1`         | string | No        | SHA1 hash.                   |
| `md5`          | string | No        | MD5 hash.                    |

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

- SHA3-256
- SHA256
- SHA1
- MD5

### Maven preset

This preset contains the hash algorithms utilized by Maven:

- SHA512
- SHA256
- SHA1
- MD5

## Signature and verification

All signature types are optional and must not depend on each other.
Multiple signatures can exist for the same Manifest, provided that they are of different
algorithms.

The example below is **valid**

``` json
{
  "signature": {
    "ed25519": "ED25519_SIGNATURE",
    "minisign": "MINISIGN_SIGNATURE",
    "pgp": "PGP_SIGNATURE"
  }
}
```

but this one is **not**.

``` json
{
  "signature": {
    "pgp": "PGP_SIGNATURE",
    "pgp": "ANOTHER_PGP_SIGNATURE"
  }
}
```

Ed25519 is the preferred algorithm for Manifest signature.

### Generating signatures

When generating the Manifest signature, the Manifest **must not** contain a `signature` field.
The generated signature is that of the contents of the Manifest as they would be stored on disk.
All generated signatures must then be [added to the Manifest](#signature-object).

For additional compatibility, it is recommended to store the Minisign signature in its native
format, and the PGP signature in a detached, ASCII-armored signature file.
As an example, a Manifest with the default name (`medhash.json`) would have its Minisign signature
and its PGP signature stored inside the Manifest and in `medhash.json.minisig` and
`medhash.json.asc`.

### Verifying signatures

When one signature is present, it is verified against the Manifest with the `signature` field
stripped.
When multiple signatures are present, the preferred signature is verified against the stripped
Manifest.
The user should be able to specific which signature to verify.
Verifying multiple signatures should be supported, but it **must not** be the default.

[Media]: #media-object
[Signature]: #signature-object
[Hash]: #hash-object
[Signature and verification]: #signature-and-verification
[xxHash]: https://xxhash.com/
[MedHash Manifest Specification v0.4.0]: https://github.com/Ghifari160/medhash-tools/tree/0b85f13fbabd6e724efe4ea872e08b60ef48da89/spec/0.4.0

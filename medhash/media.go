package medhash

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"
)

// Media stores metadata about the media.
type Media struct {
	Path string `json:"path"`
	Hash Hash   `json:"hash"`
}

// Check checks hashes for media.
// Hashes for the media are verified at the same time.
func (media Media) Check(config Config) error {
	return mediaErrOrNil(config, media, chkHash(config, media))
}

// Add adds media to man and generates the appropriate hashes as configured.
// Add also sorts the man.Media slice.
func (man *Manifest) Add(media string) error {
	med, err := genHash(man.Config, media)
	if err != nil {
		return err
	}

	man.Media = append(man.Media, med)
	man.sortMedia()

	return nil
}

// Check checks hashes for media.
// Hashes for the media are verified at the same time.
func (man *Manifest) Check(media string) error {
	med, err := man.searchMedia(media)
	if err != nil {
		return err
	}
	return med.Check(man.Config)
}

// sortMedia sorts man.Media.
// Sorting is done with a stable sorting algorithm, meaning that insertion order is preserved for
// equal elements.
// sortMedia abstracts the actual implementation detail, allowing us to change or otherwise modify
// the implementation in the future without breaking compatibility.
func (man *Manifest) sortMedia() {
	slices.SortStableFunc(man.Media, mediaCmp)
}

// searchMedia searches for media in man.Media using binary search.
// man.Media must be sorted for binary search to work (see sortMedia).
// searchMedia abstracts the actual implementation detail of the searching logic, allowing us to
// change or otherwise modify the implementation in the future without breaking compatibility.
func (man *Manifest) searchMedia(media string) (med Media, err error) {
	target := Media{Path: media}
	index, found := slices.BinarySearchFunc(man.Media, target, mediaCmp)
	if !found {
		err = fmt.Errorf("media %s not in manifest", media)
	} else {
		med = man.Media[index]
	}
	return
}

// mediaCmp compares a.Path and b.Path.
func mediaCmp(a, b Media) int {
	return strings.Compare(a.Path, b.Path)
}

// mediaErr wraps any error for a Media.
type mediaErr struct {
	path string
	err  error
}

func (err mediaErr) Error() string {
	return err.path + ": " + err.err.Error()
}

func (err mediaErr) Unwrap() error {
	return err.err
}

// mediaErrOrNil wraps err with mediaErr, ignoring nil err.
// That is, mediaErrOrNil returns nil if err is nil.
func mediaErrOrNil(config Config, media Media, err error) error {
	if err == nil {
		return nil
	} else {
		return mediaErr{path: filepath.Join(config.Dir, media.Path), err: err}
	}
}

// Hash stores each hash of a Media.
type Hash struct {
	XXH3   string `json:"xxh3,omitempty"`
	SHA512 string `json:"sha512,omitempty"`
	SHA256 string `json:"sha256,omitempty"`
	SHA3   string `json:"sha3,omitempty"`
	// Deprecated: use SHA3.
	SHA3_256 string `json:"sha3-256,omitempty"`
	SHA1     string `json:"sha1,omitempty"`
	MD5      string `json:"md5,omitempty"`
}

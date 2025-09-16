package medhash

import (
	"fmt"
	"slices"
	"strings"
)

// Media stores metadata about the media.
type Media struct {
	Path string `json:"path"`
	Hash Hash   `json:"hash"`
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

	err = chkHash(man.Config, med)
	if err != nil {
		return err
	}
	return nil
}

func (man *Manifest) sortMedia() {
	slices.SortStableFunc(man.Media, mediaCmp)
}

func (man *Manifest) searchMedia(media string) (med Media, err error) {
	const errMsg = "media %s not in manifest"

	if !slices.ContainsFunc(man.Media, func(e Media) bool {
		return e.Path == media
	}) {
		err = fmt.Errorf(errMsg, media)
	} else {
		target := Media{Path: media}
		index, found := slices.BinarySearchFunc(man.Media, target, mediaCmp)
		if !found {
			err = fmt.Errorf(errMsg, media)
		} else {
			med = man.Media[index]
		}
	}
	return
}

func mediaCmp(a, b Media) int {
	return strings.Compare(a.Path, b.Path)
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

package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"io"
	"os"
)

func Sha256WithFile(filename string) (string, error) {
	h := sha256.New()
	return SumWithFile(h, filename)
}

func SumWithFile(h hash.Hash, filename string) (string, error) {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0750)
	if err != nil {
		return "", err
	}
	defer f.Close()
	return Sum(h, f)
}

func Sum(h hash.Hash, r io.Reader) (string, error) {
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

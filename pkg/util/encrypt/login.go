package encrypt

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"hash"

	"golang.org/x/crypto/pbkdf2"
)

const (
	defaultSaltLen    = 64
	defaultIterations = 10000
	defaultKeyLen     = 128
)

var defaultHashFunction = sha512.New

type Options struct {
	SaltLen      int
	Iterations   int
	KeyLen       int
	HashFunction func() hash.Hash
}

func generateSalt(length int) []byte {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	salt := make([]byte, length)
	if _, err := rand.Read(salt); err != nil {
		return salt
	}
	for key, val := range salt {
		salt[key] = alphanum[val%byte(len(alphanum))]
	}
	return salt
}

func Encode(rawPw string, options *Options) (string, string) {
	if options == nil {
		salt := generateSalt(defaultSaltLen)
		encodedPwd := pbkdf2.Key([]byte(rawPw), salt, defaultIterations, defaultKeyLen, defaultHashFunction)
		return string(salt), hex.EncodeToString(encodedPwd)
	}
	salt := generateSalt(options.SaltLen)
	encodedPwd := pbkdf2.Key([]byte(rawPw), salt, options.Iterations, options.KeyLen, options.HashFunction)
	return string(salt), hex.EncodeToString(encodedPwd)
}

func Verify(rawPw string, salt string, encodedPw string, options *Options) bool {
	if options == nil {
		return encodedPw == hex.EncodeToString(pbkdf2.Key([]byte(rawPw), []byte(salt), defaultIterations, defaultKeyLen, defaultHashFunction))
	}
	return encodedPw == hex.EncodeToString(pbkdf2.Key([]byte(rawPw), []byte(salt), options.Iterations, options.KeyLen, options.HashFunction))
}

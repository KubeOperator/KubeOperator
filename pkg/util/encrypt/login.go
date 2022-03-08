package encrypt

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"hash"

	"golang.org/x/crypto/pbkdf2"
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
	if length < 1 {
		length = 1
	}
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
		salt := generateSalt(64)
		encodedPwd := pbkdf2.Key([]byte(rawPw), salt, 10000, 128, defaultHashFunction)
		return string(salt), hex.EncodeToString(encodedPwd)
	}
	salt := generateSalt(options.SaltLen)
	encodedPwd := pbkdf2.Key([]byte(rawPw), salt, options.Iterations, options.KeyLen, options.HashFunction)
	return string(salt), hex.EncodeToString(encodedPwd)
}

func Verify(rawPw []byte, salt []byte, encodedPw []byte, options *Options) bool {
	defer func() {
		for i := range rawPw {
			rawPw[i] = 0
		}
		for i := range salt {
			salt[i] = 0
		}
		for i := range encodedPw {
			encodedPw[i] = 0
		}
	}()
	if options == nil {
		return string(encodedPw) == hex.EncodeToString(pbkdf2.Key(rawPw, salt, 10000, 128, defaultHashFunction))
	}
	return string(encodedPw) == hex.EncodeToString(pbkdf2.Key(rawPw, salt, options.Iterations, options.KeyLen, options.HashFunction))
}

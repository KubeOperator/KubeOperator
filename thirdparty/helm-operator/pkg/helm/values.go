package helm

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/ghodss/yaml"
)

// Values represents a collection of (Helm) values.
// We define our own type to avoid working with two `chartutil`
// versions.
type Values map[string]interface{}

// YAML encodes the values into YAML bytes.
func (v Values) YAML() ([]byte, error) {
	b, err := yaml.Marshal(v)
	return b, err
}

// Checksum calculates and returns the SHA256 checksum of the YAML
// encoded values.
func (v Values) Checksum() string {
	b, _ := v.YAML()

	hasher := sha256.New()
	hasher.Write(b)
	return hex.EncodeToString(hasher.Sum(nil))
}

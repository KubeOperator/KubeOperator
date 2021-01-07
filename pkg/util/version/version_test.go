package version

import (
	"fmt"
	"testing"
)

func TestIsNewerThan(t *testing.T) {
	v1 := "v10.2.14"
	v2 := "v10.2.13"

	r := IsNewerThan(v1, v2)
	fmt.Println(r)
}

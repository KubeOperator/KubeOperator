package grafana

import (
	"fmt"
	"testing"
)

func TestNewDashboard(t *testing.T) {
	dash := NewDashboard("test")
	fmt.Println(dash)
}

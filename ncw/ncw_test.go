package ncw

import (
	"testing"
)

func TestLogin(t *testing.T) {
	err := Login("LXTY0I", "Alexandre@2021")
	if err != nil {
		t.Fatal(err)
	}
}

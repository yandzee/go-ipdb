package tests

import (
	"testing"

	"github.com/yandzee/go-ipdb"
)

func TestSmth(t *testing.T) {
	country, err := ipdb.LookupString("31.146.63.220")
	if err != nil || country != "GE" {
		t.Fatalf("Expected GE, got: '%v' (err: %v)", country, err)
	}
}

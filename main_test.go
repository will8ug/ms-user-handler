package main

import (
	"testing"
)

func TestInitGraphClient(t *testing.T) {
	client, err := initGraphClient()
	if err != nil || client == nil {
		t.Fatalf(`initGraphClient() failed with %q`, err)
	}
}
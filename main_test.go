package main

import (
	"testing"
)

func TestInitGraphClient(t *testing.T) {
	tc := &TenantCredential{"tenantId", "clientId", "username", "pwd"}
	client, err := initGraphClient(tc)
	if err != nil || client == nil {
		t.Fatalf(`initGraphClient() failed with %q`, err)
	}
}
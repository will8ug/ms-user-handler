package main

import (
	"os"
	"testing"
)

func TestParseArguments(t *testing.T) {
	os.Args = append([]string{os.Args[0]}, 
		"-dryrun=false", 
		"-extappid=test-ext-app-id",
		"-tid=test-tenant-id",
		"-cid=test-client-id",
		"-csec=test-client-secret",
	)
	t.Logf("\nos args = %v\n", os.Args)

	parseArguments()
	if isDryRun {
		t.Log(`isDryRun should be false`)
		t.Fail()
	}
	if b2cExtensionAppId != "test-ext-app-id" {
		t.Error("Incorrect b2cExtensionAppId")
	}
	if tenantCredential == nil || 
		tenantCredential.clientId != "test-client-id" ||
		tenantCredential.clientSecret != "test-client-secret" ||
		tenantCredential.tenantId != "test-tenant-id" {
		t.Error("Incorrect tenant credential")
	}
}

func TestInitGraphClient(t *testing.T) {
	tenantCredential = &TenantCredential{"tenantId", "clientId", "clientSecret"}
	client, err := initGraphClient()
	if err != nil || client == nil {
		t.Errorf(`initGraphClient() failed with %q`, err)
	}
}
package main

import (
	"flag"
	"log"

	azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
)

type TenantCredential struct {
	tenantId string
	clientId string
	username string
	password string
}

func main() {
	tenantCredential := parseArguments()
	log.Println(*tenantCredential)

	graphClient, err := initGraphClient(tenantCredential)
	log.Println(graphClient)
	log.Println(err)
}

func parseArguments() *TenantCredential {
	tenantId := flag.String("tid", "tenant_id", "Tenant name")
	clientId := flag.String("cid", "client_id", "Client ID")
	username := flag.String("user", "username", "Username")
	password := flag.String("pwd", "password", "Password")
	flag.Parse()
	
	return &TenantCredential{*tenantId, *clientId, *username, *password}
}

func initGraphClient(tc *TenantCredential) (*msgraphsdk.GraphServiceClient, error) {
	cred, _ := azidentity.NewUsernamePasswordCredential(
		tc.tenantId,
		tc.clientId,
		tc.username,
		tc.password,
		nil,
	)

	graphClient, err := msgraphsdk.NewGraphServiceClientWithCredentials(cred, []string{"User.Read"})
	return graphClient, err
}
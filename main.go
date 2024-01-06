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
	clientSecret string
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
	clientSecret := flag.String("csec", "client_secret", "Client secret")
	flag.Parse()
	
	return &TenantCredential{*tenantId, *clientId, *clientSecret}
}

func initGraphClient(tc *TenantCredential) (*msgraphsdk.GraphServiceClient, error) {
	cred, _ := azidentity.NewClientSecretCredential(
		tc.tenantId,
		tc.clientId,
		tc.clientSecret,
		nil,
	)

	graphClient, err := msgraphsdk.NewGraphServiceClientWithCredentials(cred, []string{"User.Read"})
	return graphClient, err
}
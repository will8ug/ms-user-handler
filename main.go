package main

import (
	"context"
	"flag"
	"log"

	azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	graphusers "github.com/microsoftgraph/msgraph-sdk-go/users"
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
	if err != nil {
		log.Panicln(err)
	}

	requestTop := int32(2)
	requestParameters := &graphusers.UsersRequestBuilderGetQueryParameters{
		Top: &requestTop,
	}
	config := &graphusers.UsersRequestBuilderGetRequestConfiguration{
		QueryParameters: requestParameters,
	}

	users, err := graphClient.Users().Get(context.Background(), config)
	if err != nil {
		log.Panicln(err)
	}

	log.Println(users.GetValue())
	log.Println(users.GetOdataNextLink())
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

	scopes := []string{"https://graph.microsoft.com/.default"}
	graphClient, err := msgraphsdk.NewGraphServiceClientWithCredentials(cred, scopes)
	return graphClient, err
}

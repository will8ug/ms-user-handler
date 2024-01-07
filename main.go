package main

import (
	"context"
	"flag"
	"log"

	azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	msgraphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	graphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
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

	results, err := graphClient.Users().Get(context.Background(), config)
	if err != nil {
		log.Panicln(err)
	}

	pageInterator, err := msgraphcore.NewPageIterator[graphmodels.Userable](
		results,
		graphClient.GetAdapter(),
		graphmodels.CreateUserCollectionResponseFromDiscriminatorValue,
	)
	if err != nil {
		log.Fatalf("Error creating page iterator: %v\n", err)
	}

	err = pageInterator.Iterate(context.Background(), func(user graphmodels.Userable) bool {
		log.Printf("%s\n", *user.GetDisplayName())
		// return true to continue the iteration
		return true
	})
	if err != nil {
		log.Panicln(err)
	}
	log.Println("Happy Ending.")
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

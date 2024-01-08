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
		log.Fatalln(err)
	}

	travelUsersWithPaging(graphClient, int32(2), func(user graphmodels.Userable) bool {
		log.Println()
		log.Printf("%v\n", *user.GetId())
		log.Printf("%s\n", *user.GetDisplayName())
		if (user.GetJobTitle() != nil) {
			log.Printf("%s\n", *user.GetJobTitle())
		} else {
			log.Println("Changing jobTitle when it's null")
			userPatch := graphmodels.NewUser()
			newJobTitle := "PatchedJobTitle"
			userPatch.SetJobTitle(&newJobTitle)
			graphClient.Users().ByUserId(*user.GetId()).Patch(context.Background(), userPatch, nil)
		}
		if (user.GetMobilePhone() != nil) {
			log.Printf("%s\n", *user.GetMobilePhone())
		} else {
			log.Println("mobilePhone is null")
		}
		log.Println()

		// return true to continue the iteration
		return true
	})

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

func travelUsersWithPaging(client *msgraphsdk.GraphServiceClient, pageSize int32, 
		callback func(pageItem graphmodels.Userable) bool) {
	reqParameters := &graphusers.UsersRequestBuilderGetQueryParameters{
		Select: []string{"id", "displayName", "jobTitle", "mobilePhone"},
		Top: &pageSize,
	}
	config := &graphusers.UsersRequestBuilderGetRequestConfiguration{
		QueryParameters: reqParameters,
	}

	results, err := client.Users().Get(context.Background(), config)
	if err != nil {
		log.Panicln(err)
	}

	pageInterator, err := msgraphcore.NewPageIterator[graphmodels.Userable](
		results,
		client.GetAdapter(),
		graphmodels.CreateUserCollectionResponseFromDiscriminatorValue,
	)
	if err != nil {
		log.Panicf("Error creating page iterator: %v\n", err)
	}

	err = pageInterator.Iterate(context.Background(), callback)
	if err != nil {
		log.Panicln(err)
	}
}

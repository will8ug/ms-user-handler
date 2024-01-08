package main

import (
	"context"
	"flag"
	"fmt"
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

var (
	b2cExtensionAppId string = ""
	isDryRun bool = true
	tenantCredential *TenantCredential = nil
)

func main() {
	parseArguments()
	log.Println(*tenantCredential)
	log.Printf("Is dry run: %v", isDryRun)

	graphClient, err := initGraphClient()
	if err != nil {
		log.Fatalln(err)
	}

	travelUsersWithPaging(graphClient, int32(2), func(user graphmodels.Userable) bool {
		log.Println()
		log.Printf("%v\n", *user.GetId())
		log.Printf("%s\n", *user.GetDisplayName())
		if (user.GetJobTitle() != nil) {
			log.Printf("%s\n", *user.GetJobTitle())
		} else if (!isDryRun) {
			log.Println("Changing jobTitle when it's null")
			updateJobTitle(graphClient, user)
		} else {
			log.Println("jobTitle is null")
		}
		if (user.GetMobilePhone() != nil) {
			log.Printf("%s\n", *user.GetMobilePhone())
		} else {
			log.Println("mobilePhone is null")
		}

		if (!isDryRun) {
			// just don't want to make the changes in my subsequent tests
			// so put the following lines into this "if"
			updateExtensionProperties(graphClient, user)
		}
		
		log.Println()
		// return true to continue the iteration
		return true
	})

	log.Println("Happy Ending.")
}

func parseArguments() {
	dryRun := flag.Bool("dryrun", true, "Dry run: true/false; default to true")
	tenantId := flag.String("tid", "", "Tenant ID")
	clientId := flag.String("cid", "", "Client ID")
	clientSecret := flag.String("csec", "", "Client secret")
	extappid := flag.String("extappid", "", "B2C extension application ID")
	flag.Parse()
	
	isDryRun = *dryRun
	b2cExtensionAppId = *extappid
	tenantCredential = &TenantCredential{*tenantId, *clientId, *clientSecret}
}

func initGraphClient() (*msgraphsdk.GraphServiceClient, error) {
	cred, _ := azidentity.NewClientSecretCredential(
		tenantCredential.tenantId,
		tenantCredential.clientId,
		tenantCredential.clientSecret,
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

func updateJobTitle(client *msgraphsdk.GraphServiceClient, user graphmodels.Userable) (err error) {
	userPatch := graphmodels.NewUser()
	newJobTitle := "PatchedJobTitle"
	userPatch.SetJobTitle(&newJobTitle)
	_, err = client.Users().ByUserId(*user.GetId()).Patch(context.Background(), userPatch, nil)
	log.Printf(">> New jobTitle: %v", newJobTitle)
	return
}

func updateExtensionProperties(client *msgraphsdk.GraphServiceClient, user graphmodels.Userable) (err error) {
	extUserPatch := graphmodels.NewUser()
	additionalData := map[string]interface{} {
		fmt.Sprintf("extension_%s_optIn", b2cExtensionAppId): "false",
	}
	extUserPatch.SetAdditionalData(additionalData)
	log.Printf("Trying to update extension properties: %v", additionalData)
	_, err = client.Users().ByUserId(*user.GetId()).Patch(context.Background(), extUserPatch, nil)
	return
}

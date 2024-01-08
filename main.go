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
	graphClient *msgraphsdk.GraphServiceClient = nil
)

func main() {
	parseArguments()
	log.Printf("Is dry run: %v", isDryRun)
	log.Println(*tenantCredential)

	err := initGraphClient()
	if err != nil {
		log.Fatalln(err)
	}

	travelUsersWithPaging(int32(2))

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

func initGraphClient() (err error) {
	cred, _ := azidentity.NewClientSecretCredential(
		tenantCredential.tenantId,
		tenantCredential.clientId,
		tenantCredential.clientSecret,
		nil,
	)

	scopes := []string{"https://graph.microsoft.com/.default"}
	graphClient, err = msgraphsdk.NewGraphServiceClientWithCredentials(cred, scopes)
	return err
}

func travelUsersWithPaging(pageSize int32) {
	reqParameters := &graphusers.UsersRequestBuilderGetQueryParameters{
		Select: []string{"id", "displayName", "jobTitle", "mobilePhone"},
		Top: &pageSize,
	}
	config := &graphusers.UsersRequestBuilderGetRequestConfiguration{
		QueryParameters: reqParameters,
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
		log.Panicf("Error creating page iterator: %v\n", err)
	}

	err = pageInterator.Iterate(context.Background(), handleSingleUser)
	if err != nil {
		log.Panicln(err)
	}
}

func handleSingleUser(user graphmodels.Userable) bool {
	log.Println()
	log.Printf("%v\n", *user.GetId())
	log.Printf("%s\n", *user.GetDisplayName())
	if (user.GetJobTitle() != nil) {
		log.Printf("%s\n", *user.GetJobTitle())
	} else if (!isDryRun) {
		log.Println("Changing jobTitle when it's null")
		updateJobTitle(user)
	} else {
		log.Println("jobTitle is null")
	}
	if (user.GetMobilePhone() != nil) {
		log.Printf("%s\n", *user.GetMobilePhone())
	} else {
		log.Println("mobilePhone is null")
	}

	if (!isDryRun) {
		updateExtensionProperties(user)
	}
	
	log.Println()
	// return true to continue the iteration
	return true
}

func updateJobTitle(user graphmodels.Userable) (err error) {
	userPatch := graphmodels.NewUser()
	newJobTitle := "PatchedJobTitle"
	userPatch.SetJobTitle(&newJobTitle)
	_, err = graphClient.Users().ByUserId(*user.GetId()).Patch(context.Background(), userPatch, nil)
	log.Printf(">> New jobTitle: %v", newJobTitle)
	return
}

func updateExtensionProperties(user graphmodels.Userable) (err error) {
	extUserPatch := graphmodels.NewUser()
	additionalData := map[string]interface{} {
		fmt.Sprintf("extension_%s_optIn", b2cExtensionAppId): "false",
	}
	extUserPatch.SetAdditionalData(additionalData)
	log.Printf("Trying to update extension properties: %v", additionalData)
	_, err = graphClient.Users().ByUserId(*user.GetId()).Patch(context.Background(), extUserPatch, nil)
	return
}

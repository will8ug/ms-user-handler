package main

import (
	"log"

	azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
)

func main() {
	graphClient, err := initGraphClient()
	log.Println(graphClient)
	log.Println(err)
}

func initGraphClient() (*msgraphsdk.GraphServiceClient, error) {
	cred, _ := azidentity.NewUsernamePasswordCredential(
		"tenant_id",
		"client_id",
		"user_name",
		"password",
		nil,
	)

	graphClient, err := msgraphsdk.NewGraphServiceClientWithCredentials(cred, []string{"User.Read"})
	return graphClient, err
}
package test

import (
	"os"
	"testing"

	awstest "github.com/pip-services3-go/pip-services3-aws-go/test"
	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
)

func TestDummyLambdaClient(t *testing.T) {

	lambdaArn := os.Getenv("LAMBDA_ARN")
	awsAccessId := os.Getenv("AWS_ACCESS_ID")
	awsAccessKey := os.Getenv("AWS_ACCESS_KEY")

	if lambdaArn == "" || awsAccessId == "" || awsAccessKey == "" {
		panic("AWS keys not sets!")
	}

	lambdaConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "aws",
		"connection.arn", lambdaArn,
		"credential.access_id", awsAccessId,
		"credential.access_key", awsAccessKey,
		"options.connection_timeout", 30000,
	)

	var client *DummyLambdaClient
	var fixture *awstest.DummyClientFixture

	client = NewDummyLambdaClient()
	client.Configure(lambdaConfig)

	fixture = awstest.NewDummyClientFixture(client)

	client.Open("")

	defer client.Close("")

	t.Run("DummyLambdaClient.CrudOperations", fixture.TestCrudOperations)
}

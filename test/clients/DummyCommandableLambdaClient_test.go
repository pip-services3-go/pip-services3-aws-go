package test

import (
	"os"
	"testing"

	awstest "github.com/pip-services3-go/pip-services3-aws-go/test"
	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	"github.com/stretchr/testify/assert"
)

func TestDummyCommandableLambdaClient(t *testing.T) {

	lambdaArn := os.Getenv("LAMBDA_ARN")
	awsAccessId := os.Getenv("AWS_ACCESS_ID")
	awsAccessKey := os.Getenv("AWS_ACCESS_KEY")

	if lambdaArn == "" || awsAccessId == "" || awsAccessKey == "" {
		return
	}

	lambdaConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "aws",
		"connection.arn", lambdaArn,
		"credential.access_id", awsAccessId,
		"credential.access_key", awsAccessKey,
		"options.connection_timeout", 30000,
	)

	var client *DummyCommandableLambdaClient
	var fixture *awstest.DummyClientFixture

	client = NewDummyCommandableLambdaClient()
	client.Configure(lambdaConfig)

	fixture = awstest.NewDummyClientFixture(client)

	err := client.Open("")
	assert.Nil(t, err)

	defer client.Close("")

	t.Run("DummyCommandableLambdaClient.CrudOperations", fixture.TestCrudOperations)
}

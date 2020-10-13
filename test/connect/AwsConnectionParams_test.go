package test

import (
	"testing"

	awscon "github.com/pip-services3-go/pip-services3-aws-go/connect"
	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	"github.com/stretchr/testify/assert"
)

func TestAwsConnectionParams(t *testing.T) {

	t.Run("TestAwsConnectionParams.Empty_Connection", EmptyConnection)
	t.Run("TestAwsConnectionParams.Parse_ARN", ParseARN)
	t.Run("TestAwsConnectionParams.Compose_ARN", ComposeARN)
}

func EmptyConnection(t *testing.T) {
	connection := awscon.NewEmptyAwsConnectionParams()
	assert.Equal(t, "arn:aws::::", connection.GetArn())
}

func ParseARN(t *testing.T) {
	connection := awscon.NewEmptyAwsConnectionParams()

	connection.SetArn("arn:aws:lambda:us-east-1:12342342332:function:pip-services-dummies")
	assert.Equal(t, "lambda", connection.GetService())
	assert.Equal(t, "us-east-1", connection.GetRegion())
	assert.Equal(t, "12342342332", connection.GetAccount())
	assert.Equal(t, "function", connection.GetResourceType())
	assert.Equal(t, "pip-services-dummies", connection.GetResource())

	connection.SetArn("arn:aws:s3:us-east-1:12342342332:pip-services-dummies")
	assert.Equal(t, "s3", connection.GetService())
	assert.Equal(t, "us-east-1", connection.GetRegion())
	assert.Equal(t, "12342342332", connection.GetAccount())
	assert.Equal(t, "", connection.GetResourceType())
	assert.Equal(t, "pip-services-dummies", connection.GetResource())

	connection.SetArn("arn:aws:lambda:us-east-1:12342342332:function/pip-services-dummies")
	assert.Equal(t, "lambda", connection.GetService())
	assert.Equal(t, "us-east-1", connection.GetRegion())
	assert.Equal(t, "12342342332", connection.GetAccount())
	assert.Equal(t, "function", connection.GetResourceType())
	assert.Equal(t, "pip-services-dummies", connection.GetResource())

}

func ComposeARN(t *testing.T) {
	connection := awscon.NewAwsConnectionParamsFromConfig(
		cconf.NewConfigParamsFromTuples(
			"connection.service", "lambda",
			"connection.region", "us-east-1",
			"connection.account", "12342342332",
			"connection.resource_type", "function",
			"connection.resource", "pip-services-dummies",
			"credential.access_id", "1234",
			"credential.access_key", "ABCDEF",
		))

	assert.Equal(t, "arn:aws:lambda:us-east-1:12342342332:function:pip-services-dummies", connection.GetArn())
	assert.Equal(t, "1234", connection.GetAccessId())
	assert.Equal(t, "ABCDEF", connection.GetAccessKey())

}

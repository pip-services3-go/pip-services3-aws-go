package test_container

import (
	"encoding/json"
	"testing"

	awstest "github.com/pip-services3-go/pip-services3-aws-go/test"
	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	"github.com/stretchr/testify/assert"
)

func TestDummyLambdaFunction(t *testing.T) {

	restConfig := cconf.NewConfigParamsFromTuples(
		"logger.descriptor", "pip-services:logger:console:default:1.0",
		"controller.descriptor", "pip-services-dummies:controller:default:default:1.0",
	)

	var _dummy1 awstest.Dummy
	var _dummy2 awstest.Dummy
	var lambda *DummyLambdaFunction
	ctrl := awstest.NewDummyController()

	lambda = NewDummyLambdaFunction()
	lambda.Configure(restConfig)

	var references *cref.References = cref.NewReferencesFromTuples(
		cref.NewDescriptor("pip-services-dummies", "controller", "default", "default", "1.0"), ctrl,
	)
	lambda.SetReferences(references)
	opnErr := lambda.Open("")
	assert.Nil(t, opnErr)
	defer lambda.Close("")

	_dummy1 = awstest.Dummy{Id: "", Key: "Key 1", Content: "Content 1"}
	_dummy2 = awstest.Dummy{Id: "", Key: "Key 2", Content: "Content 2"}

	var dummy1 awstest.Dummy

	params := make(map[string]interface{})

	// Create one dummy
	params["dummy"] = _dummy1
	params["cmd"] = "create_dummy"

	resBody, bodyErr := lambda.Act(params)
	assert.Nil(t, bodyErr)

	var dummy awstest.Dummy
	jsonErr := json.Unmarshal([]byte(resBody), &dummy)

	assert.Nil(t, jsonErr)
	assert.NotNil(t, dummy)
	assert.Equal(t, dummy.Content, _dummy1.Content)
	assert.Equal(t, dummy.Key, _dummy1.Key)

	dummy1 = dummy

	// Create another dummy
	params["dummy"] = _dummy2
	params["cmd"] = "create_dummy"

	resBody, bodyErr = lambda.Act(params)
	assert.Nil(t, bodyErr)

	jsonErr = json.Unmarshal([]byte(resBody), &dummy)

	assert.Nil(t, jsonErr)
	assert.NotNil(t, dummy)
	assert.Equal(t, dummy.Content, _dummy2.Content)
	assert.Equal(t, dummy.Key, _dummy2.Key)
	//dummy2 = dummy

	// Get all dummies
	delete(params, "dummy")
	params["cmd"] = "get_dummies"
	resBody, bodyErr = lambda.Act(params)
	assert.Nil(t, bodyErr)

	var dummies awstest.DummyDataPage
	jsonErr = json.Unmarshal([]byte(resBody), &dummies)
	assert.Nil(t, jsonErr)
	assert.NotNil(t, dummies)
	assert.Len(t, dummies.Data, 2)

	// Update the dummy

	dummy1.Content = "Updated Content 1"

	params["dummy"] = dummy1
	params["cmd"] = "update_dummy"

	resBody, bodyErr = lambda.Act(params)
	assert.Nil(t, bodyErr)
	jsonErr = json.Unmarshal([]byte(resBody), &dummy)
	assert.Nil(t, jsonErr)
	assert.NotNil(t, dummy)

	assert.Equal(t, dummy.Content, "Updated Content 1")
	assert.Equal(t, dummy.Key, _dummy1.Key)
	dummy1 = dummy

	// Delete dummy
	delete(params, "dummy")
	params["dummy_id"] = dummy1.Id
	params["cmd"] = "delete_dummy"
	resBody, bodyErr = lambda.Act(params)
	assert.Nil(t, bodyErr)

	// Try to get delete dummy
	dummies.Data = dummies.Data[:0]
	*dummies.Total = 0

	params["dummy_id"] = dummy1.Id
	params["cmd"] = "get_dummy_by_id"

	resBody, bodyErr = lambda.Act(params)
	assert.Nil(t, bodyErr)
	jsonErr = json.Unmarshal([]byte(resBody), &dummies)
	assert.Nil(t, jsonErr)
	assert.NotNil(t, dummies)
	assert.Len(t, dummies.Data, 0)
}

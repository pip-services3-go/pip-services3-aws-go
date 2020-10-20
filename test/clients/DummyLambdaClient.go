package test

import (
	"reflect"

	awsclient "github.com/pip-services3-go/pip-services3-aws-go/clients"
	awstest "github.com/pip-services3-go/pip-services3-aws-go/test"
	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
)

var (
	dummyDataPageType = reflect.TypeOf(&awstest.DummyDataPage{})
	dummyType         = reflect.TypeOf(&awstest.Dummy{})
)

type DummyLambdaClient struct {
	*awsclient.LambdaClient
}

func NewDummyLambdaClient() *DummyLambdaClient {
	c := &DummyLambdaClient{
		LambdaClient: awsclient.NewLambdaClient(),
	}
	return c
}
func (c *DummyLambdaClient) GetDummies(correlationId string, filter *cdata.FilterParams,
	paging *cdata.PagingParams) (result *awstest.DummyDataPage, err error) {

	params := cdata.NewEmptyAnyValueMap()
	params.SetAsObject("filter", filter)
	params.SetAsObject("paging", paging)

	calValue, calErr := c.Call(dummyDataPageType, "get_dummies", correlationId, params.Value())
	if calErr != nil {
		return nil, calErr
	}

	result, _ = calValue.(*awstest.DummyDataPage)
	c.Instrument(correlationId, "dummy.get_dummies")
	return result, nil
}

func (c *DummyLambdaClient) GetDummyById(correlationId string, dummyId string) (result *awstest.Dummy, err error) {

	params := cdata.NewEmptyAnyValueMap()
	params.SetAsObject("dummy_id", dummyId)

	calValue, calErr := c.Call(dummyType, "get_dummy_by_id", correlationId, params.Value())

	if calErr != nil {
		return nil, calErr
	}

	result, _ = calValue.(*awstest.Dummy)
	c.Instrument(correlationId, "dummy.get_one_by_id")
	return result, nil
}

func (c *DummyLambdaClient) CreateDummy(correlationId string, dummy awstest.Dummy) (result *awstest.Dummy, err error) {

	params := cdata.NewEmptyAnyValueMap()
	params.SetAsObject("dummy", dummy)

	calValue, calErr := c.Call(dummyType, "create_dummy", correlationId, params.Value())
	if calErr != nil {
		return nil, calErr
	}

	result, _ = calValue.(*awstest.Dummy)
	c.Instrument(correlationId, "dummy.create_dummy")
	return result, nil
}

func (c *DummyLambdaClient) UpdateDummy(correlationId string, dummy awstest.Dummy) (result *awstest.Dummy, err error) {

	params := cdata.NewEmptyAnyValueMap()
	params.SetAsObject("dummy", dummy)

	calValue, calErr := c.Call(dummyType, "update_dummy", correlationId, params.Value())
	if calErr != nil {
		return nil, calErr
	}

	result, _ = calValue.(*awstest.Dummy)
	c.Instrument(correlationId, "dummy.update_dummy")
	return result, nil
}

func (c *DummyLambdaClient) DeleteDummy(correlationId string, dummyId string) (result *awstest.Dummy, err error) {

	params := cdata.NewEmptyAnyValueMap()
	params.SetAsObject("dummy_id", dummyId)
	calValue, calErr := c.Call(dummyType, "delete_dummy", correlationId, params.Value())
	if calErr != nil {
		return nil, calErr
	}

	result, _ = calValue.(*awstest.Dummy)
	c.Instrument(correlationId, "dummy.delete_dummy")
	return result, nil
}

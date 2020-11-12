package test

import (
	awsclient "github.com/pip-services3-go/pip-services3-aws-go/clients"
	awstest "github.com/pip-services3-go/pip-services3-aws-go/test"
	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
)

type DummyCommandableLambdaClient struct {
	*awsclient.CommandableLambdaClient
}

func NewDummyCommandableLambdaClient() *DummyCommandableLambdaClient {
	c := &DummyCommandableLambdaClient{
		CommandableLambdaClient: awsclient.NewCommandableLambdaClient("dummy"),
	}
	return c
}
func (c *DummyCommandableLambdaClient) GetDummies(correlationId string, filter *cdata.FilterParams,
	paging *cdata.PagingParams) (result *awstest.DummyDataPage, err error) {

	params := cdata.NewEmptyAnyValueMap()
	params.SetAsObject("filter", filter.Value())
	params.SetAsObject("paging", paging)

	calValue, calErr := c.CallCommand(dummyDataPageType, "get_dummies", correlationId, params)
	if calErr != nil {
		return nil, calErr
	}

	result, _ = calValue.(*awstest.DummyDataPage)
	c.Instrument(correlationId, "dummy.get_dummies")
	return result, nil
}

func (c *DummyCommandableLambdaClient) GetDummyById(correlationId string, dummyId string) (result *awstest.Dummy, err error) {

	params := cdata.NewEmptyAnyValueMap()
	params.SetAsObject("dummy_id", dummyId)

	calValue, calErr := c.CallCommand(dummyType, "get_dummy_by_id", correlationId, params)

	if calErr != nil {
		return nil, calErr
	}

	result, _ = calValue.(*awstest.Dummy)
	c.Instrument(correlationId, "dummy.get_one_by_id")
	return result, nil
}

func (c *DummyCommandableLambdaClient) CreateDummy(correlationId string, dummy awstest.Dummy) (result *awstest.Dummy, err error) {

	params := cdata.NewEmptyAnyValueMap()
	params.SetAsObject("dummy", dummy)

	calValue, calErr := c.CallCommand(dummyType, "create_dummy", correlationId, params)
	if calErr != nil {
		return nil, calErr
	}

	result, _ = calValue.(*awstest.Dummy)
	c.Instrument(correlationId, "dummy.create_dummy")
	return result, nil
}

func (c *DummyCommandableLambdaClient) UpdateDummy(correlationId string, dummy awstest.Dummy) (result *awstest.Dummy, err error) {

	params := cdata.NewEmptyAnyValueMap()
	params.SetAsObject("dummy", dummy)

	calValue, calErr := c.CallCommand(dummyType, "update_dummy", correlationId, params)
	if calErr != nil {
		return nil, calErr
	}

	result, _ = calValue.(*awstest.Dummy)
	c.Instrument(correlationId, "dummy.update_dummy")
	return result, nil
}

func (c *DummyCommandableLambdaClient) DeleteDummy(correlationId string, dummyId string) (result *awstest.Dummy, err error) {

	params := cdata.NewEmptyAnyValueMap()
	params.SetAsObject("dummy_id", dummyId)
	calValue, calErr := c.CallCommand(dummyType, "delete_dummy", correlationId, params)
	if calErr != nil {
		return nil, calErr
	}

	result, _ = calValue.(*awstest.Dummy)
	c.Instrument(correlationId, "dummy.delete_dummy")
	return result, nil
}

package test

import cdata "github.com/pip-services3-go/pip-services3-commons-go/data"

type IDummyClient interface {
	GetDummies(correlationId string, filter *cdata.FilterParams, paging *cdata.PagingParams) (result *DummyDataPage, err error)
	GetDummyById(correlationId string, dummyId string) (result *Dummy, err error)
	CreateDummy(correlationId string, dummy Dummy) (result *Dummy, err error)
	UpdateDummy(correlationId string, dummy Dummy) (result *Dummy, err error)
	DeleteDummy(correlationId string, dummyId string) (result *Dummy, err error)
}

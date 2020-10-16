package test

import (
	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
)

type IDummyController interface {
	GetPageByFilter(correlationId string, filter *cdata.FilterParams, paging *cdata.PagingParams) (result *DummyDataPage, err error)
	GetOneById(correlationId string, id string) (result *Dummy, err error)
	Create(correlationId string, entity Dummy) (result *Dummy, err error)
	Update(correlationId string, entity Dummy) (result *Dummy, err error)
	DeleteById(correlationId string, id string) (result *Dummy, err error)
}

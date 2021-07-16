package test_services

import (
	awsserv "github.com/pip-services3-go/pip-services3-aws-go/services"
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
)

type DummyCommandableLambdaService struct {
	*awsserv.CommandableLambdaService
}

func NewDummyCommandableLambdaService() *DummyCommandableLambdaService {
	c := &DummyCommandableLambdaService{}
	c.CommandableLambdaService = awsserv.InheritCommandableLambdaService(c, "dummy")
	c.DependencyResolver.Put("controller", cref.NewDescriptor("pip-services-dummies", "controller", "default", "*", "*"))
	return c
}

package test

import (
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	cbuild "github.com/pip-services3-go/pip-services3-components-go/build"
)

type DummyFactory struct {
	cbuild.Factory
	Descriptor           *cref.Descriptor
	ControllerDescriptor *cref.Descriptor
}

func NewDummyFactory() *DummyFactory {

	c := &DummyFactory{
		Factory:              *cbuild.NewFactory(),
		Descriptor:           cref.NewDescriptor("pip-services-dummies", "factory", "default", "default", "1.0"),
		ControllerDescriptor: cref.NewDescriptor("pip-services-dummies", "controller", "default", "*", "1.0"),
	}

	c.RegisterType(c.ControllerDescriptor, NewDummyController)
	return c
}

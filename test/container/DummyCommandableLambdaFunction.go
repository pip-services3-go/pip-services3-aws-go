package test_container

import (
	awscont "github.com/pip-services3-go/pip-services3-aws-go/container"
	awstest "github.com/pip-services3-go/pip-services3-aws-go/test"
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
)

type DummyCommandableLambdaFunction struct {
	*awscont.CommandableLambdaFunction
}

func NewDummyCommandableLambdaFunction() *DummyCommandableLambdaFunction {
	c := &DummyCommandableLambdaFunction{}
	c.CommandableLambdaFunction = awscont.NewCommandableLambdaFunction("dummy", "Dummy lambda function")

	c.DependencyResolver.Put("controller", cref.NewDescriptor("pip-services-dummies", "controller", "default", "*", "*"))
	c.AddFactory(awstest.NewDummyFactory())
	return c
}

// func main() {
// 	lambda.Start(NewDummyCommandableLambdaFunction().GetHandler())
// }

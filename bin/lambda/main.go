package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	awstest "github.com/pip-services3-go/pip-services3-aws-go/test"
	awscont "github.com/pip-services3-go/pip-services3-aws-go/test/container"
	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
)

func main() {

	restConfig := cconf.NewConfigParamsFromTuples(
		"logger.descriptor", "pip-services:logger:console:default:1.0",
		"controller.descriptor", "pip-services-dummies:controller:default:default:1.0",
	)

	var container *awscont.DummyLambdaFunction
	ctrl := awstest.NewDummyController()

	container = awscont.NewDummyLambdaFunction()
	container.Configure(restConfig)

	var references *cref.References = cref.NewReferencesFromTuples(
		cref.NewDescriptor("pip-services-dummies", "controller", "default", "default", "1.0"), ctrl,
		cref.NewDescriptor("pip-services-dummies", "service", "rest", "default", "1.0"), container,
	)
	container.SetReferences(references)
	defer container.Close("")
	opnErr := container.Open("")
	if opnErr == nil {
		lambda.Start(container.GetHandler())
	}

}

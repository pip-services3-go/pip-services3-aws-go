package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	awscont "github.com/pip-services3-go/pip-services3-aws-go/test/container"
	//cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
)

func main() {

	// config := cconf.NewConfigParamsFromTuples(
	// 	"logger.descriptor", "pip-services:logger:console:default:1.0",
	// )

	var container *awscont.DummyLambdaFunction

	container = awscont.NewDummyLambdaFunction()
	// container.Configure(config)

	// defer container.Close("")
	// opnErr := container.Open("")
	// if opnErr == nil {
	lambda.Start(container.GetHandler())
	// }

}

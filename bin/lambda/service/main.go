package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	awsserv "github.com/pip-services3-go/pip-services3-aws-go/test/services"
	//cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
)

func main() {

	// config := cconf.NewConfigParamsFromTuples(
	// 	"logger.descriptor", "pip-services:logger:console:default:1.0",
	// )

	var container *awsserv.DummyLambdaFunction

	container = awsserv.NewDummyLambdaFunction()
	// container.Configure(config)

	// defer container.Close("")
	// opnErr := container.Open("")
	// if opnErr == nil {
	lambda.Start(container.GetHandler())
	// }

}

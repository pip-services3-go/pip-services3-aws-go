package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	awsserv "github.com/pip-services3-go/pip-services3-aws-go/test/services"
)

func main() {

	var container *awsserv.DummyLambdaFunction

	container = awsserv.NewDummyLambdaFunction()

	defer container.Close("")
	opnErr := container.Run()
	if opnErr == nil {
		lambda.Start(container.GetHandler())
	}

}

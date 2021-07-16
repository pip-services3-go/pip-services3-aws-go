package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	awscont "github.com/pip-services3-go/pip-services3-aws-go/test/container"
)

func main() {
	var container *awscont.DummyLambdaFunction

	container = awscont.NewDummyLambdaFunction()

	defer container.Close("")
	err := container.Run()
	if err != nil {
		panic(err)
	}
	lambda.Start(container.GetHandler())
}

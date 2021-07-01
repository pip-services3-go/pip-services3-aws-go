package services

import (
	cvalid "github.com/pip-services3-go/pip-services3-commons-go/validate"
)

type LambdaAction struct {

	//Command to call the action
	Cmd string

	//Schema to validate action parameters
	Schema *cvalid.Schema

	//Action to be executed
	Action func(params map[string]interface{}) (interface{}, error)
}

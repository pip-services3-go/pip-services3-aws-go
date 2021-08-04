package container

import cref "github.com/pip-services3-go/pip-services3-commons-go/refer"

type ILambdaFunctionOverrides interface {
	cref.IReferenceable
	// Perform required registration steps.
	Register()
}

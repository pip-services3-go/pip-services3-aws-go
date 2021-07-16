package container

import cref "github.com/pip-services3-go/pip-services3-commons-go/refer"

/*
 IRegisterable is interface to perform on-demand registrations.
*/
type IRegisterable interface {
	// Perform required registration steps.
	Register()
}

/*
 IRegisterable is interface to perform on-demand registrations.
*/
type IContainerable interface {
	IRegisterable
	cref.IReferenceable
}

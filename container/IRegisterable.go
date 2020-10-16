package container

/*
 IRegisterable is interface to perform on-demand registrations.
*/
type IRegisterable interface {
	// Perform required registration steps.
	Register()
}

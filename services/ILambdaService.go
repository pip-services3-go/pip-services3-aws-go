package services

/**
 * An interface that allows to integrate lambda services into lambda function containers
 * and connect their actions to the function calls.
 */
type ILambdaService interface {
	/**
	 * Get all actions supported by the service.
	 * Returns an array with supported actions.
	 */
	getActions() []*LambdaAction
}

package services

import (
	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cerr "github.com/pip-services3-go/pip-services3-commons-go/errors"
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	cvalid "github.com/pip-services3-go/pip-services3-commons-go/validate"
	ccount "github.com/pip-services3-go/pip-services3-components-go/count"
	clog "github.com/pip-services3-go/pip-services3-components-go/log"
	ctrace "github.com/pip-services3-go/pip-services3-components-go/trace"
	rpcserv "github.com/pip-services3-go/pip-services3-rpc-go/services"
)

type ILambdaServiceOverrides interface {
	Register()
}

/*
Abstract service that receives remove calls via AWS Lambda protocol.

This service is intended to work inside LambdaFunction container that
exploses registered actions externally.

### Configuration parameters ###

- dependencies:
  - controller:            override for Controller dependency

### References ###

- *:logger:*:*:1.0               (optional) [[ILogger]] components to pass log messages
- *:counters:*:*:1.0             (optional) [[ICounters]] components to pass collected measurements

See [[LambdaClient]]

### Example ###

    struct MyLambdaService struct  {
       LambdaService
       controller IMyController
    }
       ...
       func NewMyLambdaService()* MyLambdaService {
          c:= &MyLambdaService{}
          c.LambdaService = NewLambdaService("v1.myservice")
          c.DependencyResolver.Put(
              "controller",
              cref.NewDescriptor("mygroup","controller","*","*","1.0")
          )
          return c
       }

       func (c * LambdaService)  SetReferences(references: IReferences){
          c.LambdaService.setReferences(references);
          ref := c.DependencyResolver.GetRequired("controller")
          c.controller = ref.(IMyController)
       }

       func (c * LambdaService)  Register(){
           c.RegisterAction("get_mydata", nil,  func(params map[string]interface{})(interface{}, error) {
                correlationId := params.GetAsString("correlation_id")
                id := params.GetAsString("id")
               return  c.controller.GetMyData(correlationId, id);
           });
           ...
       }

     service := NewMyLambdaService();
    service.Configure(NewConfigParamsFromTuples(
        "connection.protocol", "http",
        "connection.host", "localhost",
        "connection.port", 8080
    ));
    service.SetReferences(cref.NewReferencesFromTuples(
       cref.NewDescriptor("mygroup","controller","default","default","1.0"), controller
    ));

    service.Open("123");
    fmt.Println("The Lambda 'v1.myservice' service is running on port 8080");
*/
type LambdaService struct { // ILambdaService, IOpenable, IConfigurable, IReferenceable

	name         string
	actions      []*LambdaAction
	interceptors []func(params map[string]interface{}, next func(params map[string]interface{}) (interface{}, error)) (interface{}, error)
	opened       bool

	Overrides ILambdaServiceOverrides

	// The dependency resolver.
	DependencyResolver *cref.DependencyResolver
	// The logger.
	Logger *clog.CompositeLogger
	//The performance counters.
	Counters *ccount.CompositeCounters
	//The tracer.
	Tracer *ctrace.CompositeTracer
}

// Creates an instance of this service.
// -  name a service name to generate action cmd. RestService()
func InheritLambdaService(overrides ILambdaServiceOverrides, name string) *LambdaService {
	return &LambdaService{
		Overrides:          overrides,
		name:               name,
		actions:            make([]*LambdaAction, 0),
		interceptors:       make([]func(params map[string]interface{}, next func(params map[string]interface{}) (interface{}, error)) (interface{}, error), 0),
		DependencyResolver: cref.NewDependencyResolver(),
		Logger:             clog.NewCompositeLogger(),
		Counters:           ccount.NewCompositeCounters(),
		Tracer:             ctrace.NewCompositeTracer(nil),
	}
}

// Configures component by passing configuration parameters.
// -  config    configuration parameters to be set.
func (c *LambdaService) Configure(config *cconf.ConfigParams) {
	c.DependencyResolver.Configure(config)
}

// Sets references to dependent components.
// -  references 	references to locate the component dependencies.
func (c *LambdaService) SetReferences(references cref.IReferences) {
	c.Logger.SetReferences(references)
	c.Counters.SetReferences(references)
	c.Tracer.SetReferences(references)
	c.DependencyResolver.SetReferences(references)
}

// Get all actions supported by the service.
// Returns an array with supported actions.
func (c *LambdaService) GetActions() []*LambdaAction {
	return c.actions
}

/*
   Adds instrumentation to log calls and measure call time.
   It returns a Timing object that is used to end the time measurement.
    *
   -  correlationId     (optional) transaction id to trace execution through call chain.
   -  name              a method name.
   @returns Timing object to end the time measurement.
*/
func (c *LambdaService) Instrument(correlationId string, name string) *rpcserv.InstrumentTiming {
	c.Logger.Trace(correlationId, "Executing %s method", name)
	c.Counters.IncrementOne(name + ".exec_count")

	counterTiming := c.Counters.BeginTiming(name + ".exec_time")
	traceTiming := c.Tracer.BeginTrace(correlationId, name, "")
	return rpcserv.NewInstrumentTiming(correlationId, name, "exec",
		c.Logger, c.Counters, counterTiming, traceTiming)
}

//    Checks if the component is opened.
//    Returns true if the component has been opened and false otherwise.
func (c *LambdaService) IsOpen() bool {
	return c.opened
}

//    Opens the component.
//    -  correlationId 	(optional) transaction id to trace execution through call chain.
func (c *LambdaService) Open(correlationId string) error {
	if c.opened {
		return nil
	}

	c.Register()

	c.opened = true
	return nil
}

// Closes component and frees used resources.
// -  correlationId 	(optional) transaction id to trace execution through call chain.
func (c *LambdaService) Close(correlationId string) error {
	if !c.opened {
		return nil
	}

	c.opened = false
	c.actions = make([]*LambdaAction, 0)
	c.interceptors = make([]func(params map[string]interface{}, next func(params map[string]interface{}) (interface{}, error)) (interface{}, error), 0)
	return nil
}

func (c *LambdaService) ApplyValidation(schema *cvalid.Schema, action func(params map[string]interface{}) (interface{}, error)) func(map[string]interface{}) (interface{}, error) {
	// Create an action function
	actionWrapper := func(params map[string]interface{}) (interface{}, error) {
		// Validate object
		if schema != nil && params != nil {
			// Perform validation
			correlationId, _ := params["correlation_id"].(string)
			err := schema.ValidateAndReturnError(correlationId, params, false)
			if err != nil {
				return nil, err
			}
		}
		return action(params)
	}

	return actionWrapper
}

func (c *LambdaService) ApplyInterceptors(action func(map[string]interface{}) (interface{}, error)) func(map[string]interface{}) (interface{}, error) {
	actionWrapper := action

	for index := len(c.interceptors) - 1; index >= 0; index-- {
		interceptor := c.interceptors[index]
		actionWrapper = (func(action func(map[string]interface{}) (interface{}, error)) func(map[string]interface{}) (interface{}, error) {
			return func(params map[string]interface{}) (interface{}, error) {
				return interceptor(params, action)
			}
		})(actionWrapper)
	}

	return actionWrapper
}

func (c *LambdaService) GenerateActionCmd(name string) string {
	cmd := name
	if c.name != "" {
		cmd = c.name + "." + cmd
	}
	return cmd
}

// Registers a action in AWS Lambda function.
// -  name          an action name
// -  schema        a validation schema to validate received parameters.
// -  action        an action function that is called when operation is invoked.
func (c *LambdaService) RegisterAction(name string, schema *cvalid.Schema, action func(params map[string]interface{}) (interface{}, error)) {
	actionWrapper := c.ApplyValidation(schema, action)
	actionWrapper = c.ApplyInterceptors(actionWrapper)

	registeredAction := &LambdaAction{
		Cmd:    c.GenerateActionCmd(name),
		Schema: schema,
		Action: func(params map[string]interface{}) (interface{}, error) { return actionWrapper(params) },
	}
	c.actions = append(c.actions, registeredAction)
}

// Registers an action with authorization.
// -  name          an action name
// -  schema        a validation schema to validate received parameters.
// -  authorize     an authorization interceptor
// -  action        an action function that is called when operation is invoked.
func (c *LambdaService) RegisterActionWithAuth(name string, schema *cvalid.Schema,
	authorize func(params map[string]interface{}, next func(map[string]interface{}) (interface{}, error)) (interface{}, error),
	action func(params map[string]interface{}) (interface{}, error)) {

	actionWrapper := c.ApplyValidation(schema, action)
	// Add authorization just before validation
	actionWrapper = func(params map[string]interface{}) (interface{}, error) {
		return authorize(params, actionWrapper)
	}
	actionWrapper = c.ApplyInterceptors(actionWrapper)

	registeredAction := &LambdaAction{
		Cmd:    c.GenerateActionCmd(name),
		Schema: schema,
		Action: func(params map[string]interface{}) (interface{}, error) { return actionWrapper(params) },
	}
	c.actions = append(c.actions, registeredAction)
}

// Registers a middleware for actions in AWS Lambda service.
// -  action        an action function that is called when middleware is invoked.
func (c *LambdaService) RegisterInterceptor(action func(params map[string]interface{}, next func(params map[string]interface{}) (interface{}, error)) (interface{}, error)) {
	c.interceptors = append(c.interceptors, action)
}

// Registers all service routes in HTTP endpoint.
// This method is called by the service and must be overriden
// in child classes.
func (c *LambdaService) Register() {
	c.Overrides.Register()
}

// Calls registered action in this lambda function.
// "cmd" parameter in the action parameters determin
// what action shall be called.
// This method shall only be used in testing.
// -  params action parameters.
func (c *LambdaService) Act(params map[string]interface{}) (interface{}, error) {
	cmd, ok := params["cmd"].(string)
	correlationId, _ := params["correlation_id"].(string)

	if !ok || cmd == "" {
		return nil, cerr.NewBadRequestError(
			correlationId,
			"NO_COMMAND",
			"Cmd parameter is missing",
		)
	}

	var action *LambdaAction
	for _, act := range c.actions {
		if act.Cmd == cmd {
			action = act
			break
		}
	}

	if action == nil {
		return nil, cerr.NewBadRequestError(
			correlationId,
			"NO_ACTION",
			"Action "+cmd+" was not found",
		).
			WithDetails("command", cmd)
	}

	return action.Action(params)
}

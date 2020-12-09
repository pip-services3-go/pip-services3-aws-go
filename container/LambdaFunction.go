package container

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"

	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cerr "github.com/pip-services3-go/pip-services3-commons-go/errors"
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	cvalid "github.com/pip-services3-go/pip-services3-commons-go/validate"
	ccount "github.com/pip-services3-go/pip-services3-components-go/count"
	"github.com/pip-services3-go/pip-services3-components-go/log"
	cproc "github.com/pip-services3-go/pip-services3-container-go/container"
)

/*
Abstract AWS Lambda function, that acts as a container to instantiate and run components
and expose them via external entry point.

When handling calls "cmd" parameter determines which what action shall be called, while
other parameters are passed to the action itself.

Container configuration for this Lambda function is stored in "./config/config.yml" file.
But this path can be overriden by CONFIG_PATH environment variable.

### Configuration parameters ###

 - dependencies:
    - controller:                  override for Controller dependency
 - connections:
    - discovery_key:               (optional) a key to retrieve the connection from IDiscovery
    - region:                      (optional) AWS region
 - credentials:
    - store_key:                   (optional) a key to retrieve the credentials from ICredentialStore
    - access_id:                   AWS access/client id
    - access_key:                  AWS access/client id

### References ###

 - \*:logger:\*:\*:1.0            (optional) ILogger components to pass log messages
 - \*:counters:\*:\*:1.0          (optional) ICounters components to pass collected measurements
 - \*:discovery:\*:\*:1.0         (optional) IDiscovery services to resolve connection
 - \*:credential-store:\*:\*:1.0  (optional) Credential stores to resolve credentials

See LambdaClient

 ### Example ###
 
    class MyLambdaFunction extends LambdaFunction {
        func (c* LambdaFunction) _controller: IMyController;
        ...
        func (c* LambdaFunction) constructor() {
            base("mygroup", "MyGroup lambda function");
            c.dependencyResolver.put(
                "controller",
                new Descriptor("mygroup","controller","*","*","1.0")
            );
        }
 
        func (c* LambdaFunction) setReferences(references: IReferences){
            base.setReferences(references);
            c.controller = c.dependencyResolver.getRequired<IMyController>("controller");
        }
 
        func (c* LambdaFunction) register(){
            registerAction("get_mydata", null, (params, callback) => {
                let correlationId = params.correlation_id;
                let id = params.id;
                c.controller.getMyData(correlationId, id, callback);
            });
            ...
        }
    }
 
    let lambda = new MyLambdaFunction();
 
    service.run((err) => {
        console.log("MyLambdaFunction is started");
    });
*/
type LambdaFunction struct {
	*cproc.Container
	IRegisterable
	/*
	   The performanc counters.
	*/
	counters *ccount.CompositeCounters
	/*
	   The dependency resolver.
	*/
	DependencyResolver *cref.DependencyResolver
	/*
	   The map of registred validation schemas.
	*/
	schemas map[string]*cvalid.Schema
	/*
	   The map of registered actions.
	*/
	actions map[string]func(map[string]interface{}) (interface{}, error)
	/*
	   The default path to config file.
	*/
	configPath string
}

/*
Creates a new instance of this lambda function.
   - name          (optional) a container name (accessible via ContextInfo)
   - description   (optional) a container description (accessible via ContextInfo)
*/
func InheriteLambdaFunction(name string, description string, register IRegisterable) *LambdaFunction {
	c := &LambdaFunction{
		counters:           ccount.NewCompositeCounters(),
		DependencyResolver: cref.NewDependencyResolver(),
		schemas:            make(map[string]*cvalid.Schema, 0),
		actions:            make(map[string]func(map[string]interface{}) (interface{}, error), 0),
		configPath:         "./config/config.yml",
		IRegisterable:      register,
	}
	c.Container = cproc.InheritContainer(name, description, c)
	c.SetLogger(log.NewConsoleLogger())
	return c
}

func (c *LambdaFunction) getConfigPath() string {
	res := os.Getenv("CONFIG_PATH")
	if res == "" {
		return c.configPath
	}
	return res
}

func (c *LambdaFunction) getParameters() *cconf.ConfigParams {
	parameters := cconf.NewConfigParamsFromValue(os.Environ())
	return parameters
}

func (c *LambdaFunction) captureErrors(correlationId string) {
	if r := recover(); r != nil {
		err, _ := r.(error)
		c.Logger().Fatal(correlationId, err, "Process is terminated")
		os.Exit(1)
	}
}

func (c *LambdaFunction) captureExit(correlationId string) {
	c.Logger().Info(correlationId, "Press Control-C to stop the microservice...")

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)

	go func() {
		select {
		case <-ch:
			c.Close(correlationId)
			c.Logger().Info(correlationId, "Googbye!")
			os.Exit(0)
		}
	}()
}

/*
Sets references to dependent components.
  - references 	references to locate the component dependencies.
*/
func (c *LambdaFunction) SetReferences(references cref.IReferences) {
	//c.Container.SetReferences(references)
	c.counters.SetReferences(references)
	c.DependencyResolver.SetReferences(references)
	c.Register()
}

/*
Adds instrumentation to log calls and measure call time.
It returns a Timing object that is used to end the time measurement.
   - correlationId     (optional) transaction id to trace execution through call chain.
   - name              a method name.
Returns Timing object to end the time measurement.
*/
func (c *LambdaFunction) Instrument(correlationId string, name string) *ccount.Timing {

	c.Logger().Trace(correlationId, "Executing %s method", name)
	return c.counters.BeginTiming(name + ".exec_time")
}

// InstrumentError method are adds instrumentation to error handling.
// Parameters:
//    - correlationId  string  (optional) transaction id to trace execution through call chain.
//    - name    string         a method name.
//    - err     error          an occured error
//    - result  interface{}    (optional) an execution result
// Returns:  result interface{}, err error
// (optional) an execution callback
func (c *LambdaFunction) InstrumentError(correlationId string, name string, errIn error,
	resIn interface{}) (result interface{}, err error) {
	if errIn != nil {
		c.Logger().Error(correlationId, errIn, "Failed to execute %s method", name)
		c.counters.IncrementOne(name + ".exec_errors")
	}
	return resIn, errIn
}

/*
Runs this lambda function, loads container configuration,
instantiate components and manage their lifecycle,
makes this function ready to access action calls.
  - callback callback function that receives error or null for success.
*/
func (c *LambdaFunction) Run() error {
	correlationId := c.Info().Name

	path := c.getConfigPath()
	parameters := c.getParameters()
	c.ReadConfigFromFile(correlationId, path, parameters)

	c.captureErrors(correlationId)
	c.captureExit(correlationId)
	return c.Open(correlationId)
}

/*
Registers an action in this lambda function.
   - cmd           a action/command name.
   - schema        a validation schema to validate received parameters.
   - action        an action function that is called when action is invoked.
*/
func (c *LambdaFunction) RegisterAction(cmd string, schema *cvalid.Schema,
	action func(params map[string]interface{}) (result interface{}, err error)) error {

	if cmd == "" {
		return cerr.NewUnknownError("", "NO_COMMAND", "Missing command")
	}

	if action == nil {
		return cerr.NewUnknownError("", "NO_ACTION", "Missing action")
	}

	// Hack!!! Wrapping action to preserve prototyping context
	actionCurl := func(params map[string]interface{}) (interface{}, error) {
		// Perform validation
		if schema != nil {
			correlationId, _ := params["correlaton_id"].(string)
			err := schema.ValidateAndReturnError(correlationId, params, false)
			if err != nil {
				return nil, err
			}
		}

		return action(params)
	}

	c.actions[cmd] = actionCurl
	return nil
}

func (c *LambdaFunction) execute(ctx context.Context, params map[string]interface{}) (string, error) {

	cmd, ok := params["cmd"].(string)
	correlationId, _ := params["correlation_id"].(string)

	if !ok || cmd == "" {
		err := cerr.NewBadRequestError(
			correlationId,
			"NO_COMMAND",
			"Cmd parameter is missing")
		ctx.Done()
		return "ERROR", err
	}

	action := c.actions[cmd]
	if action == nil {
		err := cerr.NewBadRequestError(
			correlationId,
			"NO_ACTION",
			"Action "+cmd+" was not found").
			WithDetails("command", cmd)

		ctx.Done()
		return "ERROR", err
	}

	res, err := action(params)
	ctx.Done()
	resStr := "ERROR"
	if res != nil {
		convRes, convErr := json.Marshal(res)
		if convRes == nil || convErr != nil {
			err = convErr
		} else {
			resStr = (string)(convRes)
		}
	}
	return resStr, err
}

func (c *LambdaFunction) Handler(ctx context.Context, event map[string]interface{}) (string, error) { //handler(event: any, context: any) {
	// If already started then execute
	if c.IsOpen() {
		if event != nil {
			return c.execute(ctx, event)
		}
	} else { // Start before execute
		err := c.Run()
		if err != nil {
			ctx.Done()
			return "", err
		}
		if event != nil {
			return c.execute(ctx, event)
		}
	}
	err := cerr.NewBadRequestError(
		"Lambda",
		"NO_EVENT",
		"Event is empty")
	return "ERROR", err
}

/*
Gets entry point into this lambda function.
   - event     an incoming event object with invocation parameters.
   - context   a context object with local references.
*/

func (c *LambdaFunction) GetHandler() func(ctx context.Context, event map[string]interface{}) (string, error) {

	// Return plugin function
	return func(ctx context.Context, event map[string]interface{}) (string, error) {
		// Calling run with changed context
		return c.Handler(ctx, event)
	}
}

/*
Calls registered action in this lambda function.
"cmd" parameter in the action parameters determin
what action shall be called.

This method shall only be used in testing.
   - params action parameters.
   - callback callback function that receives action result or error.
*/

func (c *LambdaFunction) Act(params map[string]interface{}) (string, error) {
	ctx := context.TODO()
	return c.GetHandler()(ctx, params)
}

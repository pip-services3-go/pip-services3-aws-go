package container

import (
	ccomands "github.com/pip-services3-go/pip-services3-commons-go/commands"
	crun "github.com/pip-services3-go/pip-services3-commons-go/run"
)

/*
Abstract AWS Lambda function, that acts as a container to instantiate and run components
and expose them via external entry point. All actions are automatically generated for commands
defined in ICommandable components. Each command is exposed as an action defined by "cmd" parameter.

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

    class MyLambdaFunction extends CommandableLambdaFunction {
        private _controller: IMyController;
        ...
        func (c* CommandableLambdaFunction) constructor() {
            base("mygroup", "MyGroup lambda function");
            c.dependencyResolver.put(
                "controller",
                new Descriptor("mygroup","controller","*","*","1.0")
            );
        }
    }

    let lambda = new MyLambdaFunction();

    service.run((err) => {
        console.log("MyLambdaFunction is started");
    });
*/
type CommandableLambdaFunction struct {
	*LambdaFunction
}

/*
Creates a new instance of this lambda function.

   - name          (optional) a container name (accessible via ContextInfo)
   - description   (optional) a container description (accessible via ContextInfo)
*/
func NewCommandableLambdaFunction(name string, description string) *CommandableLambdaFunction {
	c := &CommandableLambdaFunction{}
	c.LambdaFunction = InheriteLambdaFunction(name, description, c)
	c.DependencyResolver.Put("controller", "none")
	return c
}

func (c *CommandableLambdaFunction) registerCommandSet(commandSet *ccomands.CommandSet) {
	commands := commandSet.Commands()
	for index := 0; index < len(commands); index++ {
		command := commands[index]

		c.RegisterAction(command.Name(), nil, func(params map[string]interface{}) (result interface{}, err error) {

			correlationId, _ := params["correlation_id"].(string)

			args := crun.NewParametersFromValue(params)
			timing := c.Instrument(correlationId, c.Info().Name+"."+command.Name())
			execRes, execErr := command.Execute(correlationId, args)
			timing.EndTiming()
			instrRes, instrErr := c.InstrumentError(correlationId,
				c.Info().Name+"."+command.Name(),
				execErr, execRes)
			return instrRes, instrErr
		})
	}
}

/*
Registers all actions in this lambda function.
*/
func (c *CommandableLambdaFunction) Register() {

	ref, _ := c.DependencyResolver.GetOneRequired("controller")
	controller, ok := ref.(ccomands.ICommandable)
	if ok {
		commandSet := controller.GetCommandSet()
		c.registerCommandSet(commandSet)
	}
}

package clients

import (
	"encoding/json"
	"reflect"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	awscon "github.com/pip-services3-go/pip-services3-aws-go/connect"
	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	cerr "github.com/pip-services3-go/pip-services3-commons-go/errors"
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	ccount "github.com/pip-services3-go/pip-services3-components-go/count"
	clog "github.com/pip-services3-go/pip-services3-components-go/log"
	ctrace "github.com/pip-services3-go/pip-services3-components-go/trace"
)

/*
Abstract client that calls AWS Lambda Functions.

When making calls "cmd" parameter determines which what action shall be called, while
other parameters are passed to the action itself.

### Configuration parameters ###

 - connections:
     - discovery_key:               (optional) a key to retrieve the connection from IDiscovery
     - region:                      (optional) AWS region
 - credentials:
     - store_key:                   (optional) a key to retrieve the credentials from ICredentialStore
     - access_id:                   AWS access/client id
     - access_key:                  AWS access/client id
 - options:
     - connect_timeout:             (optional) connection timeout in milliseconds (default: 10 sec)

### References ###

 - \*:logger:\*:\*:1.0            (optional) ILogger components to pass log messages
 - \*:counters:\*:\*:1.0          (optional) ICounters components to pass collected measurements
 - \*:discovery:\*:\*:1.0         (optional) IDiscovery services to resolve connection
 - \*:credential-store:\*:\*:1.0  (optional) Credential stores to resolve credentials

 See LambdaFunction
 See CommandableLambdaClient

### Example ###

     type MyLambdaClient struct  {
        *LambdaClient
         ...
	 }
         func (c* MyLambdaClient) getData(correlationId string, id string)(result MyData, err error){

            timing := c.Instrument(correlationId, "myclient.get_data");
            callRes, callErr := c.Call(MyDataPageType ,"get_data" correlationId, map[string]interface{ "id": id })
            timing.EndTiming();
            return callRes, callErr
         }
         ...


    client = NewMyLambdaClient();
    client.Configure(NewConfigParamsFromTuples(
        "connection.region", "us-east-1",
        "connection.access_id", "XXXXXXXXXXX",
        "connection.access_key", "XXXXXXXXXXX",
        "connection.arn", "YYYYYYYYYYYYY"
    ));

    data, err := client.GetData("123", "1",)
        ...
*/
type LambdaClient struct {
	// The reference to AWS Lambda Function.
	Lambda *lambda.Lambda
	// The opened flag.
	Opened bool
	// The AWS connection parameters
	Connection     *awscon.AwsConnectionParams
	connectTimeout int
	// The dependencies resolver.
	DependencyResolver *cref.DependencyResolver
	// The connection resolver.
	ConnectionResolver *awscon.AwsConnectionResolver
	// The logger.
	Logger *clog.CompositeLogger
	//The performance counters.
	Counters *ccount.CompositeCounters
	// The tracer.
	Tracer *ctrace.CompositeTracer
}

func NewLambdaClient() *LambdaClient {
	c := &LambdaClient{
		Opened:             false,
		connectTimeout:     10000,
		DependencyResolver: cref.NewDependencyResolver(),
		ConnectionResolver: awscon.NewAwsConnectionResolver(),
		Logger:             clog.NewCompositeLogger(),
		Counters:           ccount.NewCompositeCounters(),
	}
	return c
}

// Configures component by passing configuration parameters.
//   - config    configuration parameters to be set.
func (c *LambdaClient) Configure(config *cconf.ConfigParams) {
	c.ConnectionResolver.Configure(config)
	c.DependencyResolver.Configure(config)
	c.connectTimeout = config.GetAsIntegerWithDefault("options.connect_timeout", c.connectTimeout)
}

/*
 Sets references to dependent components.

 - references 	references to locate the component dependencies.
*/
func (c *LambdaClient) SetReferences(references cref.IReferences) {
	c.Logger.SetReferences(references)
	c.Counters.SetReferences(references)
	c.ConnectionResolver.SetReferences(references)
	c.DependencyResolver.SetReferences(references)
}

// Adds instrumentation to log calls and measure call time.
// It returns a Timing object that is used to end the time measurement.
//   - correlationId     (optional) transaction id to trace execution through call chain.
//   - name              a method name.
//  Returns Timing object to end the time measurement.
func (c *LambdaClient) Instrument(correlationId string, name string) *ccount.CounterTiming {
	c.Logger.Trace(correlationId, "Executing %s method", name)
	c.Counters.IncrementOne(name + ".exec_count")
	return c.Counters.BeginTiming(name + ".exec_time")
}

//  Checks if the component is opened.
//  Returns true if the component has been opened and false otherwise.
func (c *LambdaClient) IsOpen() bool {
	return c.Opened
}

// Opens the component.
//   - correlationId 	(optional) transaction id to trace execution through call chain.
//   - Return 			 error or nil no errors occured.
func (c *LambdaClient) Open(correlationId string) error {
	if c.IsOpen() {
		return nil
	}

	wg := sync.WaitGroup{}
	var errGlobal error

	wg.Add(1)
	go func() {
		defer wg.Done()
		connection, err := c.ConnectionResolver.Resolve(correlationId)
		c.Connection = connection
		errGlobal = err

		awsCred := credentials.NewStaticCredentials(c.Connection.GetAccessId(), c.Connection.GetAccessKey(), "")
		sess := session.Must(session.NewSession(&aws.Config{
			MaxRetries:  aws.Int(3),
			Region:      aws.String(c.Connection.GetRegion()),
			Credentials: awsCred,
		}))
		// Create new cloudwatch client.
		c.Lambda = lambda.New(sess)
		c.Lambda.Config.HTTPClient.Timeout = time.Duration((int64)(c.connectTimeout)) * time.Millisecond
		c.Logger.Debug(correlationId, "Lambda client connected to %s", c.Connection.GetArn())

	}()
	wg.Wait()
	if errGlobal != nil {
		c.Opened = false
		return errGlobal
	}
	return nil
}

// Closes component and frees used resources.
//   - correlationId 	(optional) transaction id to trace execution through call chain.
//   - Returns 			 error or null no errors occured.
func (c *LambdaClient) Close(correlationId string) error {
	// Todo: close listening?
	c.Opened = false
	c.Lambda = nil
	return nil
}

// Performs AWS Lambda Function invocation.
// 	 - prototype reflect.Type type for convert result. Set nil for return raw []byte
//   - invocationType    an invocation type: "RequestResponse" or "Event"
//   - cmd               an action name to be called.
//   - correlationId 	(optional) transaction id to trace execution through call chain.
//   - args              action arguments
// Returns           result or error.

func (c *LambdaClient) Invoke(prototype reflect.Type, invocationType string, cmd string, correlationId string, args map[string]interface{}) (result interface{}, err error) {

	if cmd == "" {
		err = cerr.NewUnknownError("", "NO_COMMAND", "Missing cmd")
		c.Logger.Error(correlationId, err, "Failed to call %s", cmd)
		return nil, err
	}

	//args = _.clone(args)

	args["cmd"] = cmd
	if correlationId != "" {
		args["correlation_id"] = correlationId
	} else {
		args["correlation_id"] = cdata.IdGenerator.NextLong()
	}
	payloads, jsonErr := json.Marshal(args)

	if jsonErr != nil {
		c.Logger.Error(correlationId, jsonErr, "Failed to call %s", cmd)
		return nil, jsonErr
	}

	params := &lambda.InvokeInput{
		FunctionName:   aws.String(c.Connection.GetArn()),
		InvocationType: aws.String(invocationType),
		LogType:        aws.String("None"),
		Payload:        payloads,
	}

	data, lambdaErr := c.Lambda.Invoke(params)

	if lambdaErr != nil {
		err = cerr.NewInvocationError(
			correlationId,
			"CALL_FAILED",
			"Failed to invoke lambda function").WithCause(err)
		return nil, err
	}

	if prototype != nil {
		return ConvertComandResult(data.Payload, prototype)
	}
	return data.Payload, nil

}

// Calls a AWS Lambda Function action.
// 	 - prototype reflect.Type type for convert result. Set nil for return raw []byte
//   - cmd               an action name to be called.
//   - correlationId     (optional) transaction id to trace execution through call chain.
//   - params            (optional) action parameters.
//   - Returns           result and error.
func (c *LambdaClient) Call(prototype reflect.Type, cmd string, correlationId string, params map[string]interface{}) (result interface{}, err error) {
	return c.Invoke(prototype, "RequestResponse", cmd, correlationId, params)
}

// Calls a AWS Lambda Function action asynchronously without waiting for response.
// 	 - prototype reflect.Type type for convert result. Set nil for return raw []byte
//   - cmd               an action name to be called.
//   - correlationId     (optional) transaction id to trace execution through call chain.
//   - params            (optional) action parameters.
//   - Returns           error or null for success.
func (c *LambdaClient) CallOneWay(prototype reflect.Type, cmd string, correlationId string, params map[string]interface{}) error {
	_, err := c.Invoke(prototype, "Event", cmd, correlationId, params)
	return err
}

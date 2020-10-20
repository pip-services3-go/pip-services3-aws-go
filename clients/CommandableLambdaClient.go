package clients

import "reflect"

/*
 Abstract client that calls commandable AWS Lambda Functions.

 Commandable services are generated automatically for [[ICommandable objects]].
 Each command is exposed as action determined by "cmd" parameter.

 ### Configuration parameters ###

 - connections:
     - discovery_key:               (optional) a key to retrieve the connection from [[IDiscovery]]
     - region:                      (optional) AWS region
 - credentials:
     - store_key:                   (optional) a key to retrieve the credentials from [[ICredentialStore]]
     - access_id:                   AWS access/client id
     - access_key:                  AWS access/client id
 - options:
     - connect_timeout:             (optional) connection timeout in milliseconds (default: 10 sec)

 ### References ###

 - \*:logger:\*:\*:1.0            (optional) [[ILogger]] components to pass log messages
 - \*:counters:\*:\*:1.0          (optional) [[ICounters]] components to pass collected measurements
 - \*:discovery:\*:\*:1.0         (optional) [[IDiscovery]] services to resolve connection
 - \*:credential-store:\*:\*:1.0  (optional) Credential stores to resolve credentials

 See [[LambdaFunction]]

 ### Example ###

     type MyLambdaClient struct {
		 *CommandableLambdaClient
	 }
         ...

         func (c* MyLambdaClient) GetData(correlationId string, id string)(result MyDataPage, err error) {

           return c.callCommand(MyDataPageType,
                 "get_data",
                 correlationId,
                 map[string]interface{}{ "id": id })

         }
         ...


      client := NewMyLambdaClient();
     client.Configure(NewConfigParamsFromTuples(
         "connection.region", "us-east-1",
         "connection.access_id", "XXXXXXXXXXX",
         "connection.access_key", "XXXXXXXXXXX",
         "connection.arn", "YYYYYYYYYYYYY"
     ));

     res, err := client.GetData("123", "1")
         ...

*/
type CommandableLambdaClient struct {
	*LambdaClient
	name string
}

//  Creates a new instance of this client.
//  - name a service name.
func NewCommandableLambdaClient(name string) *CommandableLambdaClient {
	c := &CommandableLambdaClient{
		LambdaClient: NewLambdaClient(),
	}
	c.name = name
	return c
}

//    Calls a remote action in AWS Lambda function.
//    The name of the action is added as "cmd" parameter
//    to the action parameters.
//    - prototype reflect.Type type for convert result. Set nil for return raw []byte
//    - cmd               an action name
//    - correlationId     (optional) transaction id to trace execution through call chain.
//    - params            command parameters.
//    - Return           result or error.
func (c *CommandableLambdaClient) CallCommand(prototype reflect.Type, cmd string, correlationId string, params map[string]interface{}) (result interface{}, err error) {
	timing := c.Instrument(correlationId, c.name+"."+cmd)
	callRes, callErr := c.Call(prototype, cmd, correlationId, params)
	timing.EndTiming()
	return callRes, callErr

}

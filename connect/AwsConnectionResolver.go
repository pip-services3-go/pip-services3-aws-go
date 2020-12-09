package connect

import (
	"sync"

	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	cauth "github.com/pip-services3-go/pip-services3-components-go/auth"
	ccon "github.com/pip-services3-go/pip-services3-components-go/connect"
)

/*
Helper class to retrieve AWS connection and credential parameters,
validate them and compose a AwsConnectionParams value.

### Configuration parameters ###

 - connections:
     - discovery_key:               (optional) a key to retrieve the connection from IDiscovery
     - region:                      (optional) AWS region
     - partition:                   (optional) AWS partition
     - service:                     (optional) AWS service
     - resource_type:               (optional) AWS resource type
     - resource:                    (optional) AWS resource id
     - arn:                         (optional) AWS resource ARN
 - credentials:
     - store_key:                   (optional) a key to retrieve the credentials from ICredentialStore
     - access_id:                   AWS access/client id
     - access_key:                  AWS access/client id

### References ###

 - \*:discovery:\*:\*:1.0         (optional) IDiscovery services to resolve connections
 - \*:credential-store:\*:\*:1.0  (optional) Credential stores to resolve credentials

 See ConnectionParams (in the Pip.Services components package)
 See IDiscovery (in the Pip.Services components package)

 ### Example ###

    config := NewConfigParamsFromTuples(
         "connection.region", "us-east1",
         "connection.service", "s3",
         "connection.bucket", "mybucket",
         "credential.access_id", "XXXXXXXXXX",
         "credential.access_key", "XXXXXXXXXX"
     );

    connectionResolver := NewAwsConnectionResolver();
    connectionResolver.Configure(config);
    connectionResolver.SetReferences(references);

    err, connection :=connectionResolver.Resolve("123")
        // Now use connection...
*/
type AwsConnectionResolver struct {

	//The connection resolver.
	connectionResolver *ccon.ConnectionResolver

	//The credential resolver.
	credentialResolver *cauth.CredentialResolver
}

func NewAwsConnectionResolver() *AwsConnectionResolver {
	return &AwsConnectionResolver{
		connectionResolver: ccon.NewEmptyConnectionResolver(),
		credentialResolver: cauth.NewEmptyCredentialResolver(),
	}
}

// Configures component by passing configuration parameters.
//   - config    configuration parameters to be set.
func (c *AwsConnectionResolver) Configure(config *cconf.ConfigParams) {
	c.connectionResolver.Configure(config)
	c.credentialResolver.Configure(config)
}

//  Sets references to dependent components.
//    - references 	references to locate the component dependencies.
func (c *AwsConnectionResolver) SetReferences(references cref.IReferences) {
	c.connectionResolver.SetReferences(references)
	c.credentialResolver.SetReferences(references)
}

/*
Resolves connection and credental parameters and generates a single
AWSConnectionParams value.
   - correlationId     (optional) transaction id to trace execution through call chain.
   - callback 			callback function that receives AWSConnectionParams value or error.
See IDiscovery (in the Pip.Services components package)
*/
func (c *AwsConnectionResolver) Resolve(correlationId string) (connection *AwsConnectionParams, err error) {
	connection = NewEmptyAwsConnectionParams()
	//var credential *cauth.CredentialParams
	var globalErr error

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		data, connErr := c.connectionResolver.Resolve(correlationId)
		if connErr == nil && data != nil {
			connection.Append(data.Value())
		}
		globalErr = connErr
	}()
	wg.Wait()

	if globalErr != nil {
		return nil, globalErr
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		data, credErr := c.credentialResolver.Lookup(correlationId)
		if credErr == nil && data != nil {
			connection.Append(data.Value())
		}
		globalErr = credErr
	}()
	wg.Wait()

	if globalErr != nil {
		return nil, globalErr
	}
	// Force ARN parsing
	connection.SetArn(connection.GetArn())
	// Perform validation
	validErr := connection.Validate(correlationId)

	if validErr != nil {
		return nil, validErr
	}
	return connection, nil
}

package log

import (
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	awsconn "github.com/pip-services3-go/pip-services3-aws-go/connect"
	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cerr "github.com/pip-services3-go/pip-services3-commons-go/errors"
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	cinfo "github.com/pip-services3-go/pip-services3-components-go/info"
	clog "github.com/pip-services3-go/pip-services3-components-go/log"
)

/*
 Logger that writes log messages to AWS Cloud Watch Log.

 ### Configuration parameters ###

 - stream:                        (optional) Cloud Watch Log stream (default: context name)
 - group:                         (optional) Cloud Watch Log group (default: context instance ID or hostname)
 - connections:
     - discovery_key:               (optional) a key to retrieve the connection from [[IDiscovery]]
     - region:                      (optional) AWS region
 - credentials:
     - store_key:                   (optional) a key to retrieve the credentials from [[ICredentialStore]]
     - access_id:                   AWS access/client id
     - access_key:                  AWS access/client id
 - options:
     - interval:        interval in milliseconds to save current counters measurements (default: 5 mins)
     - reset_timeout:   timeout in milliseconds to reset the counters. 0 disables the reset (default: 0)

 ### References ###

 - \*:context-info:\*:\*:1.0      (optional) [[ContextInfo]] to detect the context id and specify counters source
 - \*:discovery:\*:\*:1.0         (optional) [[IDiscovery]] services to resolve connections
 - \*:credential-store:\*:\*:1.0  (optional) Credential stores to resolve credentials

 See [[Counter]] (in the Pip.Services components package)
 See [[CachedCounters]] (in the Pip.Services components package)
 See [[CompositeLogger]] (in the Pip.Services components package)


 ### Example ###

     logger := NewLogger();
     logger.Config(NewConfigParamsFromTuples(
         "stream", "mystream",
         "group", "mygroup",
         "connection.region", "us-east-1",
         "connection.access_id", "XXXXXXXXXXX",
         "connection.access_key", "XXXXXXXXXXX",
     ));
     logger.SetReferences(NewReferencesFromTuples(
         NewDescriptor("pip-services", "logger", "console", "default", "1.0"),
         NewConsoleLogger()
     ));

     err:= logger.Open("123")
         ...

     logger.SetLevel(Debug);

     logger.Error("123", ex, "Error occured: %s", ex.Message);
     logger.Debug("123", "Everything is OK.");
*/
type CloudWatchLogger struct {
	*clog.CachedLogger

	timer chan bool

	connectionResolver *awsconn.AwsConnectionResolver
	client             *cloudwatchlogs.CloudWatchLogs //AmazonCloudWatchLogsClient
	connection         *awsconn.AwsConnectionParams
	connectTimeout     int

	group     string
	stream    string
	lastToken string

	logger *clog.CompositeLogger
}

/*
   Creates a new instance of this logger.
*/
func NewCloudWatchLogger() *CloudWatchLogger {
	c := &CloudWatchLogger{
		connectionResolver: awsconn.NewAwsConnectionResolver(),
		connectTimeout:     30000,
		group:              "undefined",
		stream:             "",
		lastToken:          "",
		logger:             clog.NewCompositeLogger(),
	}
	c.CachedLogger = clog.InheritCachedLogger(c)
	return c
}

//  Configure method configures component by passing configuration parameters.
//  - config    configuration parameters to be set.
func (c *CloudWatchLogger) Configure(config *cconf.ConfigParams) {
	c.CachedLogger.Configure(config)
	c.connectionResolver.Configure(config)

	c.group = config.GetAsStringWithDefault("group", c.group)
	c.stream = config.GetAsStringWithDefault("stream", c.stream)
	c.connectTimeout = config.GetAsIntegerWithDefault("options.connect_timeout", c.connectTimeout)
}

//  SetReferences method sets references to dependent components.
//  - references 	references to locate the component dependencies.
//  See [[IReferences]] (in the Pip.Services commons package)
func (c *CloudWatchLogger) SetReferences(references cref.IReferences) {
	c.CachedLogger.SetReferences(references)
	c.logger.SetReferences(references)

	ref := references.GetOneOptional(cref.NewDescriptor("pip-services", "context-info", "default", "*", "1.0"))

	contextInfo, ok := ref.(*cinfo.ContextInfo)
	if ok && c.stream == "" {
		c.stream = contextInfo.Name
	}
	if ok && c.group == "" {
		c.group = contextInfo.ContextId
	}
}

//  Writes a log message to the logger destination.
//  - level             a log level.
//  - correlationId     (optional) transaction id to trace execution through call chain.
//  - error             an error object associated with this message.
//  - message           a human-readable message to log.
func (c *CloudWatchLogger) Write(level int, correlationId string, ex error, message string) {
	if c.Level() < level {
		return
	}
	c.CachedLogger.Write(level, correlationId, ex, message)
}

//  Checks if the component is opened.
//  Returns true if the component has been opened and false otherwise.
func (c *CloudWatchLogger) IsOpen() bool {
	return c.timer != nil
}

/*
	 Opens the component.
	 *
	 - correlationId 	(optional) transaction id to trace execution through call chain.
     - Returns 			 error or nil no errors occured.
*/
func (c *CloudWatchLogger) Open(correlationId string) error {
	if c.IsOpen() {
		return nil
	}

	var globalErr error

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		connection, err := c.connectionResolver.Resolve(correlationId)
		c.connection = connection
		if err != nil {
			globalErr = err
			return
		}

		awsCred := credentials.NewStaticCredentials(c.connection.GetAccessId(), c.connection.GetAccessKey(), "")
		sess := session.Must(session.NewSession(&aws.Config{
			MaxRetries:  aws.Int(3),
			Region:      aws.String(c.connection.GetRegion()),
			Credentials: awsCred,
		}))
		// Create new cloudwatch client.
		c.client = cloudwatchlogs.New(sess)
		c.client.APIVersion = "2014-03-28"
		c.client.Config.HTTPClient.Timeout = time.Duration((int64)(c.connectTimeout)) * time.Millisecond

		groupParam := &cloudwatchlogs.CreateLogGroupInput{
			LogGroupName: aws.String(c.group),
		}
		_, groupErr := c.client.CreateLogGroup(groupParam)
		if _, ok := groupErr.(*cloudwatchlogs.ResourceAlreadyExistsException); !ok {
			globalErr = groupErr
			return
		}

		streamParam := &cloudwatchlogs.CreateLogStreamInput{
			LogGroupName:  aws.String(c.group),
			LogStreamName: aws.String(c.stream),
		}
		_, streamErr := c.client.CreateLogStream(streamParam)

		if streamErr != nil {
			if _, ok := streamErr.(*cloudwatchlogs.ResourceAlreadyExistsException); ok {

				params := &cloudwatchlogs.DescribeLogStreamsInput{
					LogGroupName:        aws.String(c.group),
					LogStreamNamePrefix: aws.String(c.stream),
				}

				descData, describeErr := c.client.DescribeLogStreams(params)
				if describeErr != nil {
					globalErr = describeErr
					return
				}
				if len(descData.LogStreams) > 0 {
					if descData.LogStreams[0].UploadSequenceToken != nil {
						c.lastToken = *descData.LogStreams[0].UploadSequenceToken
					}
				}
			} else {
				globalErr = streamErr
				return
			}
		} else {
			c.lastToken = ""
		}

		if c.timer == nil {
			c.timer = setInterval(func() { c.Dump() }, c.Interval, true)
		}
	}()
	wg.Wait()

	if globalErr != nil {
		return globalErr
	}
	return nil
}

/*
	 Closes component and frees used resources.
	 - correlationId 	(optional) transaction id to trace execution through call chain.
     - Returns 		   error or nil no errors occured.
*/
func (c *CloudWatchLogger) Close(correlationId string) error {
	err := c.Save(c.Cache)

	if c.timer != nil {
		c.timer <- true
	}

	c.Cache = make([]*clog.LogMessage, 0)
	c.timer = nil
	c.client = nil

	return err
}

func (c *CloudWatchLogger) formatMessageText(message *clog.LogMessage) string {

	result := "["

	if message.Source != "" {
		result += message.Source
	} else {
		result += "---"
	}

	result += ":"

	if message.CorrelationId != "" {
		result += message.CorrelationId
	} else {
		result += "---"
	}
	result += ":" + clog.LogLevelConverter.ToString(message.Level) + "] " + message.Message
	if message.Error.Message != "" || message.Error.Code != "" {
		if message.Message == "" {
			result += "Error: "
		} else {
			result += ": "
		}
		result += message.Error.Message
		if message.Error.StackTrace != "" {
			result += " StackTrace: " + message.Error.StackTrace
		}
	}
	return result
}

/*
 Saves log messages from the cache.
 *
 - messages  a list with log messages
 - Returns   error or nil for success.
*/
func (c *CloudWatchLogger) Save(messages []*clog.LogMessage) error {
	if !c.IsOpen() || messages == nil || len(messages) == 0 {
		return nil
	}

	if c.client == nil {
		err := cerr.NewConfigError("cloudwatch_logger", "NOT_OPENED", "CloudWatchLogger is not opened")
		if err != nil {
			return err
		}
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		var events []*cloudwatchlogs.InputLogEvent
		events = make([]*cloudwatchlogs.InputLogEvent, 0)

		for _, message := range messages {
			events = append(events, &cloudwatchlogs.InputLogEvent{
				Timestamp: aws.Int64(message.Time.UnixNano() / (int64)(time.Millisecond)),
				Message:   aws.String(c.formatMessageText(message)),
			})
		}

		// get token again if saving log from another container
		describeParams := &cloudwatchlogs.DescribeLogStreamsInput{
			LogGroupName:        aws.String(c.group),
			LogStreamNamePrefix: aws.String(c.stream),
		}

		data, _ := c.client.DescribeLogStreams(describeParams)
		if len(data.LogStreams) > 0 {
			if data.LogStreams[0].UploadSequenceToken != nil {
				c.lastToken = *data.LogStreams[0].UploadSequenceToken
			}
		}
		var token *string = nil
		if c.lastToken != "" {
			token = &c.lastToken
		}

		params := &cloudwatchlogs.PutLogEventsInput{
			LogEvents:     events,
			LogGroupName:  aws.String(c.group),
			LogStreamName: aws.String(c.stream),
			SequenceToken: token,
		}

		putRes, putErr := c.client.PutLogEvents(params)
		if putErr != nil {
			if c.logger != nil {
				c.logger.Error("cloudwatch_logger", putErr, "putLogEvents error")
			}
		} else {
			if putRes.NextSequenceToken != nil {
				c.lastToken = *putRes.NextSequenceToken
			}
		}
	}()

	wg.Wait()

	return nil
}

func setInterval(someFunc func(), milliseconds int, async bool) chan bool {

	interval := time.Duration(milliseconds) * time.Millisecond
	ticker := time.NewTicker(interval)
	clear := make(chan bool)
	go func() {
		for {
			select {
			case <-ticker.C:
				if async {
					go someFunc()
				} else {
					someFunc()
				}
			case <-clear:
				ticker.Stop()
				return
			}

		}
	}()

	return clear
}

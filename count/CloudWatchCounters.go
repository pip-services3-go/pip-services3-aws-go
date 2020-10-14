package count

import (
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	awsconn "github.com/pip-services3-go/pip-services3-aws-go/connect"
	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	ccount "github.com/pip-services3-go/pip-services3-components-go/count"
	cinfo "github.com/pip-services3-go/pip-services3-components-go/info"
	clog "github.com/pip-services3-go/pip-services3-components-go/log"
)

/*
 Performance counters that periodically dumps counters to AWS Cloud Watch Metrics.

 ### Configuration parameters ###

 - connections:
     - discovery_key:         (optional) a key to retrieve the connection from [[IDiscovery]]
     - region:                (optional) AWS region
 - credentials:
     - store_key:             (optional) a key to retrieve the credentials from [[ICredentialStore]]
     - access_id:             AWS access/client id
     - access_key:            AWS access/client id
 - options:
     - interval:              interval in milliseconds to save current counters measurements (default: 5 mins)
     - reset_timeout:         timeout in milliseconds to reset the counters. 0 disables the reset (default: 0)

 ### References ###

 - \*:context-info:\*:\*:1.0      (optional) [[ContextInfo]] to detect the context id and specify counters source
 - \*:discovery:\*:\*:1.0         (optional) [[IDiscovery]] services to resolve connections
 - \*:credential-store:\*:\*:1.0  (optional) Credential stores to resolve credentials

 See [[Counter]] (in the Pip.Services components package)
 See [[CachedCounters]] (in the Pip.Services components package)
 See [[CompositeLogger]] (in the Pip.Services components package)

 ### Example ###

      counters := NewCloudWatchCounters();
     counters.Config(ConfigParams.fromTuples(
         "connection.region", "us-east-1",
         "connection.access_id", "XXXXXXXXXXX",
         "connection.access_key", "XXXXXXXXXXX"
     ));
     counters.SetReferences(NewReferencesFromTuples(
         NewDescriptor("pip-services", "logger", "console", "default", "1.0"),
         NewConsoleLogger()
     ));

     err := counters.Open("123")
         ...

     counters.Increment("mycomponent.mymethod.calls");
      timing:= counters.BeginTiming("mycomponent.mymethod.exec_time");

         ...

         timing.endTiming();

     counters.Dump();
*/
type CloudWatchCounters struct {
	*ccount.CachedCounters
	logger *clog.CompositeLogger

	connectionResolver *awsconn.AwsConnectionResolver
	connection         *awsconn.AwsConnectionParams
	connectTimeout     int
	client             *cloudwatch.CloudWatch //AmazonCloudWatchClient
	source             string
	instance           string
	opened             bool
}

// Creates a new instance of this counters.
func NewCloudWatchCounters() *CloudWatchCounters {
	c := &CloudWatchCounters{
		logger:             clog.NewCompositeLogger(),
		connectionResolver: awsconn.NewAwsConnectionResolver(),
		connectTimeout:     30000,
		opened:             false,
	}
	c.CachedCounters = ccount.InheritCacheCounters(c)
	return c
}

//  Configures component by passing configuration parameters.
//  - config    configuration parameters to be set.
func (c *CloudWatchCounters) Configure(config *cconf.ConfigParams) {
	c.CachedCounters.Configure(config)
	c.connectionResolver.Configure(config)

	c.source = config.GetAsStringWithDefault("source", c.source)
	c.instance = config.GetAsStringWithDefault("instance", c.instance)
	c.connectTimeout = config.GetAsIntegerWithDefault("options.connect_timeout", c.connectTimeout)
}

/*
 Sets references to dependent components.
 - references 	references to locate the component dependencies.
 See [[IReferences]] (in the Pip.Services commons package)
*/
func (c *CloudWatchCounters) SetReferences(references cref.IReferences) {
	c.logger.SetReferences(references)
	c.connectionResolver.SetReferences(references)
	ref := references.GetOneOptional(
		cref.NewDescriptor("pip-services", "context-info", "default", "*", "1.0"))
	contextInfo, ok := ref.(*cinfo.ContextInfo)

	if ok && c.source == "" {
		c.source = contextInfo.Name
	}

	if ok && c.instance == "" {
		c.instance = contextInfo.ContextId
	}

}

// Checks if the component is opened.
// Returns true if the component has been opened and false otherwise.
func (c *CloudWatchCounters) IsOpen() bool {
	return c.opened
}

/*
 Opens the component.
 - correlationId 	(optional) transaction id to trace execution through call chain.
 - Returns 			 error or null no errors occured.
*/
func (c *CloudWatchCounters) Open(correlationId string) error {
	if c.opened {
		return nil
	}

	c.opened = true

	wg := sync.WaitGroup{}
	var errGlobal error

	wg.Add(1)
	go func() {
		defer wg.Done()
		connection, err := c.connectionResolver.Resolve(correlationId)
		c.connection = connection
		errGlobal = err
	}()
	wg.Wait()
	if errGlobal != nil {
		c.opened = false
		return errGlobal
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		awsCred := credentials.NewStaticCredentials(c.connection.GetAccessId(), c.connection.GetAccessKey(), "")
		sess := session.Must(session.NewSession(&aws.Config{
			MaxRetries:  aws.Int(3),
			Region:      aws.String(c.connection.GetRegion()),
			Credentials: awsCred,
		}))
		// Create new cloudwatch client.
		c.client = cloudwatch.New(sess)
		c.client.APIVersion = "2010-08-01"
		c.client.Config.HTTPClient.Timeout = time.Duration((int64)(c.connectTimeout)) * time.Millisecond

	}()
	wg.Wait()
	if errGlobal != nil {
		c.opened = false
		return errGlobal
	}
	return nil
}

/*
 Closes component and frees used resources.
 - correlationId 	(optional) transaction id to trace execution through call chain.
 - Return 			 error or nil no errors occured.
*/
func (c *CloudWatchCounters) Close(correlationId string) error {
	c.opened = false
	c.client = nil
	return nil
}

func (c *CloudWatchCounters) getCounterData(counter *ccount.Counter, now time.Time, dimensions []*cloudwatch.Dimension) *cloudwatch.MetricDatum {

	value := &cloudwatch.MetricDatum{
		MetricName: aws.String(counter.Name),
		Unit:       aws.String(None),
		Dimensions: dimensions,
	}
	tm := counter.Time
	if tm.IsZero() {
		tm = time.Now()
	}
	value.SetTimestamp(tm)

	switch counter.Type {
	case ccount.Increment:
		value.Value = aws.Float64((float64)(counter.Count))
		value.Unit = aws.String(Count)
		break
	case ccount.Interval:
		value.Unit = aws.String(Milliseconds)
		//value.Value = counter.average;
		value.StatisticValues = &cloudwatch.StatisticSet{
			SampleCount: aws.Float64((float64)(counter.Count)),
			Maximum:     aws.Float64((float64)(counter.Max)),
			Minimum:     aws.Float64((float64)(counter.Min)),
			Sum:         aws.Float64((float64)(counter.Count) * (float64)(counter.Average)),
		}
		break
	case ccount.Statistics:
		//value.Value = counter.average;
		value.StatisticValues = &cloudwatch.StatisticSet{
			SampleCount: aws.Float64((float64)(counter.Count)),
			Maximum:     aws.Float64((float64)(counter.Max)),
			Minimum:     aws.Float64((float64)(counter.Min)),
			Sum:         aws.Float64((float64)(counter.Count) * (float64)(counter.Average)),
		}
		break
	case ccount.LastValue:
		value.Value = aws.Float64((float64)(counter.Last))
		break
	case ccount.Timestamp:
		value.Value = aws.Float64((float64)(counter.Time.UnixNano()) / (float64)(time.Millisecond.Milliseconds())) // Convert to milliseconds UnixTimeStamp
		break
	}

	return value
}

/*
 Saves the current counters measurements.
 - counters      current counters measurements to be saves.
*/
func (c *CloudWatchCounters) Save(counters []*ccount.Counter) error {
	if c.client == nil {
		return nil
	}

	var dimensions []*cloudwatch.Dimension
	dimensions = make([]*cloudwatch.Dimension, 0)
	dimensions = append(dimensions, &cloudwatch.Dimension{
		Name:  aws.String("InstanceID"),
		Value: aws.String(c.instance),
	})

	now := time.Now()

	var data []*cloudwatch.MetricDatum
	data = make([]*cloudwatch.MetricDatum, 0)

	params := &cloudwatch.PutMetricDataInput{
		MetricData: data,
		Namespace:  aws.String(c.source),
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		for _, counter := range counters {
			data = append(data, c.getCounterData(counter, now, dimensions))

			if len(data) >= 20 {
				params.MetricData = data
				_, err := c.client.PutMetricData(params)
				if err != nil && c.logger != nil {
					c.logger.Error("cloudwatch_counters", err, "putMetricData error")
				}
				data = make([]*cloudwatch.MetricDatum, 0)
			}
		}

	}()

	wg.Wait()

	params.MetricData = data

	if len(data) > 0 {
		_, err := c.client.PutMetricData(params)
		if err != nil && c.logger != nil {
			c.logger.Error("cloudwatch_counters", err, "putMetricData error")
		}
	}
	return nil
}

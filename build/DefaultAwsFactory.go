package build

import (
	awscount "github.com/pip-services3-go/pip-services3-aws-go/count"
	awslog "github.com/pip-services3-go/pip-services3-aws-go/log"
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	cbuild "github.com/pip-services3-go/pip-services3-components-go/build"
)

/*
 Creates AWS components by their descriptors.
 *
 See [[CloudWatchLogger]]
 See [[CloudWatchCounters]]
*/
type DefaultAwsFactory struct {
	cbuild.Factory

	Descriptor                   *cref.Descriptor
	CloudWatchLoggerDescriptor   *cref.Descriptor
	CloudWatchCountersDescriptor *cref.Descriptor
}

// NewDefaultAwsFactory method are create a new instance of the factory.
func NewDefaultAwsFactory() *DefaultAwsFactory {

	c := &DefaultAwsFactory{
		Factory:                      *cbuild.NewFactory(),
		Descriptor:                   cref.NewDescriptor("pip-services", "factory", "aws", "default", "1.0"),
		CloudWatchLoggerDescriptor:   cref.NewDescriptor("pip-services", "logger", "cloudwatch", "*", "1.0"),
		CloudWatchCountersDescriptor: cref.NewDescriptor("pip-services", "counters", "cloudwatch", "*", "1.0"),
	}

	c.RegisterType(c.CloudWatchLoggerDescriptor, awslog.NewCloudWatchLogger)
	c.RegisterType(c.CloudWatchCountersDescriptor, awscount.NewCloudWatchCounters)
	return c
}

package test

import (
	"os"
	"testing"

	awslog "github.com/pip-services3-go/pip-services3-aws-go/log"
	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	cinfo "github.com/pip-services3-go/pip-services3-components-go/info"
)

func TestCloudWatchLogger(t *testing.T) {

	var loggers *awslog.CloudWatchLogger
	var fixture *LoggerFixture

	AWS_REGION := os.Getenv("AWS_REGION")
	AWS_ACCESS_ID := os.Getenv("AWS_ACCESS_ID")
	AWS_ACCESS_KEY := os.Getenv("AWS_ACCESS_KEY")

	if AWS_REGION == "" || AWS_ACCESS_ID == "" || AWS_ACCESS_KEY == "" {
		return
	}

	loggers = awslog.NewCloudWatchLogger()
	fixture = NewLoggerFixture(loggers.CachedLogger)

	config := cconf.NewConfigParamsFromTuples(
		"group", "TestGroup",
		"connection.region", AWS_REGION,
		"credential.access_id", AWS_ACCESS_ID,
		"credential.access_key", AWS_ACCESS_KEY,
	)
	loggers.Configure(config)

	contextInfo := cinfo.NewContextInfo()
	contextInfo.Name = "TestStream"
	contextInfo.Description = "This is a test container"

	var references = cref.NewReferencesFromTuples(
		cref.NewDescriptor("pip-services", "context-info", "default", "default", "1.0"), contextInfo,
		cref.NewDescriptor("pip-services", "loggers", "cloudwatch", "default", "1.0"), loggers,
	)
	loggers.SetReferences(references)
	loggers.Open("")
	defer loggers.Close("")

	t.Run("Log Level", fixture.TestLogLevel)
	t.Run("Simple Logging", fixture.TestSimpleLogging)
	t.Run("Error Logging", fixture.TestErrorLogging)
}

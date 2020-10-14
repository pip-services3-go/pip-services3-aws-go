package test

import (
	"errors"
	"testing"
	"time"

	clog "github.com/pip-services3-go/pip-services3-components-go/log"
	"github.com/stretchr/testify/assert"
)

type LoggerFixture struct {
	logger *clog.CachedLogger
}

func NewLoggerFixture(logger *clog.CachedLogger) *LoggerFixture {
	return &LoggerFixture{
		logger: logger,
	}

}

func (c *LoggerFixture) TestLogLevel(t *testing.T) {
	assert.True(t, c.logger.Level() >= clog.None)
	assert.True(t, c.logger.Level() <= clog.Trace)
}

func (c *LoggerFixture) TestSimpleLogging(t *testing.T) {
	c.logger.SetLevel(clog.Trace)

	c.logger.Fatal("", nil, "Fatal error message")
	c.logger.Error("", nil, "Error message")
	c.logger.Warn("", "Warning message")
	c.logger.Info("", "Information message")
	c.logger.Debug("", "Debug message")
	c.logger.Trace("", "Trace message")

	c.logger.Dump()
	time.Sleep(1000 * time.Millisecond)
}

func (c *LoggerFixture) TestErrorLogging(t *testing.T) {
	var ex error = errors.New("Testing error throw")
	c.logger.Fatal("123", ex, "Fatal error")
	c.logger.Error("123", ex, "Recoverable error")

	assert.NotNil(t, ex)

	c.logger.Dump()
	time.Sleep(1000 * time.Millisecond)
}

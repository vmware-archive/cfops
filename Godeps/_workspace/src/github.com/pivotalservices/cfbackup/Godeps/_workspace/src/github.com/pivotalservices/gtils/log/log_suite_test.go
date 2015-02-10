package log_test

import (
	"bytes"
	"encoding/json"
	"io"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/pivotal-golang/lager"
	"github.com/pivotalservices/gtils/log"

	"testing"
)

func TestLog(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Log Suite")
}

type TestLogger struct {
	log.Logger
	*TestSink
}

type TestSink struct {
	buffer *gbytes.Buffer
}

func NewTestLogger(component string) *TestLogger {
	buffer := gbytes.NewBuffer()
	testSink := NewTestSink(buffer)
	logger := log.LogFactory(component, log.Lager, buffer)
	return &TestLogger{logger, testSink}
}

func NewTestSink(buffer *gbytes.Buffer) *TestSink {
	return &TestSink{
		buffer: buffer,
	}
}

func (s *TestSink) Buffer() *gbytes.Buffer {
	return s.buffer
}

func (s *TestSink) Logs() []lager.LogFormat {
	logs := []lager.LogFormat{}

	decoder := json.NewDecoder(bytes.NewBuffer(s.buffer.Contents()))
	for {
		var log lager.LogFormat
		if err := decoder.Decode(&log); err == io.EOF {
			return logs
		} else if err != nil {
			panic(err)
		}
		logs = append(logs, log)
	}

	return logs
}

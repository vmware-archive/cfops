package log_test

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-golang/lager"
	"github.com/pivotalservices/gtils/log"
)

var _ = Describe("Logger", func() {
	var logger *TestLogger
	var testSink *TestSink

	var component = "my-component"
	var action = "my-action"
	var logData = log.Data{
		"foo":      "bar",
		"a-number": 7,
	}
	var anotherLogData = log.Data{
		"baz":      "quux",
		"b-number": 43,
	}

	BeforeEach(func() {
		log.SetLogLevel("debug")
		logger = NewTestLogger(component)
		testSink = logger.TestSink
	})

	var TestLagerFormat = func(level lager.LogLevel) {
		var log lager.LogFormat

		BeforeEach(func() {
			log = testSink.Logs()[0]
		})

		It("outputs a properly-formatted message", func() {
			Ω(log.Message).Should(Equal(fmt.Sprintf("%s.%s", component, action)))
		})

		It("has a timestamp", func() {
			expectedTime := float64(time.Now().UnixNano()) / 1e9
			parsedTimestamp, err := strconv.ParseFloat(log.Timestamp, 64)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(parsedTimestamp).Should(BeNumerically("~", expectedTime, 1.0))
		})

		It("sets the proper output level", func() {
			Ω(log.LogLevel).Should(Equal(level))
		})
	}

	var TestLagerData = func() {
		var log lager.LogFormat

		BeforeEach(func() {
			log = testSink.Logs()[0]
		})

		It("data contains custom user data", func() {
			Ω(log.Data["foo"]).Should(Equal("bar"))
			Ω(log.Data["a-number"]).Should(BeNumerically("==", 7))
			Ω(log.Data["baz"]).Should(Equal("quux"))
			Ω(log.Data["b-number"]).Should(BeNumerically("==", 43))
		})
	}

	Describe("Debug", func() {
		Context("with log data", func() {
			BeforeEach(func() {
				logger.Debug(action, logData, anotherLogData)
			})

			TestLagerFormat(lager.DEBUG)
			TestLagerData()
		})

		Context("with no log data", func() {
			BeforeEach(func() {
				logger.Debug(action)
			})

			TestLagerFormat(lager.DEBUG)
		})
	})

	Describe("Info", func() {
		Context("with log data", func() {
			BeforeEach(func() {
				logger.Info(action, logData, anotherLogData)
			})

			TestLagerFormat(lager.INFO)
			TestLagerData()
		})

		Context("with no log data", func() {
			BeforeEach(func() {
				logger.Info(action)
			})

			TestLagerFormat(lager.INFO)
		})
	})

	Describe("Error", func() {
		var err = errors.New("oh noes!")
		Context("with log data", func() {
			BeforeEach(func() {
				logger.Error(action, err, logData, anotherLogData)
			})

			TestLagerFormat(lager.ERROR)
			TestLagerData()

			It("data contains error message", func() {
				Ω(testSink.Logs()[0].Data["error"]).Should(Equal(err.Error()))
			})
		})

		Context("with no log data", func() {
			BeforeEach(func() {
				logger.Error(action, err)
			})

			TestLagerFormat(lager.ERROR)

			It("data contains error message", func() {
				Ω(testSink.Logs()[0].Data["error"]).Should(Equal(err.Error()))
			})
		})

		Context("with no error", func() {
			BeforeEach(func() {
				logger.Error(action, nil)
			})

			TestLagerFormat(lager.ERROR)

			It("does not contain the error message", func() {
				Ω(testSink.Logs()[0].Data).ShouldNot(HaveKey("error"))
			})
		})
	})

	Describe("Fatal", func() {
		var err = errors.New("oh noes!")
		var fatalErr interface{}

		Context("with log data", func() {
			BeforeEach(func() {
				defer func() {
					fatalErr = recover()
				}()

				logger.Fatal(action, err, logData, anotherLogData)
			})

			TestLagerFormat(lager.FATAL)
			TestLagerData()

			It("data contains error message", func() {
				Ω(testSink.Logs()[0].Data["error"]).Should(Equal(err.Error()))
			})

			It("data contains stack trace", func() {
				Ω(testSink.Logs()[0].Data["trace"]).ShouldNot(BeEmpty())
			})

			It("panics with the provided error", func() {
				Ω(fatalErr).Should(Equal(err))
			})
		})

		Context("with no log data", func() {
			BeforeEach(func() {
				defer func() {
					fatalErr = recover()
				}()

				logger.Fatal(action, err)
			})

			TestLagerFormat(lager.FATAL)

			It("data contains error message", func() {
				Ω(testSink.Logs()[0].Data["error"]).Should(Equal(err.Error()))
			})

			It("data contains stack trace", func() {
				Ω(testSink.Logs()[0].Data["trace"]).ShouldNot(BeEmpty())
			})

			It("panics with the provided error", func() {
				Ω(fatalErr).Should(Equal(err))
			})
		})

		Context("with no error", func() {
			BeforeEach(func() {
				defer func() {
					fatalErr = recover()
				}()

				logger.Fatal(action, nil)
			})

			TestLagerFormat(lager.FATAL)

			It("does not contain the error message", func() {
				Ω(testSink.Logs()[0].Data).ShouldNot(HaveKey("error"))
			})

			It("data contains stack trace", func() {
				Ω(testSink.Logs()[0].Data["trace"]).ShouldNot(BeEmpty())
			})

			It("panics with the provided error (i.e. nil)", func() {
				Ω(fatalErr).Should(BeNil())
			})
		})
	})
})

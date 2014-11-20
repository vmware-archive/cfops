package config_test

import (
	. "github.com/pivotalservices/cfops/config"
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type CustomConfig struct {
	TestConfig
	CustomProperty string
}

var _ = Describe("Config", func() {
	Context("when the config parses succesfully", func() {
		It("can be loaded from a json file", func() {
			file, err := ioutil.TempFile("", "config")
			defer func() {
				os.Remove(file.Name())
			}()
			Expect(err).NotTo(HaveOccurred())
			_, err = file.Write([]byte(`{"TestString":"MyString"}`))
			Expect(err).NotTo(HaveOccurred())

			err = file.Close()
			Expect(err).NotTo(HaveOccurred())

			config := &TestConfig{}
			err = LoadConfig(config, file.Name())
			Expect(err).NotTo(HaveOccurred())

			Expect(config.TestString).To(Equal("MyString"))
		})

		It("can be loaded into a custom config", func() {
			file, err := ioutil.TempFile("", "config")
			defer func() {
				os.Remove(file.Name())
			}()
			Expect(err).NotTo(HaveOccurred())
			_, err = file.Write([]byte(`{"TestString":"MyString", "CustomProperty":"MyValue"}`))
			Expect(err).NotTo(HaveOccurred())

			err = file.Close()
			Expect(err).NotTo(HaveOccurred())

			config := &CustomConfig{}
			err = LoadConfig(config, file.Name())
			Expect(err).NotTo(HaveOccurred())

			Expect(config.TestString).To(Equal("MyString"))
			Expect(config.CustomProperty).To(Equal("MyValue"))
		})
	})

	Context("when the config fails to parse", func() {
		It("should return error if file not found", func() {
			config := &TestConfig{}
			err := LoadConfig(config, "/foo/config.json")
			Expect(err).To(HaveOccurred())
		})

		It("should return error if json is invalid", func() {
			file, err := ioutil.TempFile("", "config")
			defer func() {
				os.Remove(file.Name())
			}()
			Expect(err).NotTo(HaveOccurred())
			_, err = file.Write([]byte(`NotJson`))
			Expect(err).NotTo(HaveOccurred())

			err = file.Close()
			Expect(err).NotTo(HaveOccurred())

			config := &TestConfig{}
			err = LoadConfig(config, file.Name())
			Expect(err).To(HaveOccurred())
		})
	})
})

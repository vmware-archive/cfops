package config_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotalservices/cfops/config"
)

var _ = Describe("Config initialized", func() {
	var config Config
	var err error

	BeforeEach(func() {
		config, err = NewConfig("./fixtures/")
	})
	Context("when GetOpsManagerHost is called", func() {
		It("then it should return 'theOpsManagerHost'", func() {
			Ω(err).ShouldNot(HaveOccurred())
			Ω(config.GetOpsManagerHost()).Should(Equal("theOpsManagerHost"))
		})
	})
	Context("when GetOpsManagerUser is called", func() {
		It("then it should return 'theOpsManagerUser'", func() {
			Ω(err).ShouldNot(HaveOccurred())
			Ω(config.GetOpsManagerUser()).Should(Equal("theOpsManagerUser"))
		})
	})
	Context("when GetOpsManagerPassword is called", func() {
		It("then it should return 'theOpsManagerPassword'", func() {
			Ω(err).ShouldNot(HaveOccurred())
			Ω(config.GetOpsManagerPassword()).Should(Equal("theOpsManagerPassword"))
		})
	})
	Context("when GetAdminUser is called", func() {
		It("then it should return 'theAdminUser'", func() {
			Ω(err).ShouldNot(HaveOccurred())
			Ω(config.GetAdminUser()).Should(Equal("theAdminUser"))
		})
	})
	Context("when GetAdminPassword is called", func() {
		It("then it should return 'theAdminPassword'", func() {
			Ω(err).ShouldNot(HaveOccurred())
			Ω(config.GetAdminPassword()).Should(Equal("theAdminPassword"))
		})
	})
	Context("when GetString is called", func() {
		It("then it should return an error when called with a key that doesn't exist", func() {
			_, theError := config.GetString("madeupkey")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(theError).Should(HaveOccurred())
		})
		It("then it should return a value when called with a key that does exist", func() {
			Ω(err).ShouldNot(HaveOccurred())
			theValue, theError := config.GetString("mytestkey")
			Ω(theError).ShouldNot(HaveOccurred())
			Ω(theValue).ShouldNot(BeEmpty())
		})
		It("then it should return an error when called with a key that doesn't exist", func() {
			_, theError := config.GetString("madeupkey")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(theError).Should(HaveOccurred())
		})
		It("then it should return a value when called with a key that is in mixed case", func() {
			Ω(err).ShouldNot(HaveOccurred())
			theValue, theError := config.GetString("MyTestKey")
			Ω(theError).ShouldNot(HaveOccurred())
			Ω(theValue).ShouldNot(BeEmpty())
		})
	})
	Context("when GetSubConfig is called", func() {
		It("then it should return a value when called with a nested parent", func() {
			var theValue string
			Ω(err).ShouldNot(HaveOccurred())
			subConfig, theError := config.GetSubConfig("mysql-tile")
			Ω(theError).ShouldNot(HaveOccurred())
			Ω(subConfig).ShouldNot(BeNil())
			theValue, theError = subConfig.GetString("name")
			Ω(theError).ShouldNot(HaveOccurred())
			Ω(theValue).ShouldNot(BeEmpty())
			subConfig2, theError2 := subConfig.GetSubConfig("sub-config")
			Ω(theError2).ShouldNot(HaveOccurred())
			Ω(subConfig2).ShouldNot(BeNil())
		})
	})
	Context("when GetInt is called", func() {
		It("then it should return an error when called with a key that doesn't exist", func() {
			_, theError := config.GetInt("madeupkey")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(theError).Should(HaveOccurred())
		})
		It("then it should return an int value when called with a key that does exist", func() {
			Ω(err).ShouldNot(HaveOccurred())
			theValue, theError := config.GetInt("mytestkeyint")
			Ω(theError).ShouldNot(HaveOccurred())
			Ω(theValue).Should(Equal(4))
		})
	})
	Context("when GetBool is called", func() {
		It("then it should return an error when called with a key that doesn't exist", func() {
			_, theError := config.GetBool("madeupkey")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(theError).Should(HaveOccurred())
		})
		It("then it should return an int value when called with a key that does exist", func() {
			Ω(err).ShouldNot(HaveOccurred())
			theValue, theError := config.GetBool("mytestkeybool")
			Ω(theError).ShouldNot(HaveOccurred())
			Ω(theValue).Should(Equal(true))
		})
	})

})

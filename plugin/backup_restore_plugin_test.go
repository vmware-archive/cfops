package plugin_test

import (
	"encoding/json"
	"fmt"
	"net/rpc"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotalservices/cfops/plugin"

	"github.com/pivotalservices/cfops/plugin/fake"
)

var _ = Describe("given a plugin", func() {
	var (
		pluginName = "BRTest"
		fakePlugin = &fake.BackupRestorePlugin{
			Meta: Meta{Name: pluginName, Role: "backup-restore"},
		}
	)
	Context("when called for meta data", func() {
		var spyUI []interface{}
		UIOutput = func(a ...interface{}) (int, error) {
			spyUI = a
			return fmt.Print(a...)
		}
		var origIsPluginMetaCall = IsPluginMetaCall.Load().(func() bool)
		var newIsPluginMetaCall = func() bool {
			return true
		}

		BeforeEach(func() {
			IsPluginMetaCall.Store(newIsPluginMetaCall)
		})

		AfterEach(func() {
			IsPluginMetaCall.Store(origIsPluginMetaCall)
		})

		It("should return us the plugins metadata", func() {
			Start(fakePlugin)
			control, _ := json.Marshal(fakePlugin.GetMeta())
			Ω(spyUI[0]).Should(Equal(string(control)))
		})
	})

	Context("when called to be executed", func() {

		var origIsPluginMetaCall = IsPluginMetaCall.Load().(func() bool)
		var newIsPluginMetaCall = func() bool {
			return false
		}

		BeforeEach(func() {
			IsPluginMetaCall.Store(newIsPluginMetaCall)
		})

		AfterEach(func() {
			IsPluginMetaCall.Store(origIsPluginMetaCall)
		})

		Context("when started", func() {
			var (
				err    error
				client *rpc.Client
			)
			go Start(fakePlugin)

			BeforeEach(func() {
				client, err = rpc.DialHTTP("tcp", fmt.Sprintf("127.0.0.1:%d", PluginPort))
			})

			AfterEach(func() {
				client.Close()
			})

			It("should be serving rpc", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})

			Context("when the rpc is running", func() {
				var err error
				var activity = "backup"
				var controlProducts = []Product{Product{}, Product{}, Product{}}
				var controlCredentials = []Credential{Credential{}, Credential{}, Credential{}}
				var fakePCF = &fake.PivotalCF{
					FakeProducts:    controlProducts,
					FakeCredentials: controlCredentials,
					FakeActivity:    activity,
				}

				BeforeEach(func() {
					pcfWrapper := NewPivotalCF(fakePCF)
					err = client.Call(pluginName+".Run", pcfWrapper, &os.Args)
				})

				It("should be Runnable", func() {
					Ω(err).ShouldNot(HaveOccurred())
					Ω(fakePlugin.RunCallCount).Should(BeNumerically(">", 0))
				})

				It("should give the plugin method a PivotalCF", func() {
					Ω(err).ShouldNot(HaveOccurred())
					Ω(fakePlugin.SpyPivotalCF.GetActivity()).Should(Equal(activity))
					Ω(fakePlugin.SpyPivotalCF.GetProducts()).Should(Equal(controlProducts))
					Ω(fakePlugin.SpyPivotalCF.GetCredentials()).Should(Equal(controlCredentials))
				})
			})
		})
	})
})

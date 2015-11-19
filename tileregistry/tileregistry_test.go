package tileregistry_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotalservices/cfops/tileregistry"
	"github.com/pivotalservices/cfops/tileregistry/fake"
)

var _ = Describe("tileregistry", func() {
	var (
		controlTileGeneratorKey = "myRegisteredTileGenerator"
		myTileGenerator         = new(fake.TileGenerator)
	)
	Describe("given: a Register() method", func() {
		Context("when: passed a name and a TileGenerator interface", func() {
			BeforeEach(func() {
				Register(controlTileGeneratorKey, myTileGenerator)
			})
			AfterEach(func() {
				Repo = make(map[string]TileGenerator)
			})
			It("then: it should add the given TileGenerator under the given name in the registry", func() {
				registry := GetRegistry()
				Ω(registry).ShouldNot(BeEmpty())
				Ω(registry[controlTileGeneratorKey]).Should(Equal(myTileGenerator))
			})
		})
	})

	Describe("given: a GetRegistry() method", func() {
		Context("when: called without any registrerd TileGenerators", func() {
			var registry map[string]TileGenerator
			BeforeEach(func() {
				registry = GetRegistry()
			})
			AfterEach(func() {
				Repo = make(map[string]TileGenerator)
			})
			It("then: it should return an empty map of TileGenerator interfaces", func() {
				Ω(registry).Should(BeEmpty())
				Ω(registry).ShouldNot(BeNil())
			})
		})
		Context("when: called containing registrerd TileGenerators", func() {
			var registry map[string]TileGenerator
			BeforeEach(func() {
				Register(controlTileGeneratorKey, myTileGenerator)
				registry = GetRegistry()
			})
			AfterEach(func() {
				Repo = make(map[string]TileGenerator)
			})
			It("then: it should return the map of registered TileGenerator interfaces", func() {
				Ω(registry).ShouldNot(BeEmpty())
				Ω(registry[controlTileGeneratorKey]).Should(Equal(myTileGenerator))
			})
		})
	})
})

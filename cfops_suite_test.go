package cfops_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotalservices/cfops"

	"testing"
)

func TestCfops(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cfops")
}

func testPipelineExecutionError(action string) {
	var errMock = errors.New("random execution mock error")
	Context("when an error is returned from the builtin pipeline", func() {
		var fs *mockFlagSet
		BeforeEach(func() {
			m := mockBuiltinPipeline{
				ErrReturned: errMock,
			}
			BuiltinPipelineExecution[action] = m.action

			fs = &mockFlagSet{
				tileListFlag: "",
			}
		})

		It("should return the error", func() {
			Ω(RunPipeline(fs, action)).ShouldNot(BeNil())
			Ω(RunPipeline(fs, action)).Should(Equal(errMock))
		})
	})

}

func testWithTilelist(action string) {
	Context("with tile list flag", func() {
		var fs *mockFlagSet
		BeforeEach(func() {
			m := mockBuiltinPipeline{}
			BuiltinPipelineExecution[action] = m.action

			fs = &mockFlagSet{
				tileListFlag: "opsmanager, er",
			}
		})

		It("should return a not implemented yet error", func() {
			Ω(RunPipeline(fs, action)).ShouldNot(BeNil())
		})
	})
}

func testWithoutTilelist(action string) {
	Context("without tile list flag", func() {
		var fs *mockFlagSet
		BeforeEach(func() {
			m := mockBuiltinPipeline{}
			BuiltinPipelineExecution[action] = m.action

			fs = &mockFlagSet{
				tileListFlag: "",
			}
		})

		It("should run the builtin pipeline", func() {
			Ω(RunPipeline(fs, action)).Should(BeNil())
		})
	})
}

type mockFlagSet struct {
	tileListFlag string
}

func (s *mockFlagSet) Host() (r string) {
	return
}

func (s *mockFlagSet) User() (r string) {
	return
}

func (s *mockFlagSet) Pass() (r string) {
	return
}

func (s *mockFlagSet) Tpass() (r string) {
	return
}

func (s *mockFlagSet) Dest() (r string) {
	return
}

func (s *mockFlagSet) Tilelist() (r string) {
	r = s.tileListFlag
	return
}

type mockBuiltinPipeline struct {
	ErrReturned error
}

func (s *mockBuiltinPipeline) action(string, string, string, string, string) error {
	return s.ErrReturned
}

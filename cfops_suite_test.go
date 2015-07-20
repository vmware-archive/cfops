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

func testWithValidTilelist(action string) {
	Context("with invalid tile list flag", func() {
		var (
			fs    *mockFlagSet
			mTile *mockTile
		)

		BeforeEach(func() {
			mTile = &mockTile{}
			SupportedTiles = map[string]func() (Tile, error){
				"TESTFLAG": func() (Tile, error) {
					return mTile, nil
				},
			}
			m := mockBuiltinPipeline{}
			BuiltinPipelineExecution[action] = m.action

			fs = &mockFlagSet{
				tileListFlag: "badflag",
			}
		})

		It("should return an error", func() {
			Ω(RunPipeline(fs, action)).ShouldNot(BeNil())
		})

		It("should not call the action on the tiles in your list", func() {
			RunPipeline(fs, action)
			Ω(mTile.RunCount).Should(Equal(0))
		})
	})
}

func testWithInvalidTileList(action string) {
	Context("with valid tile list flag", func() {
		var (
			fs    *mockFlagSet
			mTile *mockTile
		)

		BeforeEach(func() {
			mTile = &mockTile{}
			SupportedTiles = map[string]func() (Tile, error){
				"TESTFLAG": func() (Tile, error) {
					return mTile, nil
				},
			}
			m := mockBuiltinPipeline{}
			BuiltinPipelineExecution[action] = m.action

			fs = &mockFlagSet{
				tileListFlag: "testflag",
			}
		})

		It("should not return an error", func() {
			Ω(RunPipeline(fs, action)).Should(BeNil())
		})

		It("should call the action on the tiles in your list", func() {
			RunPipeline(fs, action)
			Ω(mTile.RunCount).Should(BeNumerically(">", 0))
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

func (s *mockFlagSet) AdminUser() (r string) {
	return
}

func (s *mockFlagSet) AdminPass() (r string) {
	return
}

func (s *mockFlagSet) OpsManagerUser() (r string) {
	return
}

func (s *mockFlagSet) OpsManagerPass() (r string) {
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

func (s *mockBuiltinPipeline) action(string, string, string, string, string, string) error {
	return s.ErrReturned
}

type mockTile struct {
	ErrReturned error
	RunCount    int
}

func (s *mockTile) Restore() (err error) {
	s.RunCount++
	return
}

func (s *mockTile) Backup() (err error) {
	s.RunCount++
	return
}

package e2e

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/otiai10/copy"
	"kcl-lang.io/kpm/pkg/utils"
)

var _ = ginkgo.Describe("Kpm CLI Testing", func() {
	ginkgo.Context("testing...", func() {
		testSuitesHome := filepath.Join(GetWorkDir(), TEST_SUITES_DIR)
		testSuites := LoadAllTestSuites(testSuitesHome)

		for _, ts := range testSuites {
			// In the for loop, the variable ts is defined outside the scope of the ginkgo.It function.
			// This means that when the ginkgo.It function is executed,
			// it will always use the value of ts from the last iteration of the for loop.
			// To fix this issue, create a new variable inside the loop with the same value as ts,
			// and use that variable inside the ginkgo.It function.
			ts := ts
			ginkgo.Describe(ts.GetTestSuiteInfo(), func() {
				testSpace := filepath.Join(testSuitesHome, ts.Name, "test_space")
				tmp := filepath.Join(testSuitesHome, ts.Name, "tmp")
				ginkgo.BeforeEach(func() {
					if !utils.DirExists(testSpace) {
						err := os.MkdirAll(testSpace, 0755)
						gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
					}
					// backup the test space
					err := copy.Copy(testSpace, tmp)
					gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
				})

				ginkgo.AfterEach(func() {
					// restore the test space
					err := os.RemoveAll(testSpace)
					gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
					err = copy.Copy(tmp, testSpace)
					gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
					err = os.RemoveAll(tmp)
					gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
				})

				ginkgo.It(ts.GetTestSuiteInfo(), func() {
					if ts.ignore || ts.Input == "" {
						fmt.Printf("skipped: %s\n", ts.Name)
						ginkgo.Skip(ts.Name)
					}
					stdout, stderr, err := ExecKpmWithWorkDir(ts.Input, ts.Conf.Cwd)
					gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
					if stdout == ts.ExpectStdout && stderr == ts.ExpectStderr {
						gomega.Expect(ts.ExpectStdout).Should(gomega.Equal(stdout))
						gomega.Expect(ts.ExpectStderr).Should(gomega.Equal(stderr))
					} else {
						fmt.Printf("\n%s : check expected output contains: %s", color.YellowString("[warning]"), ts.Name)
						gomega.Expect(stdout).Should(gomega.ContainSubstring(ts.ExpectStdout))
						gomega.Expect(stderr).Should(gomega.ContainSubstring(ts.ExpectStderr))
					}
				})
			})
		}
	})
})

package e2e

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestE2e(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "E2e Suite")
}

var _ = ginkgo.BeforeEach(func() {
	ginkgo.By("create kpm test workspace", func() {
		_ = CreateTestWorkspace()
	})
})

var _ = ginkgo.AfterEach(func() {
	ginkgo.By("clean up kpm test workspace", func() {
		_ = CleanUpTestWorkspace()
	})
})

var _ = ginkgo.AfterSuite(func() {
	ginkgo.By("clean up kpm bin", func() {
		path := filepath.Join(GetWorkDir(), "../..", "bin")
		cli := fmt.Sprintf("rm -rf %s", path)
		stdout, stderr, err := Exec(cli)
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		gomega.Expect(stdout).To(gomega.BeEmpty())
		gomega.Expect(stderr).To(gomega.BeEmpty())
	})
})

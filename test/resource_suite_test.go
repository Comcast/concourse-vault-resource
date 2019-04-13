package test

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/gexec"
)

var (
	checkPath  string
	inPath     string
	vaultAddr  string
	vaultToken string
)

func TestResource(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Resource Suite")
}

var _ = BeforeSuite(func() {
	var err error

	By("Getting VAULT_ADDR from environment variables")
	vaultAddr = os.Getenv("VAULT_ADDR")
	Expect(vaultAddr).NotTo(BeEmpty(), "$VAULT_ADDR must be provided")

	By("Getting VAULT_TOKEN from environment variables")
	vaultToken = os.Getenv("VAULT_TOKEN")
	Expect(vaultAddr).NotTo(BeEmpty(), "$VAULT_TOKEN must be provided")

	By("Compiling check binary")
	checkPath, err = gexec.Build("github.com/comcast/concourse-vault-resource/cmd/check", "-race")
	Expect(err).NotTo(HaveOccurred())

	By("Compiling in binary")
	inPath, err = gexec.Build("github.com/comcast/concourse-vault-resource/cmd/in", "-race")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

func run(command *exec.Cmd, stdinContents []byte) *gexec.Session {
	fmt.Fprintf(GinkgoWriter, "input: %s\n", stdinContents)

	stdin, err := command.StdinPipe()
	Expect(err).ShouldNot(HaveOccurred())

	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())

	_, err = io.WriteString(stdin, string(stdinContents))
	Expect(err).ShouldNot(HaveOccurred())

	err = stdin.Close()
	Expect(err).ShouldNot(HaveOccurred())

	return session
}

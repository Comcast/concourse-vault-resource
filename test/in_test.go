package test

import (
	"encoding/json"
	"io/ioutil"
	"os/exec"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/comcast/concourse-vault-resource/pkg/resource/models"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

const (
	inTimeout = 40 * time.Second
)

var _ = Describe("In", func() {
	var (
		command       *exec.Cmd
		inRequest     models.Request
		stdinContents []byte
		destDirectory string
	)

	BeforeEach(func() {
		var err error

		By("Creating temp directory")
		destDirectory, err = ioutil.TempDir("", "concourse-vault-resource")
		Expect(err).NotTo(HaveOccurred())

		By("Creating command object")
		command = exec.Command(inPath, destDirectory)

		By("Creating default request")
		inRequest = models.Request{
			Source: models.Source{
				VaultPaths: map[string]int{
					"kv2/data/atu/foo": 2,
				},
				VaultAddr:  vaultAddr,
				VaultToken: vaultToken,
			},
		}

		stdinContents, err = json.Marshal(inRequest)
		Expect(err).ShouldNot(HaveOccurred())
	})

	Describe("successful behavior", func() {
		It("writes secret(s) to destination file", func() {
			By("Running the command")
			session := run(command, stdinContents)

			Eventually(session, inTimeout).Should(gexec.Exit(0))

			files, err := ioutil.ReadDir(destDirectory)
			Expect(err).NotTo(HaveOccurred())

			Expect(len(files)).To(BeNumerically(">", 0))
			for _, file := range files {
				Expect(file.Name()).To(MatchRegexp("secrets"))
				Expect(file.Size()).To(BeNumerically(">", 0))
			}
		})

		It("returns valid json", func() {
			By("Running the command")
			session := run(command, stdinContents)
			Eventually(session, inTimeout).Should(gexec.Exit(0))

			By("Outputting a valid json rsponse")
			response := models.Response{}
			err := json.Unmarshal(session.Out.Contents(), &response)
			Expect(err).ShouldNot(HaveOccurred())

			By("Validating output contains versions")
			Expect(len(response.Version.Version)).To(BeNumerically(">", 0))
			Expect(response.Version.Version).NotTo(BeEmpty())
		})

	})

	Context("when validation fails", func() {
		BeforeEach(func() {
			inRequest.Source.VaultPaths = make(map[string]int, 0)

			var err error
			stdinContents, err = json.Marshal(inRequest)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("exits with error", func() {
			By("Running the command")
			session := run(command, stdinContents)

			By("Validating command exited with error")
			Eventually(session, inTimeout).Should(gexec.Exit(1))
			Expect(session.Err).Should(gbytes.Say("error validating resource configuration"))
		})
	})
})

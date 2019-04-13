package test

import (
	"encoding/json"
	"os"
	"os/exec"
	"time"

	"github.com/comcast/concourse-vault-resource/pkg/resource/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

const (
	checkTimeout = 40 * time.Second
)

var _ = Describe("Check()", func() {
	var (
		command       *exec.Cmd
		checkRequest  models.Request
		err           error
		stdinContents []byte
	)

	BeforeEach(func() {
		By("Creating command object")
		command = exec.Command(checkPath)

		By("Creating default request")
		checkRequest = models.Request{
			Source: models.Source{
				VaultPaths: map[string]int{
					"kv2/data/atu/foo": 2,
				},
				VaultAddr:  vaultAddr,
				VaultToken: vaultToken,
			},
		}

		stdinContents, err = json.Marshal(checkRequest)
		Expect(err).ShouldNot(HaveOccurred())
	})

	Describe("successful behavior", func() {
		It("returns secret(s) without error", func() {
			By("Running the Check() command")
			session := run(command, stdinContents)

			By("Validating command exited without error")
			Eventually(session, checkTimeout).Should(gexec.Exit(0))

			var resp []models.Version
			err := json.Unmarshal(session.Out.Contents(), &resp)
			Expect(err).NotTo(HaveOccurred())

			Expect(len(resp)).To(BeNumerically(">", 0))
			Expect(resp).NotTo(BeEmpty())
		})

		Context("vault address not provided", func() {
			BeforeEach(func() {
				err = os.Setenv("VAULT_ADDR", checkRequest.Source.VaultAddr)
				Expect(err).ShouldNot(HaveOccurred())

				checkRequest.Source.VaultAddr = ""

				stdinContents, err = json.Marshal(checkRequest)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Context("vault token not provided", func() {
			BeforeEach(func() {
				err = os.Setenv("VAULT_TOKEN", checkRequest.Source.VaultToken)
				Expect(err).ShouldNot(HaveOccurred())

				checkRequest.Source.VaultToken = ""

				stdinContents, err = json.Marshal(checkRequest)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		It("returns secret(s) without error", func() {
			By("Running the command")
			session := run(command, stdinContents)

			By("Validating command exited without error")
			Eventually(session, checkTimeout).Should(gexec.Exit(0))

			var resp []models.Version
			err := json.Unmarshal(session.Out.Contents(), &resp)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("when resource configuration validation fails", func() {
		BeforeEach(func() {
			checkRequest.Source.VaultPaths = make(map[string]int, 0)
			stdinContents, err = json.Marshal(checkRequest)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("exits with error", func() {
			By("Running the command")
			session := run(command, stdinContents)

			By("Validating command exited with error")
			Eventually(session, checkTimeout).Should(gexec.Exit(1))
			Expect(session.Err).Should(gbytes.Say("error validating resource configuration"))
		})
	})
})

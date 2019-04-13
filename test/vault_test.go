package test

import (
	"bytes"
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/comcast/concourse-vault-resource/pkg/resource/models"
	"github.com/comcast/concourse-vault-resource/test/fakes"
)

var _ = Describe("Vault", func() {
	var (
		check   []models.Version
		err     error
		request models.Request
		stdin   *bytes.Buffer
		v       *fakes.FakeVault
	)

	BeforeEach(func() {
		v = new(fakes.FakeVault)
		check = make([]models.Version, 0)
		stdin = &bytes.Buffer{}
		request = models.Request{
			Source: models.Source{
				VaultPaths: map[string]int{
					"kv2/data/foo/bar": 1,
				},
				VaultAddr:  "http://vault.example.com",
				VaultToken: "faKev4ultT0k3n",
			},
		}

		err = json.NewEncoder(stdin).Encode(request)
		Expect(err).ShouldNot(HaveOccurred())
	})

	Describe("when Check() is called", func() {
		Context("checks the resource for versions and a version is found", func() {
			It("should return []models.Version containing a new version", func() {
				check = append(check, models.Version{
					Version: "1",
				})
				v.CheckReturns(check)
				Expect(v.Check()).To(Equal(check))
			})
		})

		Context("checks the resource for versions and no version is found", func() {
			It("should return an empty []models.Version", func() {
				v.CheckReturns(check)
				Expect(v.Check()).To(Equal(check))
			})
		})
	})

	Describe("when In() is called", func() {
		Context("retrieves a secret(s) from vault", func() {
			It("should write the secret(s) to resource/secrets and no error should occur", func() {
				v.InReturns(err)
				Expect(v.In()).ShouldNot(HaveOccurred())
			})
		})
	})
})

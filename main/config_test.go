package main_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	fakesys "github.com/cloudfoundry/bosh-utils/system/fakes"

	. "github.com/frodenas/bosh-openstack-cpi/main"

	bocaction "github.com/frodenas/bosh-openstack-cpi/action"
	bocconfig "github.com/frodenas/bosh-openstack-cpi/openstack/config"

	"github.com/frodenas/bosh-registry/client"
)

var validOpenStackConfig = bocconfig.Config{
	IdentityEndpoint: "fake-identity-endpoint",
	Username:         "fake-username",
	Password:         "fake-password",
	TenantName:       "fake-tenant-name",
}

var validActionsOptions = bocaction.ConcreteFactoryOptions{
	Agent: registry.AgentOptions{
		Mbus: "fake-mbus",
		Ntp:  []string{},
		Blobstore: registry.BlobstoreOptions{
			Type: "fake-blobstore-type",
		},
	},
	Registry: registry.ClientOptions{
		Protocol: "http",
		Host:     "fake-host",
		Port:     5555,
		Username: "fake-username",
		Password: "fake-password",
	},
}

var validConfig = Config{
	OpenStack: validOpenStackConfig,
	Actions:   validActionsOptions,
}

var _ = Describe("NewConfigFromPath", func() {
	var (
		fs *fakesys.FakeFileSystem
	)

	BeforeEach(func() {
		fs = fakesys.NewFakeFileSystem()
	})

	It("returns error if config is not valid", func() {
		err := fs.WriteFileString("/config.json", "{}")
		Expect(err).ToNot(HaveOccurred())

		_, err = NewConfigFromPath("/config.json", fs)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("Validating config"))
	})

	It("returns error if file contains invalid json", func() {
		err := fs.WriteFileString("/config.json", "-")
		Expect(err).ToNot(HaveOccurred())

		_, err = NewConfigFromPath("/config.json", fs)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("Unmarshalling config"))
	})

	It("returns error if file cannot be read", func() {
		err := fs.WriteFileString("/config.json", "{}")
		Expect(err).ToNot(HaveOccurred())

		fs.ReadFileError = errors.New("fake-read-err")

		_, err = NewConfigFromPath("/config.json", fs)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("fake-read-err"))
	})
})

var _ = Describe("Config", func() {
	var (
		config Config
	)

	Describe("Validate", func() {
		BeforeEach(func() {
			config = validConfig
		})

		It("does not return error if all openstack and actions sections are valid", func() {
			err := config.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns error if openstack section is not valid", func() {
			config.OpenStack = bocconfig.Config{}

			err := config.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Validating OpenStack configuration"))
		})

		It("returns error if actions section is not valid", func() {
			config.Actions.Agent = registry.AgentOptions{}
			config.Actions.Registry = registry.ClientOptions{}

			err := config.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Validating Actions configuration"))
		})
	})
})

package config_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/frodenas/bosh-openstack-cpi/openstack/config"
)

var validConfig = Config{
	IdentityEndpoint: "fake-identity-endpoint",
	Username:         "fake-username",
	UserID:           "fake-user-id",
	Password:         "fake-password",
	APIKey:           "fake-api-key",
	TenantName:       "fake-tenant-name",
	TenantID:         "fake-tenant-id",
}

var _ = Describe("Config", func() {
	var (
		config Config
	)

	Describe("Validate", func() {
		BeforeEach(func() {
			config = validConfig
		})

		It("does not return error if all fields are valid", func() {
			err := config.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns error if Identity Endpoint is empty", func() {
			config.IdentityEndpoint = ""

			err := config.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Must provide a non-empty Identity Endpoint"))
		})

		It("returns error if Username and UserID are empty", func() {
			config.Username = ""
			config.UserID = ""

			err := config.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Must provide a non-empty Username or UserID"))
		})

		It("does returns error if Username is empty but UserID is not empty", func() {
			config.Username = ""

			err := config.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("does returns error if UserID is empty but Username is not empty", func() {
			config.UserID = ""

			err := config.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns error if Password and APIKey are empty", func() {
			config.Password = ""
			config.APIKey = ""

			err := config.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Must provide a non-empty Password or APIKey"))
		})

		It("does returns error if Password is empty but APIKey is not empty", func() {
			config.Password = ""

			err := config.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("does returns error if just APIKey is empty but Password is not empty", func() {
			config.APIKey = ""

			err := config.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns error if TenantName and TenantID are empty", func() {
			config.TenantName = ""
			config.TenantID = ""

			err := config.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Must provide a non-empty TenantName or TenantID"))
		})

		It("does returns error if TenantName is empty but TenantID is not empty", func() {
			config.TenantName = ""

			err := config.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("does returns error if just TenantID is empty but TenantName is not empty", func() {
			config.TenantID = ""

			err := config.Validate()
			Expect(err).ToNot(HaveOccurred())
		})
	})
})

package config

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type Config struct {
	IdentityEndpoint             string   `json:"identity_endpoint,omitempty"`
	Username                     string   `json:"username,omitempty"`
	UserID                       string   `json:"user_id,omitempty"`
	Password                     string   `json:"password,omitempty"`
	APIKey                       string   `json:"api_key,omitempty"`
	TenantName                   string   `json:"tenant_name,omitempty"`
	TenantID                     string   `json:"tenant_id,omitempty"`
	DomainName                   string   `json:"domain_name,omitempty"`
	DomainID                     string   `json:"domain_id,omitempty"`
	Region                       string   `json:"region,omitempty"`
	DefaultKeyPair               string   `json:"default_keypair,omitempty"`
	DefaultSecurityGroups        []string `json:"default_security_groups,omitempty"`
	DisableConfigDrive           bool     `json:"disable_config_drive,omitempty"`
	DisableNeutron               bool     `json:"disable_neutron,omitempty"`
	IgnoreServerAvailabilityZone bool     `json:"ignore_server_availability_zone,omitempty"`
}

func (c Config) Validate() error {
	if c.IdentityEndpoint == "" {
		return bosherr.Error("Must provide a non-empty Identity Endpoint")
	}

	if c.Username == "" && c.UserID == "" {
		return bosherr.Error("Must provide a non-empty Username or UserID")
	}

	if c.Password == "" && c.APIKey == "" {
		return bosherr.Error("Must provide a non-empty Password or APIKey")
	}

	if c.TenantName == "" && c.TenantID == "" {
		return bosherr.Error("Must provide a non-empty TenantName or TenantID")
	}

	return nil
}

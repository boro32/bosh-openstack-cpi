package client

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/frodenas/bosh-openstack-cpi/openstack/config"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
)

type OpenStackClient struct {
	config              config.Config
	computeService      *gophercloud.ServiceClient
	blockstorageService *gophercloud.ServiceClient
	networkService      *gophercloud.ServiceClient
}

func NewOpenStackClient(
	config config.Config,
) (OpenStackClient, error) {
	authOptions := gophercloud.AuthOptions{
		IdentityEndpoint: config.IdentityEndpoint,
		Username:         config.Username,
		UserID:           config.UserID,
		Password:         config.Password,
		APIKey:           config.APIKey,
		TenantName:       config.TenantName,
		TenantID:         config.TenantID,
		DomainID:         config.DomainID,
		DomainName:       config.DomainName,
		AllowReauth:      true,
	}

	providerClient, err := openstack.AuthenticatedClient(authOptions)
	if err != nil {
		return OpenStackClient{}, bosherr.WrapError(err, "Authenticating client")
	}

	endpointOpts := gophercloud.EndpointOpts{
		Region: config.Region,
	}

	computeService, err := openstack.NewComputeV2(providerClient, endpointOpts)
	if err != nil {
		return OpenStackClient{}, bosherr.WrapError(err, "Creating Compute Service")
	}

	blockstorageService, err := openstack.NewBlockStorageV1(providerClient, endpointOpts)
	if err != nil {
		return OpenStackClient{}, bosherr.WrapError(err, "Creating Block Storage Service")
	}

	var networkService *gophercloud.ServiceClient
	if !config.DisableNeutron {
		networkService, err = openstack.NewNetworkV2(providerClient, endpointOpts)
		if err != nil {
			return OpenStackClient{}, bosherr.WrapError(err, "Creating Network Service")
		}
	}

	return OpenStackClient{
		config:              config,
		computeService:      computeService,
		blockstorageService: blockstorageService,
		networkService:      networkService,
	}, nil
}

func (c OpenStackClient) ComputeService() *gophercloud.ServiceClient {
	return c.computeService
}

func (c OpenStackClient) BlockStorageService() *gophercloud.ServiceClient {
	return c.blockstorageService
}

func (c OpenStackClient) NetworkService() *gophercloud.ServiceClient {
	return c.networkService
}

func (c OpenStackClient) DefaultKeyPair() string {
	return c.config.DefaultKeyPair
}
func (c OpenStackClient) DefaultSecurityGroups() []string {
	return c.config.DefaultSecurityGroups
}

func (c OpenStackClient) DisableConfigDrive() bool {
	return c.config.DisableConfigDrive
}

func (c OpenStackClient) IgnoreServerAvailabilityZone() bool {
	return c.config.IgnoreServerAvailabilityZone
}

package client

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	"github.com/frodenas/bosh-openstack-cpi/openstack/config"
	"github.com/frodenas/bosh-openstack-cpi/openstack/extension_service"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
)

type OpenStackClient struct {
	config              config.Config
	computeService      *gophercloud.ServiceClient
	blockstorageService *gophercloud.ServiceClient
	networkService      *gophercloud.ServiceClient
	computeExtensions   extension.OpenStackComputeExtensions
	networkExtensions   extension.OpenStackNetworkExtensions
	logger              boshlog.Logger
}

func NewOpenStackClient(
	config config.Config,
	logger boshlog.Logger,
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

	computeExtensionService := extension.NewOpenStackComputeExtensionService(computeService, logger)
	computeExtensions, err := computeExtensionService.List()
	if err != nil {
		return OpenStackClient{}, bosherr.WrapError(err, "Creating Compute Service")
	}

	blockstorageService, err := openstack.NewBlockStorageV1(providerClient, endpointOpts)
	if err != nil {
		return OpenStackClient{}, bosherr.WrapError(err, "Creating Block Storage Service")
	}

	var networkService *gophercloud.ServiceClient
	var networkExtensions map[string]struct{}
	if !config.DisableNeutron {
		networkService, err = openstack.NewNetworkV2(providerClient, endpointOpts)
		if err != nil {
			return OpenStackClient{}, bosherr.WrapError(err, "Creating Network Service")
		}

		networkExtensionService := extension.NewOpenStackNetworkExtensionService(networkService, logger)
		networkExtensions, err = networkExtensionService.List()
		if err != nil {
			return OpenStackClient{}, bosherr.WrapError(err, "Creating Compute Service")
		}
	}

	return OpenStackClient{
		config:              config,
		computeService:      computeService,
		blockstorageService: blockstorageService,
		networkService:      networkService,
		computeExtensions:   computeExtensions,
		networkExtensions:   networkExtensions,
		logger:              logger,
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

func (c OpenStackClient) ComputeExtensions() extension.OpenStackComputeExtensions {
	return c.computeExtensions
}

func (c OpenStackClient) NetworkExtensions() extension.OpenStackNetworkExtensions {
	return c.networkExtensions
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

func (c OpenStackClient) DisableNeutron() bool {
	return c.config.DisableNeutron
}

func (c OpenStackClient) IgnoreServerAvailabilityZone() bool {
	return c.config.IgnoreServerAvailabilityZone
}

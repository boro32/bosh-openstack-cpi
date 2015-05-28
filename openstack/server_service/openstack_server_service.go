package server

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshuuid "github.com/cloudfoundry/bosh-utils/uuid"

	"github.com/frodenas/bosh-openstack-cpi/openstack/floating_ip_service"
	"github.com/frodenas/bosh-openstack-cpi/openstack/network_service"

	"github.com/rackspace/gophercloud"
)

const openstackServerServiceLogTag = "OpenStackServerService"
const openstackServerNamePrefix = "vm"
const openstackServerDescription = "Server managed by BOSH"

type OpenStackServerService struct {
	computeService    *gophercloud.ServiceClient
	floatingIPService floatingip.Service
	networkService    network.Service
	uuidGen           boshuuid.Generator
	logger            boshlog.Logger
}

func NewOpenStackServerService(
	computeService *gophercloud.ServiceClient,
	floatingIPService floatingip.Service,
	networkService network.Service,
	uuidGen boshuuid.Generator,
	logger boshlog.Logger,
) OpenStackServerService {
	return OpenStackServerService{
		computeService:    computeService,
		floatingIPService: floatingIPService,
		networkService:    networkService,
		uuidGen:           uuidGen,
		logger:            logger,
	}
}

type OpenStackUserData struct {
	Server   OpenStackUserDataServerName       `json:"server"`
	Registry OpenStackUserDataRegistryEndpoint `json:"registry"`
	DNS      OpenStackUserDataDNSItems         `json:"dns,omitempty"`
}

type OpenStackUserDataServerName struct {
	Name string `json:"name"`
}

type OpenStackUserDataRegistryEndpoint struct {
	Endpoint string `json:"endpoint"`
}

type OpenStackUserDataDNSItems struct {
	NameServer []string `json:"nameserver,omitempty"`
}

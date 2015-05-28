package fakes

import (
	"github.com/frodenas/bosh-openstack-cpi/openstack/client"
)

func NewFakeOpenStackClient() client.OpenStackClient { return client.OpenStackClient{} }

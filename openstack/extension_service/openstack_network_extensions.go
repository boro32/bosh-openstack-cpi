package extension

const openStackNetworkFloatingIpsExtension = "floatingips"
const openStackNetworkSecurityGroupsExtension = "security-group"

type OpenStackNetworkExtensions map[string]struct{}

func (ext OpenStackNetworkExtensions) HasFloatingIpsSupport() bool {
	_, found := ext[openStackNetworkFloatingIpsExtension]

	return found
}

func (ext OpenStackNetworkExtensions) HasSecurityGroupsSupport() bool {
	_, found := ext[openStackNetworkSecurityGroupsExtension]

	return found
}

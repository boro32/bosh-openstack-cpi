package extension

const openStackComputeBlockDeviceMappingExtension = "os-block-device-mapping-v2-boot"
const openStackComputeConfigDriveExtension = "os-config-drive"
const openStackComputeFloatingIpsExtension = "os-floating-ips"
const openStackComputeKeyPairsExtension = "os-keypairs"
const openStackComputeSchedulerHintsExtension = "OS-SCH-HNT"
const openStackComputeSecurityGroupsExtension = "os-security-groups"
const openStackComputeTenantNetworksExtension = "os-tenant-networks"

type OpenStackComputeExtensions map[string]struct{}

func (ext OpenStackComputeExtensions) HasBlockDeviceMappingSupport() bool {
	_, found := ext[openStackComputeBlockDeviceMappingExtension]

	return found
}

func (ext OpenStackComputeExtensions) HasConfigDriveSupport() bool {
	_, found := ext[openStackComputeConfigDriveExtension]

	return found
}

func (ext OpenStackComputeExtensions) HasFloatingIpsSupport() bool {
	_, found := ext[openStackComputeFloatingIpsExtension]

	return found
}

func (ext OpenStackComputeExtensions) HasKeyPairsSupport() bool {
	_, found := ext[openStackComputeKeyPairsExtension]

	return found
}

func (ext OpenStackComputeExtensions) HasSchedulerHintsSupport() bool {
	_, found := ext[openStackComputeSchedulerHintsExtension]

	return found
}

func (ext OpenStackComputeExtensions) HasSecurityGroupsSupport() bool {
	_, found := ext[openStackComputeSecurityGroupsExtension]

	return found
}

func (ext OpenStackComputeExtensions) HasTenantNetworksSupport() bool {
	_, found := ext[openStackComputeTenantNetworksExtension]

	return found
}

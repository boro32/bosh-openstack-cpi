package action_test

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	fakeuuid "github.com/cloudfoundry/bosh-utils/uuid/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/frodenas/bosh-openstack-cpi/action"

	clientfakes "github.com/frodenas/bosh-openstack-cpi/openstack/client/fakes"

	"github.com/frodenas/bosh-openstack-cpi/openstack/client"
	"github.com/frodenas/bosh-openstack-cpi/openstack/flavor_service"
	"github.com/frodenas/bosh-openstack-cpi/openstack/floating_ip_service"
	"github.com/frodenas/bosh-openstack-cpi/openstack/image_service"
	"github.com/frodenas/bosh-openstack-cpi/openstack/keypair_service"
	"github.com/frodenas/bosh-openstack-cpi/openstack/network_service"
	"github.com/frodenas/bosh-openstack-cpi/openstack/security_group_service"
	"github.com/frodenas/bosh-openstack-cpi/openstack/server_service"
	"github.com/frodenas/bosh-openstack-cpi/openstack/snapshot_service"
	"github.com/frodenas/bosh-openstack-cpi/openstack/volume_service"
	"github.com/frodenas/bosh-openstack-cpi/openstack/volume_type_service"

	"github.com/frodenas/bosh-registry/client"
)

var _ = Describe("ConcreteFactory", func() {
	var (
		uuidGen         *fakeuuid.FakeGenerator
		openstackClient client.OpenStackClient
		logger          boshlog.Logger

		options = ConcreteFactoryOptions{
			Registry: registry.ClientOptions{
				Protocol: "http",
				Host:     "fake-host",
				Port:     5555,
				Username: "fake-username",
				Password: "fake-password",
			},
		}

		factory Factory
	)

	var (
		flavorService        flavor.Service
		floatingIPService    floatingip.Service
		imageService         image.Service
		keypairService       keypair.Service
		networkService       network.Service
		registryClient       registry.Client
		securityGroupService securitygroup.Service
		serverService        server.Service
		snapshotService      snapshot.Service
		volumeService        volume.Service
		volumeTypeService    volumetype.Service
	)

	BeforeEach(func() {
		openstackClient = clientfakes.NewFakeOpenStackClient()
		uuidGen = &fakeuuid.FakeGenerator{}
		logger = boshlog.NewLogger(boshlog.LevelNone)

		factory = NewConcreteFactory(
			openstackClient,
			uuidGen,
			options,
			logger,
		)
	})

	BeforeEach(func() {
		flavorService = flavor.NewOpenStackFlavorService(
			openstackClient.ComputeService(),
			logger,
		)

		if openstackClient.DisableNeutron() {
			floatingIPService = floatingip.NewOpenStackComputeFloatingIPService(
				openstackClient.ComputeService(),
				logger,
			)
		} else {
			floatingIPService = floatingip.NewOpenStackNetworkFloatingIPService(
				openstackClient.NetworkService(),
				openstackClient.ComputeService(),
				logger,
			)
		}

		imageService = image.NewOpenStackImageService(
			openstackClient.ComputeService(),
			uuidGen,
			logger,
		)

		keypairService = keypair.NewOpenStackKeyPairService(
			openstackClient.ComputeService(),
			logger,
		)

		if openstackClient.DisableNeutron() {
			networkService = network.NewOpenStackComputeNetworkService(
				openstackClient.ComputeService(),
				logger,
			)
		} else {
			networkService = network.NewOpenStackNetworkNetworkService(
				openstackClient.NetworkService(),
				logger,
			)
		}

		registryClient = registry.NewHTTPClient(
			options.Registry,
			logger,
		)

		if openstackClient.DisableNeutron() {
			securityGroupService = securitygroup.NewOpenStackComputeSecurityGroupService(
				openstackClient.ComputeService(),
				logger,
			)
		} else {
			securityGroupService = securitygroup.NewOpenStackNetworkSecurityGroupService(
				openstackClient.NetworkService(),
				logger,
			)
		}

		serverService = server.NewOpenStackServerService(
			openstackClient.ComputeService(),
			floatingIPService,
			networkService,
			securityGroupService,
			openstackClient.DefaultSecurityGroups(),
			openstackClient.DisableConfigDrive(),
			openstackClient.DisableNeutron(),
			uuidGen,
			logger,
		)

		snapshotService = snapshot.NewOpenStackSnapshotService(
			openstackClient.BlockStorageService(),
			uuidGen,
			logger,
		)

		volumeService = volume.NewOpenStackVolumeService(
			openstackClient.BlockStorageService(),
			uuidGen,
			logger,
		)

		volumeTypeService = volumetype.NewOpenStackVolumeTypeService(
			openstackClient.BlockStorageService(),
			logger,
		)
	})

	It("returns error if action cannot be created", func() {
		action, err := factory.Create("fake-unknown-action")
		Expect(err).To(HaveOccurred())
		Expect(action).To(BeNil())
	})

	It("create_disk", func() {
		action, err := factory.Create("create_disk")
		Expect(err).ToNot(HaveOccurred())
		Expect(action).To(Equal(NewCreateDisk(
			volumeService,
			volumeTypeService,
			serverService,
			openstackClient.IgnoreServerAvailabilityZone(),
		)))
	})

	It("delete_disk", func() {
		action, err := factory.Create("delete_disk")
		Expect(err).ToNot(HaveOccurred())
		Expect(action).To(Equal(NewDeleteDisk(volumeService)))
	})

	It("snapshot_disk", func() {
		action, err := factory.Create("snapshot_disk")
		Expect(err).ToNot(HaveOccurred())
		Expect(action).To(Equal(NewSnapshotDisk(snapshotService)))
	})

	It("delete_snapshot", func() {
		action, err := factory.Create("delete_snapshot")
		Expect(err).ToNot(HaveOccurred())
		Expect(action).To(Equal(NewDeleteSnapshot(snapshotService)))
	})

	It("create_stemcell", func() {
		action, err := factory.Create("create_stemcell")
		Expect(err).ToNot(HaveOccurred())
		Expect(action).To(Equal(NewCreateStemcell(imageService)))
	})

	It("delete_stemcell", func() {
		action, err := factory.Create("delete_stemcell")
		Expect(err).ToNot(HaveOccurred())
		Expect(action).To(Equal(NewDeleteStemcell(imageService)))
	})

	It("create_vm", func() {
		action, err := factory.Create("create_vm")
		Expect(err).ToNot(HaveOccurred())
		Expect(action).To(Equal(NewCreateVM(
			serverService,
			flavorService,
			imageService,
			keypairService,
			volumeService,
			registryClient,
			options.Registry,
			options.Agent,
			openstackClient.DefaultKeyPair(),
		)))
	})

	It("configure_networks", func() {
		action, err := factory.Create("configure_networks")
		Expect(err).ToNot(HaveOccurred())
		Expect(action).To(Equal(NewConfigureNetworks(serverService, registryClient)))
	})

	It("delete_vm", func() {
		action, err := factory.Create("delete_vm")
		Expect(err).ToNot(HaveOccurred())
		Expect(action).To(Equal(NewDeleteVM(serverService, registryClient)))
	})

	It("reboot_vm", func() {
		action, err := factory.Create("reboot_vm")
		Expect(err).ToNot(HaveOccurred())
		Expect(action).To(Equal(NewRebootVM(serverService)))
	})

	It("set_vm_metadata", func() {
		action, err := factory.Create("set_vm_metadata")
		Expect(err).ToNot(HaveOccurred())
		Expect(action).To(Equal(NewSetVMMetadata(serverService)))
	})

	It("has_vm", func() {
		action, err := factory.Create("has_vm")
		Expect(err).ToNot(HaveOccurred())
		Expect(action).To(Equal(NewHasVM(serverService)))
	})

	It("attach_disk", func() {
		action, err := factory.Create("attach_disk")
		Expect(err).ToNot(HaveOccurred())
		Expect(action).To(Equal(NewAttachDisk(serverService, registryClient)))
	})

	It("detach_disk", func() {
		action, err := factory.Create("detach_disk")
		Expect(err).ToNot(HaveOccurred())
		Expect(action).To(Equal(NewDetachDisk(serverService, registryClient)))
	})

	It("get_disks", func() {
		action, err := factory.Create("get_disks")
		Expect(err).ToNot(HaveOccurred())
		Expect(action).To(Equal(NewGetDisks(serverService)))
	})

	It("ping", func() {
		action, err := factory.Create("ping")
		Expect(err).ToNot(HaveOccurred())
		Expect(action).To(Equal(NewPing()))
	})

	It("when action is current_vm_id returns an error because this CPI does not implement the method", func() {
		action, err := factory.Create("current_vm_id")
		Expect(err).To(HaveOccurred())
		Expect(action).To(BeNil())
	})

	It("when action is wrong returns an error because it is not an official CPI method", func() {
		action, err := factory.Create("wrong")
		Expect(err).To(HaveOccurred())
		Expect(action).To(BeNil())
	})
})

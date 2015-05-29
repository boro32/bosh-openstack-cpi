package action_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/frodenas/bosh-openstack-cpi/action"

	flavorfakes "github.com/frodenas/bosh-openstack-cpi/openstack/flavor_service/fakes"
	imagefakes "github.com/frodenas/bosh-openstack-cpi/openstack/image_service/fakes"
	keypairfakes "github.com/frodenas/bosh-openstack-cpi/openstack/keypair_service/fakes"
	serverfakes "github.com/frodenas/bosh-openstack-cpi/openstack/server_service/fakes"
	volumefakes "github.com/frodenas/bosh-openstack-cpi/openstack/volume_service/fakes"
	registryfakes "github.com/frodenas/bosh-registry/client/fakes"

	"github.com/frodenas/bosh-openstack-cpi/api"
	"github.com/frodenas/bosh-openstack-cpi/openstack/flavor_service"
	"github.com/frodenas/bosh-openstack-cpi/openstack/image_service"
	"github.com/frodenas/bosh-openstack-cpi/openstack/keypair_service"
	"github.com/frodenas/bosh-openstack-cpi/openstack/server_service"
	"github.com/frodenas/bosh-openstack-cpi/openstack/volume_service"

	"github.com/frodenas/bosh-registry/client"
)

var _ = Describe("CreateVM", func() {
	var (
		err                      error
		vmCID                    VMCID
		networks                 Networks
		cloudProps               VMCloudProperties
		disks                    []DiskCID
		env                      Environment
		registryOptions          registry.ClientOptions
		agentOptions             registry.AgentOptions
		defaultKeyPair           string
		disableConfigDrive       bool
		expectedServerProperties *server.Properties
		expectedServerNetworks   server.Networks
		expectedAgentSettings    registry.AgentSettings

		serverService  *serverfakes.FakeServerService
		flavorService  *flavorfakes.FakeFlavorService
		imageService   *imagefakes.FakeImageService
		keypairService *keypairfakes.FakeKeyPairService
		volumeService  *volumefakes.FakeVolumeService
		registryClient *registryfakes.FakeClient

		createVM CreateVM
	)

	BeforeEach(func() {
		serverService = &serverfakes.FakeServerService{}
		flavorService = &flavorfakes.FakeFlavorService{}
		imageService = &imagefakes.FakeImageService{}
		keypairService = &keypairfakes.FakeKeyPairService{}
		volumeService = &volumefakes.FakeVolumeService{}
		registryClient = &registryfakes.FakeClient{}
		registryOptions = registry.ClientOptions{
			Protocol: "http",
			Host:     "fake-registry-host",
			Port:     25777,
			Username: "fake-registry-username",
			Password: "fake-registry-password",
		}
		agentOptions = registry.AgentOptions{
			Mbus: "http://fake-mbus",
			Blobstore: registry.BlobstoreOptions{
				Type: "fake-blobstore-type",
			},
		}
		defaultKeyPair = "fake-default-keypair"
		disableConfigDrive = true
		createVM = NewCreateVM(
			serverService,
			flavorService,
			imageService,
			keypairService,
			volumeService,
			registryClient,
			registryOptions,
			agentOptions,
			defaultKeyPair,
			disableConfigDrive,
		)
	})

	Describe("Run", func() {
		BeforeEach(func() {
			serverService.CreateID = "fake-server-id"
			flavorService.FindByNameFound = true
			imageService.FindFound = true
			keypairService.FindFound = true

			volumeService.FindVolume = volume.Volume{AvailabilityZone: "fake-volume-availability-zone"}
			imageService.FindImage = image.Image{ID: "fake-image-id"}
			flavorService.FindByNameFlavor = flavor.Flavor{ID: "fake-flavor-id"}
			keypairService.FindKeyPair = keypair.KeyPair{Name: "fake-keypair"}

			cloudProps = VMCloudProperties{
				Flavor:           "fake-flavor",
				AvailabilityZone: "",
				SchedulerHints: VMSchedulerHintsProperties{
					Group:           "fake-scheduler-hints-group",
					DifferentHost:   []string{"fake-scheduler-hints-different-host"},
					SameHost:        []string{"fake-scheduler-hints-same-host"},
					Query:           []interface{}{"fake-scheduler-hints-query"},
					TargetCell:      "fake-scheduler-hints-target-cell",
					BuildNearHostIP: "fake-scheduler-hints-build-near-host-ip",
				},
			}

			networks = Networks{
				"fake-network-name": Network{
					Type:    "dynamic",
					IP:      "fake-network-ip",
					Gateway: "fake-network-gateway",
					Netmask: "fake-network-netmask",
					DNS:     []string{"fake-network-dns"},
					Default: []string{"fake-network-default"},
					CloudProperties: NetworkCloudProperties{
						Network:        "fake-network-cloud-network",
						SecurityGroups: []string{"fake-network-cloud-security-group"},
					},
				},
			}

			expectedServerProperties = &server.Properties{
				ImageID:            "fake-image-id",
				FlavorID:           "fake-flavor-id",
				AvailabilityZone:   "",
				KeyPair:            defaultKeyPair,
				DisableConfigDrive: disableConfigDrive,
				SchedulerHints: server.SchedulerHintsProperties{
					Group:           "fake-scheduler-hints-group",
					DifferentHost:   []string{"fake-scheduler-hints-different-host"},
					SameHost:        []string{"fake-scheduler-hints-same-host"},
					Query:           []interface{}{"fake-scheduler-hints-query"},
					TargetCell:      "fake-scheduler-hints-target-cell",
					BuildNearHostIP: "fake-scheduler-hints-build-near-host-ip",
				},
			}

			expectedServerNetworks = networks.AsServerServiceNetworks()

			expectedAgentSettings = registry.AgentSettings{
				AgentID: "fake-agent-id",
				Blobstore: registry.BlobstoreSettings{
					Provider: "fake-blobstore-type",
				},
				Disks: registry.DisksSettings{
					System:     "/dev/sda",
					Persistent: map[string]registry.PersistentSettings{},
				},
				Mbus: "http://fake-mbus",
				Networks: registry.NetworksSettings{
					"fake-network-name": registry.NetworkSettings{
						Type:    "dynamic",
						IP:      "fake-network-ip",
						Gateway: "fake-network-gateway",
						Netmask: "fake-network-netmask",
						DNS:     []string{"fake-network-dns"},
						Default: []string{"fake-network-default"},
					},
				},
				VM: registry.VMSettings{
					Name: "fake-server-id",
				},
			}
		})

		It("creates the vm", func() {
			vmCID, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
			Expect(err).NotTo(HaveOccurred())
			Expect(volumeService.FindCalled).To(BeFalse())
			Expect(imageService.FindCalled).To(BeTrue())
			Expect(flavorService.FindByNameCalled).To(BeTrue())
			Expect(keypairService.FindCalled).To(BeTrue())
			Expect(serverService.CreateCalled).To(BeTrue())
			Expect(serverService.CleanUpCalled).To(BeFalse())
			Expect(serverService.AddNetworkConfigurationCalled).To(BeTrue())
			Expect(registryClient.UpdateCalled).To(BeTrue())
			Expect(registryClient.UpdateSettings).To(Equal(expectedAgentSettings))
			Expect(vmCID).To(Equal(VMCID("fake-server-id")))
			Expect(serverService.CreateServerProps).To(Equal(expectedServerProperties))
			Expect(serverService.CreateNetworks).To(Equal(expectedServerNetworks))
			Expect(serverService.CreateRegistryEndpoint).To(Equal("http://fake-registry-host:25777"))
		})

		It("returns an error if imageService find call returns an error", func() {
			imageService.FindErr = errors.New("fake-image-service-error")

			_, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-image-service-error"))
			Expect(volumeService.FindCalled).To(BeFalse())
			Expect(imageService.FindCalled).To(BeTrue())
			Expect(flavorService.FindByNameCalled).To(BeFalse())
			Expect(keypairService.FindCalled).To(BeFalse())
			Expect(serverService.CreateCalled).To(BeFalse())
			Expect(serverService.CleanUpCalled).To(BeFalse())
			Expect(serverService.AddNetworkConfigurationCalled).To(BeFalse())
			Expect(registryClient.UpdateCalled).To(BeFalse())
		})

		It("returns an error if stemcell is not found", func() {
			imageService.FindFound = false

			_, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Stemcell 'fake-stemcell-id' does not exists"))
			Expect(volumeService.FindCalled).To(BeFalse())
			Expect(imageService.FindCalled).To(BeTrue())
			Expect(flavorService.FindByNameCalled).To(BeFalse())
			Expect(keypairService.FindCalled).To(BeFalse())
			Expect(serverService.CreateCalled).To(BeFalse())
			Expect(serverService.CleanUpCalled).To(BeFalse())
			Expect(serverService.AddNetworkConfigurationCalled).To(BeFalse())
			Expect(registryClient.UpdateCalled).To(BeFalse())
		})

		It("returns an error if flavor is not set", func() {
			cloudProps.Flavor = ""

			_, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("'flavor' must be provided"))
			Expect(volumeService.FindCalled).To(BeFalse())
			Expect(imageService.FindCalled).To(BeTrue())
			Expect(flavorService.FindByNameCalled).To(BeFalse())
			Expect(keypairService.FindCalled).To(BeFalse())
			Expect(serverService.CreateCalled).To(BeFalse())
			Expect(serverService.CleanUpCalled).To(BeFalse())
			Expect(serverService.AddNetworkConfigurationCalled).To(BeFalse())
			Expect(registryClient.UpdateCalled).To(BeFalse())
		})

		It("returns an error if flavorService find call returns an error", func() {
			flavorService.FindByNameErr = errors.New("fake-flavor-service-error")

			_, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-flavor-service-error"))
			Expect(volumeService.FindCalled).To(BeFalse())
			Expect(imageService.FindCalled).To(BeTrue())
			Expect(flavorService.FindByNameCalled).To(BeTrue())
			Expect(keypairService.FindCalled).To(BeFalse())
			Expect(serverService.CreateCalled).To(BeFalse())
			Expect(serverService.CleanUpCalled).To(BeFalse())
			Expect(serverService.AddNetworkConfigurationCalled).To(BeFalse())
			Expect(registryClient.UpdateCalled).To(BeFalse())
		})

		It("returns an error if flavor is not found", func() {
			flavorService.FindByNameFound = false

			_, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Flavor 'fake-flavor' does not exists"))
			Expect(volumeService.FindCalled).To(BeFalse())
			Expect(imageService.FindCalled).To(BeTrue())
			Expect(flavorService.FindByNameCalled).To(BeTrue())
			Expect(keypairService.FindCalled).To(BeFalse())
			Expect(serverService.CreateCalled).To(BeFalse())
			Expect(serverService.CleanUpCalled).To(BeFalse())
			Expect(serverService.AddNetworkConfigurationCalled).To(BeFalse())
			Expect(registryClient.UpdateCalled).To(BeFalse())
		})

		It("returns an error if keypairService find call returns an error", func() {
			keypairService.FindErr = errors.New("fake-keypair-service-error")

			_, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-keypair-service-error"))
			Expect(volumeService.FindCalled).To(BeFalse())
			Expect(imageService.FindCalled).To(BeTrue())
			Expect(flavorService.FindByNameCalled).To(BeTrue())
			Expect(keypairService.FindCalled).To(BeTrue())
			Expect(serverService.CreateCalled).To(BeFalse())
			Expect(serverService.CleanUpCalled).To(BeFalse())
			Expect(serverService.AddNetworkConfigurationCalled).To(BeFalse())
			Expect(registryClient.UpdateCalled).To(BeFalse())
		})

		It("returns an error if keypair is not found", func() {
			keypairService.FindFound = false

			_, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("KeyPair 'fake-default-keypair' does not exists"))
			Expect(volumeService.FindCalled).To(BeFalse())
			Expect(imageService.FindCalled).To(BeTrue())
			Expect(flavorService.FindByNameCalled).To(BeTrue())
			Expect(keypairService.FindCalled).To(BeTrue())
			Expect(serverService.CreateCalled).To(BeFalse())
			Expect(serverService.CleanUpCalled).To(BeFalse())
			Expect(serverService.AddNetworkConfigurationCalled).To(BeFalse())
			Expect(registryClient.UpdateCalled).To(BeFalse())
		})

		It("returns an error if networks are not valid", func() {
			_, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, Networks{}, disks, env)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Creating VM"))
			Expect(volumeService.FindCalled).To(BeFalse())
			Expect(imageService.FindCalled).To(BeTrue())
			Expect(flavorService.FindByNameCalled).To(BeTrue())
			Expect(keypairService.FindCalled).To(BeTrue())
			Expect(serverService.CreateCalled).To(BeFalse())
			Expect(serverService.CleanUpCalled).To(BeFalse())
			Expect(serverService.AddNetworkConfigurationCalled).To(BeFalse())
			Expect(registryClient.UpdateCalled).To(BeFalse())
		})

		It("returns an error if serverService create call returns an error", func() {
			serverService.CreateErr = errors.New("fake-server-service-error")

			_, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-server-service-error"))
			Expect(volumeService.FindCalled).To(BeFalse())
			Expect(imageService.FindCalled).To(BeTrue())
			Expect(flavorService.FindByNameCalled).To(BeTrue())
			Expect(keypairService.FindCalled).To(BeTrue())
			Expect(serverService.CreateCalled).To(BeTrue())
			Expect(serverService.CleanUpCalled).To(BeFalse())
			Expect(serverService.AddNetworkConfigurationCalled).To(BeFalse())
			Expect(registryClient.UpdateCalled).To(BeFalse())
		})

		It("returns an error if serverService add network configuration call returns an error", func() {
			serverService.AddNetworkConfigurationErr = errors.New("fake-server-service-error")

			_, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-server-service-error"))
			Expect(volumeService.FindCalled).To(BeFalse())
			Expect(imageService.FindCalled).To(BeTrue())
			Expect(flavorService.FindByNameCalled).To(BeTrue())
			Expect(keypairService.FindCalled).To(BeTrue())
			Expect(serverService.CreateCalled).To(BeTrue())
			Expect(serverService.CleanUpCalled).To(BeTrue())
			Expect(serverService.AddNetworkConfigurationCalled).To(BeTrue())
			Expect(registryClient.UpdateCalled).To(BeFalse())
		})

		It("returns an error if registryClient update call returns an error", func() {
			registryClient.UpdateErr = errors.New("fake-registry-client-error")

			_, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-registry-client-error"))
			Expect(volumeService.FindCalled).To(BeFalse())
			Expect(imageService.FindCalled).To(BeTrue())
			Expect(flavorService.FindByNameCalled).To(BeTrue())
			Expect(keypairService.FindCalled).To(BeTrue())
			Expect(serverService.CreateCalled).To(BeTrue())
			Expect(serverService.CleanUpCalled).To(BeTrue())
			Expect(serverService.AddNetworkConfigurationCalled).To(BeTrue())
			Expect(registryClient.UpdateCalled).To(BeTrue())
		})

		Context("when availability zone is set at cloud properties", func() {
			BeforeEach(func() {
				cloudProps.AvailabilityZone = "fake-availability-zone"
				expectedServerProperties.AvailabilityZone = "fake-availability-zone"
			})

			It("creates the vm at the right zone", func() {
				vmCID, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
				Expect(err).NotTo(HaveOccurred())
				Expect(volumeService.FindCalled).To(BeFalse())
				Expect(imageService.FindCalled).To(BeTrue())
				Expect(flavorService.FindByNameCalled).To(BeTrue())
				Expect(keypairService.FindCalled).To(BeTrue())
				Expect(serverService.CreateCalled).To(BeTrue())
				Expect(serverService.CleanUpCalled).To(BeFalse())
				Expect(serverService.AddNetworkConfigurationCalled).To(BeTrue())
				Expect(registryClient.UpdateCalled).To(BeTrue())
				Expect(registryClient.UpdateSettings).To(Equal(expectedAgentSettings))
				Expect(vmCID).To(Equal(VMCID("fake-server-id")))
				Expect(serverService.CreateServerProps).To(Equal(expectedServerProperties))
				Expect(serverService.CreateNetworks).To(Equal(expectedServerNetworks))
				Expect(serverService.CreateRegistryEndpoint).To(Equal("http://fake-registry-host:25777"))
			})
		})

		Context("when keyname is set at cloud properties", func() {
			BeforeEach(func() {
				cloudProps.KeyPair = "fake-keypair"
				expectedServerProperties.KeyPair = "fake-keypair"
			})

			It("creates the vm at the right zone", func() {
				vmCID, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
				Expect(err).NotTo(HaveOccurred())
				Expect(volumeService.FindCalled).To(BeFalse())
				Expect(imageService.FindCalled).To(BeTrue())
				Expect(flavorService.FindByNameCalled).To(BeTrue())
				Expect(keypairService.FindCalled).To(BeTrue())
				Expect(serverService.CreateCalled).To(BeTrue())
				Expect(serverService.CleanUpCalled).To(BeFalse())
				Expect(serverService.AddNetworkConfigurationCalled).To(BeTrue())
				Expect(registryClient.UpdateCalled).To(BeTrue())
				Expect(registryClient.UpdateSettings).To(Equal(expectedAgentSettings))
				Expect(vmCID).To(Equal(VMCID("fake-server-id")))
				Expect(serverService.CreateServerProps).To(Equal(expectedServerProperties))
				Expect(serverService.CreateNetworks).To(Equal(expectedServerNetworks))
				Expect(serverService.CreateRegistryEndpoint).To(Equal("http://fake-registry-host:25777"))
			})
		})

		Context("when DiskCIDs is set", func() {
			BeforeEach(func() {
				volumeService.FindFound = true
				disks = []DiskCID{"fake-volume-1"}
				expectedServerProperties.AvailabilityZone = "fake-volume-availability-zone"
			})

			It("creates the vm at the right zone", func() {
				vmCID, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
				Expect(err).NotTo(HaveOccurred())
				Expect(volumeService.FindCalled).To(BeTrue())
				Expect(imageService.FindCalled).To(BeTrue())
				Expect(flavorService.FindByNameCalled).To(BeTrue())
				Expect(keypairService.FindCalled).To(BeTrue())
				Expect(serverService.CreateCalled).To(BeTrue())
				Expect(serverService.CleanUpCalled).To(BeFalse())
				Expect(serverService.AddNetworkConfigurationCalled).To(BeTrue())
				Expect(registryClient.UpdateCalled).To(BeTrue())
				Expect(registryClient.UpdateSettings).To(Equal(expectedAgentSettings))
				Expect(vmCID).To(Equal(VMCID("fake-server-id")))
				Expect(serverService.CreateServerProps).To(Equal(expectedServerProperties))
				Expect(serverService.CreateNetworks).To(Equal(expectedServerNetworks))
				Expect(serverService.CreateRegistryEndpoint).To(Equal("http://fake-registry-host:25777"))
			})

			It("returns an error if volumeService find call returns an error", func() {
				volumeService.FindErr = errors.New("fake-volume-service-error")

				_, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-volume-service-error"))
				Expect(volumeService.FindCalled).To(BeTrue())
				Expect(imageService.FindCalled).To(BeFalse())
				Expect(flavorService.FindByNameCalled).To(BeFalse())
				Expect(keypairService.FindCalled).To(BeFalse())
				Expect(serverService.CreateCalled).To(BeFalse())
				Expect(serverService.CleanUpCalled).To(BeFalse())
				Expect(serverService.AddNetworkConfigurationCalled).To(BeFalse())
				Expect(registryClient.UpdateCalled).To(BeFalse())
			})

			It("returns an error if disk is not found", func() {
				volumeService.FindFound = false

				_, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(api.NewDiskNotFoundError("fake-volume-1", false).Error()))
				Expect(volumeService.FindCalled).To(BeTrue())
				Expect(imageService.FindCalled).To(BeFalse())
				Expect(flavorService.FindByNameCalled).To(BeFalse())
				Expect(keypairService.FindCalled).To(BeFalse())
				Expect(serverService.CreateCalled).To(BeFalse())
				Expect(serverService.CleanUpCalled).To(BeFalse())
				Expect(serverService.AddNetworkConfigurationCalled).To(BeFalse())
				Expect(registryClient.UpdateCalled).To(BeFalse())
			})

			Context("and zone is set", func() {
				BeforeEach(func() {
					cloudProps.AvailabilityZone = "fake-availability-zone"
				})

				It("returns an error if zone and disk are different", func() {
					_, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("can't use multiple availability zones:"))
					Expect(volumeService.FindCalled).To(BeTrue())
					Expect(imageService.FindCalled).To(BeFalse())
					Expect(flavorService.FindByNameCalled).To(BeFalse())
					Expect(keypairService.FindCalled).To(BeFalse())
					Expect(serverService.CreateCalled).To(BeFalse())
					Expect(serverService.CleanUpCalled).To(BeFalse())
					Expect(serverService.AddNetworkConfigurationCalled).To(BeFalse())
					Expect(registryClient.UpdateCalled).To(BeFalse())
				})
			})
		})
	})
})

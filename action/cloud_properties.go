package action

type DiskCloudProperties struct {
	VolumeType       string `json:"volume_type,omitempty"`
	AvailabilityZone string `json:"availability_zone,omitempty"`
}

type Environment map[string]interface{}

type NetworkCloudProperties struct {
	Network        string   `json:"network,omitempty"`
	SecurityGroups []string `json:"security_groups,omitempty"`
}

type SnapshotMetadata struct {
	Deployment string `json:"deployment,omitempty"`
	Job        string `json:"job,omitempty"`
	Index      string `json:"index,omitempty"`
}

type StemcellCloudProperties struct {
	Name           string `json:"name,omitempty"`
	Version        string `json:"version,omitempty"`
	Infrastructure string `json:"infrastructure,omitempty"`
	ImageUUID      string `json:"image_uuid,omitempty"`
}

type VMCloudProperties struct {
	Flavor           string                     `json:"flavor,omitempty"`
	AvailabilityZone string                     `json:"availability_zone,omitempty"`
	KeyPair          string                     `json:"keypair,omitempty"`
	RootDiskSizeGb   int                        `json:"root_disk_size_gb,omitempty"`
	SchedulerHints   VMSchedulerHintsProperties `json:"scheduler_hints,omitempty"`
}

type VMSchedulerHintsProperties struct {
	Group           string        `json:"group,omitempty"`
	DifferentHost   []string      `json:"different_host,omitempty"`
	SameHost        []string      `json:"same_host,omitempty"`
	Query           []interface{} `json:"query,omitempty"`
	TargetCell      string        `json:"target_cell,omitempty"`
	BuildNearHostIP string        `json:"build_near_host_ip,omitempty"`
}

type VMMetadata map[string]interface{}

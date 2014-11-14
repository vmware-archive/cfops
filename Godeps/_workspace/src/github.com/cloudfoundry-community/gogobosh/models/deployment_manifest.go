package models

// DeploymentManifest describes all the configuration for any BOSH deployment
type DeploymentManifest struct {
	Meta          map[string]interface{}
	Name          string
	DirectorUUID  string `yaml:"director_uuid"`
	Releases      []*NameVersion
	Compilation   *manifestCompilation
	Update        *manifestUpdate
	Networks      []*manifestNetwork
	ResourcePools []*manifestResourcePool `yaml:"resource_pools"`
	Jobs          []*ManifestJob
	Properties    *map[string]interface{}
}

type manifestCompilation struct {
	Workers             int                     `yaml:"workers"`
	NetworkName         string                  `yaml:"network"`
	ReuseCompilationVMs bool                    `yaml:"reuse_compilation_vms"`
	CloudProperties     *map[string]interface{} `yaml:"cloud_properties"`
}

type manifestUpdate struct {
	Canaries        int
	MaxInFlight     int    `yaml:"max_in_flight"`
	CanaryWatchTime string `yaml:"canary_watch_time"`
	UpdateWatchTime string `yaml:"update_watch_time"`
	Serial          bool
}

type manifestNetwork struct {
	Name            string
	Type            string
	CloudProperties *map[string]interface{} `yaml:"cloud_properties"`
	Subnets         interface{}
}

type manifestResourcePool struct {
	Name            string
	NetworkName     string `yaml:"network"`
	Stemcell        string
	CloudProperties string `yaml:"cloud_properties"`
}

// ManifestJob describes a cluster of VMs each running the same set of job templates
type ManifestJob struct {
	Name             string
	JobTemplates     []*ManifestJobTemplate `yaml:"templates"`
	Instances        int                    `yaml:"instances,omitempty"`
	ResourcePoolName string                 `yaml:"resource_pool"`
	PersistentDisk   int                    `yaml:"persistent_disk,omitempty"`
	Lifecycle        string                 `yaml:"lifecycle,omitempty"`
	Update           *manifestUpdate        `yaml:"update"`
	Networks         []*manifestJobNetwork
	Properties       *map[string]interface{} `yaml:"properties"`
}

// ManifestJobTemplate describes a job template included in a ManifestJob
type ManifestJobTemplate struct {
	Name    string
	Release string
}

type manifestJobNetwork struct {
	Name      string
	Default   *[]string `yaml:"default"`
	StaticIPs *[]string `yaml:"static_ips"`
}

// FindByJobTemplates returns the subnet of ManifestJobs that include a specific job template
func (manifest *DeploymentManifest) FindByJobTemplates(jobTemplateName string) (jobs []*ManifestJob) {
	jobs = []*ManifestJob{}
	for _, job := range manifest.Jobs {
		for _, jobTemplate := range job.JobTemplates {
			if jobTemplate.Name == jobTemplateName {
				jobs = append(jobs, job)
			}
		}
	}
	return
}

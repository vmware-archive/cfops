package models

// Deployments is a collection of deployments in the Director
type Deployments []*Deployment

// Deployment describes a running BOSH deployment and the
// Releases and Stemcells it is using.
type Deployment struct {
	Name      string
	Releases  []NameVersion
	Stemcells []NameVersion
}

// DeploymentVM describes the association of a running server
// within a Deployment
type DeploymentVM struct {
	JobName string
	Index   int
	VMCid   string
	AgentID string
}

// FindByRelease returns a list of deployments that use a release
func (deployments Deployments) FindByRelease(releaseName string) Deployments {
	subset := []*Deployment{}
	for _, deployment := range deployments {
		for _, release := range deployment.Releases {
			if release.Name == releaseName {
				subset = append(subset, deployment)
			}
		}
	}
	return Deployments(subset)
}

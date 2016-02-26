package cfbackup

import (
	"fmt"
	"strings"
)

//GetIPsByJob - get array of ips for a job
func (s *Products) GetIPsByJob(jobname string) (ips []string) {
	for vmName, ipList := range s.IPS {
		if strings.HasPrefix(vmName, jobname+"-") {
			ips = append(ips, ipList...)
		}
	}
	return
}

//GetPropertiesByJob - get array of []Properties for a job
func (s *Products) GetPropertiesByJob(jobname string) (properties []Properties, err error) {
	var job Jobs
	job, err = s.GetJob(jobname)
	if err != nil {
		return
	}
	properties = job.Properties
	return
}

//GetJob - get Job by name
func (s *Products) GetJob(jobName string) (job Jobs, err error) {
	var jobFound = false

	for _, theJob := range s.Jobs {
		if theJob.Identifier == jobName {
			job = theJob
			jobFound = true
			break
		}
	}
	if !jobFound {
		err = fmt.Errorf("job %s not found for product %s", jobName, s.Identifier)
	}
	return
}

//GetVMCredentialsByJob - returns VMCredentials for a job
func (s *Products) GetVMCredentialsByJob(jobName string) (vmCredentials VMCredentials, err error) {
	var job Jobs
	if job, err = s.GetJob(jobName); err == nil {

		if s.isLegacyFormat(job) {
			vmCredentials = s.extractLegacyCredentials(job)
		} else {
			vmCredentials = s.extractCredentials(job)
		}
	}
	return
}

func (s *Products) extractLegacyCredentials(job Jobs) (vmCredentials VMCredentials) {
	propMap := s.GetPropertyValues(job, "vm_credentials")
	vmCredentials.UserID = propMap["identity"]
	vmCredentials.Password = propMap["password"]
	return
}

func (s *Products) extractCredentials(job Jobs) (vmCredentials VMCredentials) {
	vmCredentials.UserID = job.VMCredentials["identity"]
	vmCredentials.Password = job.VMCredentials["password"]
	return
}

func (s *Products) isLegacyFormat(job Jobs) bool {
	return job.VMCredentials == nil
}

//GetPropertyValues = returns a map of property values for an identifier
func (s *Products) GetPropertyValues(job Jobs, identifier string) (propertyMap map[string]string) {
	properties := job.Properties
	propertyMap = make(map[string]string)
	for _, property := range properties {

		if property.Identifier == identifier {
			pMap := property.Value.(map[string]interface{})
			for key, value := range pMap {
				propertyMap[key] = fmt.Sprintf("%v", value)
			}
		}
	}
	return
}

//GetAvailabilityZoneNames - returns a list of availability zones
func (s *Products) GetAvailabilityZoneNames() (availablityZoneNames []string) {
	if len(s.AZReference) > 0 {
		availablityZoneNames = s.AZReference
	} else {
		availablityZoneNames = append(availablityZoneNames, s.SingletonAvailabilityZoneReference)
	}
	return
}

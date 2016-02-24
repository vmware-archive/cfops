package cfbackup

import "fmt"

//FindIPsByProductGUIDAndJobGUIDAndAvailabilityZoneGUID - returns array of IPs based on product and job guids
func (s *IPAssignments) FindIPsByProductGUIDAndJobGUIDAndAvailabilityZoneGUID(productGUID, jobGUID string, azGUID string) (ips []string, err error) {
	var assignmentsJob AssignmentsJob
	if assignmentsJob, err = s.Assignments.getAssignmentsJob(productGUID); err == nil {
		var assignmentsAZ AssignmentsAZ
		if assignmentsAZ, err = assignmentsJob.getAssignmentsAZ(jobGUID); err == nil {
			ips, err = assignmentsAZ.getIPs(azGUID)
		}
	}
	return
}

func (s AssignmentsProduct) getAssignmentsJob(productGUID string) (assignmentsJob AssignmentsJob, err error) {
	var ok bool
	if assignmentsJob, ok = s[productGUID]; !ok {
		err = fmt.Errorf("Product guid not found %s", productGUID)
	}
	return
}

func (s AssignmentsAZ) getIPs(azGUID string) (ips []string, err error) {
	var ok bool
	if ips, ok = s[azGUID]; !ok {
		err = fmt.Errorf("AZ guid not found %s", azGUID)
	}
	return
}

func (s AssignmentsJob) getAssignmentsAZ(jobGUID string) (assignmentsAZ AssignmentsAZ, err error) {
	var ok bool
	if assignmentsAZ, ok = s[jobGUID]; !ok {
		err = fmt.Errorf("Job guid not found %s", jobGUID)
	}
	return
}

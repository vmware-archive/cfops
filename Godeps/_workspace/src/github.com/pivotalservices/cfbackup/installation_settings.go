package cfbackup

import (
	"fmt"

	"github.com/pivotalservices/gtils/persistence"
)

var (
	boshVersionNames = map[string]string{
		"1.5": legacyBoshName,
		"1.4": legacyBoshName,
		"":    legacyBoshName,
	}
	pgDumpBin = map[string]string{
		"1.5": legacyPGDumpBin,
		"1.4": legacyPGDumpBin,
		"":    legacyPGDumpBin,
	}
	pgRestoreBin = map[string]string{
		"1.5": legacyPGRestoreBin,
		"1.4": legacyPGRestoreBin,
		"":    legacyPGRestoreBin,
	}
)

const (
	defaultBoshName     = "p-bosh"
	legacyBoshName      = "microbosh"
	defaultPGDumpBin    = "/var/vcap/packages/postgres-9.4.2/bin/pg_dump"
	defaultPGRestoreBin = "/var/vcap/packages/postgres-9.4.2/bin/pg_restore"
	legacyPGDumpBin     = "/var/vcap/packages/postgres/bin/pg_dump"
	legacyPGRestoreBin  = "/var/vcap/packages/postgres/bin/pg_restore"
)

//FindPropertyValues - returns a map of property values for a given product, job and identifier
func (s *InstallationSettings) FindPropertyValues(productName, jobName, identifier string) (propertyMap map[string]string, err error) {
	var product Products
	var job Jobs
	if product, err = s.FindByProductID(productName); err == nil {
		if job, err = s.FindJobByProductAndJobName(productName, jobName); err == nil {
			propertyMap = product.GetPropertyValues(job, identifier)
		}
	}
	return
}

//FindJobByProductAndJobName gets job for a given product and jobName
func (s *InstallationSettings) FindJobByProductAndJobName(productName, jobName string) (job Jobs, err error) {
	var product Products
	if product, err = s.FindByProductID(productName); err == nil {
		job, err = product.GetJob(jobName)
	}
	return
}

//FindVMCredentialsByProductAndJob gets VMCredentials for a given product and job
func (s *InstallationSettings) FindVMCredentialsByProductAndJob(productName, jobName string) (vmCredentials VMCredentials, err error) {
	var product Products
	if product, err = s.FindByProductID(productName); err == nil {
		vmCredentials, err = product.GetVMCredentialsByJob(jobName)
	}
	return
}

// FindIPsByProductAndJob finds a product and jobName
func (s *InstallationSettings) FindIPsByProductAndJob(productName string, jobName string) (IPs []string, err error) {

	if s.isLegacyFormat() {
		IPs, err = s.extractLegacyIPsForProductAndJob(productName, jobName)
	} else {
		IPs, err = s.extractIPsForProductAndJob(productName, jobName)
	}
	return
}

func (s *InstallationSettings) extractLegacyIPsForProductAndJob(productName, jobName string) (IPs []string, err error) {
	var product Products
	if product, err = s.FindByProductID(productName); err == nil {
		IPs = product.GetIPsByJob(jobName)
	}
	return
}

func (s *InstallationSettings) extractIPsForProductAndJob(productName, jobName string) (IPs []string, err error) {
	var product Products
	if product, err = s.FindByProductID(productName); err == nil {
		var job Jobs
		if job, err = product.GetJob(jobName); err == nil {
			IPs, err = s.findIPs(product, job)
		}
	}
	return
}

func (s *InstallationSettings) findIPs(product Products, job Jobs) (IPs []string, err error) {
	var IPsResponse []string
	for _, azGUID := range product.GetAvailabilityZoneNames() {
		if IPsResponse, err = s.IPAssignments.FindIPsByProductGUIDAndJobGUIDAndAvailabilityZoneGUID(product.GUID, job.GUID, azGUID); err == nil {
			for _, ip := range IPsResponse {
				IPs = append(IPs, ip)
			}
		}
	}
	return
}

// FindByProductID finds a product by product id
func (s *InstallationSettings) FindByProductID(id string) (productResponse Products, err error) {
	var found bool
	for _, product := range s.Products {
		identifier := product.Identifier
		if identifier == id {
			productResponse = product
			found = true
			break
		}
	}
	if !found {
		err = fmt.Errorf("Product not found %s", id)
	}

	return
}

// FindJobsByProductID finds all the jobs in an installation by product id
func (s *InstallationSettings) FindJobsByProductID(id string) []Jobs {
	cfJobs := []Jobs{}

	for _, product := range s.Products {
		identifier := product.Identifier
		if identifier == id {
			for _, job := range product.Jobs {
				cfJobs = append(cfJobs, job)
			}
		}
	}
	return cfJobs
}

// FindCFPostgresJobs finds all the postgres jobs in the cf product
func (s *InstallationSettings) FindCFPostgresJobs() (jobs []Jobs) {

	jobsList := s.FindJobsByProductID("cf")
	for _, job := range jobsList {
		if isPostgres(job.Identifier, job.Instances) {
			jobs = append(jobs, job)
		}
	}

	return jobs
}

func isPostgres(job string, instances []Instances) bool {
	pgdbs := []string{"ccdb", "uaadb", "consoledb"}

	for _, pgdb := range pgdbs {
		if pgdb == job {
			for _, instances := range instances {
				val := instances.Value
				if val >= 1 {
					return true
				}
			}
		}
	}
	return false
}

func (s *InstallationSettings) isLegacyFormat() bool {
	return s.IPAssignments.Assignments == nil
}

//GetBoshName - returns name of BoshProduct
func (s *InstallationSettings) GetBoshName() (boshName string) {
	var ok bool
	if boshName, ok = boshVersionNames[s.Version]; !ok {
		boshName = defaultBoshName
	}
	return
}

//SetPGDumpUtilVersions - initializes the correct dump commands
func (s *InstallationSettings) SetPGDumpUtilVersions() {
	var ok bool
	var pgDump, pgRestore string
	if pgDump, ok = pgDumpBin[s.Version]; !ok {
		pgDump = defaultPGDumpBin
	}
	persistence.PGDmpDumpBin = pgDump

	if pgRestore, ok = pgRestoreBin[s.Version]; !ok {
		pgRestore = defaultPGRestoreBin
	}
	persistence.PGDmpRestoreBin = pgRestore
}

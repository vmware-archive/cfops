package cfbackup

type connectionBucket struct {
	hostname           string
	adminUsername      string
	adminPassword      string
	opsManagerUsername string
	opsManagerPassword string
	destination        string
}

func (s connectionBucket) Host() string {
	return s.hostname
}

func (s connectionBucket) AdminUser() string {
	return s.adminUsername
}

func (s connectionBucket) AdminPass() string {
	return s.adminPassword
}

func (s connectionBucket) OpsManagerUser() string {
	return s.opsManagerUsername
}

func (s connectionBucket) OpsManagerPass() string {
	return s.opsManagerPassword
}

func (s connectionBucket) Destination() string {
	return s.destination
}

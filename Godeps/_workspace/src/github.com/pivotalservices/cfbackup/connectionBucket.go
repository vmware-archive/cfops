package cfbackup

type connectionBucket struct {
	hostname           string
	directorUsername   string
	directorPassword   string
	opsManagerUser     string
	opsManagerPassword string
	destination        string
}

func (s connectionBucket) Host() string {
	return s.hostname
}

func (s connectionBucket) DirectorUser() string {
	return s.directorUsername
}

func (s connectionBucket) DirectorPass() string {
	return s.directorPassword
}

func (s connectionBucket) OpsManagerUser() string {
	return s.opsManagerUser
}

func (s connectionBucket) OpsManagerPass() string {
	return s.opsManagerPassword
}

func (s connectionBucket) Destination() string {
	return s.destination
}

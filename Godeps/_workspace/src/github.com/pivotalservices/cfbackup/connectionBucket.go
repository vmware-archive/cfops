package cfbackup

type connectionBucket struct {
	hostname        string
	username        string
	password        string
	tempestPassword string
	destination     string
}

func (s connectionBucket) Host() string {
	return s.hostname
}

func (s connectionBucket) User() string {
	return s.username
}

func (s connectionBucket) Pass() string {
	return s.password
}

func (s connectionBucket) TempestPass() string {
	return s.tempestPassword
}

func (s connectionBucket) Destination() string {
	return s.destination
}

package cfbackup

import "os"

//BoshName - function which returns proper bosh component name for given version
func BoshName() (bosh string) {
	switch os.Getenv(ERVersionEnvFlag) {
	case ERVersion16:
		bosh = "p-bosh"
	default:
		bosh = "microbosh"
	}
	return
}

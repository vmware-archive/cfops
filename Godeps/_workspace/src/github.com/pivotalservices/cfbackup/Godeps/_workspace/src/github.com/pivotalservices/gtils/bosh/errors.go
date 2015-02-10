package bosh

import (
	"errors"
)

var (
	TaskStatusCodeError         error = errors.New("The resp code from return task should return 200")
	ManifestStatusCodeError     error = errors.New("The retriveing bosh manifest API response code is not equal to 200")
	TaskRedirectStatusCodeError error = errors.New("The resp code after task creation should return 302")
	TaskResultUnknown           error = errors.New("TASK processed result is unknown")
)

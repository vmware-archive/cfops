package gogobosh

import "github.com/cloudfoundry-community/gogobosh/models"

// NewDirector constructs a Director
func NewDirector(targetURL string, username string, password string) (director models.Director) {
	director = models.Director{}
	director.TargetURL = targetURL
	director.Username = username
	director.Password = password
	return
}

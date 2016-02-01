package plugin

import "github.com/pivotalservices/cfbackup"

func (s *defaultPivotalCF) SetActivity(activity string) {
	s.activity = activity
}

func (s *defaultPivotalCF) GetActivity() string {
	return s.activity
}

func (s *defaultPivotalCF) GetProducts() (products map[string]cfbackup.Products) {
	products = make(map[string]cfbackup.Products)

	for _, product := range s.installationSettings.GetProducts() {
		products[product.Identifier] = product
	}
	return
}

func (s *defaultPivotalCF) GetCredentials() (creds map[string]map[string][]cfbackup.Properties) {
	creds = make(map[string]map[string][]cfbackup.Properties)

	for _, product := range s.installationSettings.GetProducts() {
		creds[product.Identifier] = make(map[string][]cfbackup.Properties)

		for _, job := range product.Jobs {
			creds[product.Identifier][job.Identifier] = job.Properties
		}
	}
	return
}

//DefaultPivotalCF - creates the default pivotacf
var DefaultPivotalCF = func(installationSettings *cfbackup.ConfigurationParser) PivotalCF {

	return &defaultPivotalCF{
		installationSettings: installationSettings,
	}
}

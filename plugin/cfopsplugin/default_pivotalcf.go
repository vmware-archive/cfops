package cfopsplugin

import "github.com/pivotalservices/cfbackup"

//GetProducts - gets a products object from the given pivotalcf
func (s *DefaultPivotalCF) GetProducts() (products map[string]cfbackup.Products) {
	products = make(map[string]cfbackup.Products)

	for _, product := range s.InstallationSettings.GetProducts() {
		products[product.Identifier] = product
	}
	return
}

//GetCredentials - gets a credentials object from the given pivotalcf
func (s *DefaultPivotalCF) GetCredentials() (creds map[string]map[string][]cfbackup.Properties) {
	creds = make(map[string]map[string][]cfbackup.Properties)

	for _, product := range s.InstallationSettings.GetProducts() {
		creds[product.Identifier] = make(map[string][]cfbackup.Properties)

		for _, job := range product.Jobs {
			creds[product.Identifier][job.Identifier] = job.Properties
		}
	}
	return
}

//NewPivotalCF - creates the default pivotacf
var NewPivotalCF = func(installationSettings *cfbackup.ConfigurationParser) PivotalCF {

	return &DefaultPivotalCF{
		InstallationSettings: installationSettings,
	}
}

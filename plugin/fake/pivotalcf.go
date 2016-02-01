package fake

import "github.com/pivotalservices/cfbackup"

//PivotalCF --
type PivotalCF struct {
	FakeProducts    map[string]cfbackup.Products
	FakeCredentials map[string]map[string][]cfbackup.Properties
	FakeActivity    string
}

//SetActivity --
func (s *PivotalCF) SetActivity(activity string) {
	s.FakeActivity = activity
}

//GetActivity --
func (s *PivotalCF) GetActivity() string {
	return s.FakeActivity
}

//GetProducts --
func (s *PivotalCF) GetProducts() (products map[string]cfbackup.Products) {
	return s.FakeProducts
}

//GetCredentials --
func (s *PivotalCF) GetCredentials() (credentials map[string]map[string][]cfbackup.Properties) {
	return s.FakeCredentials
}

package fake

import (
	"github.com/pivotalservices/cfops/plugin"
)

//PivotalCF --
type PivotalCF struct {
	FakeProducts    []plugin.Product
	FakeCredentials []plugin.Credential
	FakeActivity    string
}

//GetActivity --
func (s *PivotalCF) GetActivity() string {
	return s.FakeActivity
}

//GetProducts --
func (s *PivotalCF) GetProducts() (products []plugin.Product) {
	return s.FakeProducts
}

//GetCredentials --
func (s *PivotalCF) GetCredentials() (credentials []plugin.Credential) {
	return s.FakeCredentials
}

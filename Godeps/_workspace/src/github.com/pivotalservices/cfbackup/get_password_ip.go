package cfbackup

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/xchapter7x/goutil/itertools"
)

type (
	InstallationCompareObject struct {
		Guid                 string
		Installation_Version string
		Products             []productCompareObject
		Infrastructure       infrastructure
	}

	infrastructure struct {
		Type string
	}

	productCompareObject struct {
		Type              string
		Identifier        string
		Installation_name string
		Jobs              []jobCompare
		IPs               ipCompare
	}

	ipCompare map[string][]string

	jobCompare struct {
		Type       string
		Identifier string
		Properties []propertyCompare
	}

	propertyCompare struct {
		Value interface{}
	}

	IpPasswordParser struct {
		Product   string
		Component string
		Username  string
		ip        string
		password  string
	}
)

func filterERProductsVersion13(v interface{}, product string) bool {
	return v.(productCompareObject).Type == product
}

func filterERProductsVersion14(v interface{}, product string) bool {
	return v.(productCompareObject).Identifier == product
}

func filterJobsVersion13(v interface{}, product string) bool {
	return v.(jobCompare).Type == product
}

func filterJobsVersion14(v interface{}, product string) bool {
	return v.(jobCompare).Identifier == product
}

func filterERProducts(i, v interface{}) bool {
	product := "cf"
	return filterERProductsVersion13(v, product) || filterERProductsVersion14(v, product)
}

func GetDeploymentName(jsonObj InstallationCompareObject) (deploymentName string, err error) {

	if o := itertools.Filter(jsonObj.Products, filterERProducts); len(o) > 0 {
		var (
			idx  interface{}
			prod productCompareObject
		)
		itertools.PairUnPack(<-o, &idx, &prod)
		deploymentName = prod.Installation_name

	} else {
		err = fmt.Errorf("could not find a cf install to pull name from")
	}
	return
}

func GetPasswordAndIP(jsonObj InstallationCompareObject, product, component, username string) (ip, password string, err error) {
	parser := &IpPasswordParser{
		Product:   product,
		Component: component,
		Username:  username,
	}
	return parser.Parse(jsonObj)
}

func (s *IpPasswordParser) Parse(jsonObj InstallationCompareObject) (ip, password string, err error) {
	if err = s.setupAndRun(jsonObj); err == nil {
		ip = s.ip
		password = s.password
	}
	return
}

func ReadAndUnmarshal(src io.Reader) (jsonObj InstallationCompareObject, err error) {
	var contents []byte

	if contents, err = ioutil.ReadAll(src); err == nil {
		err = json.Unmarshal(contents, &jsonObj)
	}
	return
}

func (s *IpPasswordParser) setupAndRun(jsonObj InstallationCompareObject) (err error) {
	var productObj productCompareObject
	s.modifyProductTypeName(jsonObj.Infrastructure.Type)

	if err = jsonFilter(jsonObj.Products, s.productFilter, &productObj); err == nil {
		err = s.ipPasswordSet(productObj)
	}
	return
}

func (s *IpPasswordParser) ipPasswordSet(productObj productCompareObject) (err error) {

	if err = s.setPassword(productObj); err == nil {
		err = s.setIp(productObj)
	}
	return
}

func (s *IpPasswordParser) setIp(productObj productCompareObject) (err error) {
	var iplist []string

	if err = jsonFilter(productObj.IPs, s.ipsFilter, &iplist); err == nil {

		s.ip = iplist[0]
	}
	return
}

func (s *IpPasswordParser) setPassword(productObj productCompareObject) (err error) {
	var jobObj jobCompare
	var property propertyCompare

	if err = jsonFilter(productObj.Jobs, s.jobsFilter, &jobObj); err == nil {

		if err = jsonFilter(jobObj.Properties, s.propertiesFilter, &property); err == nil {
			switch v := property.Value.(type) {
			case map[string]interface{}:
				s.password = property.Value.(map[string]interface{})["password"].(string)

			default:
				err = fmt.Errorf("unable to cast: map[string]interface{} :", v)
			}
		}
	}
	return
}

func (s *IpPasswordParser) productFilter(i, v interface{}) bool {
	return filterERProductsVersion13(v, s.Product) || filterERProductsVersion14(v, s.Product)
}

func (s *IpPasswordParser) jobsFilter(i, v interface{}) bool {
	return filterJobsVersion13(v, s.Component) || filterJobsVersion14(v, s.Component)
}

func (s *IpPasswordParser) propertiesFilter(i, v interface{}) (ok bool) {
	var identity interface{}

	switch v.(propertyCompare).Value.(type) {
	case map[string]interface{}:
		val := v.(propertyCompare).Value.(map[string]interface{})

		if identity, ok = val["identity"]; ok {
			ok = identity.(string) == s.Username
		}
	default:
		ok = false
	}
	return
}

func (s *IpPasswordParser) ipsFilter(i, v interface{}) bool {
	name := i.(string)
	val := v.([]string)
	return strings.Contains(name, fmt.Sprintf("%s-", s.Component)) && len(val) > 0
}

func (s *IpPasswordParser) modifyProductTypeName(typeval string) {
	typename := "vlcoud"
	productname := "microbosh"

	if typeval == typename && s.Product == productname {
		s.Product = fmt.Sprintf("%s-%s", productname, typename)
	}
}

func jsonFilter(list interface{}, filter func(i, v interface{}) bool, unpack interface{}) (err error) {
	var idx interface{}

	if o := itertools.Filter(list, filter); len(o) > 0 {
		itertools.PairUnPack(<-o, &idx, unpack)

	} else {
		err = fmt.Errorf("no matches in list for filter")
	}
	return
}

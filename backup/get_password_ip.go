package backup

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
		Type string
		Jobs []jobCompare
		IPs  ipCompare
	}

	ipCompare map[string][]string

	jobCompare struct {
		Type       string
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
		err = jsonFilter(jobObj.Properties, s.propertiesFilter, &property)
		switch v := property.Value.(type) {
		case map[string]interface{}:
			s.password = property.Value.(map[string]interface{})["password"].(string)

		default:
			err = fmt.Errorf("unable to cast: map[string]interface{}", v)
		}
	}
	return
}

func (s *IpPasswordParser) productFilter(i, v interface{}) bool {
	return v.(productCompareObject).Type == s.Product
}

func (s *IpPasswordParser) jobsFilter(i, v interface{}) bool {
	return v.(jobCompare).Type == s.Component
}

func (s *IpPasswordParser) propertiesFilter(i, v interface{}) bool {
	return v.(propertyCompare).Value.(map[string]interface{})["identity"].(string) == s.Username
}

func (s *IpPasswordParser) ipsFilter(i, v interface{}) bool {
	return strings.Contains(i.(string), s.Component)
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

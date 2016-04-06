package cfbackup

import "encoding/json"

//UnmarshalJSON Custom handling of PropertyValue unmarshal
//UnmarshalJSON Custom handling of PropertyValue unmarshal
func (p *PropertyValue) UnmarshalJSON(b []byte) (err error) {
	var arrayValue []interface{}
	var mapValue map[string]interface{}
	var stringValue string
	var intValue uint64
	var boolValue bool

	// if the value is string
	if err = json.Unmarshal(b, &stringValue); err == nil {
		p.StringValue = stringValue
		return
	}

	// if the value is int
	if err = json.Unmarshal(b, &intValue); err == nil {
		p.IntValue = intValue
		return
	}

	// if the value is an array
	if err = json.Unmarshal(b, &arrayValue); err == nil {
		p.ArrayValue = arrayValue
		return
	}

	// if the value is a map
	if err = json.Unmarshal(b, &mapValue); err == nil {
		p.MapValue = mapValue
		return
	}

	// if the value is a bool
	if err = json.Unmarshal(b, &boolValue); err == nil {
		p.BoolValue = boolValue
		return
	}

	return
}

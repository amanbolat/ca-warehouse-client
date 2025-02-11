// Code generated by "enumer -type=TransportMethod -json -sql -transform=snake"; DO NOT EDIT.

package logistics

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

const _TransportMethodName = "airautotrainexpresssealocal"

var _TransportMethodIndex = [...]uint8{0, 3, 7, 12, 19, 22, 27}

func (i TransportMethod) String() string {
	if i < 0 || i >= TransportMethod(len(_TransportMethodIndex)-1) {
		return fmt.Sprintf("TransportMethod(%d)", i)
	}
	return _TransportMethodName[_TransportMethodIndex[i]:_TransportMethodIndex[i+1]]
}

var _TransportMethodValues = []TransportMethod{0, 1, 2, 3, 4, 5}

var _TransportMethodNameToValueMap = map[string]TransportMethod{
	_TransportMethodName[0:3]:   0,
	_TransportMethodName[3:7]:   1,
	_TransportMethodName[7:12]:  2,
	_TransportMethodName[12:19]: 3,
	_TransportMethodName[19:22]: 4,
	_TransportMethodName[22:27]: 5,
}

// TransportMethodString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func TransportMethodString(s string) (TransportMethod, error) {
	if val, ok := _TransportMethodNameToValueMap[s]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to TransportMethod values", s)
}

// TransportMethodValues returns all values of the enum
func TransportMethodValues() []TransportMethod {
	return _TransportMethodValues
}

// IsATransportMethod returns "true" if the value is listed in the enum definition. "false" otherwise
func (i TransportMethod) IsATransportMethod() bool {
	for _, v := range _TransportMethodValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for TransportMethod
func (i TransportMethod) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for TransportMethod
func (i *TransportMethod) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("TransportMethod should be a string, got %s", data)
	}

	var err error
	*i, err = TransportMethodString(s)
	return err
}

func (i TransportMethod) Value() (driver.Value, error) {
	return i.String(), nil
}

func (i *TransportMethod) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	str, ok := value.(string)
	if !ok {
		bytes, ok := value.([]byte)
		if !ok {
			return fmt.Errorf("value is not a byte slice")
		}

		str = string(bytes[:])
	}

	val, err := TransportMethodString(str)
	if err != nil {
		return err
	}

	*i = val
	return nil
}

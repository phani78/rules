package rules

import (
	"strconv"
)

func ProcessNotEquals(value interface{}) (bool, interface{}) {
	switch value.(type) {
	case bool:
		boolVal := value.(bool)
		boolVal = !(boolVal)
		return boolVal, boolVal
	case string:
		strVal := value.(string)

		// Try to convert to bool
		if boolValue, err := strconv.ParseBool(strVal); err == nil {
			return boolValue, boolValue
		}
	}
	return false, value
}

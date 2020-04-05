package util

import "fmt"

func IsTrue(i interface{}) (bool, error) {
	switch t := i.(type) {
	case int:
		return i != 0, nil
	case string:
		return i != "", nil
	case bool:
		return i.(bool), nil
	default:
		return false, fmt.Errorf("unable to test type '%s'", t)
	}
}

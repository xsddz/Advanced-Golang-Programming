package helper

import "reflect"

func InArray(val interface{}, arr interface{}) bool {
	switch arrS := arr.(type) {
	case []int:
		return InIntSlice(val, arrS)
	case []int64:
		return InInt64Slice(val, arrS)
	case []string:
		return InStringSlice(val, arrS)
	case []interface{}:
		return InInterfaceSlice(val, arrS)
	default:
		return false
	}
}

func InInterfaceSlice(v interface{}, arr []interface{}) bool {
	for _, item := range arr {
		if reflect.DeepEqual(v, item) {
			return true
		}
	}
	return false
}

func InStringSlice(v interface{}, arr []string) bool {
	for _, item := range arr {
		if reflect.DeepEqual(v, item) {
			return true
		}
	}
	return false
}

func InInt64Slice(v interface{}, arr []int64) bool {
	for _, item := range arr {
		if reflect.DeepEqual(v, item) {
			return true
		}
	}
	return false
}

func InIntSlice(v interface{}, arr []int) bool {
	for _, item := range arr {
		if reflect.DeepEqual(v, item) {
			return true
		}
	}
	return false
}

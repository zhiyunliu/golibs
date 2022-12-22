package internal

import "reflect"

type XMapScanError struct {
	Message string
	Type    reflect.Type
}

func (e *XMapScanError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	if e.Type == nil {
		return "XMap: Scan(nil)"
	}

	// if e.Type.Kind() != reflect.Pointer {
	// 	return "XMap: Scan(non-pointer " + e.Type.String() + ")"
	// }
	return "XMap: Scan(nil " + e.Type.String() + ")"
}

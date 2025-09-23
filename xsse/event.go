package xsse

import (
	"encoding"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

type Event struct {
	Event string
	Id    string
	Retry uint
	Data  any
	Err   error
}

var (
	// default marshal callback
	DefaultMarshalCallback MarshalCallback = json.Marshal
	// default unmarshal callback
	DefaultUnmarshalCallback UnmarshalCallback = func(bytes []byte) (data any, err error) {
		return string(bytes), nil
	}
)

type MarshalCallback func(data any) ([]byte, error)
type UnmarshalCallback func(bytes []byte) (data any, err error)

func (e *Event) Encode(w stringWriter) error {
	e.writeId(w)
	e.writeEvent(w)
	e.writeRetry(w)
	return e.writeData(w)
}

func (e *Event) writeId(w stringWriter) {
	if len(e.Id) > 0 {
		_, _ = w.WriteString("id:")
		_, _ = fieldReplacer.WriteString(w, e.Id)
		_, _ = w.WriteString("\n")
	}
}

func (e *Event) writeEvent(w stringWriter) {
	if len(e.Event) > 0 {
		_, _ = w.WriteString("event:")
		_, _ = fieldReplacer.WriteString(w, e.Event)
		_, _ = w.WriteString("\n")
	}
}

func (e *Event) writeRetry(w stringWriter) {
	if e.Retry > 0 {
		_, _ = w.WriteString("retry:")
		_, _ = w.WriteString(strconv.FormatUint(uint64(e.Retry), 10))
		_, _ = w.WriteString("\n")
	}
}

func (e *Event) writeData(w stringWriter) error {
	_, _ = w.WriteString("data:")

	dataKind, custom := kindOfData(e.Data)
	switch dataKind {
	case reflect.Struct, reflect.Slice, reflect.Map:
		callback := DefaultMarshalCallback
		if custom {
			callback = customMarshal
		}
		databytes, err := callback(e.Data)
		if err != nil {
			return err
		}
		_, _ = w.Write(databytes)
		_, _ = w.WriteString("\n\n")
	default:
		_, _ = dataReplacer.WriteString(w, fmt.Sprint(e.Data))
		_, _ = w.WriteString("\n\n")
	}
	return nil
}

func customMarshal(data any) ([]byte, error) {
	return data.(encoding.BinaryMarshaler).MarshalBinary()
}

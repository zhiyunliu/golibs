package xsse

import (
	"encoding"
	"fmt"
	"reflect"
	"strconv"
)

type SSEEvent interface {
	Encode(w StringWriter) error
}

type Event struct {
	Event string
	Id    string
	Retry uint
	Data  any
	Err   error
}

type MarshalCallback func(data any) ([]byte, error)
type UnmarshalCallback func(bytes []byte) (data any, err error)

func (e *Event) Encode(w StringWriter) error {
	e.writeId(w)
	e.writeEvent(w)
	e.writeRetry(w)
	return e.writeData(w)
}

func (e *Event) writeId(w StringWriter) {
	if len(e.Id) > 0 {
		_, _ = w.Write([]byte("id:"))
		_, _ = fieldReplacer.WriteString(w, e.Id)
		_, _ = w.Write([]byte("\n"))
	}
}

func (e *Event) writeEvent(w StringWriter) {
	if len(e.Event) > 0 {
		_, _ = w.Write([]byte("event:"))
		_, _ = fieldReplacer.WriteString(w, e.Event)
		_, _ = w.Write([]byte("\n"))
	}
}

func (e *Event) writeRetry(w StringWriter) {
	if e.Retry > 0 {
		_, _ = w.Write([]byte("retry:"))
		_, _ = w.Write([]byte(strconv.FormatUint(uint64(e.Retry), 10)))
		_, _ = w.Write([]byte("\n"))
	}
}

func (e *Event) writeData(w StringWriter) error {
	_, _ = w.Write([]byte("data:"))

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
		_, _ = w.Write([]byte("\n\n"))
	default:
		_, _ = dataReplacer.WriteString(w, fmt.Sprint(e.Data))
		_, _ = w.Write([]byte("\n\n"))
	}
	return nil
}

func customMarshal(data any) ([]byte, error) {
	return data.(encoding.BinaryMarshaler).MarshalBinary()
}

type HeartbeatEvent struct {
	heartbeat string
}

func (e *HeartbeatEvent) Encode(w StringWriter) (err error) {
	if len(e.heartbeat) > 0 {
		_, err = w.Write([]byte(":" + e.heartbeat + "\n\n"))
		return
	}
	_, err = w.Write([]byte(":hrtb\n\n"))
	return
}

// NewHeartbeatEvent creates a new heartbeat event.
// If heartbeat is empty, the default heartbeat is used.
// If heartbeat is not empty, it is used as the heartbeat.
func defaultHeartbeatEvent(heartbeat ...string) SSEEvent {
	if len(heartbeat) > 0 {
		return &HeartbeatEvent{heartbeat: heartbeat[0]}
	}
	return &HeartbeatEvent{heartbeat: ""}
}

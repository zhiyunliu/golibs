package xsse

import (
	"encoding"
	"io"
	"net/http"
	"reflect"
	"strings"
)

// Server-Sent Events
// W3C Working Draft 29 October 2009
// http://www.w3.org/TR/2009/WD-eventsource-20091029/

const (
	ContentType = "text/event-stream"
	NoCache     = "no-cache"
)

var (
	//contentType = []string{ContentType}
	//noCache     = []string{NoCache}

	fieldReplacer = strings.NewReplacer(
		"\n", "\\n",
		"\r", "\\r")

	dataReplacer = strings.NewReplacer(
		"\n", "\ndata:",
		"\r", "\\r")
)

func HttpEncode(writer http.ResponseWriter, event Event) error {
	header := writer.Header()
	header["Content-Type"] = []string{ContentType}
	if _, exist := header["Cache-Control"]; !exist {
		header["Cache-Control"] = []string{NoCache}
	}

	w := wrapWriter(writer)
	return event.Encode(w)
}

func Encode(writer io.Writer, event Event) error {
	w := wrapWriter(writer)
	return event.Encode(w)
}

var (
	encodingBinaryMarshalerType = reflect.TypeOf((*encoding.BinaryMarshaler)(nil)).Elem()
)

func kindOfData(data interface{}) (dataKind reflect.Kind, custom bool) {
	value := reflect.ValueOf(data)
	dataKind = value.Kind()

	if value.Type().Implements(encodingBinaryMarshalerType) {
		custom = true
	}
	if dataKind == reflect.Ptr {
		dataKind = value.Elem().Kind()
	}
	return
}

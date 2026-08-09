package xsse

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

type sseEncodePayload struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type sseBinaryPayload struct {
	Value string
}

func (payload sseBinaryPayload) MarshalBinary() ([]byte, error) {
	return []byte("binary:" + payload.Value), nil
}

type sseWriteOnlyBuffer struct {
	data bytes.Buffer
}

func (writer *sseWriteOnlyBuffer) Write(bytes []byte) (int, error) {
	return writer.data.Write(bytes)
}

func (writer *sseWriteOnlyBuffer) String() string {
	return writer.data.String()
}

func TestEncodeWritesFieldsAndEscapesText(t *testing.T) {
	var buffer bytes.Buffer
	event := &Event{
		Id:    "id\n1\r2",
		Event: "update\nname",
		Retry: 15,
		Data:  "line1\nline2\rend",
	}

	err := Encode(&buffer, event)

	if err != nil {
		t.Fatalf("Encode() error = %v, want nil", err)
	}
	want := "id:id\\n1\\r2\nevent:update\\nname\nretry:15\ndata:line1\ndata:line2\\rend\n\n"
	if buffer.String() != want {
		t.Errorf("Encode() output = %q, want %q", buffer.String(), want)
	}
}

func TestEncodeMarshalsStructuredData(t *testing.T) {
	tests := []struct {
		name string
		data any
		want string
	}{
		{name: "struct", data: sseEncodePayload{Name: "alice", Count: 2}, want: "data:{\"name\":\"alice\",\"count\":2}\n\n"},
		{name: "slice", data: []string{"a", "b"}, want: "data:[\"a\",\"b\"]\n\n"},
		{name: "map", data: map[string]int{"count": 3}, want: "data:{\"count\":3}\n\n"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var buffer bytes.Buffer
			err := Encode(&buffer, &Event{Data: test.data})

			if err != nil {
				t.Fatalf("Encode() error = %v, want nil", err)
			}
			if buffer.String() != test.want {
				t.Errorf("Encode() output = %q, want %q", buffer.String(), test.want)
			}
		})
	}
}

func TestEncodeUsesBinaryMarshaler(t *testing.T) {
	var buffer bytes.Buffer
	event := &Event{Data: sseBinaryPayload{Value: "payload"}}

	err := Encode(&buffer, event)

	if err != nil {
		t.Fatalf("Encode() error = %v, want nil", err)
	}
	if buffer.String() != "data:binary:payload\n\n" {
		t.Errorf("Encode() output = %q, want %q", buffer.String(), "data:binary:payload\n\n")
	}
}

func TestEncodeReturnsMarshalError(t *testing.T) {
	wantErr := errors.New("marshal failed")
	originalCallback := DefaultMarshalCallback
	DefaultMarshalCallback = func(any) ([]byte, error) {
		return nil, wantErr
	}
	t.Cleanup(func() {
		DefaultMarshalCallback = originalCallback
	})

	err := Encode(&bytes.Buffer{}, &Event{Data: sseEncodePayload{Name: "alice"}})

	if !errors.Is(err, wantErr) {
		t.Fatalf("Encode() error = %v, want %v", err, wantErr)
	}
}

func TestEncodeWrapsPlainWriter(t *testing.T) {
	var writer sseWriteOnlyBuffer

	err := Encode(&writer, &Event{Data: "plain"})

	if err != nil {
		t.Fatalf("Encode() error = %v, want nil", err)
	}
	if writer.String() != "data:plain\n\n" {
		t.Errorf("Encode() output = %q, want %q", writer.String(), "data:plain\n\n")
	}
}

func TestHttpEncodeSetsSSEHeaders(t *testing.T) {
	recorder := httptest.NewRecorder()
	recorder.Header().Set("Cache-Control", "private")

	err := HttpEncode(recorder, &Event{Data: "hello"})

	if err != nil {
		t.Fatalf("HttpEncode() error = %v, want nil", err)
	}
	if got := recorder.Header().Get("Content-Type"); got != ContentType {
		t.Errorf("Content-Type = %q, want %q", got, ContentType)
	}
	if got := recorder.Header().Get("Cache-Control"); got != "private" {
		t.Errorf("Cache-Control = %q, want %q", got, "private")
	}
	if recorder.Body.String() != "data:hello\n\n" {
		t.Errorf("HttpEncode() body = %q, want %q", recorder.Body.String(), "data:hello\n\n")
	}
}

func TestHttpEncodeSetsDefaultCacheControl(t *testing.T) {
	recorder := httptest.NewRecorder()

	err := HttpEncode(recorder, &Event{Data: "hello"})

	if err != nil {
		t.Fatalf("HttpEncode() error = %v, want nil", err)
	}
	if got := recorder.Header().Get("Cache-Control"); got != NoCache {
		t.Errorf("Cache-Control = %q, want %q", got, NoCache)
	}
}

func TestHeartbeatEventEncode(t *testing.T) {
	tests := []struct {
		name  string
		event SSEEvent
		want  string
	}{
		{name: "default", event: NewHeartbeatEvent(), want: ":hrtb\n\n"},
		{name: "custom", event: NewHeartbeatEvent("ping"), want: ":ping\n\n"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var buffer bytes.Buffer
			err := test.event.Encode(&buffer)

			if err != nil {
				t.Fatalf("HeartbeatEvent.Encode() error = %v, want nil", err)
			}
			if buffer.String() != test.want {
				t.Errorf("HeartbeatEvent.Encode() output = %q, want %q", buffer.String(), test.want)
			}
		})
	}
}

func TestBuildEventReturnsHeartbeatForNilItem(t *testing.T) {
	event := buildEvent[any](nil, false, 1)

	var buffer bytes.Buffer
	err := event.Encode(&buffer)

	if err != nil {
		t.Fatalf("heartbeat Encode() error = %v, want nil", err)
	}
	if buffer.String() != ":hrtb\n\n" {
		t.Errorf("heartbeat Encode() output = %q, want %q", buffer.String(), ":hrtb\n\n")
	}
}

func TestBuildEventReturnsExistingSSEEvent(t *testing.T) {
	wantEvent := &Event{Id: "custom", Data: "value"}

	gotEvent := buildEvent[SSEEvent](wantEvent, false, 1)

	if gotEvent != wantEvent {
		t.Fatalf("buildEvent() = %#v, want same event %#v", gotEvent, wantEvent)
	}
}

func TestDecodeParsesEvents(t *testing.T) {
	input := strings.Join([]string{
		": ignored",
		"id: 42",
		"event: update",
		"data: hello",
		"data: world",
		"",
		"data: fallback",
		"",
	}, "\n")
	events, err := Decode(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Decode() error = %v, want nil", err)
	}

	gotEvents := collectEvents(events)
	wantEvents := []*Event{
		{Id: "42", Event: "update", Data: "hello\nworld"},
		{Event: "message", Data: "fallback"},
	}
	if !reflect.DeepEqual(gotEvents, wantEvents) {
		t.Errorf("Decode() events = %#v, want %#v", gotEvents, wantEvents)
	}
}

func TestDecodeUsesCustomUnmarshal(t *testing.T) {
	events, err := Decode(strings.NewReader("data: 21\n\n"), WithDecoderUnmarshal(func(bytes []byte) (any, error) {
		return strconv.Atoi(string(bytes))
	}))
	if err != nil {
		t.Fatalf("Decode() error = %v, want nil", err)
	}

	gotEvents := collectEvents(events)
	wantEvents := []*Event{{Event: "message", Data: 21}}
	if !reflect.DeepEqual(gotEvents, wantEvents) {
		t.Errorf("Decode() events = %#v, want %#v", gotEvents, wantEvents)
	}
}

func TestDecodeStoresUnmarshalErrorOnEvent(t *testing.T) {
	wantErr := errors.New("decode failed")
	events, err := Decode(strings.NewReader("data: value\n\n"), WithDecoderUnmarshal(func([]byte) (any, error) {
		return nil, wantErr
	}))
	if err != nil {
		t.Fatalf("Decode() error = %v, want nil", err)
	}

	gotEvents := collectEvents(events)
	if len(gotEvents) != 1 {
		t.Fatalf("Decode() event count = %d, want 1", len(gotEvents))
	}
	if !errors.Is(gotEvents[0].Err, wantErr) {
		t.Fatalf("Decode() event error = %v, want %v", gotEvents[0].Err, wantErr)
	}
}

func TestDecodeCanSkipUnmarshal(t *testing.T) {
	events, err := Decode(strings.NewReader("event: ping\n\n"), WithDecoderUnmarshal(nil))
	if err != nil {
		t.Fatalf("Decode() error = %v, want nil", err)
	}

	gotEvents := collectEvents(events)
	wantEvents := []*Event{{Event: "ping"}}
	if !reflect.DeepEqual(gotEvents, wantEvents) {
		t.Errorf("Decode() events = %#v, want %#v", gotEvents, wantEvents)
	}
}

func collectEvents(events <-chan *Event) []*Event {
	var result []*Event
	for event := range events {
		result = append(result, event)
	}
	return result
}

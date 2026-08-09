package xsse

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

type sseTestPayload struct {
	ID    string
	Value string
}

func (payload sseTestPayload) GetEventId() string {
	return payload.ID
}

type ssePlainPayload struct {
	Value string
}

func TestCheckEventIdGetter(t *testing.T) {
	tests := []struct {
		name string
		got  bool
		want bool
	}{
		{name: "value implements getter", got: checkEventIdGetter[sseTestPayload](), want: true},
		{name: "pointer implements getter", got: checkEventIdGetter[*sseTestPayload](), want: true},
		{name: "plain value does not implement getter", got: checkEventIdGetter[ssePlainPayload](), want: false},
		{name: "interface value may implement getter at runtime", got: checkEventIdGetter[any](), want: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.got != test.want {
				t.Errorf("checkEventIdGetter() = %v, want %v", test.got, test.want)
			}
		})
	}
}

func TestBuildEvent(t *testing.T) {
	payload := sseTestPayload{ID: "payload-1", Value: "hello"}

	tests := []struct {
		name        string
		item        any
		hasIDGetter bool
		idx         int
		wantID      string
	}{
		{name: "uses item event id when getter is enabled", item: payload, hasIDGetter: true, idx: 9, wantID: "payload-1"},
		{name: "uses index when getter is disabled", item: payload, hasIDGetter: false, idx: 9, wantID: "9"},
		{name: "keeps empty id when runtime item has no getter", item: ssePlainPayload{Value: "plain"}, hasIDGetter: true, idx: 3, wantID: ""},
		{name: "uses index for plain value", item: "plain", hasIDGetter: false, idx: 4, wantID: "4"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			event := buildEvent(test.item, test.hasIDGetter, test.idx)
			assertEvent(t, event, test.wantID, test.item)
		})
	}
}

func TestResultListGetEventV2(t *testing.T) {
	resultList := NewResultList(context.Background(), "first", "second")

	firstEvent, err := resultList.GetEventV2()
	if err != nil {
		t.Fatalf("ResultList.GetEventV2() first error = %v", err)
	}
	assertEvent(t, firstEvent, "1", "first")

	secondEvent, err := resultList.GetEventV2()
	if err != nil {
		t.Fatalf("ResultList.GetEventV2() second error = %v", err)
	}
	assertEvent(t, secondEvent, "2", "second")

	finalEvent, err := resultList.GetEventV2()
	if !errors.Is(err, ErrChanIsEmpty) {
		t.Fatalf("ResultList.GetEventV2() final error = %v, want %v", err, ErrChanIsEmpty)
	}
	if finalEvent != nil {
		t.Errorf("ResultList.GetEventV2() final event = %#v, want nil", finalEvent)
	}
}

func TestResultListGetEventUsesEventIdGetter(t *testing.T) {
	payload := sseTestPayload{ID: "custom-id", Value: "data"}
	resultList := NewResultList(context.Background(), payload)

	event, ok := resultList.GetEvent()
	if !ok {
		t.Fatal("ResultList.GetEvent() ok = false, want true")
	}
	assertEvent(t, event, "custom-id", payload)
}

func TestResultListGetEventWithInterfaceItems(t *testing.T) {
	payload := sseTestPayload{ID: "runtime-id", Value: "data"}
	resultList := NewResultList[any](context.Background(), payload, ssePlainPayload{Value: "plain"})

	firstEvent, ok := resultList.GetEvent()
	if !ok {
		t.Fatal("ResultList.GetEvent() first ok = false, want true")
	}
	assertEvent(t, firstEvent, "runtime-id", payload)

	secondEvent, ok := resultList.GetEvent()
	if !ok {
		t.Fatal("ResultList.GetEvent() second ok = false, want true")
	}
	assertEvent(t, secondEvent, "", ssePlainPayload{Value: "plain"})
}

func TestResultListGetEventReturnsFalseWhenEmpty(t *testing.T) {
	resultList := NewResultList[string](context.Background())

	event, ok := resultList.GetEvent()
	if ok {
		t.Fatal("ResultList.GetEvent() ok = true, want false")
	}
	if event != nil {
		t.Errorf("ResultList.GetEvent() event = %#v, want nil", event)
	}
}

func TestResultListGetEventV2ReturnsContextError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	resultList := NewResultList(ctx, "value")

	event, err := resultList.GetEventV2()
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("ResultList.GetEventV2() error = %v, want %v", err, context.Canceled)
	}
	if event != nil {
		t.Errorf("ResultList.GetEventV2() event = %#v, want nil", event)
	}
}

func TestResultChanGetEventV2(t *testing.T) {
	dataChan := make(chan int, 2)
	dataChan <- 10
	dataChan <- 20
	close(dataChan)
	resultChan := NewResultChan(context.Background(), dataChan)

	firstEvent, err := resultChan.GetEventV2()
	if err != nil {
		t.Fatalf("ResultChan.GetEventV2() first error = %v", err)
	}
	assertEvent(t, firstEvent, "1", 10)

	secondEvent, err := resultChan.GetEventV2()
	if err != nil {
		t.Fatalf("ResultChan.GetEventV2() second error = %v", err)
	}
	assertEvent(t, secondEvent, "2", 20)

	finalEvent, err := resultChan.GetEventV2()
	if !errors.Is(err, ErrChanIsEmpty) {
		t.Fatalf("ResultChan.GetEventV2() final error = %v, want %v", err, ErrChanIsEmpty)
	}
	if finalEvent != nil {
		t.Errorf("ResultChan.GetEventV2() final event = %#v, want nil", finalEvent)
	}
}

func TestResultChanGetEventUsesEventIdGetter(t *testing.T) {
	payload := sseTestPayload{ID: "chan-id", Value: "data"}
	dataChan := make(chan sseTestPayload, 1)
	dataChan <- payload
	close(dataChan)
	resultChan := NewResultChan(context.Background(), dataChan)

	event, ok := resultChan.GetEvent()
	if !ok {
		t.Fatal("ResultChan.GetEvent() ok = false, want true")
	}
	assertEvent(t, event, "chan-id", payload)
}

func TestResultChanGetEventReturnsFalseWhenClosed(t *testing.T) {
	dataChan := make(chan string)
	close(dataChan)
	resultChan := NewResultChan(context.Background(), dataChan)

	event, ok := resultChan.GetEvent()
	if ok {
		t.Fatal("ResultChan.GetEvent() ok = true, want false")
	}
	if event != nil {
		t.Errorf("ResultChan.GetEvent() event = %#v, want nil", event)
	}
}

func TestResultChanGetEventV2ReturnsContextError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	dataChan := make(chan string)
	resultChan := NewResultChan(ctx, dataChan)

	event, err := resultChan.GetEventV2()
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("ResultChan.GetEventV2() error = %v, want %v", err, context.Canceled)
	}
	if event != nil {
		t.Errorf("ResultChan.GetEventV2() event = %#v, want nil", event)
	}
}

func TestResultChanGetEventReturnsFalseOnContextError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	dataChan := make(chan string)
	resultChan := NewResultChan(ctx, dataChan)

	event, ok := resultChan.GetEvent()
	if ok {
		t.Fatal("ResultChan.GetEvent() ok = true, want false")
	}
	if event != nil {
		t.Errorf("ResultChan.GetEvent() event = %#v, want nil", event)
	}
}

func assertEvent(t *testing.T, event SSEEvent, wantID string, wantData any) {
	t.Helper()
	if event == nil {
		t.Fatal("event = nil, want non-nil")
	}
	evt, ok := event.(*Event)
	if !ok {
		t.Fatalf("event = %T, want *Event", event)
	}
	if evt.Id != wantID {
		t.Errorf("event.Id = %q, want %q", evt.Id, wantID)
	}
	if !reflect.DeepEqual(evt.Data, wantData) {
		t.Errorf("event.Data = %#v, want %#v", evt.Data, wantData)
	}
}

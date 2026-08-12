package xsse

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"
)

var (
	ErrChanIsEmpty = errors.New("chan is empty")
)

// EventIdGetter interface
type EventIdGetter interface {
	GetEventId() string
}
type EventTypeGetter interface {
	GetEventType() string
}

var (
	eventIdGetterType   = reflect.TypeOf((*EventIdGetter)(nil)).Elem()
	eventTypeGetterType = reflect.TypeOf((*EventTypeGetter)(nil)).Elem()
)

// needCheckEventIdGetter 判定类型 T 是否可能实现 EventIdGetter。
// T 为接口类型时，实际类型只能在运行时确定，需逐个元素断言。
func checkEventIdGetter[T any]() bool {
	typ := reflect.TypeFor[T]()
	if typ.Kind() == reflect.Interface {
		return true
	}
	return typ.Implements(eventIdGetterType)
}

// needCheckEventIdGetter 判定类型 T 是否可能实现 EventIdGetter。
// T 为接口类型时，实际类型只能在运行时确定，需逐个元素断言。
func checkEventTypeGetter[T any]() bool {
	typ := reflect.TypeFor[T]()
	if typ.Kind() == reflect.Interface {
		return true
	}
	return typ.Implements(eventTypeGetterType)
}

// ServerSentEvents interface
type ServerSentEvents interface {
	GetEvent() (evt SSEEvent, ok bool)
}

// ServerSentEvents interface
type ServerSentEventv2 interface {
	GetEventV2() (evt SSEEvent, err error)
}

// ResultChan is a channel that returns events
type ResultChan[T any] struct {
	ctx           context.Context
	timer         *time.Timer
	dataChan      chan T
	idx           int
	hasIdGetter   bool
	hasTypeGetter bool
}

var (
	_ ServerSentEvents  = (*ResultChan[any])(nil)
	_ ServerSentEventv2 = (*ResultChan[any])(nil)
)

// NewResultChan creates a new ResultChan
func NewResultChan[T any](ctx context.Context, dataChan chan T) *ResultChan[T] {
	return &ResultChan[T]{
		ctx:           ctx,
		timer:         time.NewTimer(HeartbeatInterval),
		dataChan:      dataChan,
		hasIdGetter:   checkEventIdGetter[T](),
		hasTypeGetter: checkEventTypeGetter[T](),
	}
}

// GetEvent returns the next event from the channel
func (s *ResultChan[T]) GetEvent() (evt SSEEvent, ok bool) {
	evt, err := s.GetEventV2()
	if errors.Is(err, ErrChanIsEmpty) {
		return evt, false
	}
	if err != nil {
		return nil, false
	}
	return evt, true
}

func (s *ResultChan[T]) GetEventV2() (evt SSEEvent, err error) {
	var (
		item T
		ok   bool
	)

	select {
	case <-s.ctx.Done():
		stopTimer(s.timer)
		return nil, s.ctx.Err()
	case item, ok = <-s.dataChan:
		if !ok {
			stopTimer(s.timer)
			return nil, ErrChanIsEmpty
		}
		resetTimer(s.timer, HeartbeatInterval)
	case <-s.timer.C:
		resetTimer(s.timer, HeartbeatInterval)
		return GetHeartBeatEvent(), nil
	}
	s.idx++
	return buildEvent(item, s.hasIdGetter, s.hasTypeGetter, s.idx), nil
}

func resetTimer(timer *time.Timer, duration time.Duration) {
	if !timer.Stop() {
		select {
		case <-timer.C:
		default:
		}
	}
	timer.Reset(duration)
}

func stopTimer(timer *time.Timer) {
	if !timer.Stop() {
		select {
		case <-timer.C:
		default:
		}
	}
}

// ResultList is a list that returns events
type ResultList[T any] struct {
	ctx           context.Context
	dataList      []T
	idx           int
	hasIdGetter   bool
	hasTypeGetter bool
}

var (
	_ ServerSentEvents  = (*ResultList[any])(nil)
	_ ServerSentEventv2 = (*ResultList[any])(nil)
)

// NewResultList creates a new ResultList
func NewResultList[T any](ctx context.Context, dataList ...T) *ResultList[T] {
	return &ResultList[T]{
		ctx:           ctx,
		dataList:      dataList,
		hasIdGetter:   checkEventIdGetter[T](),
		hasTypeGetter: checkEventTypeGetter[T](),
	}
}

// GetEvent returns the next event from the datalist
func (s *ResultList[T]) GetEvent() (evt SSEEvent, ok bool) {
	evt, err := s.GetEventV2()
	if errors.Is(err, ErrChanIsEmpty) {
		return evt, false
	}
	if err != nil {
		return nil, false
	}
	return evt, true
}

func (s *ResultList[T]) GetEventV2() (evt SSEEvent, err error) {
	select {
	case <-s.ctx.Done():
		return nil, s.ctx.Err()
	default:
	}

	if s.idx >= len(s.dataList) {
		return nil, ErrChanIsEmpty
	}
	item := s.dataList[s.idx]
	s.idx++

	return buildEvent(item, s.hasIdGetter, s.hasTypeGetter, s.idx), nil
}

func buildEvent[T any](item T, hasIdGetter bool, hasTypeGetter bool, idx int) (evt SSEEvent) {
	anyItem := any(item)

	if anyItem == nil {
		return GetHeartBeatEvent()
	}
	if evt, ok := anyItem.(SSEEvent); ok {
		return evt
	}

	var id string
	if hasIdGetter {
		if getter, ok := any(item).(EventIdGetter); ok {
			id = getter.GetEventId()
		}
	} else {
		id = fmt.Sprint(idx)
	}

	var eventType string
	if hasTypeGetter {
		if getter, ok := any(item).(EventTypeGetter); ok {
			eventType = getter.GetEventType()
		}
	}

	return &Event{
		Event: eventType,
		Id:    id,
		Data:  item,
	}
}

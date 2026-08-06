package xsse

import (
	"fmt"
	"reflect"
)

// EventIdGetter interface
type EventIdGetter interface {
	GetEventId() string
}

var eventIdGetterType = reflect.TypeOf((*EventIdGetter)(nil)).Elem()

// needCheckEventIdGetter 判定类型 T 是否可能实现 EventIdGetter。
// T 为接口类型时，实际类型只能在运行时确定，需逐个元素断言。
func checkEventIdGetter[T any]() bool {
	typ := reflect.TypeFor[T]()
	if typ.Kind() == reflect.Interface {
		return true
	}
	return typ.Implements(eventIdGetterType)
}

// ServerSentEvents interface
type ServerSentEvents interface {
	GetEvent() (evt *Event, ok bool)
}

// ResultChan is a channel that returns events
type ResultChan[T any] struct {
	dataChan    chan T
	idx         int
	hasIdGetter bool
}

// NewResultChan creates a new ResultChan
func NewResultChan[T any](dataChan chan T) *ResultChan[T] {

	return &ResultChan[T]{
		dataChan:    dataChan,
		hasIdGetter: checkEventIdGetter[T](),
	}
}

// GetEvent returns the next event from the channel
func (s *ResultChan[T]) GetEvent() (evt *Event, ok bool) {
	item, ok := <-s.dataChan
	if !ok {
		return nil, false
	}
	s.idx++

	var id string
	if s.hasIdGetter {
		if getter, ok := any(item).(EventIdGetter); ok {
			id = getter.GetEventId()
		}
	} else {
		id = fmt.Sprint(s.idx)
	}

	return &Event{
		Id:   id,
		Data: item,
	}, ok
}

// ResultList is a list that returns events
type ResultList[T any] struct {
	dataList    []T
	idx         int
	hasIdGetter bool
}

// NewResultList creates a new ResultList
func NewResultList[T any](dataList ...T) *ResultList[T] {
	return &ResultList[T]{
		dataList:    dataList,
		hasIdGetter: checkEventIdGetter[T](),
	}
}

// GetEvent returns the next event from the datalist
func (s *ResultList[T]) GetEvent() (evt *Event, ok bool) {
	if s.idx >= len(s.dataList) {
		return nil, false
	}
	item := s.dataList[s.idx]
	s.idx++

	var id string
	if s.hasIdGetter {
		if getter, ok := any(item).(EventIdGetter); ok {
			id = getter.GetEventId()
		}
	} else {
		id = fmt.Sprint(s.idx)
	}

	return &Event{
		Id:   id,
		Data: item,
	}, true
}

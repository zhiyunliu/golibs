package xsse

import "fmt"

// ServerSentEvents interface
type ServerSentEvents interface {
	GetEvent() (evt *Event, ok bool)
}

// ResultChan is a channel that returns events
type ResultChan[T any] struct {
	dataChan chan T
	idx      int
}

// NewResultChan creates a new ResultChan
func NewResultChan[T any](dataChan chan T) *ResultChan[T] {
	return &ResultChan[T]{
		dataChan: dataChan,
	}
}

// GetEvent returns the next event from the channel
func (s *ResultChan[T]) GetEvent() (evt *Event, ok bool) {
	item, ok := <-s.dataChan
	if !ok {
		return nil, false
	}
	s.idx++
	return &Event{
		Id:   fmt.Sprint(s.idx),
		Data: item,
	}, ok
}

// ResultList is a list that returns events
type ResultList[T any] struct {
	dataList []T
	idx      int
}

// NewResultChan creates a new ResultChan
func NewResultList[T any](dataList ...T) *ResultList[T] {
	return &ResultList[T]{
		dataList: dataList,
	}
}

// GetEvent returns the next event from the datalist
func (s *ResultList[T]) GetEvent() (evt *Event, ok bool) {
	if s.idx >= len(s.dataList) {
		return nil, false
	}
	item := s.dataList[s.idx]
	s.idx++
	return &Event{
		Id:   fmt.Sprint(s.idx),
		Data: item,
	}, true
}

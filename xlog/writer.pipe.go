package xlog

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type WriterPipe struct {
	ctx        context.Context
	ctxCancel  context.CancelFunc
	eventsChan chan *Event
	onceLock   sync.Once
	closed     bool
}

type WriterPipes []*WriterPipe

func newWriterPipe() *WriterPipe {
	if BufferSize <= 0 {
		panic(fmt.Errorf("WriterPipe xlog.BufferSize must more than 0"))
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &WriterPipe{
		ctx:        ctx,
		ctxCancel:  cancel,
		eventsChan: make(chan *Event, BufferSize),
		closed:     false,
	}
}

func (p *WriterPipe) Close() error {
	p.onceLock.Do(func() {
		p.closed = true
		close(p.eventsChan)
	})
	return nil
}

func (p *WriterPipe) complete() error {
	if p.ctxCancel != nil {
		p.ctxCancel()
	}
	return nil
}

func (ps WriterPipes) Write(evt *Event) error {
	idx := int(time.Now().UnixMicro() % int64(len(ps)))
	p := ps[0]
	if len(ps) > idx {
		p = ps[idx]
	}
	if p.closed {
		return fmt.Errorf("log writer pipe closed")
	}
	select {
	case p.eventsChan <- evt:
	default:
		//丢弃多余的日志
	}

	return nil
}

func (ps WriterPipes) Close() {
	for _, p := range ps {
		p.Close()
	}
}

func (ps WriterPipes) CloseAndWait() {
	group := &sync.WaitGroup{}
	group.Add(len(ps))
	for _, p := range ps {

		go func(w *WriterPipe) {
			<-w.ctx.Done()
			group.Done()
		}(p)
		p.Close()
	}

	group.Wait()
}

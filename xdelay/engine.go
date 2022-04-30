package xdelay

import (
	"sync"
	"time"

	"github.com/zhiyunliu/golibs/xlist"
)

type Engine struct {
	slotCount int
	//当前下标
	curIndex uint
	//环形槽
	slots []*xlist.List
	//关闭
	closed   chan struct{}
	onceLock sync.Once
}

//执行的任务函数
type TaskCallback func(args ...interface{})

//任务
type dealyTask struct {
	//循环次数
	cycleCnt int
	//执行的函数
	callback TaskCallback
	params   []interface{}
}

//创建一个延迟消息
func NewEngine(slotCount int) *Engine {
	dm := &Engine{
		slotCount: slotCount,
		curIndex:  0,
		closed:    make(chan struct{}),
		slots:     make([]*xlist.List, slotCount),
	}
	for i := 0; i < slotCount; i++ {
		dm.slots[i] = xlist.NewList()
	}
	return dm
}

//启动延迟消息
func (e *Engine) Start() {
	go e.taskLoop()
	<-e.closed
}

//关闭延迟消息
func (e *Engine) Close() {
	e.onceLock.Do(func() {
		close(e.closed)
	})
}

//处理每1秒的任务
func (e *Engine) taskLoop() {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-e.closed:
			return
		case <-ticker.C:
			list := e.slots[e.curIndex]
			if !list.IsEmpty() {
				list.Iter(func(idx int, node *xlist.Node) bool {
					task := node.Value.(*dealyTask)
					if task.cycleCnt <= 0 {
						go task.callback(task.params...)
						node.Remove()
						return true
					}
					task.cycleCnt--
					return true
				})

			}
			e.curIndex++
			e.curIndex = e.curIndex % uint(e.slotCount)
		}
	}
}

//添加任务
func (e *Engine) AddTask(seconds uint, callback TaskCallback, params ...interface{}) error {
	slotIdx := (e.curIndex + seconds) % uint(e.slotCount)
	//把任务加入tasks中
	e.slots[slotIdx].Append(&dealyTask{
		cycleCnt: int(seconds / uint(e.slotCount)),
		callback: callback,
		params:   params,
	})
	return nil
}

package xdelay

// 执行的任务函数
type TaskCallback func(args ...interface{})

type TaskInfo interface {
	GetCycleCount() int
	GetCallback() TaskCallback
	GetTag() TaskTag
	GetArgs() []any
}

type TaskTag interface {
	UniqueIdentifier() string
}

// 任务
type dealyTask struct {
	//循环次数
	cycleCnt int
	//执行的函数
	callback TaskCallback
	params   []any
	tag      TaskTag
}

func (t *dealyTask) GetCycleCount() int {
	return t.cycleCnt
}

func (t *dealyTask) GetTag() TaskTag {
	return t.tag
}

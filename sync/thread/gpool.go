package thread

import (
	"github.com/pkg/errors"
	"log"
	"sync"
	"sync/atomic"
)

//任务要包含需要执行的函数、以及函数要传的参数, 因为参数类型、个数不确定, 这里使用可变参数和空接口的形式
type Task struct {
	Handler func(v ...interface{})
	Params  []interface{}
}

//因为直接用了Mutex，注意加锁问题
type Pool struct {
	capacity       uint64     //池的容量
	runningWorkers uint64     //当前运行的 worker（goroutine）数量
	state          int64      //任务池的状态 state
	taskC          chan *Task ///任务队列（channel）
	sync.Mutex
	PanicHandler func(interface{}) //提供可订制的 panic handler
}

var ErrInvalidPoolCap = errors.New("invalid pool cap")
var ErrPoolAlreadyClosed = errors.New("pool already closed")

const (
	RUNNING = 1
	STOPED  = 0
)

//初始化容量
func NewPool(capacity uint64) (*Pool, error) {
	if capacity <= 0 {
		return nil, ErrInvalidPoolCap
	}
	return &Pool{
		capacity: capacity,
		state:    RUNNING,
		// 初始化任务队列, 队列长度为容量
		taskC: make(chan *Task, capacity),
	}, nil
}

//跑任务
func (p *Pool) run() {
	p.incRunning()

	go func() {
		defer func() {
			p.decRunning()
			if r := recover(); r != nil { // 恢复 panic
				if p.PanicHandler != nil { // 如果设置了 PanicHandler, 调用
					p.PanicHandler(r)
				} else { // 默认处理
					log.Printf("Worker panic: %s\n", r)
				}
			}
		}()

		for {
			select {
			case task, ok := <-p.taskC:
				if !ok {
					return
				}
				task.Handler(task.Params...)
			}
		}
	}()
}

func (p *Pool) Put(task *Task) error {

	if p.getState() == STOPED { // 如果任务池处于关闭状态, 再 put 任务会返回 ErrPoolAlreadyClosed 错误
		return ErrPoolAlreadyClosed
	}

	p.Lock()
	if p.GetRunningWorkers() < p.GetCap() {
		p.run()
	}
	p.Unlock()

	// 安全的推送任务, 以防在推送任务到 taskC 时 state 改变而关闭了 taskC
	p.Lock()
	if p.state == RUNNING {
		p.taskC <- task
	}
	p.Unlock()

	return nil
}

// 安全关闭 taskC
func (p *Pool) close() {
	p.Lock()
	defer p.Unlock()

	close(p.taskC)
}

func (p *Pool) Close() {

	if p.getState() == STOPED { // 如果已经关闭, 不能重复关闭
		return
	}

	p.setState(STOPED) // 设置 state 为已停止

	for len(p.taskC) > 0 { // 阻塞等待所有任务被 worker 消费
	}

	p.close()
}
func (p *Pool) incRunning() { // runningWorkers + 1
	atomic.AddUint64(&p.runningWorkers, 1)
}

func (p *Pool) decRunning() { // runningWorkers - 1
	atomic.AddUint64(&p.runningWorkers, ^uint64(0))
}

func (p *Pool) GetRunningWorkers() uint64 {
	return atomic.LoadUint64(&p.runningWorkers)
}
func (p *Pool) GetCap() uint64 {
	return p.capacity
}

func (p *Pool) getState() int64 {
	p.Lock()
	defer p.Unlock()

	return p.state
}

func (p *Pool) setState(state int64) {
	p.Lock()
	defer p.Unlock()

	p.state = state
}

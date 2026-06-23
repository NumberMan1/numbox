package timewheel

import (
	"container/list"
	"fmt"
	"time"
)

// Task 是时间轮中的一个定时任务。
type Task struct {
	Delay    time.Duration // 任务的延迟时间
	Key      any           // 任务的唯一标识符，用于取消任务
	Data     any           // 任务携带的自定义数据
	Callback func(any)     // 任务到期时执行的回调函数
	circle   int           // 任务需要在时间轮上转动的圈数
}

// TimeWheel 是一个时间轮调度器。
// 它将时间划分为多个槽（slot），每个槽代表一个时间间隔。
// 定时任务根据其到期时间被放入相应的槽中。
// 调度器通过一个指针按固定间隔扫过所有槽，来触发到期的任务。
// 此实现是线程安全的。所有操作通过内部 channel 在单一 goroutine 中完成。
type TimeWheel struct {
	interval   time.Duration    // 每个槽代表的时间间隔
	ticker     *time.Ticker     // Go原生的定时器，用于驱动时间轮指针移动
	slots      []*list.List     // 时间轮的槽位，每个槽是一个双向链表，存储任务
	timer      map[any]int      // 任务Key到其所在槽索引的映射，用于快速删除任务
	currentPos int              // 指针当前所在的槽位索引
	slotNum    int              // 槽位的总数
	taskChan   chan *Task       // 用于添加任务的channel
	removeChan chan any         // 用于移除任务的channel
	stopChan   chan struct{}    // 用于停止时间轮的channel
	OnExpired  func(task *Task) // 任务到期时的统一回调
}

// NewTimeWheel 创建一个新的时间轮。
// interval: 时间轮的滴答间隔，即指针多久移动一格。
// slotNum: 时间轮的槽位数量。
// onExpiredCallback: 一个统一的回调函数，当任何任务到期时都会被调用。
func NewTimeWheel(interval time.Duration, slotNum int, onExpiredCallback func(task *Task)) (*TimeWheel, error) {
	if interval <= 0 || slotNum <= 0 {
		return nil, fmt.Errorf("interval and slotNum must be positive")
	}
	tw := &TimeWheel{
		interval:   interval,
		slots:      make([]*list.List, slotNum),
		timer:      make(map[any]int),
		currentPos: 0,
		slotNum:    slotNum,
		taskChan:   make(chan *Task),
		removeChan: make(chan any),
		stopChan:   make(chan struct{}),
		OnExpired:  onExpiredCallback,
	}

	// 初始化每个槽位
	for i := 0; i < slotNum; i++ {
		tw.slots[i] = list.New()
	}
	return tw, nil
}

// Start 启动时间轮。
func (tw *TimeWheel) Start() {
	tw.ticker = time.NewTicker(tw.interval)
	go tw.run()
}

// Stop 停止时间轮。
func (tw *TimeWheel) Stop() {
	close(tw.stopChan)
}

// AddTask 添加一个定时任务。
// 任务的Key必须是可比较的类型，并且是唯一的。这是一个线程安全的操作。
func (tw *TimeWheel) AddTask(task *Task) {
	tw.taskChan <- task
}

// RemoveTask 根据任务的Key来移除一个还未执行的定时任务。这是一个线程安全的操作。
func (tw *TimeWheel) RemoveTask(key any) {
	tw.removeChan <- key
}

// run 是时间轮的主循环 Goroutine。
func (tw *TimeWheel) run() {
	defer tw.ticker.Stop()
	for {
		select {
		case <-tw.ticker.C:
			tw.tick()
		case task := <-tw.taskChan:
			tw.addTask(task)
		case key := <-tw.removeChan:
			tw.removeTask(key)
		case <-tw.stopChan:
			return
		}
	}
}

// tick 是每次指针移动时执行的核心逻辑。
func (tw *TimeWheel) tick() {
	// 获取当前槽位的任务列表
	l := tw.slots[tw.currentPos]
	// 遍历列表
	for e := l.Front(); e != nil; {
		task := e.Value.(*Task)
		if task.circle > 0 {
			// 如果任务还需要转几圈，则圈数减一，继续等待
			task.circle--
			e = e.Next()
			continue
		}

		// 任务到期，执行回调。使用 goroutine 以防止回调阻塞时间轮。
		if tw.OnExpired != nil {
			go tw.OnExpired(task)
		} else if task.Callback != nil {
			go task.Callback(task.Data)
		}

		// 从链表和映射中移除该任务
		next := e.Next()
		l.Remove(e)
		if task.Key != nil {
			delete(tw.timer, task.Key)
		}
		e = next
	}

	// 移动指针到下一个槽位
	tw.currentPos = (tw.currentPos + 1) % tw.slotNum
}

// addTask 将任务添加到正确的槽位。
func (tw *TimeWheel) addTask(task *Task) {
	if task == nil || task.Key == nil {
		return
	}
	// 计算任务应该在多少个滴答之后执行
	delayTicks := int(task.Delay / tw.interval)
	// 计算任务需要转动的圈数
	task.circle = delayTicks / tw.slotNum
	// 计算任务最终落脚的槽位索引
	pos := (tw.currentPos + delayTicks) % tw.slotNum

	// 将任务添加到槽位的链表中
	tw.slots[pos].PushBack(task)
	// 记录任务Key和槽位索引的映射关系
	tw.timer[task.Key] = pos
}

// removeTask 移除任务。
func (tw *TimeWheel) removeTask(key any) {
	pos, ok := tw.timer[key]
	if !ok {
		return
	}
	l := tw.slots[pos]
	for e := l.Front(); e != nil; e = e.Next() {
		task := e.Value.(*Task)
		if task.Key == key {
			delete(tw.timer, key)
			l.Remove(e)
			break
		}
	}
}

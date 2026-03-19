package scheduler

import (
	"container/heap"
	"log"
	"sync"
	"time"
)

type Schedule struct {
	subjectID   uint64
	executeTime time.Time
}

type ScheduleHeap []*Schedule

func (h *ScheduleHeap) Len() int { return len(*h) }

func (h *ScheduleHeap) Less(i, j int) bool {
	return (*h)[i].executeTime.Before((*h)[j].executeTime)
}

func (h *ScheduleHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *ScheduleHeap) Push(x any) {
	task, ok := x.(*Schedule)
	if !ok {
		return
	}
	*h = append(*h, task)
}

func (h *ScheduleHeap) Pop() any {
	old := *h
	n := len(old)
	task := old[n-1]
	*h = old[0 : n-1]
	return task
}

func NewScheduleService(executeFunction func(uint64) error, getScheduleFunction func(page, pageSize int) (int64, []*Schedule)) *ScheduleService {
	scheduleHeap := &ScheduleHeap{}
	heap.Init(scheduleHeap)
	return &ScheduleService{
		executeFunction:     executeFunction,
		addChan:             make(chan *Schedule),
		updateChan:          make(chan *Schedule),
		reworkChan:          make(chan *Schedule),
		deleteChan:          make(chan uint64),
		schedules:           scheduleHeap,
		getScheduleFunction: getScheduleFunction,
	}
}

type ScheduleService struct {
	executeFunction     func(uint64) error
	getScheduleFunction func(page, pageSize int) (int64, []*Schedule) //在堆为空时调用
	addChan             chan *Schedule
	updateChan          chan *Schedule
	deleteChan          chan uint64
	stopChan            chan bool
	schedules           *ScheduleHeap
	reworkChan          chan *Schedule
	rwLock              sync.RWMutex
}

func (s *ScheduleService) Start() {
	var trigger *time.Timer
	<-trigger.C

	processDue := func() {
		now := time.Now()
		s.rwLock.RLock()
		defer s.rwLock.RUnlock()

		for s.schedules.Len() > 0 {
			top := (*s.schedules)[0]
			if top.executeTime.After(now) {
				break
			}
			// 弹出并执行
			item := heap.Pop(s.schedules).(*Schedule)
			err := s.executeFunction(item.subjectID)
			if err != nil {
				log.Printf("schedule service execute function error: %v", err)
				go s.rework(item.subjectID)
			}
		}
		if s.schedules.Len() <= 0 {
			var total int64
			page := 1
			pageSize := 200
			total, *s.schedules = s.getScheduleFunction(page, pageSize)
			if total > 0xFFFFFFF { //在数据量大于2^28(268,435,455)时计划崩溃并提示升级，采用更优雅的数据库版本
				log.Fatal("current scheduler solution is too weak in this large data situation,please try to update it to database version")
			}
			if int(total) > pageSize {
				_, *s.schedules = s.getScheduleFunction(page, int(total))
			}
			//检查是否有更新不及时等等
		}
	}
	// 根据当前堆顶重置定时器
	adjustTimer := func() {
		if trigger != nil {
			if !trigger.Stop() {
				// 清空可能已触发的 channel
				select {
				case <-trigger.C:
				default:
				}
			}
			trigger = nil
		}
		if s.schedules.Len() == 0 {
			return
		}
		top := (*s.schedules)[0]
		duration := time.Until(top.executeTime)
		if duration < 0 {
			duration = 0
		}
		trigger = time.NewTimer(duration)
	}

	for {
		processDue()
		adjustTimer()

		select {
		case addItem := <-(s.addChan):
			s.schedules.Push(&addItem)
			processDue()
			adjustTimer()

		case updateItem := <-(s.updateChan):
			var found = false
			for i, item := range *s.schedules {
				if item.subjectID == updateItem.subjectID {
					newItem := &Schedule{
						subjectID:   updateItem.subjectID,
						executeTime: updateItem.executeTime,
					}
					(*s.schedules)[i] = newItem
					heap.Fix(s.schedules, i)
					found = true
					break
				}
			}
			if !found {
				s.schedules.Push(&updateItem)
			}
			processDue()
			adjustTimer()
		case ID := <-s.deleteChan:
			for i, item := range *s.schedules {
				if item.subjectID == ID {
					heap.Remove(s.schedules, i)
				}
			}
		case <-s.stopChan:
			return
		case <-trigger.C: //这时候到期了，执行即可
		}
	}
}

func (s *ScheduleService) rework(ID uint64) {
	retryCountDown := 3
	for retryCountDown > 0 {
		time.Sleep(300 * time.Millisecond)
		s.rwLock.RLock()
		err := s.executeFunction(ID)
		if err == nil {
			return
		}
		log.Printf("schedule service retry(%d) execute function error: %v", retryCountDown, err)
		s.rwLock.RUnlock()
		retryCountDown--
	}
}

func (s *ScheduleService) Stop() {
	s.stopChan <- true
}

func (s *ScheduleService) AddRefreshSchedule(subjectID uint64, executeTime time.Time) {
	s.addChan <- &Schedule{
		subjectID:   subjectID,
		executeTime: executeTime,
	}
}

func (s *ScheduleService) DeleteRefreshSchedule(subjectID uint64) {
	s.deleteChan <- subjectID
}

func (s *ScheduleService) UpdateRefreshSchedule(subjectID uint64, executeTime time.Time) {
	s.updateChan <- &Schedule{
		subjectID:   subjectID,
		executeTime: executeTime,
	}
	return
}

func (s *ScheduleService) SetGetScheduleFunction(getScheduleFunction func(page, pageSize int) (int64, []*Schedule)) {
	s.rwLock.Lock()
	defer s.rwLock.Unlock()
	s.getScheduleFunction = getScheduleFunction

}

func (s *ScheduleService) SetExecuteFunction(executeFunction func(ID uint64) error) error {
	s.rwLock.Lock()
	defer s.rwLock.Unlock()
	s.executeFunction = executeFunction
	return nil
}

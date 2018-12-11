package containers

import (
	"container/list"
	"sync"
)

type SyncContiner struct {
	slist  *list.List
	locker sync.RWMutex
}

func New() *SyncContiner {
	return &SyncContiner{
		slist: list.New(),
	}
}

func NewArray(size int) []*SyncContiner {
	sync := []*SyncContiner{}
	for i := 0; i < size; i++ {
		sync = append(sync, New())
	}

	return sync
}

func (sc *SyncContiner) Push(item interface{}) {
	sc.locker.Lock()
	defer sc.locker.Unlock()

	sc.slist.PushBack(item)
}

func (sc *SyncContiner) Range(fn func(i interface{})) {
	sc.locker.RLock()
	defer sc.locker.RUnlock()

	for e := sc.slist.Front(); e != nil; e = e.Next() {
		fn(e.Value)
	}
}

func (sc *SyncContiner) Search(fn func(i interface{}) bool) (interface{}, bool) {
	sc.locker.RLock()
	defer sc.locker.RUnlock()

	for e := sc.slist.Front(); e != nil; e = e.Next() {
		if fn(e.Value) {
			return e.Value, true
		}
	}

	return nil, false
}

func (sc *SyncContiner) Remove(fn func(i interface{}) bool) {
	sc.locker.Lock()
	defer sc.locker.Unlock()

	for e := sc.slist.Front(); e != nil; e = e.Next() {
		if fn(e.Value) {
			sc.slist.Remove(e)
			return
		}
	}
}

package grapeGuard

import (
	"container/list"
	"sync"
	"time"
)

////////////////////////////////////////////////////////////
//  add by koangel
//	date: 2017/09/30
//  在一个时间周期内超过指定数量的计数则直接报错，未超过 过时重置
////////////////////////////////////////////////////////////

var (
	locker    sync.Mutex
	syncRange = list.New()
)

func procCoudMemory() {
	ticks := time.NewTicker(20 * time.Minute)
	for {
		select {
		case <-ticks.C:
			locker.Lock()
			for e := syncRange.Front(); e != nil; e = e.Next() {
				e.Value.(*TimeGroup).ClearMemory()
			}
			locker.Unlock()
		}
	}
}

type timeCount struct {
	count    int
	nextTime int64
	hotTime  int64
}

func (c *timeCount) tickAdd(limit int, t time.Duration) bool {
	if c.IsExpired() {
		c.count = 0
		c.next(t)
	}

	if c.IsExpired() == false && c.count >= limit {
		return false
	}

	c.hotTime = time.Now().Add(40 * time.Minute).Unix()
	c.count += 1
	return true
}

func (c *timeCount) next(t time.Duration) {
	c.nextTime = time.Now().Add(t).Unix()
}

func (c *timeCount) IsExpired() bool {
	if time.Now().Unix() >= c.nextTime {
		return true
	}

	return false
}

func (c *timeCount) IsNotHot() bool {
	if time.Now().Unix() >= c.hotTime {
		return true
	}

	return false
}

type TimeGroup struct {
	mux      sync.RWMutex
	mapData  map[interface{}]*timeCount
	limit    int
	lootTime time.Duration
	once     *timeCount
}

func init() {
	go procCoudMemory()
}

func NewTimeGroup(loopTime time.Duration, limit int) *TimeGroup {
	ret := &TimeGroup{
		mapData:  map[interface{}]*timeCount{},
		limit:    limit,
		lootTime: loopTime,
		once:     nil,
	}

	return ret
}

func (t *TimeGroup) AddLimitMap(key interface{}) bool {

	if t.once == nil {
		c := &timeCount{
			count:    0,
			nextTime: 0,
		}
		c.next(t.lootTime)
		t.once = c
	}

	if t.once.IsExpired() {
		t.mux.Lock()
		t.mapData = map[interface{}]*timeCount{}
		t.mux.Unlock()
		t.once.next(t.lootTime)
	}

	if t.once.IsExpired() == false && len(t.mapData) >= t.limit {
		return false
	}

	t.mux.RLock()
	_, ok := t.mapData[key]
	t.mux.RUnlock()
	if !ok {
		t.mux.Lock()
		t.mapData[key] = &timeCount{
			count:    0,
			nextTime: 0,
		}
		t.mux.Unlock()
	}

	return true
}

func (t *TimeGroup) ClearMemory() {
	delKey := []interface{}{}
	t.mux.RLock()
	for k, v := range t.mapData {
		if v.IsNotHot() {
			delKey = append(delKey, k)
		}
	}
	t.mux.RUnlock()

	if len(delKey) == 0 {
		return
	}

	t.mux.Lock()
	for _, v := range delKey {
		delete(t.mapData, v)
	}
	t.mux.Unlock()
}

func (t *TimeGroup) AddCount(key interface{}) bool {

	t.mux.RLock()
	mv, ok := t.mapData[key]
	t.mux.RUnlock()
	if !ok {
		c := &timeCount{
			count:    0,
			nextTime: 0,
		}

		t.mux.Lock()
		t.mapData[key] = c
		t.mux.Unlock()

		c.next(t.lootTime)
		return c.tickAdd(t.limit, t.lootTime)
	}

	return mv.tickAdd(t.limit, t.lootTime)
}

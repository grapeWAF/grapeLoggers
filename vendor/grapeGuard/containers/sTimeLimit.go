package containers

import (
	"sync/atomic"
	"time"
)

type STimeLimit struct {
	time     time.Time
	loopTime time.Duration
	count    int32
	limit    int32
}

func NewSTL(loop time.Duration, limit int32) *STimeLimit {
	return &STimeLimit{
		time:     time.Now().Add(loop),
		loopTime: loop,
		count:    0,
		limit:    limit,
	}
}

func NewSTLOnce() *STimeLimit {
	return &STimeLimit{
		time:     time.Now(),
		loopTime: 0,
		count:    0,
		limit:    0,
	}
}
func (s *STimeLimit) AddCount() {
	if s.loopTime != 0 {
		s.Reset(s.loopTime)
	}

	atomic.AddInt32(&s.count, 1)
}

func (s *STimeLimit) IsOverflow(limit int32, limitTime time.Duration) bool {
	s.Reset(limitTime)
	return (s.count >= limit)
}

func (s *STimeLimit) IsTimeout(limit int32, ltime time.Duration) bool {
	return (time.Now().Unix() >= s.time.Add(ltime).Unix() && s.count < limit)
}

func (s *STimeLimit) Reset(limitTime time.Duration) {
	s.loopTime = limitTime
	if time.Now().Unix() >= s.time.Add(limitTime).Unix() {
		atomic.StoreInt32(&s.count, 0)
		s.time = time.Now()
	}
}

var (
	HostGuards = &SUrlMaps{}
)

func AddGuardCount(host string) {
	val, has := HostGuards.LookupS("", host)
	if !has {
		st := NewSTLOnce()
		st.AddCount()
		HostGuards.Map(host, st)
		return
	}

	st := val.(*STimeLimit)
	st.AddCount()
}

func IsOverflow(host string, limit int32, ltime time.Duration) bool {
	val, has := HostGuards.LookupS("", host)
	if !has {
		return false
	}

	return val.(*STimeLimit).IsOverflow(limit, ltime)
}

func IsTimeout(host string, limit int32, ltime time.Duration) bool {
	val, has := HostGuards.LookupS("", host)
	if !has {
		return true
	}

	return val.(*STimeLimit).IsTimeout(limit, ltime)
}

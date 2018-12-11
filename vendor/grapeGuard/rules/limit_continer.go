package rules

import (
	"container/list"
	"strings"
	"sync"
)

////////////////////////////////////////////////////////////
//  add by koangel
//	date: 2017/12/05
//  一种限制性的容器，用于做数据统计，例如相同PATH的匹配数
////////////////////////////////////////////////////////////

type LimitContiner struct {
	limit  int
	con    *list.List
	locker sync.RWMutex
}

func NewLC(lc int) *LimitContiner {
	return &LimitContiner{
		limit: lc,
		con:   list.New(),
	}
}

func (l *LimitContiner) Add(val interface{}) {
	defer l.locker.Unlock()
	l.locker.Lock()

	l.con.PushBack(val)

	if l.con.Len() >= l.limit {
		l.con.Remove(l.con.Front()) // 删除第一个
	}
}

func (l *LimitContiner) Clear() {
	defer l.locker.Unlock()
	l.locker.Lock()

	l.con = list.New()
}

func (l *LimitContiner) Foreach(fn func(n interface{})) {
	l.locker.RLock()
	defer l.locker.RUnlock()

	for e := l.con.Front(); e != nil; e = e.Next() {
		fn(e.Value)
	}
}

func (l *LimitContiner) Search(fn func(n interface{}) bool) bool {
	l.locker.RLock()
	defer l.locker.RUnlock()

	for e := l.con.Front(); e != nil; e = e.Next() {
		if fn(e.Value) {
			return true
		}
	}

	return false
}

func (l *LimitContiner) MatchLimit(fn func(n interface{}) bool, limit int) bool {
	l.locker.RLock()
	defer l.locker.RUnlock()

	lcount := 0
	for e := l.con.Front(); e != nil; e = e.Next() {
		if fn(e.Value) {
			lcount++
		}

		if lcount >= limit {
			return true
		}
	}

	return false
}

type StringLimit struct {
	lc *LimitContiner
}

func NewSL(ls int) *StringLimit {
	return &StringLimit{
		lc: NewLC(ls),
	}
}

func (l *StringLimit) Push(v string) {
	l.lc.Add(strings.ToLower(v))
}

func (l *StringLimit) MatchPrefix(src string, limit int) bool {
	pathPos := strings.LastIndex(src, "/")
	if pathPos == -1 {
		return false
	}

	path := src[:pathPos]

	return l.lc.MatchLimit(func(v interface{}) bool {
		return strings.HasPrefix(v.(string), path)
	}, limit)
}

func (l *StringLimit) Match(val string, limit int) bool {
	return l.lc.MatchLimit(func(v interface{}) bool {
		return (strings.ToLower(val) == v.(string))
	}, limit)
}

func (l *StringLimit) LineMatch(limit int) bool {
	pathData := l.Strings()
	count := 0
	for i := len(pathData) - 1; i > 0; i-- {
		if (i - 1) < 0 {
			break
		}

		if pathData[i] != pathData[i-1] {
			break
		}

		count++
	}

	return (count >= limit)
}

func (l *StringLimit) Strings() []string {
	rv := []string{}
	l.lc.Foreach(func(v interface{}) {
		rv = append(rv, v.(string))
	})

	return rv
}

func (l *StringLimit) Search(val string) bool {
	return l.lc.Search(func(v interface{}) bool {
		return (strings.ToLower(val) == v.(string))
	})
}

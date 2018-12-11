package containers

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"sync"

	"github.com/koangel/grapeNet/Utils"
)

func Contain(obj interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}

	return false
}

// 一个随随便便的小容器，模拟基本的VECTOR
type SVector struct {
	c   sync.RWMutex
	vec []interface{}
}

func (s *SVector) Len() int {
	s.c.RLock()
	defer s.c.RUnlock()

	return len(s.vec)
}

func (s *SVector) Append(item ...interface{}) {
	s.c.Lock()
	defer s.c.Unlock()

	s.vec = append(s.vec, item...)
}

func (s *SVector) I(index int) (val interface{}, ok bool) {
	s.c.RLock()
	defer s.c.RUnlock()

	if index < 0 || index >= len(s.vec) {
		return nil, false
	}

	return s.vec[index], true
}

func (s *SVector) Range(fn func(index int, val interface{}) bool) {
	s.c.RLock()
	defer s.c.RUnlock()

	for i := 0; i < len(s.vec); i++ {
		if fn(i, s.vec[i]) == false {
			break
		}
	}
}

func (s *SVector) DRange(fn func(index int, val interface{}) bool) {
	s.c.RLock()
	defer s.c.RUnlock()

	total := len(s.vec)

	firstPos := 0
	secondPos := total / 2 // 从尾部开始
	firstEnd := secondPos - 1
	secondEnd := total

	for {
		if firstPos >= firstEnd && secondPos >= secondEnd {
			break
		}

		if firstPos < firstEnd {
			if fn(firstPos, s.vec[firstPos]) == false {
				break
			}

			firstPos++
		}

		if secondPos < secondEnd {
			if fn(secondPos, s.vec[secondPos]) == false {
				break
			}

			secondPos++
		}
	}
}

func (s *SVector) Remove(i int) bool {
	s.c.Lock()
	defer s.c.Unlock()

	if i < 0 || i >= len(s.vec) {
		return false
	}

	s.vec = append(s.vec[:i], s.vec[i+1:]...)
	return true
}

// 一个正则容器项目，用于快速匹配正则数据
type SRegexItem struct {
	Regex string
	RM    *regexp.Regexp
	Item  interface{}
}

func (r *SRegexItem) Compile(src string) error {
	_, err := regexp.Compile(src)
	if err != nil {
		return err
	}

	r.RM = regexp.MustCompile(src)
	return nil
}

func (r *SRegexItem) MatchAll(src string) bool {
	q := r.RM.FindAllString(src, -1)
	if len(q) != 1 {
		return false
	}

	return true
}

func (r *SRegexItem) Match(src string) bool {
	return r.RM.MatchString(src)
}

func (r *SRegexItem) Extend(src, template string) string {
	dst := []byte{}
	match := r.RM.FindStringSubmatchIndex(src)
	dst = r.RM.ExpandString(dst, template, src, match)

	return string(dst)
}

type SRegexVec struct {
	SVector

	index []int
}

func NewRegVecArray(max int, extend []int) []*SRegexVec {
	ret := []*SRegexVec{}
	for i := 0; i < max; i++ {
		if Contain(i, extend) {
			ret = append(ret, nil)
			continue
		}
		ret = append(ret, new(SRegexVec))
	}

	return ret
}

func (s *SRegexVec) AddRegexS(src string, val interface{}) error {
	regp, err := regexp.Compile(src)
	if err != nil {
		return err
	}

	newRegex := new(SRegexItem)
	newRegex.Item = val
	newRegex.Regex = src
	newRegex.RM = regp

	s.Append(newRegex)
	return nil
}

func (s *SRegexVec) Lookup(url string) (vidx int, val interface{}, ok bool) {
	val = nil
	ok = false
	vidx = -1

	s.Range(func(index int, item interface{}) bool {
		regex := item.(*SRegexItem)
		if regex.Match(url) {
			vidx = index
			val = regex.Item
			ok = true
			return false
		}

		return true
	})

	return
}

func (s *SRegexVec) BuildJobs(maxgo int) {

	maxproc := 1000
	if s.Len() <= maxproc {
		return
	}

	psize := s.Len() / maxproc
	lowCount := s.Len() % maxproc
	if psize > maxgo {
		psize = maxgo
		maxproc = s.Len() / maxgo

		if (maxproc * maxgo) < s.Len() {
			lowCount = s.Len() - (maxproc * maxgo)
		}
	}

	for i := 0; i < psize; i++ {

		s.index = append(s.index, Utils.Ifn(i > 0, i*maxproc+1, i*maxproc))
		s.index = append(s.index, Utils.Ifn(i == psize-1, (i*maxproc+maxproc)+lowCount, (i*maxproc+maxproc)))
	}
}
func (s *SRegexVec) Prefix(url string) (vidx int, val interface{}, ok bool) {
	vidx = -1
	val = nil
	ok = false
	s.Range(func(index int, item interface{}) bool {
		regex := item.(*SRegexItem)
		if strings.HasPrefix(regex.Regex, url) {

			vidx = index
			val = regex.Item
			ok = true

			return false
		}

		return true
	})

	return
}

func (s *SRegexVec) LookupP(url string) (vidx int, val interface{}, ok bool) {
	if len(s.index) == 0 {
		return s.Lookup(url)
	}

	val = nil
	ok = false
	vidx = -1

	var cl sync.Mutex
	isFind := false

	var wait sync.WaitGroup
	for i := 0; i < len(s.index); i += 2 {
		wait.Add(1)

		go func(start, end int) {
			s.c.RLock()
			defer s.c.RUnlock()

			for idx := start; idx < end; idx++ {
				if isFind {
					break
				}

				if idx >= len(s.vec) {
					break
				}

				item := s.vec[idx].(*SRegexItem)
				if item.Match(url) {
					val = item.Item
					ok = true
					vidx = idx

					cl.Lock()
					isFind = true
					cl.Unlock()
					break
				}
			}

			wait.Done()
		}(s.index[i], s.index[i+1])
	}

	wait.Wait()

	return
}

func (s *SRegexVec) Extend(url, template string) (val string, ok bool) {
	i, _, has := s.Lookup(url)
	if !has {
		return "", has
	}

	s.c.RLock()
	defer s.c.RUnlock()

	regex := s.vec[i].(*SRegexItem)
	ok = true
	val = regex.Extend(url, template)
	return
}

func (s *SRegexVec) dump() {
	s.c.RLock()
	defer s.c.RUnlock()

	for _, v := range s.vec {
		fmt.Println(v)
	}
}

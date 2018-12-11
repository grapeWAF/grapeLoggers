package containers

import (
	"sort"
	"strings"
	"sync"
)

type SUrlMaps struct {
	maps  sync.Map
	hosts sync.Map
}

func NewSMArray(num int) []*SUrlMaps {
	r := []*SUrlMaps{}
	for i := 0; i < num; i++ {
		r = append(r, &SUrlMaps{})
	}

	return r
}

func (s *SUrlMaps) getHosts(scheme, url string) []string {
	keys := scheme + url
	vals, has := s.hosts.Load(keys)
	if has {
		return vals.([]string)
	}

	hosts := []string{keys}
	labels := strings.Split(url, ".")

	if len(labels) <= 2 {
		labels = append([]string{"*"}, labels...)
	}

	for i := range labels {
		labels[i] = "*"
		if (i + 2) >= len(labels) {
			break
		}

		hosts = append(hosts, scheme+strings.Join(labels[i:], "."))
	}

	s.hosts.Store(keys, hosts)
	return hosts
}

func (s *SUrlMaps) Lookup(host string) (val interface{}, has bool) {
	return s.maps.Load(strings.ToLower(host))
}

func (s *SUrlMaps) LookupS(scheme, host string) (val interface{}, has bool) {
	for _, vh := range s.getHosts(scheme, strings.ToLower(host)) {
		tv, ok := s.maps.Load(vh)
		if ok {
			return tv, ok
		}
	}

	return nil, false
}

func (s *SUrlMaps) Map(key string, val interface{}) {
	s.maps.Store(strings.ToLower(key), val)
}

func (s *SUrlMaps) DeleteS(scheme, host string) {
	for _, vh := range s.getHosts(scheme, host) {
		s.maps.Delete(vh)
	}
}

func (s *SUrlMaps) Delete(key string) {
	s.maps.Delete(key)
}

func (s *SUrlMaps) Range(fn func(key, value interface{}) bool) {
	s.maps.Range(fn)
}

func (s *SUrlMaps) SortRange(fn func(key, value interface{}) bool) {
	keys := []string{}
	s.maps.Range(func(key, value interface{}) bool {
		keys = append(keys, key.(string))
		return true
	})

	sort.Strings(keys)

	for _, vk := range keys {
		v, has := s.maps.Load(vk)
		if has && fn(vk, v) == false {
			break
		}
	}
}

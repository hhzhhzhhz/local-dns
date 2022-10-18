package cache

import "sync"

type Set interface {
	Add(value ...interface{})
	Contains(value interface{}) bool
	Values() []interface{}
	Len() int64
}

func NewSet() Set {
	return &set{
		mset: make(map[interface{}]struct{}, 5),
	}
}

type set struct {
	mset map[interface{}]struct{}
	Set []interface{}
	rmux sync.RWMutex
	size int64
}

func (s *set) Add(values ...interface{}) {
	if len(values) == 0 {
		return
	}
	s.rmux.Lock()
	for _, v := range values {
		if _, ok := s.mset[v]; ok {
			continue
		}
		s.size++
		s.mset[v] = struct{}{}
		s.Set = append(s.Set, v)
	}
	s.rmux.Unlock()
}

func (s *set) Contains(value interface{}) bool {
	s.rmux.RLock()
	_, ok := s.mset[value]
	s.rmux.RUnlock()
	return ok
}
func (s *set)  Values() (vs []interface{}) {
	s.rmux.RLock()
	for _, v := range s.Set {
		vs = append(vs, v)
	}
	s.rmux.RUnlock()
	return vs
}

func (s *set) Len() int64 {
	return s.size
}
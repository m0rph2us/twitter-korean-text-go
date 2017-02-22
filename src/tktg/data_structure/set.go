package data_structure

import (
	"fmt"
	"sync"
)

type Set interface {
	Size() int
	Add(value interface{}) bool
	Iter() chan interface{}
	ToSlice() []interface{}
	Contains(value interface{}) bool
	Remove(value interface{})
}

type linked_set struct {
	sync.RWMutex
	unordered_map map[string]*Node
	list          LinkedList
}

func NewLinkedSet(values ...interface{}) Set {
	s := &linked_set{
		unordered_map:   map[string]*Node{},
		list: NewLinkedList(),
	}

	for _, item := range values {
		s.add(item)
	}

	return s
}

func (t *linked_set) getKey(value interface{}) string {
	return fmt.Sprintf("%v", value)
}

func (t *linked_set) Size() int {
	return len(t.unordered_map)
}

func (t *linked_set) add(value interface{}) bool {
	k := t.getKey(value)

	if t.exist(k) {
		return false
	}

	n := t.list.Add(value)

	t.unordered_map[k] = n

	return true
}

func (t *linked_set) Add(value interface{}) bool {
	t.Lock()
	defer t.Unlock()

	ret := t.add(value)

	return ret
}

func (t *linked_set) exist(key string) bool {
	if _, e := t.unordered_map[key]; e == false {
		return false
	}
	return true
}

func (t *linked_set) Iter() chan interface{} {
	c := make(chan interface{})
	go func() {
		for v := range t.list.Iter() {
			c <- v
		}
		close(c)
	}()
	return c
}

func (t *linked_set) ToSlice() []interface{} {
	return t.list.ToSlice()
}

func (t *linked_set) Contains(value interface{}) bool {
	k := t.getKey(value)
	return t.exist(k)
}

func (t *linked_set) Remove(value interface{}) {
	t.Lock()
	defer t.Unlock()

	k := t.getKey(value)

	t.list.RemoveNode(t.unordered_map[k])

	delete(t.unordered_map, k)
}

package data_structure

import (
	"sync"
)

type Map interface {
	Size() int
	Put(key interface{}, value interface{}) bool
	Get(key interface{}) interface{}
	Iter() chan nodeItem
	Contains(key interface{}) bool
	Remove(key interface{})
}

type nodeItem struct {
	Key interface{}
	Value interface{}
}

type linked_map struct {
	sync.RWMutex
	unordered_map   map[interface{}]*Node
	list LinkedList
}

func NewLinkedMap() Map {
	s := &linked_map{
		unordered_map:   map[interface{}]*Node{},
		list: NewLinkedList(),
	}

	return s
}

func (t *linked_map) Size() int {
	return len(t.unordered_map)
}

func (t *linked_map) put(key interface{}, value interface{}) bool {
	if item, e := t.unordered_map[key]; e == true {
		item.value.(*nodeItem).Value = value
	} else {
		n := t.list.Add(&nodeItem{key, value})
		t.unordered_map[key] = n
	}

	return true
}

func (t *linked_map) Put(key interface{}, value interface{}) bool {
	t.Lock()
	defer t.Unlock()

	ret := t.put(key, value)

	return ret
}

func (t *linked_map) Get(key interface{}) interface{} {
	return t.unordered_map[key].value.(*nodeItem).Value
}

func (t *linked_map) Iter() chan nodeItem {
	d := make(chan nodeItem)
	go func() {
		for v := range t.list.Iter() {
			d <- *(v.(*nodeItem))
		}
		close(d)
	}()
	return d
}

func (t *linked_map) Contains(key interface{}) bool {
	_, e := t.unordered_map[key];
	return e
}

func (t *linked_map) Remove(key interface{}) {
	t.Lock()
	defer t.Unlock()

	t.list.RemoveNode(t.unordered_map[key])

	delete(t.unordered_map, key)
}

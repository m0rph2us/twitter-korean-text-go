package data_structure

import "sync"

type LinkedList interface {
	Add(value interface{}) *Node
	Size() int
	Iter() chan interface{}
	ToSlice() []interface{}
	Contains(value interface{}) bool
	Remove(value interface{})
	RemoveNode(node *Node)
}

type Node struct {
	value interface{}
	prev *Node
	next *Node
}

type linkedList struct {
	sync.RWMutex
	count int
	head *Node
	tail *Node
}

func NewLinkedList() LinkedList {
	s := &linkedList{
		count: 0,
		head: nil,
		tail: nil,
	}
	return s
}

func (t *linkedList)Add(value interface{}) *Node {
	t.Lock()
	defer t.Unlock()

	node := &Node{value, nil, nil}

	if t.head == nil && t.tail == nil {
		t.head = node
		t.tail = node
	} else {
		t.tail.next = node
		node.prev = t.tail
		t.tail = node
	}

	t.count++

	return node
}

func (t *linkedList)Size() int {
	return t.count
}

func (t *linkedList)Iter() chan interface{} {
	c := make(chan interface{})
	go func() {
		cur := t.head
		for {
			if cur == nil {
				break
			}
			c <- cur.value
			cur = cur.next
		}
		close(c)
	}()
	return c
}

func (t *linkedList)ToSlice() []interface{} {
	s := []interface{}{}
	for v := range t.Iter() {
		s = append(s, v)
	}
	return s
}

func (t *linkedList)Contains(value interface{}) bool {
	for v := range t.Iter() {
		if v == value {
			return true
		}
	}
	return false
}

func (t *linkedList)Remove(value interface{}) {
	cur := t.head
	for {
		if cur == nil {
			break
		}

		if cur.value == value {
			t.RemoveNode(cur)
		}

		cur = cur.next
	}
}

func (t *linkedList)RemoveNode(node *Node) {
	t.Lock()
	defer t.Unlock()

	if node == nil {
		return
	}

	if node == t.head {
		t.head = t.head.next
	}

	if node == t.tail {
		t.tail = t.tail.prev
	}

	if node.next != nil {
		node.next.prev = node.prev
	}

	if node.prev != nil {
		node.prev.next = node.next
	}

	t.count--
}






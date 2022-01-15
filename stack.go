package adt

import (
	"sync"
)

type stack struct {
	size int
	items []interface{}
	mx sync.Mutex
}

func NewStack() *stack {
	st := stack{
		size: 0,
		items: make([]interface{}, 0),
	}
	return &st
}

func (st *stack) Size() int {
	return st.size
}

func (st *stack) Push(elm interface{}) {
	st.mx.Lock()
	defer st.mx.Unlock()
	st.size++
	st.items = append(st.items, elm)
}

func (st *stack) Pop() interface{} {
	st.mx.Lock()
	defer st.mx.Unlock()
	if st.size == 0 {
		return nil
	}
	st.size--
	elm := st.items[st.size]
	st.items = st.items[:st.size]
	return elm
}

func (st *stack) Top() interface{} {
	st.mx.Lock()
	defer st.mx.Unlock()
	if st.size == 0 {
		return nil
	}
	return st.items[st.size-1]
}

func (st *stack) Destroy() {
	st.mx.Lock()
	defer st.mx.Unlock()
	st.items = nil
	st.items = make([]interface{}, 0)
	st.size = 0
}
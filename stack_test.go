package adt

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestNewStack(t *testing.T) {
	st := NewStack()
	assert.NotNil(t, st)
}

func TestStack_Size(t *testing.T) {
	tt := []struct{
		name string
		prestate func(st *stack)
		expectedSize int
	}{
		{
			name: "initial stack should return size of zero",
			prestate: func(st *stack) {},
			expectedSize: 0,
		},
		{
			name: "initial stack should increment size value after pushes an element",
			prestate: func(st *stack) {
				st.Push("elm")
			},
			expectedSize: 1,
		},
		{
			name: "should decrease size value when pops an element",
			prestate: func(st *stack) {
				st.Push("elm")
				st.Pop()
			},
			expectedSize: 0,
		},
		{
			name: "should return size 0 when pops from an empty stack",
			prestate: func(st *stack) {
				st.Pop()
			},
			expectedSize: 0,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			st := NewStack()
			tc.prestate(st)
			assert.EqualValues(t, tc.expectedSize, st.Size())
		})
	}
}

func TestStack_Push(t *testing.T) {
	tests := []struct {
		name string
		prestate func(st *stack)
		expectedSize int
		expectedTop interface{}
	}{
		{
			name: "push one element",
			prestate: func(st *stack) {
				st.Push("elm")
			},
			expectedSize: 1,
			expectedTop: "elm",
		},
		{
			name: "push many elements",
			prestate: func(st *stack) {
				st.Push("elm1")
				st.Push("elm2")
				st.Push("elm3")
			},
			expectedSize: 3,
			expectedTop: "elm3",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			st := NewStack()
			test.prestate(st)
			assert.EqualValues(t, test.expectedSize, st.Size())
			assert.EqualValues(t, test.expectedTop, st.Top())
		})
	}
}

func TestStack_Pop(t *testing.T) {
	tests := []struct {
		name string
		prestate func(st *stack)
		expectedSize int
		expectedItem interface{}
	}{
		{
			name: "pop from an empty stack should return nil",
			prestate: func(st *stack) {},
			expectedSize: 0,
			expectedItem: nil,
		},
		{
			name: "should pop an item from stack with one element",
			prestate: func(st *stack) {
				st.Push("elm")
			},
			expectedSize: 0,
			expectedItem: "elm",
		},
		{
			name: "should pop an item from stack with many elements",
			prestate: func(st *stack) {
				st.Push("elm1")
				st.Push("elm2")
				st.Push("elm3")
			},
			expectedSize: 2,
			expectedItem: "elm3",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			st := NewStack()
			test.prestate(st)
			item := st.Pop()
			assert.EqualValues(t, test.expectedItem, item)
		})
	}
}

func TestStack_Top(t *testing.T) {
	st := NewStack()
	st.Push("elm")
	assert.EqualValues(t, 1, st.Size())
	assert.EqualValues(t, "elm", st.Top())
	assert.EqualValues(t, 1, st.Size())
}

func TestStack_TopFromAnEmptyStack(t *testing.T) {
	st := NewStack()
	assert.EqualValues(t, nil, st.Top())
	assert.EqualValues(t, 0, st.Size())
}

func TestStack_TopCalledManyTimes(t *testing.T) {
	st := NewStack()
	st.Push("elm")
	assert.EqualValues(t, 1, st.Size())
	assert.EqualValues(t, "elm", st.Top())
	assert.EqualValues(t, "elm", st.Top())
	assert.EqualValues(t, "elm", st.Top())
	assert.EqualValues(t, 1, st.Size())
}

func TestStack_ConcurrencyPush(t *testing.T) {
	st := NewStack()
	wg := sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		elm := i
		wg.Add(1)
		go func() {
			st.Push(elm)
			wg.Done()
		}()
	}
	wg.Wait()
	assert.EqualValues(t, 1000, st.Size())
}

func TestStack_ConcurrencyPop(t *testing.T) {
	st := NewStack()
	for i := 0; i < 1000; i++ {
		elm := i
		st.Push(elm)
	}
	assert.EqualValues(t, 1000, st.Size())

	wg := sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			_ = st.Pop()
			wg.Done()
		}()
	}
	wg.Wait()
	assert.EqualValues(t, 0, st.Size())
}

func TestStack_ConcurrencyTop(t *testing.T) {
	st := NewStack()
	for i := 0; i < 1000; i++ {
		elm := i
		st.Push(elm)
	}
	assert.EqualValues(t, 1000, st.Size())

	wg := sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			_ = st.Top()
			wg.Done()
		}()
	}
	wg.Wait()
	assert.EqualValues(t, 1000, st.Size())
}

func TestStack_IntegrationConcurrency(t *testing.T) {
	st := NewStack()
	wgPush := sync.WaitGroup{}
	wgPop := sync.WaitGroup{}
	wgPush.Add(1)
	wgPop.Add(1)
	go func(wg *sync.WaitGroup) {
		for i := 0; i < 1000; i++ {
			elm := i
			st.Push(elm)
		}
		wg.Done()
	}(&wgPush)
	go func(wg *sync.WaitGroup) {
		for i := 0; i < 1000; i++ {
			_ = st.Pop()
		}
		wg.Done()
	}(&wgPop)
	go func() {
		for i := 0; i < 1000; i++ {
			_ = st.Top()
		}
	}()
	wgPush.Wait()
	wgPop.Wait()
}

func TestStack_Volume(t *testing.T) {
	st := NewStack()
	count := 1000000
	for i := 0; i < count; i++ {
		elm := i
		st.Push(elm)
	}
	assert.EqualValues(t, count, st.Size())
	st.Destroy()
	assert.EqualValues(t, 0, st.Size())
	st.Push("item")
	assert.EqualValues(t, 1, st.Size())
	assert.EqualValues(t, "item", st.Top())
}

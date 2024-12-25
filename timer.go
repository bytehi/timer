package bunch

import (
	"container/heap"
	"time"
)

type heapNode struct {
	interval      time.Duration
	executionTime time.Time
	callback      func(Cancel)
	canceled      bool
}

type minHeap []*heapNode

func (h minHeap) Len() int           { return len(h) }
func (h minHeap) Less(i, j int) bool { return h[i].executionTime.Before(h[j].executionTime) }
func (h minHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *minHeap) Push(x interface{}) {
	*h = append(*h, x.(*heapNode))
}

func (h *minHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (h *minHeap) Peek() interface{} {
	return (*h)[0]
}

type Timer struct {
	h minHeap
}

func New() *Timer {
	return &Timer{}
}

type Cancel func()

func (t *Timer) Add(duration time.Duration, callback func(Cancel)) Cancel {
	node := &heapNode{
		interval:      duration,
		executionTime: time.Now().Add(duration),
		callback:      callback,
	}
	heap.Push(&t.h, node)
	return func() {
		node.canceled = true
	}
}

func (t *Timer) Timeout(now time.Time) {
	for {
		first := t.h.Peek()
		if first == nil {
			break
		}
		node := first.(*heapNode)
		if node.executionTime.After(now) {
			break
		}
		heap.Pop(&t.h)
		if node.canceled {
			continue
		}
		node.callback(func() { node.canceled = true })
		node.executionTime = now.Add(node.interval)
		heap.Push(&t.h, node)
	}
}

func (t *Timer) NextTime() (time.Time, bool) {
	for {
		first := t.h.Peek()
		if first == nil {
			return time.Time{}, false
		}
		node := first.(*heapNode)
		if node.canceled {
			heap.Pop(&t.h)
			continue
		}
		return node.executionTime, true
	}
}

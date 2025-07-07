package queue

import (
	"container/heap"
	"errors"
	job "github.com/mdavecraft/job-queue/job"
	nanoid "github.com/mdavecraft/job-queue/nanoid"
)

type ComparatorFunc[T any] func(a, b T) int

type PriorityQueueType int

const (
	MaxHeap PriorityQueueType = iota
	MinHeap
)

const (
	defaultInitialCapacity = 0
)

var (
	ErrQueueEmpty = errors.New("queue is empty")
)

type PriorityQueue struct {
	heap.Interface
	Id         string
	Jobs       []*job.Job
	Type       PriorityQueueType
	Tags       map[string]string
	Comparator ComparatorFunc[*job.Job]
}

func (p *PriorityQueue) Len() int {
	return len(p.Jobs)
}

func (p *PriorityQueue) Swap(i, j int) {
	p.Jobs[i], p.Jobs[j] = p.Jobs[j], p.Jobs[i]
	p.Jobs[i].Index = i
	p.Jobs[j].Index = j
}

func (p *PriorityQueue) Less(i, j int) bool {
	cmpResult := p.Comparator(p.Jobs[i], p.Jobs[j])
	if p.Type == MinHeap {
		return cmpResult < 0
	}
	return cmpResult > 0
}

func (p *PriorityQueue) Push(x any) {
	jobItem := x.(*job.Job)
	jobItem.Index = p.Len()
	p.Jobs = append(p.Jobs, jobItem)
}

func (p *PriorityQueue) Pop() any {
	old := p.Jobs
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	p.Jobs = old[0 : n-1]
	item.Index = -1
	return item
}

func (p *PriorityQueue) Peek() (*job.Job, error) {
	if p.Len() == 0 {
		return nil, ErrQueueEmpty
	}
	return p.Jobs[0], nil
}

func (p *PriorityQueue) IsEmpty() bool {
	return p.Len() == 0
}

func defaultComparator(a, b *job.Job) int {
	if p := a.Priority - b.Priority; p != 0 {
		return int(p)
	}
	if t := a.CreatedAt.Compare(b.CreatedAt); t != 0 {
		return t
	}
	return a.UpdatedAt.Compare(b.UpdatedAt)
}

func NewPriorityQueue(queueType PriorityQueueType, comparator ComparatorFunc[*job.Job]) (*PriorityQueue, error) {
	id, err := nanoid.New()
	if err != nil {
		return nil, nanoid.ErrFailedToGenerateID
	}

	if comparator == nil {
		comparator = defaultComparator
	}

	pq := &PriorityQueue{
		Id:         id,
		Jobs:       make([]*job.Job, defaultInitialCapacity),
		Type:       queueType,
		Tags:       make(map[string]string),
		Comparator: comparator,
	}
	heap.Init(pq)
	return pq, nil
}
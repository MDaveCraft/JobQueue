package queue

import (
	"sync"
	job "github.com/mdavecraft/job-queue/job"
	nanoid "github.com/mdavecraft/job-queue/nanoid"
)

type PriorityQueue struct {
	id string
	mu sync.Mutex
	jobs map[job.Priority][]*job.Job
	priorityOrder []job.Priority
	buffer []interface{}
	head int
	tail int
	deadLetterQueue []*job.Job
	notEmpty *sync.Cond
}

func NewQueue(priorityOrder[]job.Priority, bufferSize int) (*PriorityQueue, error) {
	if len(priorityOrder) == 0 {
		return nil, job.ErrEmptyPriorityOrder
	}
	if bufferSize <= 0 {	
		return nil, job.ErrInvalidBufferSize
	}

	id, err := nanoid.Generate(nanoid.DefaultAlphabetString, 21)
	if err != nil {
		return nil, err
	}
	pq := &PriorityQueue{
		id: id,
		jobs: make(map[job.Priority][]*job.Job),
		priorityOrder: priorityOrder,
		buffer: make([]interface{}, bufferSize),
		head: 0,
		tail: 0,
		deadLetterQueue: make([]*job.Job, 0),
	}
	pq.notEmpty = sync.NewCond(&pq.mu)
	return pq, nil
}

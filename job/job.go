package job

import (
	"errors"
	"time"

	nanoid "github.com/mdavecraft/job-queue/nanoid"
)

const (
	DefaultMaxRetryCount = 5
)

var (
	ErrInvalidBufferSize = errors.New("buffer size must be greater than zero")
	ErrEmptyPriorityOrder = errors.New("priority order cannot be empty")
)

type JobStatus int

type Payload map[string]interface{}

type Job struct {
	Id string `json:"id"`
	Type string `json:"type"`
	Payload Payload `json:"payload"`
	Priority Priority `json:"priority"`
	RetryCount int `json:"retry_count"`
	MaxRetries int `json:"max_retries"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Status JobStatus `json:"status"`
	Metadata Metadata `json:"metadata,omitempty"`
}

const (
	New JobStatus = iota
	Pending
	Queued
	Scheduled
	Completed
	Failed
	Deferred
	Cancelled
	Paused
)

func NewJob(jobType string, payload Payload, priority Priority, maxRetries int, metaData *Metadata, ) *Job {
	id, err := nanoid.Generate(nanoid.DefaultAlphabetString, 21)
	if err != nil {
		panic("Failed to generate job ID: " + err.Error())
	}
	return &Job{
		Id:       id,
		Type:     jobType,
		Payload:  payload,
		Priority: priority,
		RetryCount: 0,
		MaxRetries: maxRetries,
		Status:   New,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Metadata: *metaData,
	}
}

func (j *Job) UpdateStatus(newStatus JobStatus) {
	j.Status = newStatus
	j.UpdatedAt = time.Now()
}

func (j *Job) UpdateMetadata(key string, value any) {
	if j.Metadata == nil {
		j.Metadata = make(Metadata)
	}
	j.Metadata.Set(key, value)
}

func (j *Job) GetMetadata(key string) (any, bool) {
	if j.Metadata == nil {
		return nil, false
	}
	return j.Metadata.Get(key)
}

func (j *Job) IncrementRetryCount() {
	if j.RetryCount < j.MaxRetries {
		j.RetryCount++
		j.UpdateStatus(Pending)
	} else {
		j.UpdateStatus(Failed)
	}
	j.UpdatedAt = time.Now()
}

func (j *Job) ResetRetryCount() {
	j.RetryCount = 0
	j.UpdateStatus(New)
	j.UpdatedAt = time.Now()
}

func (j *Job) IsRetryable() bool {
	return j.RetryCount < j.MaxRetries
}

func (j *Job) ReNice(priority Priority){
	j.Priority = priority
	j.UpdatedAt = time.Now()
}
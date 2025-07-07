package job

import "errors"

type Priority int

var ErrPriority = errors.New("invalid priority value")

const (
	HighPriority Priority = iota + 1 
	MediumPriority
	LowPriority
)
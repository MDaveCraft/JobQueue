package job

type Priority int

const (
	DefaultPriority Priority = 6
	HighPriority    Priority = 1
	MediumPriority  Priority = 5
	LowPriority     Priority = 10
)
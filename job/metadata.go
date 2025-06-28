package job

type Metadata map[string]any

func (m Metadata) Set(key string, value any) {
	if m == nil {
		m = make(Metadata)
	}
	m[key] = value
}

func (m Metadata) Get(key string) (any, bool) {
	if m == nil {
		return nil, false
	}
	value, exists := m[key]
	return value, exists
}
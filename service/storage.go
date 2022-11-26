package service

type Storage[T interface{}] struct {
	storage map[string]T
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Success  bool   `json:"success"`
}

type Authorization struct {
	ChatID     int64
	Authorized bool
}

func NewStorage[T interface{}]() *Storage[T] {
	return &Storage[T]{
		storage: make(map[string]T),
	}
}

func (s *Storage[V]) Get(key string) *V {
	if val, ok := s.storage[key]; !ok {
		return nil
	} else {
		return &val
	}
}

func (s *Storage[V]) Set(key string, val V) {
	s.storage[key] = val
}

func (s *Storage[V]) List() []V {
	vs := make([]V, 0)
	for _, val := range s.storage {
		vs = append(vs, val)
	}

	return vs
}

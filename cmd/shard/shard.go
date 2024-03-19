package main

type Shard struct {
	key_value map[string]string
}

func NewShard() *Shard {
	s := &Shard{}
	s.key_value = make(map[string]string)
	return s
}

func (s *Shard) Get(key string) string {
	return s.key_value[key]
}

func (s *Shard) Set(key string, value string) {
	s.key_value[key] = value
}

func (s *Shard) Delete(key string) {
	delete(s.key_value, key)
}

func (s *Shard) GetAll() map[string]string {
	return s.key_value
}

func (s *Shard) GetKeys() []string {
	keys := make([]string, 0, len(s.key_value))
	for k := range s.key_value {
		keys = append(keys, k)
	}
	return keys
}

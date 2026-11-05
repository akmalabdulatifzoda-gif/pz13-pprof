package work

import (
	"sync"
	"time"
)

// Store — общий ресурс для демонстрации Block/Mutex-профилей.
type Store struct {
	mu   sync.Mutex
	hits map[string]int
}

func NewStore() *Store { return &Store{hits: make(map[string]int)} }

// slowWork имитирует «долгую» работу (подготовка данных, сериализация и т.п.).
func slowWork() { time.Sleep(2 * time.Millisecond) }

// HitSlow держит мьютекс во время долгой работы — горутины ждут друг друга.
func (s *Store) HitSlow(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	slowWork()
	s.hits[key]++
}

// HitFast выполняет работу вне критической секции; под мьютексом —
// только изменение map.
func (s *Store) HitFast(key string) {
	slowWork()
	s.mu.Lock()
	s.hits[key]++
	s.mu.Unlock()
}

func (s *Store) Count(key string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.hits[key]
}

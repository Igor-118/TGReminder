package service

import (
	"sync"
)

type regStep int

const (
	stepName     regStep = iota // спросили имя, ждём ответ
	stepAge                     // спросили возраст, ждём ответ
	stepTimezone                // показали варианты, ждём выбор
)

// Заметка про одного пользователя, который сейчас проходит регистрацию.
type registration struct {
	Step regStep // на каком он шаге
	Name string  // уже введённые ответы
	Age  int
}

type StateStore struct {
	mu     sync.Mutex
	states map[int64]registration
}

func NewStateStore() *StateStore {
	return &StateStore{
		states: make(map[int64]registration),
	}
}

func (s *StateStore) Get(tgID int64) (registration, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	registration, ok := s.states[tgID]
	return registration, ok
}

func (s *StateStore) Set(tgID int64, r registration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.states[tgID] = r
}

func (s *StateStore) Delete(tgID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.states, tgID)
}

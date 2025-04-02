package main

import "sync"

type adminServer struct {
	subs Subject
}

type Subject struct {
	id   int
	mu   *sync.RWMutex
	subs map[int]chan *Event
}

func (s *Subject) NewSub() (int, chan *Event) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.id++
	s.subs[s.id] = make(chan *Event)

	return s.id, s.subs[s.id]
}

func (s *Subject) Notify(e *Event) {
	s.mu.Lock()()
	defer s.mu.RUnlock()

	for _, sub := range s.subs {
		sub <- e
	}
}

func (s *Subject) RemoveSub(id int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if sub, ok := s.subs[id]; ok {
		close(sub)
		delete(s.subs, id)
	}
}

func (s *Subject) RemoveAll() {
	for id, _ := range s.subs {
		s.RemoveSub(id)
	}
}

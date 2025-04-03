package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"
)

type adminServer struct {
	UnimplementedAdminServer
	subs *Subscriber
}

func newAdminServer(subs *Subscriber) *adminServer {
	return &adminServer{
		subs: subs,
	}
}

func (a *adminServer) Logging(_ *Nothing, srv Admin_LoggingServer) error {
	id, events := a.subs.Attach()
	defer a.subs.Dettach(id)

	for e := range events {
		err := srv.Send(e)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *adminServer) Statistics(i *StatInterval, srv Admin_StatisticsServer) error {
	id, events := a.subs.Attach()
	defer a.subs.Dettach(id)

	stat := newStatCollector()

	ticker := time.NewTicker(time.Duration(i.IntervalSeconds) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case e, ok := <-events:
			if ok {
				stat.Update(e)
			} else {
				return nil
			}
		case <-ticker.C:
			err := srv.Send(stat.Collect())
			if err != nil {
				return err
			}
		}
	}
}

type Subscriber struct {
	id   int
	mu   *sync.RWMutex
	subs map[int]chan *Event
}

func newSubscriber() *Subscriber {
	return &Subscriber{
		mu:   &sync.RWMutex{},
		subs: make(map[int]chan *Event),
	}
}

func (s *Subscriber) Attach() (int, chan *Event) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.id++
	s.subs[s.id] = make(chan *Event)

	return s.id, s.subs[s.id]
}

func (s *Subscriber) Notify(e *Event) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, sub := range s.subs {
		sub <- e
	}
}

func (s *Subscriber) Dettach(id int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if sub, ok := s.subs[id]; ok {
		close(sub)
		delete(s.subs, id)
	}
}

func (s *Subscriber) DettachAll() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for id, sub := range s.subs {
		close(sub)
		delete(s.subs, id)
	}
}

type statCollector struct {
	mu   sync.RWMutex
	stat *Stat
}

func newStatCollector() *statCollector {
	s := &statCollector{}
	s.Reset()

	s.mu = sync.RWMutex{}

	return s
}

func (s *statCollector) Reset() {
	s.stat = &Stat{
		ByMethod:   make(map[string]uint64),
		ByConsumer: make(map[string]uint64),
	}
}

func (s *statCollector) Update(e *Event) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.stat.ByMethod[e.Method]++
	s.stat.ByConsumer[e.Consumer]++
}

func (s *statCollector) Collect() *Stat {
	s.mu.Lock()
	defer s.mu.Unlock()

	stat := s.stat
	stat.Timestamp = time.Now().Unix()
	s.Reset()

	return stat
}

type authACL struct {
	acl map[string][]string
}

func newAuth(aclData string) (*authACL, error) {
	aclParsed := make(map[string][]string)

	err := json.Unmarshal([]byte(aclData), &aclParsed)
	if err != nil {
		return nil, fmt.Errorf("failed to parce ACL data: %w", err)
	}

	auth := &authACL{
		acl: aclParsed,
	}

	return auth, nil
}

func (auth *authACL) Check(consumer, method string) error {
	methods := strings.Split(method, "/")

	m, ok := auth.acl[consumer]
	if !ok {
		return fmt.Errorf("failed authorization")
	}

nextMethod:
	for _, mthd := range m {
		for i, p := range strings.Split(mthd, "/") {
			if len(methods) > i && (p == methods[i] || p == "*") {
				continue
			}
			break nextMethod
		}
	}
	return nil
}

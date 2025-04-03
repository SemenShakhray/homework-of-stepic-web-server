package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type adminServer struct {
	UnimplementedAdminServer
	subs *subscriber
}

type subscriber struct {
	id   int
	mu   *sync.RWMutex
	subs map[int]chan *Event
}

type statCollector struct {
	mu   sync.RWMutex
	stat *Stat
}

type authACL struct {
	acl map[string][]string
}

type middleware struct {
	serverOptions []grpc.ServerOption
	auth          *authACL
	subs          *subscriber
}

func newAdminServer(subs *subscriber) *adminServer {
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

func newSubscriber() *subscriber {
	return &subscriber{
		mu:   &sync.RWMutex{},
		subs: make(map[int]chan *Event),
	}
}

func (s *subscriber) Attach() (int, chan *Event) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.id++
	s.subs[s.id] = make(chan *Event)

	return s.id, s.subs[s.id]
}

func (s *subscriber) Notify(e *Event) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, sub := range s.subs {
		sub <- e
	}
}

func (s *subscriber) Dettach(id int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if sub, ok := s.subs[id]; ok {
		close(sub)
		delete(s.subs, id)
	}
}

func (s *subscriber) DettachAll() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for id, sub := range s.subs {
		close(sub)
		delete(s.subs, id)
	}
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

func (m *middleware) Do(ctx context.Context, method string) error {
	var consumer, host string

	md, _ := metadata.FromIncomingContext(ctx)

	consumer = strings.Join(md.Get("consumer"), "")
	if p, ok := peer.FromContext(ctx); ok {
		host = p.Addr.String()
	}

	m.subs.Notify(&Event{
		Method:    method,
		Consumer:  consumer,
		Host:      host,
		Timestamp: time.Now().Unix(),
	})

	if err := m.auth.Check(consumer, method); err != nil {
		status.Errorf(codes.Unauthenticated, "failed authorization")
	}
	return nil
}

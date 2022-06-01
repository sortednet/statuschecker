package statuschecker

import (
	"context"
	"fmt"
	"github.com/sortednet/statuschecker/internal/store"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Status string

const (
	Down    Status = "down"
	Up      Status = "up"
	Unknown Status = "unknown"
)

type ServiceStatus struct {
	Name   string
	Status Status
}

type HttpClient interface {
	Get(url string) (*http.Response, error)
}

type StatusChecker struct {
	queries      DbQuery
	retrievers   map[string]*statusRetriever
	pollInterval time.Duration
	httpClient   HttpClient
	pollContext  context.Context
}

type DbQuery interface {
	GetServices(ctx context.Context) ([]store.Service, error)
	RegisterService(ctx context.Context, arg store.RegisterServiceParams) (store.Service, error)
	UnregisterService(ctx context.Context, name string) error
}

// Interface to allow mocking of the service
type StatusService interface {
	RegisterService(ctx context.Context, name, url string) error
	UnregisterService(ctx context.Context, name string) error
	GetAllServiceStatus(ctx context.Context) []ServiceStatus
	GetServiceStatus(ctx context.Context, name string) Status
}

func NewStatusChecker(pollCtx context.Context, queries DbQuery, pollInterval time.Duration, httpClient HttpClient) *StatusChecker {

	return &StatusChecker{
		queries:      queries,
		retrievers:   map[string]*statusRetriever{},
		pollInterval: pollInterval,
		httpClient:   httpClient,
		pollContext:  pollCtx,
	}
}

// Start a status retriever goroutine for each service
func (s *StatusChecker) StartPolling() error {

	services, err := s.queries.GetServices(s.pollContext)
	if err != nil {
		return err
	}

	for _, service := range services {
		s.startRetriever(service.Name, service.Url)
	}

	return nil
}

func (s *StatusChecker) RegisterService(ctx context.Context, name, url string) error {

	service, err := s.queries.RegisterService(ctx, store.RegisterServiceParams{Name: name, Url: url})

	if err != nil {
		return fmt.Errorf("Cannot add service %s to database %w", name, err)
	}

	s.startRetriever(name, url)

	zap.L().Info("Registered service ", zap.Any("service", service))
	return nil
}

func (s *StatusChecker) UnregisterService(ctx context.Context, name string) error {

	err := s.queries.UnregisterService(ctx, name)
	if err != nil {
		return fmt.Errorf("Cannot remove service %s from database %w", name, err)
	}
	s.stopRetriever(name)

	zap.L().Info("Unregistered service", zap.String("name", name))
	return nil
}

func (s *StatusChecker) GetAllServiceStatus(ctx context.Context) []ServiceStatus {
	allStatus := []ServiceStatus{}

	for name, retriever := range s.retrievers {
		allStatus = append(allStatus, ServiceStatus{
			Name:   name,
			Status: retriever.status,
		})
	}

	return allStatus
}

func (s *StatusChecker) GetServiceStatus(ctx context.Context, name string) Status {
	if retreiver, ok := s.retrievers[name]; ok {
		return retreiver.getStatus()
	}

	return Unknown
}

func (s *StatusChecker) startRetriever(name, url string) {
	s.stopRetriever(name) // just in case there is already one for this name already
	retriever := newStatusRetriever(s.httpClient, name, url)
	s.retrievers[name] = retriever
	go retriever.start(s.pollContext, s.pollInterval)
}

func (s *StatusChecker) stopRetriever(name string) {
	if retriever, ok := s.retrievers[name]; ok {
		close(retriever.quit)
		delete(s.retrievers, name)
	}
}

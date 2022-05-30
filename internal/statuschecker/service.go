package statuschecker

import (
	"context"
	"fmt"
	"github.com/sortednet/statuschecker/internal/store"
	"go.uber.org/zap"
	"net/http"
	"sync"
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
	queries    DbQuery
	cache      map[string]Status
	cacheLock  sync.Mutex
	ticker     *time.Ticker
	httpClient HttpClient
}

type DbQuery interface {
	GetServices(ctx context.Context) ([]store.Service, error)
	RegisterService(ctx context.Context, arg store.RegisterServiceParams) (store.Service, error)
	UnregisterService(ctx context.Context, name string) error
}

// Interface to allow mocking of the service
// NB This is really for the controller package but that causes import cycles with the mock
type StatusService interface {
	RegisterService(ctx context.Context, name, url string) error
	UnregisterService(ctx context.Context, name string) error
	GetAllServiceStatus(ctx context.Context) []ServiceStatus
	GetServiceStatus(ctx context.Context, name string) Status
}

func NewStatusChecker(queries DbQuery, pollInterval time.Duration, httpClient HttpClient) *StatusChecker {

	return &StatusChecker{
		queries:    queries,
		cache:      map[string]Status{},
		cacheLock:  sync.Mutex{},
		ticker:     time.NewTicker(pollInterval),
		httpClient: httpClient,
	}

}

func (s *StatusChecker) StartPolling(ctx context.Context) {
	err := s.pollServices(ctx) // initialise the cache now
	if err != nil {
		zap.L().Error("Error polling services for initial cache fill")
	}

	// update the status periodically
	go func() {
		for {
			select {
			case <-s.ticker.C:
				err := s.pollServices(ctx)
				if err != nil {
					zap.L().Error("Error polling services")
				}
			case <-ctx.Done():
				s.ticker.Stop()
				zap.L().Info("stopped polling")
				return
			}
		}
	}()
}

func (s *StatusChecker) pollServices(ctx context.Context) error {
	services, err := s.queries.GetServices(ctx)
	if err != nil {
		return err
	}

	for _, service := range services {
		go s.pollService(ctx, service)
	}

	return nil
}

func (s *StatusChecker) pollService(ctx context.Context, service store.Service) {

	zap.L().Info("Poll Service", zap.Any("service", service))
	resp, err := s.httpClient.Get(service.Url)

	if err != nil {
		zap.L().Error("Cannot get response from service",
			zap.String("name", service.Name),
			zap.String("url", service.Url),
			zap.Error(err))

		s.setServiceStatus(ctx, service.Name, Down) // cannot hit the service, assume it is down

	} else if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		zap.L().Error("Bad response from service",
			zap.Int("statusCode", resp.StatusCode),
			zap.String("name", service.Name),
			zap.String("url", service.Url),
			zap.Error(err))

		s.setServiceStatus(ctx, service.Name, Down)

	} else {
		s.setServiceStatus(ctx, service.Name, Up)
	}
}

func (s *StatusChecker) RegisterService(ctx context.Context, name, url string) error {

	_, err := s.queries.RegisterService(ctx, store.RegisterServiceParams{
		Name: name,
		Url:  url,
	},
	)
	if err != nil {
		return fmt.Errorf("Cannot add service %s to database %w", name, err)
	}

	s.setServiceStatus(ctx, name, Unknown)

	return nil
}

func (s *StatusChecker) UnregisterService(ctx context.Context, name string) error {

	err := s.queries.UnregisterService(ctx, name)
	if err != nil {
		return fmt.Errorf("Cannot remove service %s from database %w", name, err)
	}
	s.removeServiceStatus(ctx, name)
	return nil
}

func (s *StatusChecker) GetAllServiceStatus(ctx context.Context) []ServiceStatus {
	s.cacheLock.Lock()
	defer s.cacheLock.Unlock()

	allStatus := []ServiceStatus{}
	for name, status := range s.cache {
		allStatus = append(allStatus, ServiceStatus{
			Name:   name,
			Status: status,
		})
	}

	return allStatus
}

func (s *StatusChecker) GetServiceStatus(ctx context.Context, name string) Status {
	s.cacheLock.Lock()
	defer s.cacheLock.Unlock()

	return s.cache[name]
}

func (s *StatusChecker) setServiceStatus(ctx context.Context, name string, status Status) {
	s.cacheLock.Lock()
	defer s.cacheLock.Unlock()

	s.cache[name] = status
}

func (s *StatusChecker) removeServiceStatus(ctx context.Context, name string) {
	s.cacheLock.Lock()
	defer s.cacheLock.Unlock()

	delete(s.cache, name)
}

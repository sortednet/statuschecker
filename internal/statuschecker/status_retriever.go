package statuschecker

import (
	"context"
	"go.uber.org/zap"
	"time"
)

type statusRetriever struct {
	quit       chan struct{}
	name       string
	url        string
	status     Status
	httpClient HttpClient
	ctx        context.Context
}

func newStatusRetriever(httpClient HttpClient, name, url string) *statusRetriever {
	return &statusRetriever{
		quit:       make(chan struct{}),
		name:       name,
		url:        url,
		status:     Unknown,
		httpClient: httpClient,
	}
}

func (s *statusRetriever) start(ctx context.Context, pollInterval time.Duration) {
	ticker := time.NewTicker(pollInterval)

	for {
		select {
		case <-ticker.C:
			s.checkService(ctx)
		case <-s.quit:
			ticker.Stop()
			zap.L().Info("stopped polling - quit service", zap.String("service", s.name))
			return
		case <-ctx.Done():
			ticker.Stop()
			zap.L().Info("stopped polling - context is donw", zap.String("service", s.name))
			return
		}
	}
}

func (s *statusRetriever) getStatus() Status {
	return s.status
}

func (s *statusRetriever) checkService(ctx context.Context) {

	zap.L().Info("Check Service", zap.Any("service", s.name))
	resp, err := s.httpClient.Get(s.url)

	if err != nil {
		zap.L().Error("Cannot get response from service",
			zap.String("name", s.name),
			zap.String("url", s.url),
			zap.Error(err))

		s.status = Down

	} else if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		zap.L().Error("Bad response from service",
			zap.Int("statusCode", resp.StatusCode),
			zap.String("name", s.name),
			zap.String("url", s.url),
			zap.Error(err))

		s.status = Down

	} else {
		s.status = Up
	}

}

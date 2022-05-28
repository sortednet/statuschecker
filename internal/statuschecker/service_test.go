package statuschecker

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/sortednet/statuschecker/internal/store"
	"github.com/sortednet/statuschecker/test/mocks"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestStatusChecker_ServiceRegistration(t *testing.T) {

	ctx := context.TODO()

	googleParams := store.RegisterServiceParams{
		Name: "google",
		Url:  "http://google.com",
	}
	google := store.Service{Name: googleParams.Name, Url: googleParams.Url}

	mockCtrl := gomock.NewController(t)
	db := mocks.NewMockDbQuery(mockCtrl)
	db.EXPECT().RegisterService(gomock.Any(), googleParams).Return(google, nil)
	db.EXPECT().UnregisterService(gomock.Any(), googleParams.Name).Return(nil)
	httpClient := mocks.NewMockHttpClient(mockCtrl)

	checker, err := NewStatusChecker(ctx, db, time.Minute, httpClient)
	assert.NoError(t, err)
	assert.NotNil(t, checker)

	// Check registration
	status := checker.GetServiceStatus(ctx, "google")
	assert.Empty(t, status, "Unregistered service status is always unknown")

	checker.RegisterService(ctx, googleParams.Name, googleParams.Url)
	status = checker.GetServiceStatus(ctx, "google")
	assert.Equal(t, Unknown, status, "Unknown as poll will not have run (not empty as it would be if unregistered")

	// Check unregistration
	err = checker.UnregisterService(ctx, "google")
	assert.NoError(t, err)
	status = checker.GetServiceStatus(ctx, "google")
	assert.Empty(t, status, "Should be no status after the service has be unregistered")

}

func TestStatusChecker_pollService(t *testing.T) {
	ctx := context.TODO()

	google := store.Service{Name: "google", Url: "http://google.com"}
	down := store.Service{Name: "down", Url: "http://down.com"}

	mockCtrl := gomock.NewController(t)
	db := mocks.NewMockDbQuery(mockCtrl)

	httpClient := mocks.NewMockHttpClient(mockCtrl)
	httpClient.EXPECT().Get(google.Url).Return(&http.Response{StatusCode: 200}, nil)
	httpClient.EXPECT().Get(down.Url).Return(&http.Response{StatusCode: 404}, nil)

	checker, err := NewStatusChecker(ctx, db, time.Minute, httpClient)
	assert.NoError(t, err)
	assert.NotNil(t, checker)

	checker.pollService(ctx, google)
	status := checker.GetServiceStatus(ctx, "google")
	assert.Equal(t, Up, status)

	checker.pollService(ctx, down)
	status = checker.GetServiceStatus(ctx, "down")
	assert.Equal(t, Down, status)

}

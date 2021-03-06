package statuschecker

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/sortednet/statuschecker/internal/statuschecker/statuschecker_test"
	"github.com/sortednet/statuschecker/internal/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

func TestStatusChecker_ServiceRegistration(t *testing.T) {

	ctx := context.TODO()
	testSvcName1 := "testSvc"
	testSvcName2 := "testSvc2"

	reqParams1 := store.RegisterServiceParams{Name: testSvcName1, Url: "http://testservice.com"}
	testService1 := store.Service{Name: reqParams1.Name, Url: reqParams1.Url}
	reqParams2 := store.RegisterServiceParams{Name: testSvcName2, Url: "http://testservice2.com"}
	testService2 := store.Service{Name: reqParams2.Name, Url: reqParams2.Url}

	mockCtrl := gomock.NewController(t)
	db := statuschecker.NewMockDbQuery(mockCtrl)
	db.EXPECT().RegisterService(gomock.Any(), reqParams1).Return(testService1, nil)
	db.EXPECT().UnregisterService(gomock.Any(), reqParams1.Name).Return(nil)
	db.EXPECT().RegisterService(gomock.Any(), reqParams2).Return(testService2, nil)

	checker := NewStatusChecker(db, time.Minute, nil)
	require.NotNil(t, checker)

	// Check registration
	status := checker.GetServiceStatus(ctx, testSvcName1)
	assert.Empty(t, status, "Unregistered service status is always unknown")

	err := checker.RegisterService(ctx, reqParams1.Name, reqParams1.Url)
	require.NoError(t, err)
	status = checker.GetServiceStatus(ctx, testSvcName1)
	assert.Equal(t, Unknown, status, "Unknown as poll will not have run (not empty as it would be if unregistered")

	err = checker.RegisterService(ctx, reqParams2.Name, reqParams2.Url)
	require.NoError(t, err)
	status = checker.GetServiceStatus(ctx, testSvcName2)
	assert.Equal(t, Unknown, status, "Unknown as poll will not have run (not empty as it would be if unregistered")

	all := checker.GetAllServiceStatus(ctx)
	assert.Contains(t, all, ServiceStatus{Name: testSvcName1, Status: Unknown})
	assert.Contains(t, all, ServiceStatus{Name: testSvcName2, Status: Unknown})

	// Check unregistration
	err = checker.UnregisterService(ctx, testSvcName1)
	assert.NoError(t, err)
	status = checker.GetServiceStatus(ctx, testSvcName1)
	assert.Empty(t, status, "Should be no status after the service has be unregistered")

}

func TestStatusChecker_pollService(t *testing.T) {
	ctx := context.TODO()

	testSvcUp := store.Service{Name: "testSvcUp", Url: "http://testservice.com"}
	testSvcDown := store.Service{Name: "testSvcDown", Url: "http://down.com"}

	mockCtrl := gomock.NewController(t)
	db := statuschecker.NewMockDbQuery(mockCtrl)

	httpClient := statuschecker.NewMockHttpClient(mockCtrl)
	httpClient.EXPECT().Get(testSvcUp.Url).Return(&http.Response{StatusCode: 200}, nil)
	httpClient.EXPECT().Get(testSvcDown.Url).Return(nil, fmt.Errorf("Timeout"))

	checker := NewStatusChecker(db, time.Minute, httpClient)
	require.NotNil(t, checker)

	checker.pollService(ctx, testSvcUp)
	status := checker.GetServiceStatus(ctx, "testSvcUp")
	assert.Equal(t, Up, status)

	checker.pollService(ctx, testSvcDown)
	status = checker.GetServiceStatus(ctx, "testSvcDown")
	assert.Equal(t, Down, status)
}

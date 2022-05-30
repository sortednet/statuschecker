package web

import (
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/sortednet/statuschecker/internal/statuschecker"
	"github.com/sortednet/statuschecker/internal/web/web_test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStatusCheckerController_GetAllStatus(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	all := []statuschecker.ServiceStatus{
		statuschecker.ServiceStatus{Name: "s1", Status: statuschecker.Up},
		statuschecker.ServiceStatus{Name: "s2", Status: statuschecker.Down},
		statuschecker.ServiceStatus{Name: "s3", Status: statuschecker.Unknown},
	}

	svc := web.NewMockStatusService(mockCtrl)
	svc.EXPECT().GetAllServiceStatus(gomock.Any()).Return(all)

	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	recorder := httptest.NewRecorder()
	e := echo.New()
	webCtx := e.NewContext(req, recorder)

	controller := NewStatusCheckerController(svc)
	err := controller.GetAllStatus(webCtx)
	require.NoError(t, err)
	statusList := ServiceStatusList{}
	err = json.Unmarshal(recorder.Body.Bytes(), &statusList)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, recorder.Code)

	assert.Len(t, statusList, 3)
	assert.Contains(t, statusList, ServiceStatus{Name: "s1", Status: "up"})
	assert.Contains(t, statusList, ServiceStatus{Name: "s2", Status: "down"})
	assert.Contains(t, statusList, ServiceStatus{Name: "s3", Status: "unknown"})
}

func TestStatusCheckerController_GetHealth(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()
	e := echo.New()
	webCtx := e.NewContext(req, recorder)
	controller := NewStatusCheckerController(nil)
	err := controller.GetHealth(webCtx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestStatusCheckerController_Register_Happy(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	svc := web.NewMockStatusService(mockCtrl)
	svc.EXPECT().RegisterService(gomock.Any(), "s1", "http://s1.com").Return(nil)
	service := Service{Name: "s1", Url: "http://s1.com"}
	jsonBytes, err := json.Marshal(service)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/service", strings.NewReader(string(jsonBytes)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	e := echo.New()
	webCtx := e.NewContext(req, recorder)

	controller := NewStatusCheckerController(svc)
	err = controller.Register(webCtx)
	assert.NoError(t, err)
}

func TestStatusCheckerController_Register_BadRequest(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	svc := web.NewMockStatusService(mockCtrl)

	req := httptest.NewRequest(http.MethodPost, "/service", strings.NewReader(string("notjson")))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	e := echo.New()
	webCtx := e.NewContext(req, recorder)

	controller := NewStatusCheckerController(svc)
	err := controller.Register(webCtx)
	assert.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, httpErr.Code)
}

func TestStatusCheckerController_Register_InternalError(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	svc := web.NewMockStatusService(mockCtrl)
	svc.EXPECT().RegisterService(gomock.Any(), "s1", "http://s1.com").Return(fmt.Errorf("bang"))
	service := Service{Name: "s1", Url: "http://s1.com"}
	jsonBytes, err := json.Marshal(service)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/service", strings.NewReader(string(jsonBytes)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	e := echo.New()
	webCtx := e.NewContext(req, recorder)

	controller := NewStatusCheckerController(svc)
	err = controller.Register(webCtx)
	assert.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusInternalServerError, httpErr.Code)
}

func TestStatusCheckerController_Unregister_Happy(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	svc := web.NewMockStatusService(mockCtrl)
	svc.EXPECT().UnregisterService(gomock.Any(), "s1").Return(nil)
	service := Service{Name: "s1", Url: "http://s1.com"}
	jsonBytes, err := json.Marshal(service)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/service", strings.NewReader(string(jsonBytes)))
	//	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	e := echo.New()
	webCtx := e.NewContext(req, recorder)

	controller := NewStatusCheckerController(svc)
	err = controller.Unregister(webCtx, "s1")
	assert.NoError(t, err)
}

func TestStatusCheckerController_Unregister_InternalError(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	svc := web.NewMockStatusService(mockCtrl)
	svc.EXPECT().UnregisterService(gomock.Any(), "s1").Return(fmt.Errorf("bang"))
	service := Service{Name: "s1", Url: "http://s1.com"}
	jsonBytes, err := json.Marshal(service)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/service", strings.NewReader(string(jsonBytes)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	e := echo.New()
	webCtx := e.NewContext(req, recorder)

	controller := NewStatusCheckerController(svc)
	err = controller.Unregister(webCtx, "s1")
	assert.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusInternalServerError, httpErr.Code)
}

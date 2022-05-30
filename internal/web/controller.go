package web

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sortednet/statuschecker/internal/statuschecker"
	"go.uber.org/zap"
	"net/http"
)

//type StatusService interface {
//	RegisterService(ctx context.Context, name, url string) error
//	UnregisterService(ctx context.Context, name string) error
//	GetAllServiceStatus(ctx context.Context) []statuschecker.ServiceStatus
//	GetServiceStatus(ctx context.Context, name string) statuschecker.Status
//}

type StatusCheckerController struct {
	statusChecker statuschecker.StatusService
}

func NewStatusCheckerController(statusChecker statuschecker.StatusService) *StatusCheckerController {
	return &StatusCheckerController{
		statusChecker: statusChecker,
	}
}

// Add Service
// (POST /service)
func (s *StatusCheckerController) Register(webCtx echo.Context) error {
	var (
		err error
		ctx = webCtx.Request().Context()
	)

	log := zap.L()
	log.Info("Register service")
	service := Service{}
	err = webCtx.Bind(&service)
	if err != nil {
		log.Error("Error binding register service request", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "Cannot bind to request")
	}

	if service.Name == "" || service.Url == "" {
		log.Error("Invalid service request", zap.Any("service", service), zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "Missing field(s) in request")

	}

	err = s.statusChecker.RegisterService(ctx, service.Name, service.Url)
	if err != nil {
		log.Error("Error registering service", zap.String("name", service.Name), zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Cannot add service %s", service.Name))
	}
	log.Info("Registered service", zap.String("name", service.Name))

	return err
}

// Remove Service
// (DELETE /service/{name})
func (s *StatusCheckerController) Unregister(webCtx echo.Context, name string) error {
	var (
		err error
		ctx = webCtx.Request().Context()
	)

	log := zap.L()
	log.Info("Unregister service")

	err = s.statusChecker.UnregisterService(ctx, name)
	if err != nil {
		log.Error("Error unregistering service", zap.String("name", name), zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Cannot remove service %s", name))
	}
	log.Info("Unregistered service", zap.String("name", name))

	return err
}

// (GET /status)
func (s *StatusCheckerController) GetAllStatus(webCtx echo.Context) error {
	ctx := webCtx.Request().Context()

	serviceStatusList := ServiceStatusList{}
	allStatus := s.statusChecker.GetAllServiceStatus(ctx)

	for _, status := range allStatus {
		webStatus := ServiceStatus{
			Name:   status.Name,
			Status: string(status.Status),
		}
		serviceStatusList = append(serviceStatusList, webStatus)
	}

	webCtx.JSON(http.StatusOK, serviceStatusList)

	return nil
}

// (GET /health)
func (s *StatusCheckerController) GetHealth(webCtx echo.Context) error {
	return nil // no error = 200
}

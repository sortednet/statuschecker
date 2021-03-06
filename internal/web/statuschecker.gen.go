// Package web provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.11.0 DO NOT EDIT.
package web

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

// Service defines model for Service.
type Service struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

// ServiceStatus defines model for ServiceStatus.
type ServiceStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

// ServiceStatusList defines model for ServiceStatusList.
type ServiceStatusList = []ServiceStatus

// RegisterJSONBody defines parameters for Register.
type RegisterJSONBody = Service

// RegisterJSONRequestBody defines body for Register for application/json ContentType.
type RegisterJSONRequestBody = RegisterJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /health)
	GetHealth(ctx echo.Context) error
	// Add Service
	// (POST /service)
	Register(ctx echo.Context) error
	// Remove Service
	// (DELETE /service/{name})
	Unregister(ctx echo.Context, name string) error
	// get all service status
	// (GET /status)
	GetAllStatus(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetHealth converts echo context to params.
func (w *ServerInterfaceWrapper) GetHealth(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetHealth(ctx)
	return err
}

// Register converts echo context to params.
func (w *ServerInterfaceWrapper) Register(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.Register(ctx)
	return err
}

// Unregister converts echo context to params.
func (w *ServerInterfaceWrapper) Unregister(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "name" -------------
	var name string

	err = runtime.BindStyledParameterWithLocation("simple", false, "name", runtime.ParamLocationPath, ctx.Param("name"), &name)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter name: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.Unregister(ctx, name)
	return err
}

// GetAllStatus converts echo context to params.
func (w *ServerInterfaceWrapper) GetAllStatus(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetAllStatus(ctx)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/health", wrapper.GetHealth)
	router.POST(baseURL+"/service", wrapper.Register)
	router.DELETE(baseURL+"/service/:name", wrapper.Unregister)
	router.GET(baseURL+"/status", wrapper.GetAllStatus)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/7RUwY7TMBD9FWvgGDUBlktuu4AAidNWe1r1YJJp461jm/GkKKry78h20rRLF1oJTk2T",
	"5zfP783MHirbOmvQsIdyD75qsJXxcYm0UxWGR0fWIbHC+MHINr7l3iGU4JmU2cCQQUf6zPshA8IfnSKs",
	"oXxMpxN2lU1Y+/0JKw4cY9UlS+78FbX94cBF5Uf4XxV8U54DqWJsI/trwjWU8CqffctH0/JT7cOBWxLJ",
	"HoYgRZm1DTSVNSwrnq8EH2230dKLO0nckGzD+Rp9RcqxsgbKSZj40GC1RRJTQBmwYo1/ROyQfGIpFsXi",
	"TSC3Do10Ckp4tygWBWTgJDfxknmDUnMTHjcYRYYMZNDxtYYSPiN/SYhgrnfW+BTP26IIP6e6E7QPNd+f",
	"+/5gmgkRPcr9UefZ5P/pidu6FlKMMMFWbBGd4AZFylWsLUH2TPQ9bpRnJEgNgZ7vbN1PYaCJdaRzWlXx",
	"UP7kQ7FpJi4MP8U8dxxTh8MlLi27qkLv153WvXBkwx+sg2k35+B3shbjNV409hORpYlrGgbfta2kfnRx",
	"Vj37nu9DSw6JTyPj78z32NodXp/Bg6HjFE48ubnOk5v/4sl4sYMtYShItshIHsrHPahAEgYFsmlyx5Vy",
	"mnl21DbPN9Iqmn1YWC+N2K3W4yo53z//sm2Ptt2QpnD2ZIMspNaHrEflCYW0m6yJ6x8aZlfmeWe2xv40",
	"edgww2r4FQAA//+9cTDcaAYAAA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}

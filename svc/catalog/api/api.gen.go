// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.11.0 DO NOT EDIT.
package api

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

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse struct {
	// Error message
	Message string `json:"message"`
}

// ItemResponse defines model for ItemResponse.
type ItemResponse struct {
	// Description of the item
	Description *string `json:"description,omitempty"`

	// Unique ID of the item
	Id *uint `json:"id,omitempty"`

	// Name of the item
	Name *string `json:"name,omitempty"`

	// Price of the item
	Price *float64 `json:"price,omitempty"`

	// Currency of the price
	PriceCode *string `json:"priceCode,omitempty"`
}

// NewItemRequest defines model for NewItemRequest.
type NewItemRequest struct {
	// Description of the item
	Description string `json:"description"`

	// Name of the item
	Name string `json:"name"`

	// Price of the item
	Price float64 `json:"price"`

	// Currency of the price
	PriceCode string `json:"priceCode"`
}

// UpdateItemRequest defines model for UpdateItemRequest.
type UpdateItemRequest struct {
	// Description of the item
	Description *string `json:"description,omitempty"`

	// Name of the item
	Name *string `json:"name,omitempty"`

	// Price of the item
	Price *float64 `json:"price,omitempty"`

	// Currency of the price
	PriceCode *string `json:"priceCode,omitempty"`
}

// GetItemsParams defines parameters for GetItems.
type GetItemsParams struct {
	// Number of elements to be returned. Default 100
	PageSize *int `form:"pageSize,omitempty" json:"pageSize,omitempty"`

	// Page number.
	Page *int `form:"page,omitempty" json:"page,omitempty"`
}

// CreateItemJSONBody defines parameters for CreateItem.
type CreateItemJSONBody = NewItemRequest

// UpdateItemByIDJSONBody defines parameters for UpdateItemByID.
type UpdateItemByIDJSONBody = UpdateItemRequest

// CreateItemJSONRequestBody defines body for CreateItem for application/json ContentType.
type CreateItemJSONRequestBody = CreateItemJSONBody

// UpdateItemByIDJSONRequestBody defines body for UpdateItemByID for application/json ContentType.
type UpdateItemByIDJSONRequestBody = UpdateItemByIDJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Swagger documentation
	// (GET /api-docs)
	GetApiDocs(ctx echo.Context) error
	// Returns all items
	// (GET /api/v1/items)
	GetItems(ctx echo.Context, params GetItemsParams) error
	// Create new item
	// (POST /api/v1/items)
	CreateItem(ctx echo.Context) error
	// Removes an item by ID
	// (DELETE /api/v1/items/{id})
	DeleteItemByID(ctx echo.Context, id uint) error
	// Returns an item by ID
	// (GET /api/v1/items/{id})
	FindItemByID(ctx echo.Context, id uint) error
	// Updates an item by ID
	// (PUT /api/v1/items/{id})
	UpdateItemByID(ctx echo.Context, id uint) error
	// Health endpoint
	// (GET /healtz)
	GetHealtz(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetApiDocs converts echo context to params.
func (w *ServerInterfaceWrapper) GetApiDocs(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetApiDocs(ctx)
	return err
}

// GetItems converts echo context to params.
func (w *ServerInterfaceWrapper) GetItems(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetItemsParams
	// ------------- Optional query parameter "pageSize" -------------

	err = runtime.BindQueryParameter("form", true, false, "pageSize", ctx.QueryParams(), &params.PageSize)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter pageSize: %s", err))
	}

	// ------------- Optional query parameter "page" -------------

	err = runtime.BindQueryParameter("form", true, false, "page", ctx.QueryParams(), &params.Page)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter page: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetItems(ctx, params)
	return err
}

// CreateItem converts echo context to params.
func (w *ServerInterfaceWrapper) CreateItem(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.CreateItem(ctx)
	return err
}

// DeleteItemByID converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteItemByID(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id uint

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.DeleteItemByID(ctx, id)
	return err
}

// FindItemByID converts echo context to params.
func (w *ServerInterfaceWrapper) FindItemByID(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id uint

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.FindItemByID(ctx, id)
	return err
}

// UpdateItemByID converts echo context to params.
func (w *ServerInterfaceWrapper) UpdateItemByID(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id uint

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.UpdateItemByID(ctx, id)
	return err
}

// GetHealtz converts echo context to params.
func (w *ServerInterfaceWrapper) GetHealtz(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetHealtz(ctx)
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

	router.GET(baseURL+"/api-docs", wrapper.GetApiDocs)
	router.GET(baseURL+"/api/v1/items", wrapper.GetItems)
	router.POST(baseURL+"/api/v1/items", wrapper.CreateItem)
	router.DELETE(baseURL+"/api/v1/items/:id", wrapper.DeleteItemByID)
	router.GET(baseURL+"/api/v1/items/:id", wrapper.FindItemByID)
	router.PUT(baseURL+"/api/v1/items/:id", wrapper.UpdateItemByID)
	router.GET(baseURL+"/healtz", wrapper.GetHealtz)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+yXS4/jNgzHv4rA9ujm0aIX37qTPnJZFDPY02IOisU4KqzHUFIG3oG/eyHJedrZpIvZ",
	"xaDtzZElkub/R1J5gcooazRq76B8AVdtUPH0+CuRoXt01miHccGSsUheYnqt0DlepxcCXUXSemk0lPkc",
	"270uwLcWoQTnSeoauq4AwqcgCQWUH/dmHrsClh7VZYcnXs6dLg6/mFkzv0EmPaqh+wKkGB7/oOVTQLZc",
	"nB1eG1LcQwlBan8wJrXHGila01yN5OA9V3gtDkuyGjn6Z1y+FIYwYdUcJVUHtcpxJGt3RoxYvAtEqKt2",
	"ZzQ7HgrTFfAen7MITwGd/3oa/Iuydopz+rDixNQu5mNvkfYPVnCP/6f7H0IaMy712gxNPCBtY1iUW4hc",
	"NcjWsRVxzWupa1ZxzxtTp4BdtC19E43f5XUoYIvksrH5ZDaZxZiNRc2thBJ+SksFWO43SZ4pt/IHYar0",
	"o0Y/DOkefSDt2MMzr2skJkwVFGrPezKi2Ol5KaCE39H/YuUiGoxU5UaYjP84m41876jRroCf8+7KaI86",
	"RcWtbWSVNkz/cpmm3Orj0/eEayjhu+lhFkz7QTA9nQIp+2Ptng47CnBBKU7t5RC7IuVuup1PsxbX8sd1",
	"Fo2tyajER6/lZCyJy15fy4kr9EgOyo8D+hOVETdsMMbmmDdsFfGJPlFM2ALXPDSezWdRdhlPPQWkFnYV",
	"BZbX+CA/RVAP6RT5GJTz2awAJbVUQUE5H46PrhgUFq+R5YKZfMbnBX9XvD2OU3UzJ3utPgfMyRDv9lFw",
	"It6O8ZPEusTPXv+m6as2dhHjRlC5I+Qe96gwqS8zkvcuc9+i3H3fGdG+WtWczdHudEx4CtgNtJi/mvdT",
	"CcZTfpTxt9QwsjJM43MeK4NWMX2RosvqN+hxbCzGdZchOG8XbMUdCmZ0WlwumAvxQ1EMEMlmYqbetcvF",
	"tWaSL4479LxhfXR9BceJcShgKeCchuNyvnbnvFTGIxrnKMTbkvgeldkeFeqqZctFDPGWCXDDAPhNavHF",
	"uq3RV5tvKtt/veLPBd7jYMMIDvnSPNblby/tw8X7ixAJ6fhXZOT1x9Hwr8ZNE+kb85kT+8a61TlwPZ5x",
	"LG2QN/7T0d11cBH9I++46TLf/3mRjiW7m/YskGRrw1ALayI96UMc0naHbaAGSphC99j9HQAA//815Snf",
	"1hEAAA==",
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

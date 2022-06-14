package handler

import (
	"github.com/konrad945/eCommerce/svc/catalog/api"
	"github.com/labstack/echo/v4"
	"net/http"
)

var _ api.ServerInterface = (*handler)(nil)

type handler struct {
}

func NewHandler() *handler {
	return &handler{}
}

func (h *handler) GetApiDocs(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

func (h *handler) GetItems(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

func (h *handler) CreateItem(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

func (h *handler) DeleteItemByID(ctx echo.Context, _ int64) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

func (h *handler) FindItemByID(ctx echo.Context, _ int64) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

func (h *handler) UpdateItemByID(ctx echo.Context, _ int64) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

package handler

import (
	"fmt"
	"github.com/konrad945/eCommerce/svc/catalog/api"
	"github.com/konrad945/eCommerce/svc/catalog/internal/store"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
)

var _ api.ServerInterface = (*handler)(nil)

type CatalogStore interface {
	CreateItem(item store.Item) (uint, error)
}

type handler struct {
	log   *logrus.Logger
	store CatalogStore
}

func NewHandler(log *logrus.Logger, store CatalogStore) *handler {
	return &handler{store: store, log: log}
}

func (h *handler) GetApiDocs(ctx echo.Context) error {
	swagger, err := api.GetSwagger()
	if err != nil {
		return h.writeErrorResponse(ctx, http.StatusInternalServerError, fmt.Errorf("error while getting swagger documentation: %w", err))
	}
	json, err := swagger.MarshalJSON()
	if err != nil {
		return h.writeErrorResponse(ctx, http.StatusInternalServerError, fmt.Errorf("error while marshalling swagger documentation: %w", err))
	}
	return ctx.JSONBlob(http.StatusOK, json)
}

func (h *handler) GetItems(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

func (h *handler) CreateItem(ctx echo.Context) error {
	var item api.NewItem
	if err := ctx.Bind(&item); err != nil {
		return h.writeErrorResponse(ctx, http.StatusInternalServerError, fmt.Errorf("error while decoding request body: %w", err))
	}

	id, err := h.store.CreateItem(mapItemToItemModel(item))
	if err != nil {
		return h.writeErrorResponse(ctx, http.StatusInternalServerError, fmt.Errorf("error while creating new item: %w", err))
	}

	return ctx.JSON(http.StatusCreated, mapItemModelToItem(id, item))
}

func (h *handler) DeleteItemByID(ctx echo.Context, _ uint) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

func (h *handler) FindItemByID(ctx echo.Context, _ uint) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

func (h *handler) UpdateItemByID(ctx echo.Context, _ uint) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

func (h *handler) writeErrorResponse(ctx echo.Context, code int, err error) error {
	h.log.Errorf(err.Error())
	return ctx.JSON(code, api.ErrorResponse{Message: err.Error()})
}

func mapItemModelToItem(id uint, item api.NewItem) api.Item {
	return api.Item{
		Id:          id,
		Name:        item.Name,
		Description: item.Description,
		Price:       item.Price,
		PriceCode:   item.PriceCode,
	}
}

func mapItemToItemModel(newItem api.NewItem) store.Item {
	return store.Item{
		Name:        newItem.Name,
		Description: newItem.Description,
		Price:       newItem.Price,
		PriceCode:   newItem.PriceCode,
	}
}

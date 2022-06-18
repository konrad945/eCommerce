package handler

import (
	"errors"
	"fmt"
	"github.com/konrad945/eCommerce/svc/catalog/api"
	"github.com/konrad945/eCommerce/svc/catalog/internal/store"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
)

var _ api.ServerInterface = (*handler)(nil)

type CatalogStore interface {
	CreateItem(item store.Item) (store.Item, error)
	DeleteItem(id uint) error
	GetItem(id uint) (store.Item, error)
	GetItems(pageSize, page int) (users []store.Item, err error)
	UpdateItem(id uint, item store.Item) error
}

type handler struct {
	log   *logrus.Logger
	store CatalogStore
}

func NewHandler(log *logrus.Logger, store CatalogStore) *handler {
	return &handler{store: store, log: log}
}

// GetApiDocs returns openapi documentation
func (h *handler) GetApiDocs(ctx echo.Context) error {
	swagger, err := api.GetSwagger()
	if err != nil {
		return h.writeErrorResponse(ctx, fmt.Errorf("error while getting swagger documentation: %w", err))
	}
	json, err := swagger.MarshalJSON()
	if err != nil {
		return h.writeErrorResponse(ctx, fmt.Errorf("error while marshalling swagger documentation: %w", err))
	}
	return ctx.JSONBlob(http.StatusOK, json)
}

// GetItems returns items from underlying store
func (h *handler) GetItems(ctx echo.Context, params api.GetItemsParams) error {
	page := nvl(params.Page, 0)
	pageSize := nvl(params.PageSize, 100)

	items, err := h.store.GetItems(pageSize, page)
	if err != nil {
		return h.writeErrorResponse(ctx, fmt.Errorf("error while getting items: %w", err))
	}
	return ctx.JSON(http.StatusOK, items)
}

// CreateItem handles creation of a new item in the underlying store
func (h *handler) CreateItem(ctx echo.Context) error {
	var newItem api.NewItemRequest
	if err := ctx.Bind(&newItem); err != nil {
		return h.writeErrorResponse(ctx, fmt.Errorf("error while decoding request body: %w", err))
	}

	item, err := h.store.CreateItem(mapItemToItemModel(newItem))
	if err != nil {
		return h.writeErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusCreated, item)
}

// DeleteItemByID deletes item by ID from the underlying store
func (h *handler) DeleteItemByID(ctx echo.Context, id uint) error {
	if err := h.store.DeleteItem(id); err != nil {
		return h.writeErrorResponse(ctx, err)
	}
	return ctx.NoContent(http.StatusOK)
}

// FindItemByID returns item with ID from the underlying store
func (h *handler) FindItemByID(ctx echo.Context, id uint) error {
	item, err := h.store.GetItem(id)
	if err != nil {
		return h.writeErrorResponse(ctx, err)
	}
	return ctx.JSON(http.StatusOK, item)
}

// UpdateItemByID updates item with ID in the underlying store
func (h *handler) UpdateItemByID(ctx echo.Context, id uint) error {
	var item api.UpdateItemRequest
	if err := ctx.Bind(&item); err != nil {
		return h.writeErrorResponse(ctx, fmt.Errorf("error while decoding request body: %w", err))
	}

	if err := h.store.UpdateItem(id, store.Item{
		Name:        item.Name,
		Description: item.Description,
		Price:       item.Price,
		PriceCode:   item.PriceCode,
	}); err != nil {
		return h.writeErrorResponse(ctx, err)
	}

	return ctx.NoContent(http.StatusOK)
}

func (h *handler) writeErrorResponse(ctx echo.Context, err error) error {
	h.log.Errorf(err.Error())
	code := http.StatusInternalServerError
	if errors.Is(err, gorm.ErrRecordNotFound) {
		code = http.StatusNotFound
	}
	return ctx.JSON(code, api.ErrorResponse{Message: err.Error()})
}

func mapItemToItemModel(newItem api.NewItemRequest) store.Item {
	return store.Item{
		Name:        &newItem.Name,
		Description: &newItem.Description,
		Price:       &newItem.Price,
		PriceCode:   &newItem.PriceCode,
	}
}

// nvl returns a if a is not nil or else return b
func nvl(a *int, b int) int {
	if a == nil {
		return b
	}
	return *a
}

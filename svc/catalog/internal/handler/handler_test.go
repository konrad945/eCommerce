package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/konrad945/eCommerce/svc/catalog/api"
	"github.com/konrad945/eCommerce/svc/catalog/internal/store"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetApiDocs(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api-docs", nil)
	rec := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, rec)

	err := NewHandler(logrus.New(), &mockCatalogStore{}).GetApiDocs(ctx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, oApiBytes(t), rec.Body.Bytes())
}

func TestGetItems(t *testing.T) {
	tests := []struct {
		name             string
		queryParams      api.GetItemsParams
		storeResponse    []store.Item
		err              error
		expectedPage     int
		expectedPageSize int
		expectedStatus   int
	}{
		{
			name:             "Successful - default values taken for query params",
			queryParams:      api.GetItemsParams{},
			storeResponse:    []store.Item{{ID: 0}, {ID: 1}},
			expectedPage:     1,
			expectedPageSize: 100,
			expectedStatus:   http.StatusOK,
		},
		{
			name: "Successful - override default values",
			queryParams: api.GetItemsParams{
				PageSize: v2p(2),
				Page:     v2p(2),
			},
			storeResponse:    []store.Item{{ID: 1}, {ID: 2}},
			expectedPage:     2,
			expectedPageSize: 2,
			expectedStatus:   http.StatusOK,
		},
		{
			name:             "Internal Error - error from store",
			queryParams:      api.GetItemsParams{},
			err:              errors.New("some error"),
			expectedPage:     1,
			expectedPageSize: 100,
			expectedStatus:   http.StatusInternalServerError,
		},
		{
			name:             "Not Found - no record in db",
			queryParams:      api.GetItemsParams{},
			err:              gorm.ErrRecordNotFound,
			expectedPage:     1,
			expectedPageSize: 100,
			expectedStatus:   http.StatusNotFound,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockStore := &mockCatalogStore{
				t:                t,
				expectedPage:     test.expectedPage,
				expectedPageSize: test.expectedPageSize,
				getItemsResponse: test.storeResponse,
				err:              test.err}

			req := httptest.NewRequest(http.MethodGet, "/api/v1/items", nil)
			rec := httptest.NewRecorder()
			ctx := echo.New().NewContext(req, rec)

			err := NewHandler(logrus.New(), mockStore).GetItems(ctx, test.queryParams)
			require.NoError(t, err)

			assert.Equal(t, rec.Code, test.expectedStatus)
			if test.expectedStatus != http.StatusOK {
				var respMsg api.ErrorResponse
				err = json.NewDecoder(rec.Body).Decode(&respMsg)
				require.NoError(t, err)
				assert.Contains(t, respMsg.Message, test.err.Error())
			}

		})
	}
}

func TestFindItemByID(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedStatus int
	}{
		{
			name:           "Successful",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Unsuccessful - no record in db",
			err:            gorm.ErrRecordNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Unsuccessful - store internal error",
			err:            errors.New("some error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockStore := &mockCatalogStore{t: t, expectedID: 1, err: test.err}
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/items/1"), nil)
			rec := httptest.NewRecorder()
			ctx := echo.New().NewContext(req, rec)

			err := NewHandler(logrus.New(), mockStore).FindItemByID(ctx, 1)
			require.NoError(t, err)

			assert.Equal(t, test.expectedStatus, rec.Code)
			if test.expectedStatus != http.StatusOK {
				var respMsg api.ErrorResponse
				err = json.NewDecoder(rec.Body).Decode(&respMsg)
				require.NoError(t, err)
				assert.Contains(t, respMsg.Message, test.err.Error())
			}
		})
	}
}

func TestDeleteItem(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedStatus int
	}{
		{
			name:           "Successful",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Unsuccessful - no record in db",
			err:            gorm.ErrRecordNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Unsuccessful - store internal error",
			err:            errors.New("some error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockStore := &mockCatalogStore{t: t, expectedID: 1, err: test.err}
			req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/items/1"), nil)
			rec := httptest.NewRecorder()
			ctx := echo.New().NewContext(req, rec)

			err := NewHandler(logrus.New(), mockStore).DeleteItemByID(ctx, 1)
			require.NoError(t, err)

			assert.Equal(t, test.expectedStatus, rec.Code)
			if test.expectedStatus != http.StatusOK {
				var respMsg api.ErrorResponse
				err = json.NewDecoder(rec.Body).Decode(&respMsg)
				require.NoError(t, err)
				assert.Contains(t, respMsg.Message, test.err.Error())
			}
		})
	}
}

func TestCreateItem(t *testing.T) {
	tests := []struct {
		name           string
		newItem        api.NewItemRequest
		err            error
		expectedStatus int
	}{
		{
			name: "Successful",
			newItem: api.NewItemRequest{
				Description: "someDesc",
				Name:        "someName",
				Price:       50,
				PriceCode:   "EUR",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Unsuccessful - store internal error",
			err:            errors.New("some error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockStore := &mockCatalogStore{t: t, err: test.err, expectedItem: store.Item{
				Name:        &test.newItem.Name,
				Description: &test.newItem.Description,
				Price:       &test.newItem.Price,
				PriceCode:   &test.newItem.PriceCode,
			}}
			body, err := json.Marshal(test.newItem)
			require.NoError(t, err)
			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/items"), bytes.NewReader(body))
			req.Header.Add("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			ctx := echo.New().NewContext(req, rec)

			err = NewHandler(logrus.New(), mockStore).CreateItem(ctx)
			require.NoError(t, err)

			assert.Equal(t, test.expectedStatus, rec.Code)
			if test.expectedStatus != http.StatusCreated {
				var respMsg api.ErrorResponse
				err = json.NewDecoder(rec.Body).Decode(&respMsg)
				require.NoError(t, err)
				assert.Contains(t, respMsg.Message, test.err.Error())
			}
		})
	}
}

func TestUpdateItem(t *testing.T) {
	tests := []struct {
		name           string
		updateItemReq  api.UpdateItemRequest
		err            error
		expectedItem   store.Item
		expectedStatus int
	}{
		{
			name:           "Successful",
			updateItemReq:  api.UpdateItemRequest{Price: v2p(float64(50)), PriceCode: v2p("EUR")},
			expectedItem:   store.Item{Price: v2p(float64(50)), PriceCode: v2p("EUR")},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Unsuccessful - no record in db",
			err:            gorm.ErrRecordNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Unsuccessful - store internal error",
			err:            errors.New("some error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockStore := &mockCatalogStore{t: t, err: test.err, expectedItem: test.expectedItem, expectedID: 1}
			body, err := json.Marshal(test.updateItemReq)
			require.NoError(t, err)
			req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/items/1"), bytes.NewReader(body))
			req.Header.Add("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			ctx := echo.New().NewContext(req, rec)

			err = NewHandler(logrus.New(), mockStore).UpdateItemByID(ctx, 1)
			require.NoError(t, err)

			assert.Equal(t, test.expectedStatus, rec.Code)
			if test.expectedStatus != http.StatusOK {
				var respMsg api.ErrorResponse
				err = json.NewDecoder(rec.Body).Decode(&respMsg)
				require.NoError(t, err)
				assert.Contains(t, respMsg.Message, test.err.Error())
			}
		})
	}
}

type mockCatalogStore struct {
	t                *testing.T
	expectedID       uint
	expectedPage     int
	expectedPageSize int
	expectedItem     store.Item
	err              error
	getItemsResponse []store.Item
	findItemResponse store.Item
}

func (m *mockCatalogStore) CreateItem(item store.Item) (store.Item, error) {
	assert.Equal(m.t, m.expectedItem, item)
	item.ID = 1
	return item, m.err
}

func (m *mockCatalogStore) DeleteItem(id uint) error {
	assert.Equal(m.t, m.expectedID, id)
	return m.err
}

func (m *mockCatalogStore) GetItem(id uint) (store.Item, error) {
	assert.Equal(m.t, m.expectedID, id)
	return store.Item{ID: 1}, m.err
}

func (m *mockCatalogStore) GetItems(pageSize, page int) ([]store.Item, error) {
	assert.Equal(m.t, m.expectedPageSize, pageSize)
	assert.Equal(m.t, m.expectedPage, page)
	return m.getItemsResponse, m.err
}

func (m *mockCatalogStore) UpdateItem(id uint, item store.Item) error {
	assert.Equal(m.t, m.expectedID, id)
	assert.Equal(m.t, m.expectedItem, item)
	return m.err
}

func v2p[V int | float64 | string](val V) *V {
	return &val
}

func oApiBytes(t *testing.T) []byte {
	oApi, err := api.GetSwagger()
	require.NoError(t, err)
	expectedResp, err := oApi.MarshalJSON()
	require.NoError(t, err)
	return expectedResp
}

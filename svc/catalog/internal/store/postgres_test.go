package store

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

const (
	getItem    = `SELECT \* FROM "items" WHERE "items"\."id" = \$1 ORDER BY "items"\."id" LIMIT 1`
	deleteItem = `DELETE FROM "items" WHERE "items"\."id" = \$1`
	createItem = `INSERT INTO "items"`
	updateItem = `UPDATE "items"`
)

var (
	itemID        = uint(1)
	itemName      = "some name"
	itemDesc      = "some desc"
	itemPrice     = float64(50)
	itemPriceCode = "EUR"
)

func TestGetItem(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}))
	require.NoError(t, err)
	store := &CatalogStore{db: db}

	tests := []struct {
		name        string
		rows        *sqlmock.Rows
		expectedRes Item
		expectedErr error
	}{
		{
			name: "Successful",
			rows: sqlmock.
				NewRows([]string{"id", "name", "description", "price", "price_code"}).
				AddRow(itemID, itemName, itemDesc, itemPrice, itemPriceCode),
			expectedRes: Item{
				ID:          itemID,
				Name:        &itemName,
				Description: &itemDesc,
				Price:       &itemPrice,
				PriceCode:   &itemPriceCode,
			},
		},
		{
			name: "Unsuccessful - error",
			rows: sqlmock.
				NewRows([]string{"id", "name", "description", "price", "price_code"}).
				AddRow(itemID, itemName, itemDesc, itemPrice, itemPriceCode),
			expectedErr: errors.New("some err"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.expectedErr != nil {
				mock.ExpectQuery(getItem).WithArgs(itemID).WillReturnError(test.expectedErr)
			} else {
				mock.ExpectQuery(getItem).WithArgs(itemID).WillReturnRows(test.rows)
			}

			item, err := store.GetItem(1)

			if test.expectedErr != nil {
				assert.ErrorContains(t, err, "some err")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedRes, item)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetItems(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}))
	require.NoError(t, err)
	store := &CatalogStore{db: db}

	tests := []struct {
		name        string
		rows        *sqlmock.Rows
		page        int
		pageSize    int
		expectedErr error
	}{
		{
			name: "Successful",
			rows: sqlmock.
				NewRows([]string{"id", "name", "description", "price", "price_code"}).
				AddRow(itemID, itemName, itemDesc, itemPrice, itemPriceCode),
			page:     1,
			pageSize: 10,
		},
		{
			name: "Successful",
			rows: sqlmock.
				NewRows([]string{"id", "name", "description", "price", "price_code"}).
				AddRow(itemID, itemName, itemDesc, itemPrice, itemPriceCode),
			page:     2,
			pageSize: 10,
		},
		{
			name: "Unsuccessful - error",
			rows: sqlmock.
				NewRows([]string{"id", "name", "description", "price", "price_code"}).
				AddRow(itemID, itemName, itemDesc, itemPrice, itemPriceCode),
			page:        1,
			pageSize:    10,
			expectedErr: errors.New("some err"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			query := fmt.Sprintf(`SELECT \* FROM "items" LIMIT %d`, test.pageSize)
			if test.page > 1 {
				query = query + fmt.Sprintf(" OFFSET %d", (test.page-1)*test.pageSize)
			}
			if test.expectedErr != nil {
				mock.ExpectQuery(query).WillReturnError(test.expectedErr)
			} else {
				mock.ExpectQuery(query).WillReturnRows(test.rows)
			}

			items, err := store.GetItems(test.pageSize, test.page)

			if test.expectedErr != nil {
				assert.ErrorContains(t, err, "some err")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, 1, len(items))
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetItems_inputValidation(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		pageSize int
	}{
		{
			name:     "Unsuccessful - page invalid",
			page:     0,
			pageSize: 10,
		},
		{
			name:     "Unsuccessful - pageSize invalid",
			page:     2,
			pageSize: 0,
		},
		{
			name:     "Unsuccessful - both params invalid",
			page:     0,
			pageSize: 0,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := (&CatalogStore{}).GetItems(test.pageSize, test.page)
			assert.ErrorIs(t, err, ErrInvalidPageParams)
		})
	}
}

func TestDeleteItem(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}))
	require.NoError(t, err)
	store := &CatalogStore{db: db}

	tests := []struct {
		name        string
		expectedErr error
	}{
		{
			name: "Successful",
		},
		{
			name:        "Unsuccessful - error",
			expectedErr: errors.New("some err"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock.ExpectBegin()
			if test.expectedErr != nil {
				mock.ExpectExec(deleteItem).WithArgs(1).WillReturnError(test.expectedErr)
				mock.ExpectRollback()
			} else {
				mock.ExpectExec(deleteItem).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			}

			err := store.DeleteItem(1)

			if test.expectedErr != nil {
				assert.ErrorContains(t, err, "some err")
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCreateItem(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}))
	require.NoError(t, err)
	store := &CatalogStore{db: db}

	tests := []struct {
		name        string
		rows        *sqlmock.Rows
		expectedRes Item
		expectedErr error
	}{
		{
			name: "Successful",
			rows: sqlmock.
				NewRows([]string{"id", "name", "description", "price", "price_code"}).
				AddRow(itemID, itemName, itemDesc, itemPrice, itemPriceCode),
			expectedRes: Item{
				ID:          itemID,
				Name:        &itemName,
				Description: &itemDesc,
				Price:       &itemPrice,
				PriceCode:   &itemPriceCode,
			},
		},
		{
			name: "Unsuccessful - error",
			rows: sqlmock.
				NewRows([]string{"id", "name", "description", "price", "price_code"}).
				AddRow(itemID, itemName, itemDesc, itemPrice, itemPriceCode),
			expectedErr: errors.New("some err"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock.ExpectBegin()
			if test.expectedErr != nil {
				mock.ExpectQuery(createItem).WithArgs(itemName, itemDesc, itemPrice, itemPriceCode).WillReturnError(test.expectedErr)
				mock.ExpectRollback()
			} else {
				mock.ExpectQuery(createItem).WithArgs(itemName, itemDesc, itemPrice, itemPriceCode).WillReturnRows(test.rows)
				mock.ExpectCommit()
			}

			item, err := store.CreateItem(Item{
				Name:        &itemName,
				Description: &itemDesc,
				Price:       &itemPrice,
				PriceCode:   &itemPriceCode,
			})

			if test.expectedErr != nil {
				assert.ErrorContains(t, err, "some err")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedRes, item)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUpdateItem(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}))
	require.NoError(t, err)
	store := &CatalogStore{db: db}

	tests := []struct {
		name      string
		result    driver.Result
		mockErr   error
		expectErr error
	}{
		{
			name:   "Successful",
			result: sqlmock.NewResult(1, 1),
		},
		{
			name:      "Unsuccessful - error",
			mockErr:   errors.New("some err"),
			expectErr: errors.New("some err"),
		},
		{
			name:      "Unsuccessful - no rows affected",
			result:    sqlmock.NewResult(1, 0),
			expectErr: gorm.ErrRecordNotFound,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock.ExpectBegin()
			if test.mockErr != nil {
				mock.ExpectExec(updateItem).WithArgs(itemName, itemDesc, itemPrice, itemPriceCode, itemID).WillReturnError(test.mockErr)
				mock.ExpectRollback()
			} else {
				mock.ExpectExec(updateItem).WithArgs(itemName, itemDesc, itemPrice, itemPriceCode, itemID).WillReturnResult(test.result)
				mock.ExpectCommit()
			}

			err := store.UpdateItem(itemID, Item{
				Name:        &itemName,
				Description: &itemDesc,
				Price:       &itemPrice,
				PriceCode:   &itemPriceCode,
			})

			if test.expectErr != nil {
				assert.ErrorContains(t, err, test.expectErr.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

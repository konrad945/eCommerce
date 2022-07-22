package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var ErrInvalidPageParams = errors.New("page and pageSize parameters should be greater than or equal to 1")

type CatalogStore struct {
	db *gorm.DB
}

// dbCredentials represents vault credential file struct
type dbCredentials = struct {
	DBConnection string `json:"db_connection"`
}

// NewCatalogStore initialize connection to underlying db. Connection parameters should be passed by env variables.
func NewCatalogStore(dbCredPath string) (*CatalogStore, error) {
	file, err := os.Open(dbCredPath)
	if err != nil {
		return nil, fmt.Errorf("error while opening file at %s: %w", dbCredPath, err)
	}
	var credentials = &dbCredentials{}
	if err := json.NewDecoder(file).Decode(credentials); err != nil {
		return nil, fmt.Errorf("error while decoding json: %w", err)
	}

	db, err := gorm.Open(postgres.Open(credentials.DBConnection), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error shile opening connection to db: %w", err)
	}

	return &CatalogStore{db: db}, nil
}

// CreateItem persists Item in db and returns its ID
func (s *CatalogStore) CreateItem(item Item) (Item, error) {
	if err := s.db.Create(&item).Error; err != nil {
		return Item{}, fmt.Errorf("error while adding item to db: %w", err)
	}
	return item, nil
}

// DeleteItem deletes item with ID from db
func (s *CatalogStore) DeleteItem(id uint) error {
	if err := s.db.Delete(&Item{}, id).Error; err != nil {
		return fmt.Errorf("error while deleting item with id %d: %w", id, err)
	}
	return nil
}

// GetItem returns item with provided ID from db
func (s *CatalogStore) GetItem(id uint) (Item, error) {
	var item Item
	if err := s.db.First(&item, id).Error; err != nil {
		return Item{}, fmt.Errorf("error while getting item with id %d: %w", id, err)
	}
	return item, nil
}

func (s *CatalogStore) GetItems(pageSize, page int) (items []Item, err error) {
	if page < 1 || pageSize < 1 {
		return nil, ErrInvalidPageParams
	}
	err = s.db.Model(&Item{}).Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error
	return
}

// UpdateItem updates item with ID in db.
func (s *CatalogStore) UpdateItem(id uint, item Item) error {
	resp := s.db.Model(Item{}).Where("id = ?", id).Updates(&item)
	if err := resp.Error; err != nil {
		return fmt.Errorf("error while updating item with id %d: %w", id, err)
	}
	if resp.RowsAffected != 1 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

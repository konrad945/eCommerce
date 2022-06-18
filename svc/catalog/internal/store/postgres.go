package store

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type config struct {
	Username string `envconfig:"DB_USERNAME" default:"catalog_user"`
	Password string `envconfig:"DB_PASSWORD" default:"password"`
	Name     string `envconfig:"DB_NAME" default:"catalog"`
	Port     int    `envconfig:"DB_PORT" default:"5432"`
	Host     string `envconfig:"DB_HOST" default:"localhost"`
}

type CatalogStore struct {
	db *gorm.DB
}

// NewCatalogStore initialize connection to underlying db. Connection parameters should be passed by env variables.
func NewCatalogStore() (*CatalogStore, error) {
	var conf config
	if err := envconfig.Process("", &conf); err != nil {
		return nil, fmt.Errorf("error while processing env variables: %w", err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", conf.Host, conf.Username, conf.Password, conf.Name, conf.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
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
	err = s.db.Model(&Item{}).Offset(page - 1).Limit(pageSize).Find(&items).Error
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

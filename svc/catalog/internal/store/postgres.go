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

func (s *CatalogStore) CreateItem(item Item) (uint, error) {
	if res := s.db.Create(&item); res.Error != nil {
		return 0, fmt.Errorf("error while adding item to db: %w", res.Error)
	}
	return item.ID, nil
}

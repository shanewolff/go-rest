package db

import (
	"time"

	"gorm.io/gorm"

	"go-rest/internal/domain"
)

// This package acts as the Outbound Adapter for PostgreSQL.

// Item is the ORM model specifically tailored for GORM, separated from the pure domain model.
type Item struct {
	ID        uint    `gorm:"primaryKey"`
	Title     string  `gorm:"not null"`
	Price     float64 `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Map from pure domain item to DB item
func toDBModel(item *domain.Item) *Item {
	return &Item{
		ID:        item.ID,
		Title:     item.Title,
		Price:     item.Price,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}

// Map from DB item to pure domain item
func toDomainModel(dbItem *Item) domain.Item {
	return domain.Item{
		ID:        dbItem.ID,
		Title:     dbItem.Title,
		Price:     dbItem.Price,
		CreatedAt: dbItem.CreatedAt,
		UpdatedAt: dbItem.UpdatedAt,
	}
}

// postgresRepository implements `domain.ItemRepository`
type postgresRepository struct {
	db *gorm.DB
}

// NewPostgresRepository returns a new instance of the repository.
// It accepts a gorm.DB handle, following dependency injection principles.
func NewPostgresRepository(db *gorm.DB) domain.ItemRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) GetAll() ([]domain.Item, error) {
	var dbItems []Item
	result := r.db.Find(&dbItems)
	if result.Error != nil {
		return nil, result.Error
	}

	// Convert DB items to Domain items before returning to the core logic
	var domainItems []domain.Item
	for _, dbItem := range dbItems {
		domainItems = append(domainItems, toDomainModel(&dbItem))
	}

	return domainItems, nil
}

func (r *postgresRepository) GetByID(id uint) (*domain.Item, error) {
	var dbItem Item
	result := r.db.First(&dbItem, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return new(toDomainModel(&dbItem)), nil
}

func (r *postgresRepository) Create(item *domain.Item) error {
	dbItem := toDBModel(item)

	// GORM will populate the ID and CreatedAt fields after creation
	result := r.db.Create(dbItem)
	if result.Error != nil {
		return result.Error
	}

	// Update the original domain item with the DB-assigned ID
	item.ID = dbItem.ID
	item.CreatedAt = dbItem.CreatedAt
	item.UpdatedAt = dbItem.UpdatedAt
	return nil
}

func (r *postgresRepository) Delete(id uint) error {
	result := r.db.Delete(&Item{}, id)
	return result.Error
}

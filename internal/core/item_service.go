package core

import (
	"errors"
	"time"

	"github.com/shanewolff/go-rest/internal/domain"
)

// The Core Business Logic (Application Layer).

// This implements the `domain.ItemService` interface (Inbound Port).
type itemService struct {
	repo domain.ItemRepository
}

// NewItemService creates a new instance of the core item service.
func NewItemService(r domain.ItemRepository) domain.ItemService {
	return &itemService{
		repo: r,
	}
}

func (s *itemService) GetAllItems() ([]domain.Item, error) {
	// Add business logic here if needed before returning items from DB
	return s.repo.GetAll()
}

func (s *itemService) GetItem(id uint) (*domain.Item, error) {
	if id == 0 {
		return nil, errors.New("invalid item ID")
	}
	return s.repo.GetByID(id)
}

func (s *itemService) CreateItem(req domain.CreateItemRequest) (*domain.Item, error) {
	// Create the core domain model from the request model
	newItem := &domain.Item{
		Title:     req.Title,
		Price:     req.Price,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := s.repo.Create(newItem)
	if err != nil {
		return nil, err
	}

	return newItem, nil
}

func (s *itemService) DeleteItem(id uint) error {
	if id == 0 {
		return errors.New("invalid item ID")
	}
	return s.repo.Delete(id)
}

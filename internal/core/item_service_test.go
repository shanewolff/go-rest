package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-rest/internal/domain"
)

// MockItemRepository is a mock implementation of domain.ItemRepository
type MockItemRepository struct {
	mock.Mock
}

func (m *MockItemRepository) GetAll() ([]domain.Item, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Item), args.Error(1)
}

func (m *MockItemRepository) GetByID(id uint) (*domain.Item, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Item), args.Error(1)
}

func (m *MockItemRepository) Create(item *domain.Item) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MockItemRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestGetAllItems(t *testing.T) {
	mockRepo := new(MockItemRepository)
	service := NewItemService(mockRepo)

	expectedItems := []domain.Item{
		{ID: 1, Title: "Item 1", Price: 10.0},
		{ID: 2, Title: "Item 2", Price: 20.0},
	}

	mockRepo.On("GetAll").Return(expectedItems, nil)

	items, err := service.GetAllItems()

	assert.NoError(t, err)
	assert.Equal(t, 2, len(items))
	assert.Equal(t, "Item 1", items[0].Title)
	mockRepo.AssertExpectations(t)
}

func TestGetItem(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockItemRepository)
		service := NewItemService(mockRepo)

		expectedItem := &domain.Item{ID: 1, Title: "Item 1", Price: 10.0}
		mockRepo.On("GetByID", uint(1)).Return(expectedItem, nil)

		item, err := service.GetItem(1)

		assert.NoError(t, err)
		assert.NotNil(t, item)
		assert.Equal(t, uint(1), item.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid id", func(t *testing.T) {
		mockRepo := new(MockItemRepository)
		service := NewItemService(mockRepo)

		item, err := service.GetItem(0)

		assert.Error(t, err)
		assert.Nil(t, item)
		assert.Equal(t, "invalid item ID", err.Error())
	})
}

func TestCreateItem(t *testing.T) {
	mockRepo := new(MockItemRepository)
	service := NewItemService(mockRepo)

	req := domain.CreateItemRequest{
		Title: "New Item",
		Price: 15.0,
	}

	// We use Matcher because the item object is created inside the service and includes timestamps
	mockRepo.On("Create", mock.MatchedBy(func(item *domain.Item) bool {
		return item.Title == "New Item" && item.Price == 15.0
	})).Return(nil)

	item, err := service.CreateItem(req)

	assert.NoError(t, err)
	assert.NotNil(t, item)
	assert.Equal(t, "New Item", item.Title)
	mockRepo.AssertExpectations(t)
}

func TestDeleteItem(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockItemRepository)
		service := NewItemService(mockRepo)

		mockRepo.On("Delete", uint(1)).Return(nil)

		err := service.DeleteItem(1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid id", func(t *testing.T) {
		mockRepo := new(MockItemRepository)
		service := NewItemService(mockRepo)

		err := service.DeleteItem(0)

		assert.Error(t, err)
		assert.Equal(t, "invalid item ID", err.Error())
	})
}

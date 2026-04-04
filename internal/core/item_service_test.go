package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-rest/internal/domain"
	"go-rest/internal/mocks"
)

func TestGetAllItems(t *testing.T) {
	mockRepo := mocks.NewItemRepository(t)
	service := NewItemService(mockRepo)

	expectedItems := []domain.Item{
		{ID: 1, Title: "Item 1", Price: 10.0},
		{ID: 2, Title: "Item 2", Price: 20.0},
	}

	mockRepo.EXPECT().GetAll().Return(expectedItems, nil)

	items, err := service.GetAllItems()

	assert.NoError(t, err)
	assert.Equal(t, 2, len(items))
	assert.Equal(t, "Item 1", items[0].Title)
}

func TestGetItem(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := mocks.NewItemRepository(t)
		service := NewItemService(mockRepo)

		expectedItem := &domain.Item{ID: 1, Title: "Item 1", Price: 10.0}
		mockRepo.EXPECT().GetByID(uint(1)).Return(expectedItem, nil)

		item, err := service.GetItem(1)

		assert.NoError(t, err)
		assert.NotNil(t, item)
		assert.Equal(t, uint(1), item.ID)
	})

	t.Run("invalid id", func(t *testing.T) {
		mockRepo := mocks.NewItemRepository(t)
		service := NewItemService(mockRepo)

		item, err := service.GetItem(0)

		assert.Error(t, err)
		assert.Nil(t, item)
		assert.Equal(t, "invalid item ID", err.Error())
	})
}

func TestCreateItem(t *testing.T) {
	mockRepo := mocks.NewItemRepository(t)
	service := NewItemService(mockRepo)

	req := domain.CreateItemRequest{
		Title: "New Item",
		Price: 15.0,
	}

	// We use Matcher because the item object is created inside the service and includes timestamps
	mockRepo.EXPECT().Create(mock.MatchedBy(func(item *domain.Item) bool {
		return item.Title == "New Item" && item.Price == 15.0
	})).Return(nil)

	item, err := service.CreateItem(req)

	assert.NoError(t, err)
	assert.NotNil(t, item)
	assert.Equal(t, "New Item", item.Title)
}

func TestDeleteItem(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := mocks.NewItemRepository(t)
		service := NewItemService(mockRepo)

		mockRepo.EXPECT().Delete(uint(1)).Return(nil)

		err := service.DeleteItem(1)

		assert.NoError(t, err)
	})

	t.Run("invalid id", func(t *testing.T) {
		mockRepo := mocks.NewItemRepository(t)
		service := NewItemService(mockRepo)

		err := service.DeleteItem(0)

		assert.Error(t, err)
		assert.Equal(t, "invalid item ID", err.Error())
	})
}

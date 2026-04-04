package domain

// These are the "Ports" in Hexagonal Architecture.

// ItemRepository is an outbound port. It defines the contract for what the core logic
// needs from a database or any other persistence mechanism.
type ItemRepository interface {
	GetAll() ([]Item, error)
	GetByID(id uint) (*Item, error)
	Create(item *Item) error
	Delete(id uint) error
}

// ItemService is an inbound port. It defines the core application's public API.
// The web handlers (or any other entry point) will use this interface to interact with the core logic.
type ItemService interface {
	GetAllItems() ([]Item, error)
	GetItem(id uint) (*Item, error)
	CreateItem(req CreateItemRequest) (*Item, error)
	DeleteItem(id uint) error
}

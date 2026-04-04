package domain

import "time"

// Item is the core business model. It has no knowledge of the database or web framework.
type Item struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// CreateItemRequest is a specific model for handling incoming creation requests.
// This separates the API model from the core domain model, which is good practice.
type CreateItemRequest struct {
	Title string  `json:"title" binding:"required,min=3"`
	Price float64 `json:"price" binding:"required,gt=0"`
}

package exampledtos

import "time"

// ExampleDTO represents the data transfer object for an example entity.
type ExampleDTO struct {
	ID     int64     `json:"id"`
	UserID string    `json:"user_id" validate:"required"`
	Amount int64     `json:"amount" validate:"gt=0"`
	Date   time.Time `json:"date"`
}

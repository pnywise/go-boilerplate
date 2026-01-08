package _mock

import (
    "context"
    "go-boilerplate/internal/entities"
)

// MockExampleRepository is a mock implementation of ExampleRepository for testing purposes.
type MockExampleRepository struct {
    GetByIDFunc func(ctx context.Context, id int64) (*entities.ExampleEntity, error)
    CreateFunc  func(ctx context.Context, u *entities.ExampleEntity) (int64, error)
}

// GetByID calls the mocked GetByIDFunc.
func (m *MockExampleRepository) GetByID(ctx context.Context, id int64) (*entities.ExampleEntity, error) {
    if m.GetByIDFunc != nil {
        return m.GetByIDFunc(ctx, id)
    }
    return nil, nil
}

// Create calls the mocked CreateFunc.
func (m *MockExampleRepository) Create(ctx context.Context, u *entities.ExampleEntity) (int64, error) {
    if m.CreateFunc != nil {
        return m.CreateFunc(ctx, u)
    }
    return 0, nil
}

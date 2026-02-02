package services

import (
    "context"
    "testing"

    "go-boilerplate/internal/configs"
    exampledtos "go-boilerplate/internal/dtos/example_dtos"
    "go-boilerplate/internal/entities"
    "go-boilerplate/internal/repositories/_mock"

    "github.com/go-playground/validator/v10"
    "github.com/stretchr/testify/require"
    "go.uber.org/zap"
)

func TestCreateExampleSuccess(t *testing.T) {
    mockRepo := &_mock.MockExampleRepository{
        CreateFunc: func(ctx context.Context, u *entities.ExampleEntity) (int64, error) {
            return 7, nil
        },
    }

    svc := NewExampleService(mockRepo, zap.NewNop(), configs.Config{}, validator.New())

    dto := exampledtos.ExampleDTO{
        UserID: "u1",
        Amount: 10,
    }

    id, err := svc.CreateExample(context.Background(), dto)
    require.NoError(t, err)
    require.Equal(t, int64(7), id)
}

func TestCreateExamplePassesEntityToRepo(t *testing.T) {
    var captured *entities.ExampleEntity
    mockRepo := &_mock.MockExampleRepository{
        CreateFunc: func(ctx context.Context, u *entities.ExampleEntity) (int64, error) {
            captured = u
            return 33, nil
        },
    }

    svc := NewExampleService(mockRepo, zap.NewNop(), configs.Config{}, validator.New())

    dto := exampledtos.ExampleDTO{
        UserID: "pass-through",
        Amount: 77,
    }

    id, err := svc.CreateExample(context.Background(), dto)
    require.NoError(t, err)
    require.Equal(t, int64(33), id)
    require.NotNil(t, captured)
    require.Equal(t, "pass-through", captured.UserID)
    require.Equal(t, int64(77), captured.Amount)
}
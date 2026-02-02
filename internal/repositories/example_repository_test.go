package repositories

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"go-boilerplate/internal/entities"
)

func TestExampleRepositoryGetByIDSuccess(t *testing.T) {
    db, mock, err := sqlmock.New()
    require.NoError(t, err)
    defer db.Close()

    repo := NewExampleRepository(db)

    rows := sqlmock.NewRows([]string{"id", "user_id", "amount"}).
        AddRow(int64(42), "user1", int64(500))
    mock.ExpectQuery(`SELECT id, user_id, amount FROM users WHERE id = \?`).
        WithArgs(int64(42)).
        WillReturnRows(rows)

    res, err := repo.GetByID(context.Background(), 42)
    require.NoError(t, err)
    require.NotNil(t, res)
    require.Equal(t, "user1", res.UserID)
    require.Equal(t, int64(500), res.Amount)
    require.NoError(t, mock.ExpectationsWereMet())
}

func TestExampleRepositoryGetByIDNoRows(t *testing.T) {
    db, mock, err := sqlmock.New()
    require.NoError(t, err)
    defer db.Close()

    repo := NewExampleRepository(db)

    mock.ExpectQuery(`SELECT id, user_id, amount FROM users WHERE id = \?`).
        WithArgs(int64(999)).
        WillReturnError(sql.ErrNoRows)

    res, err := repo.GetByID(context.Background(), 999)
    require.NoError(t, err)
    require.Nil(t, res)
    require.NoError(t, mock.ExpectationsWereMet())
}

func TestExampleRepositoryCreate(t *testing.T) {
    db, mock, err := sqlmock.New()
    require.NoError(t, err)
    defer db.Close()

    repo := NewExampleRepository(db)

    mock.ExpectExec(`INSERT INTO users \(user_id, amount\) VALUES \(\?, \?\)`).
        WithArgs("userX", int64(111)).
        WillReturnResult(sqlmock.NewResult(7, 1))

    id, err := repo.Create(context.Background(), &entities.ExampleEntity{UserID: "userX", Amount: 111})
    require.NoError(t, err)
    require.Equal(t, int64(7), id)
    require.NoError(t, mock.ExpectationsWereMet())
}
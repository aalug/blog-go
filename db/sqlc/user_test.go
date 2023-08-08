package db

import (
	"context"
	"database/sql"
	"github.com/aalug/blog-go/utils"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// createRandomUser creates a random user and returns it
func createRandomUser(t *testing.T) User {
	// password not hashed for now
	params := CreateUserParams{
		Username:       utils.RandomString(5),
		HashedPassword: utils.RandomString(6),
		Email:          utils.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, params.Username, user.Username)
	require.Equal(t, params.HashedPassword, user.HashedPassword)
	require.Equal(t, params.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

// TestQueries_CreateUser tests the create user function
func TestQueries_CreateUser(t *testing.T) {
	createRandomUser(t)
}

// TestQueries_GetUser tests the get user function
func TestQueries_GetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Email)

	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

// TestQueries_DeleteUser tests the delete user function
func TestQueries_DeleteUser(t *testing.T) {
	user := createRandomUser(t)

	err := testQueries.DeleteUser(context.Background(), user.Email)
	require.NoError(t, err)
	user2, err := testQueries.GetUser(context.Background(), user.Email)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, user2)
}

func TestQueries_ListUsersContainingString(t *testing.T) {
	tme := time.Now().Format("20060102150405")
	for i := 0; i < 10; i++ {
		params := CreateUserParams{
			Username:       utils.RandomString(3),
			HashedPassword: utils.RandomString(6),
			Email:          utils.RandomEmail(),
		}
		if i%2 == 0 {
			params.Username = tme + utils.RandomString(10)
		}
		_, err := testQueries.CreateUser(context.Background(), params)
		require.NoError(t, err)
	}

	users, err := testQueries.ListUsersContainingString(context.Background(), tme)
	require.NoError(t, err)
	require.Len(t, users, 5)
}

func TestQueries_UpdateUser(t *testing.T) {
	user := createRandomUser(t)
	params := UpdateUserParams{
		HashedPassword: sql.NullString{
			String: utils.RandomString(6),
			Valid:  true,
		},
		PasswordChangedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		Username: sql.NullString{
			String: utils.RandomString(5),
			Valid:  true,
		},
		Email: user.Email,
	}

	updatedUser, err := testQueries.UpdateUser(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)

	require.Equal(t, params.Username.String, updatedUser.Username)
	require.Equal(t, params.HashedPassword.String, updatedUser.HashedPassword)
	require.Equal(t, params.Email, updatedUser.Email)
	require.WithinDuration(t, params.PasswordChangedAt.Time, updatedUser.PasswordChangedAt, time.Second)
}

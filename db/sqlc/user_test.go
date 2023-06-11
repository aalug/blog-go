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

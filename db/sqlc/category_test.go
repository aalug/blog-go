package db

import (
	"context"
	"github.com/aalug/blog-go/utils"
	"github.com/stretchr/testify/require"
	"testing"
)

// createRandomCategory creates and returns a random category
func createRandomCategory(t *testing.T) Category {
	name := utils.RandomString(6)
	category, err := testQueries.CreateCategory(context.Background(), name)

	require.NoError(t, err)
	require.NotEmpty(t, category)
	require.Equal(t, category.Name, name)
	require.NotZero(t, category.ID)
	require.NotZero(t, category.CreatedAt)

	return category
}

// TestQueries_CreateCategory tests the create category function
func TestQueries_CreateCategory(t *testing.T) {
	createRandomCategory(t)
}

// TestQueries_ListCategories tests the list categories function
func TestQueries_ListCategories(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomCategory(t)
	}

	params := ListCategoriesParams{
		Limit:  5,
		Offset: 5,
	}

	categories, err := testQueries.ListCategories(context.Background(), params)

	require.NoError(t, err)
	require.NotEmpty(t, categories)
	require.Len(t, categories, 5)

	for _, category := range categories {
		require.NotEmpty(t, category)
		require.NotZero(t, category.ID)
		require.NotZero(t, category.CreatedAt)
	}
}

// TestQueries_DeleteCategory tests the delete category function
func TestQueries_DeleteCategory(t *testing.T) {
	category := createRandomCategory(t)

	err := testQueries.DeleteCategory(context.Background(), category.Name)
	require.NoError(t, err)
}

// TestQueries_UpdateCategory tests the update category function
func TestQueries_UpdateCategory(t *testing.T) {
	category := createRandomCategory(t)

	params := UpdateCategoryParams{
		Name:   category.Name,
		Name_2: utils.RandomString(4),
	}

	updatedCategory, err := testQueries.UpdateCategory(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, updatedCategory)
	require.Equal(t, updatedCategory.Name, params.Name_2)
	require.Equal(t, updatedCategory.ID, category.ID)
	require.NotZero(t, updatedCategory.CreatedAt)
}

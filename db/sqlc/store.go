package db

import (
	"context"
	"database/sql"
)

type Store interface {
	Querier
	AddTagsToPost(ctx context.Context, params AddTagsToPostParams) error
	RemoveTagsFromPost(ctx context.Context, params RemoveTagsFromPostParams) error
}

// SQLStore provides all functions to execute db queries and transactions
type SQLStore struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new Store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

type AddTagsToPostParams struct {
	PostID int64
	Tags   []string
}

// AddTagsToPost creates tags (if they don't exist) and creates post_tags table (many to many)
func (store SQLStore) AddTagsToPost(ctx context.Context, arg AddTagsToPostParams) error {
	// get tag ids
	tagIDs, err := store.GetOrCreateTags(ctx, arg.Tags)
	if err != nil {
		return err
	}

	params := AddMultipleTagsToPostParams{
		PostID: arg.PostID,
		TagIds: tagIDs,
	}

	// add all tags to the post
	err = store.AddMultipleTagsToPost(ctx, params)
	if err != nil {
		return err
	}
	return nil
}

type RemoveTagsFromPostParams struct {
	PostID int64
	Tags   []string
}

// RemoveTagsFromPost for each tag, checks if it's used by a different post, if not - deletes the tag and post_tags table,
// if it is used be another post - just deletes the post_tags table
func (store SQLStore) RemoveTagsFromPost(ctx context.Context, params RemoveTagsFromPostParams) error {
	//TODO: implement me
	panic("implement me")
}

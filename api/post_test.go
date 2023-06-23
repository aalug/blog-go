package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	mockdb "github.com/aalug/blog-go/db/mock"
	db "github.com/aalug/blog-go/db/sqlc"
	"github.com/aalug/blog-go/token"
	"github.com/aalug/blog-go/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreatePostAPI(t *testing.T) {
	randomUser, _ := generateRandomUser(t)
	category, post, tags := generateRandomCategoryPostAndTags(int32(randomUser.ID))

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, r *http.Request, maker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"title":       post.Title,
				"description": post.Description,
				"content":     post.Content,
				"image":       post.Image,
				"tags":        tags,
				"category":    category.Name,
			},
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, randomUser.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetOrCreateCategory(gomock.Any(), gomock.Eq(category.Name)).
					Times(1).
					Return(category.ID, nil)
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(randomUser.Email)).
					Times(1).
					Return(randomUser, nil)

				params := db.CreatePostParams{
					Title:       post.Title,
					Description: post.Description,
					Content:     post.Content,
					AuthorID:    post.AuthorID,
					CategoryID:  post.CategoryID,
					Image:       post.Image,
				}
				store.EXPECT().
					CreatePost(gomock.Any(), gomock.Eq(params)).
					Times(1).
					Return(post, nil)

				postTagsParams := db.AddTagsToPostParams{
					PostID: post.ID,
					Tags:   tags,
				}
				store.EXPECT().
					AddTagsToPost(gomock.Any(), gomock.Eq(postTagsParams)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name: "Internal Server Error In CreatePost",
			body: gin.H{
				"title":       post.Title,
				"description": post.Description,
				"content":     post.Content,
				"image":       post.Image,
				"tags":        tags,
				"category":    category.Name,
			},
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, randomUser.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetOrCreateCategory(gomock.Any(), gomock.Any()).
					Times(1).
					Return(category.ID, nil)
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(randomUser, nil)
				store.EXPECT().
					CreatePost(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Post{}, sql.ErrConnDone)
				store.EXPECT().
					AddTagsToPost(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Internal Server Error In AddTagsToPost",
			body: gin.H{
				"title":       post.Title,
				"description": post.Description,
				"content":     post.Content,
				"image":       post.Image,
				"tags":        tags,
				"category":    category.Name,
			},
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, randomUser.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetOrCreateCategory(gomock.Any(), gomock.Any()).
					Times(1).
					Return(category.ID, nil)
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(randomUser, nil)
				store.EXPECT().
					CreatePost(gomock.Any(), gomock.Any()).
					Times(1).
					Return(post, nil)
				store.EXPECT().
					AddTagsToPost(gomock.Any(), gomock.Any()).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Internal Server Error In GetOrCreateCategory",
			body: gin.H{
				"title":       post.Title,
				"description": post.Description,
				"content":     post.Content,
				"image":       post.Image,
				"tags":        tags,
				"category":    category.Name,
			},
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, randomUser.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetOrCreateCategory(gomock.Any(), gomock.Any()).
					Times(1).
					Return(category.ID, sql.ErrConnDone)
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0)
				store.EXPECT().
					CreatePost(gomock.Any(), gomock.Any()).
					Times(0)
				store.EXPECT().
					AddTagsToPost(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Internal Server Error In GetUser",
			body: gin.H{
				"title":       post.Title,
				"description": post.Description,
				"content":     post.Content,
				"image":       post.Image,
				"tags":        tags,
				"category":    category.Name,
			},
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, randomUser.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetOrCreateCategory(gomock.Any(), gomock.Any()).
					Times(1).
					Return(category.ID, nil)
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
				store.EXPECT().
					CreatePost(gomock.Any(), gomock.Any()).
					Times(0)
				store.EXPECT().
					AddTagsToPost(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Invalid Body",
			body: gin.H{
				"title":       post.Title,
				"description": post.Description,
				"tags":        tags,
				"category":    category.Name,
			},
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, randomUser.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetOrCreateCategory(gomock.Any(), gomock.Any()).
					Times(0)
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0)
				store.EXPECT().
					CreatePost(gomock.Any(), gomock.Any()).
					Times(0)
				store.EXPECT().
					AddTagsToPost(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/posts"
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, req, server.tokenMaker)

			server.router.ServeHTTP(recorder, req)

			tc.checkResponse(recorder)
		})
	}
}

func TestDeletePostAPI(t *testing.T) {
	randomUser, _ := generateRandomUser(t)
	_, post, _ := generateRandomCategoryPostAndTags(int32(randomUser.ID))

	testCases := []struct {
		name          string
		postID        int64
		setupAuth     func(t *testing.T, r *http.Request, maker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			postID: post.ID,
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, randomUser.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(randomUser.Email)).
					Times(1).
					Return(randomUser, nil)

				data := db.GetMinimalPostDataRow{
					ID:       post.ID,
					AuthorID: post.AuthorID,
				}
				store.EXPECT().
					GetMinimalPostData(gomock.Any(), gomock.Eq(post.ID)).
					Times(1).
					Return(data, nil)
				store.EXPECT().
					DeletePost(gomock.Any(), gomock.Eq(post.ID)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNoContent, recorder.Code)
			},
		},
		{
			name:   "Invalid Post ID",
			postID: 0,
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, randomUser.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0)
				store.EXPECT().
					GetMinimalPostData(gomock.Any(), gomock.Any()).
					Times(0)
				store.EXPECT().
					DeletePost(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "Internal Server Error GetUser",
			postID: post.ID,
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, randomUser.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
				store.EXPECT().
					GetMinimalPostData(gomock.Any(), gomock.Any()).
					Times(0)
				store.EXPECT().
					DeletePost(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:   "Internal Server Error GetMinimalPostData",
			postID: post.ID,
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, randomUser.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(randomUser, nil)
				store.EXPECT().
					GetMinimalPostData(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.GetMinimalPostDataRow{}, sql.ErrConnDone)
				store.EXPECT().
					DeletePost(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:   "Internal Server Error DeletePost",
			postID: post.ID,
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, randomUser.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(randomUser, nil)
				data := db.GetMinimalPostDataRow{
					ID:       post.ID,
					AuthorID: post.AuthorID,
				}
				store.EXPECT().
					GetMinimalPostData(gomock.Any(), gomock.Eq(post.ID)).
					Times(1).
					Return(data, nil)
				store.EXPECT().
					DeletePost(gomock.Any(), gomock.Any()).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:   "Not Found",
			postID: post.ID,
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, randomUser.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(randomUser, nil)
				store.EXPECT().
					GetMinimalPostData(gomock.Any(), gomock.Eq(post.ID)).
					Times(1).
					Return(db.GetMinimalPostDataRow{}, sql.ErrNoRows)
				store.EXPECT().
					DeletePost(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:   "Unauthorized User",
			postID: post.ID,
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, "unauthorized@example.com", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{
						ID:       int64(utils.RandomInt(11, 20)),
						Email:    "unauthorized@example.com",
						Username: "unauthorized",
					}, nil)
				data := db.GetMinimalPostDataRow{
					ID:       post.ID,
					AuthorID: post.AuthorID,
				}
				store.EXPECT().
					GetMinimalPostData(gomock.Any(), gomock.Eq(post.ID)).
					Times(1).
					Return(data, nil)
				store.EXPECT().
					DeletePost(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/posts/%d", tc.postID)
			req, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, req, server.tokenMaker)

			server.router.ServeHTTP(recorder, req)

			tc.checkResponse(recorder)
		})
	}
}

// generateRandomCategoryPostAndTags generates random category, post and tags
func generateRandomCategoryPostAndTags(userID int32) (db.Category, db.Post, []string) {
	category := db.Category{
		ID:   int64(utils.RandomInt(1, 10)),
		Name: utils.RandomString(6),
	}

	post := db.Post{
		ID:          int64(utils.RandomInt(1, 10)),
		Title:       utils.RandomString(6),
		Description: utils.RandomString(7),
		Content:     utils.RandomString(8),
		AuthorID:    userID,
		CategoryID:  int32(category.ID),
		Image:       fmt.Sprintf("%s.jpg", utils.RandomString(3)),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	tags := []string{
		utils.RandomString(3),
		utils.RandomString(4),
	}

	return category, post, tags
}

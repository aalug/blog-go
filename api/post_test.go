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
	category := db.Category{
		ID:   int64(utils.RandomInt(1, 10)),
		Name: utils.RandomString(6),
	}

	post := db.Post{
		Title:       utils.RandomString(6),
		Description: utils.RandomString(7),
		Content:     utils.RandomString(8),
		AuthorID:    int32(randomUser.ID),
		CategoryID:  int32(category.ID),
		Image:       fmt.Sprintf("%s.jpg", utils.RandomString(3)),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	tags := []string{
		utils.RandomString(3),
		utils.RandomString(4),
	}

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

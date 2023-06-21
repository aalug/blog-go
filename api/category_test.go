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
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreateCategoryAPI(t *testing.T) {
	randomUser, _ := generateRandomUser(t)
	category := db.Category{
		Name: utils.RandomString(5),
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
				"name": category.Name,
			},
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, randomUser.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateCategory(gomock.Any(), gomock.Eq(category.Name)).
					Times(1).
					Return(category, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireBodyMatchCategory(t, recorder.Body, category)
			},
		},
		{
			name: "Internal Server Error",
			body: gin.H{
				"name": category.Name,
			},
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, randomUser.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateCategory(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Category{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Duplicated Name",
			body: gin.H{
				"name": category.Name,
			},
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, randomUser.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateCategory(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Category{}, &pq.Error{Code: "23505"})
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "Invalid Name",
			body: gin.H{
				"name": "123Invalid$#",
			},
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, randomUser.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateCategory(gomock.Any(), gomock.Any()).
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

			url := "/category"
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, req, server.tokenMaker)

			server.router.ServeHTTP(recorder, req)

			tc.checkResponse(recorder)
		})
	}
}

func TestDeleteCategoryAPI(t *testing.T) {
	randomUser, _ := generateRandomUser(t)
	category := db.Category{
		Name: utils.RandomString(5),
	}

	testCases := []struct {
		name          string
		categoryName  string
		setupAuth     func(t *testing.T, r *http.Request, maker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:         "OK",
			categoryName: category.Name,
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, randomUser.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteCategory(gomock.Any(), gomock.Eq(category.Name)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNoContent, recorder.Code)
			},
		},
		{
			name:         "Not Found",
			categoryName: category.Name,
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, randomUser.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteCategory(gomock.Any(), gomock.Any()).
					Times(1).
					Return(sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:         "Internal Server Error",
			categoryName: category.Name,
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, randomUser.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteCategory(gomock.Any(), gomock.Any()).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:         "Invalid Name",
			categoryName: "$@#^!",
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, randomUser.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteCategory(gomock.Any(), gomock.Any()).
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

			url := fmt.Sprintf("/category/%s", tc.categoryName)
			req, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, req, server.tokenMaker)

			server.router.ServeHTTP(recorder, req)

			tc.checkResponse(recorder)
		})
	}
}

func TestListCategoriesAPI(t *testing.T) {
	n := 5
	categories := make([]db.Category, n)
	for i := 0; i < n; i++ {
		categories[i] = db.Category{Name: utils.RandomString(5)}
	}

	type Query struct {
		page     int
		pageSize int
	}

	testCases := []struct {
		name          string
		query         Query
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				page:     1,
				pageSize: n,
			},
			buildStubs: func(store *mockdb.MockStore) {
				params := db.ListCategoriesParams{
					Limit:  int32(n),
					Offset: 0,
				}
				store.EXPECT().
					ListCategories(gomock.Any(), gomock.Eq(params)).
					Times(1).
					Return(categories, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchCategories(t, recorder.Body, categories)
			},
		},
		{
			name: "Internal Server Error",
			query: Query{
				page:     1,
				pageSize: n,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListCategories(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Category{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Invalid Page Size",
			query: Query{
				page:     1,
				pageSize: 4,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListCategories(gomock.Any(), gomock.Any()).
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

			url := "/category"
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// Add query params
			q := req.URL.Query()
			q.Add("page", fmt.Sprintf("%d", tc.query.page))
			q.Add("page_size", fmt.Sprintf("%d", tc.query.pageSize))
			req.URL.RawQuery = q.Encode()

			server.router.ServeHTTP(recorder, req)

			tc.checkResponse(recorder)
		})
	}
}

func requireBodyMatchCategory(t *testing.T, body *bytes.Buffer, category db.Category) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotCategory db.Category
	err = json.Unmarshal(data, &gotCategory)
	require.NoError(t, err)
	require.Equal(t, category, gotCategory)
}

func requireBodyMatchCategories(t *testing.T, body *bytes.Buffer, categories []db.Category) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotCategories []db.Category
	err = json.Unmarshal(data, &gotCategories)
	require.NoError(t, err)
	require.Equal(t, categories, gotCategories)
}

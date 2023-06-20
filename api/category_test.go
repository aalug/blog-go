package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	mockdb "github.com/aalug/blog-go/db/mock"
	db "github.com/aalug/blog-go/db/sqlc"
	"github.com/aalug/blog-go/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateCategoryAPI(t *testing.T) {
	category := db.Category{
		Name: utils.RandomString(5),
	}

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"name": category.Name,
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

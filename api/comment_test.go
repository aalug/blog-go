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

func TestCreateCommentAPI(t *testing.T) {
	randomUser, _ := generateRandomUser(t)
	_, post, _ := generateRandomCategoryPostAndTags(int32(randomUser.ID))
	comment := db.Comment{
		Content: "test comment",
		UserID:  int32(randomUser.ID),
		PostID:  utils.RandomInt(1, 100),
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
				"content": comment.Content,
				"post_id": post.ID,
			},
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, randomUser.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(randomUser.Email)).
					Times(1).
					Return(randomUser, nil)
				store.EXPECT().
					CreateComment(gomock.Any(), gomock.Any()).
					Times(1).
					Return(comment, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireBodyMatchComment(t, recorder.Body, comment)
			},
		},
		{
			name: "Internal Server Error GetUser",
			body: gin.H{
				"content": comment.Content,
				"post_id": post.ID,
			},
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, randomUser.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(randomUser.Email)).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
				store.EXPECT().
					CreateComment(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Internal Server Error CreateComment",
			body: gin.H{
				"content": comment.Content,
				"post_id": post.ID,
			},
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, randomUser.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(randomUser.Email)).
					Times(1).
					Return(randomUser, nil)
				store.EXPECT().
					CreateComment(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Comment{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Foreign Key Violation",
			body: gin.H{
				"content": comment.Content,
				"post_id": post.ID,
			},
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, randomUser.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(randomUser.Email)).
					Times(1).
					Return(randomUser, nil)
				store.EXPECT().
					CreateComment(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Comment{}, &pq.Error{Code: "23503"})
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Invalid Post ID",
			body: gin.H{
				"content": comment.Content,
				"post_id": 0,
			},
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, randomUser.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0)
				store.EXPECT().
					CreateComment(gomock.Any(), gomock.Any()).
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

			url := "/comments"
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, req, server.tokenMaker)

			server.router.ServeHTTP(recorder, req)

			tc.checkResponse(recorder)
		})
	}
}

func TestDeleteCommentAPI(t *testing.T) {
	randomUser, _ := generateRandomUser(t)
	comment := db.Comment{
		ID:      int64(utils.RandomInt(1, 100)),
		Content: "test comment",
		UserID:  int32(randomUser.ID),
		PostID:  utils.RandomInt(1, 100),
	}

	testCases := []struct {
		name          string
		commentID     int64
		setupAuth     func(t *testing.T, r *http.Request, maker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			commentID: comment.ID,
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, randomUser.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetComment(gomock.Any(), gomock.Eq(comment.ID)).
					Times(1).
					Return(db.GetCommentRow{
						ID:     comment.ID,
						UserID: int32(randomUser.ID),
					}, nil)
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(randomUser.Email)).
					Times(1).
					Return(randomUser, nil)
				store.EXPECT().
					DeleteComment(gomock.Any(), gomock.Eq(comment.ID)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNoContent, recorder.Code)
			},
		},
		{
			name:      "Not Found",
			commentID: comment.ID,
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, randomUser.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetComment(gomock.Any(), gomock.Eq(comment.ID)).
					Times(1).
					Return(db.GetCommentRow{}, sql.ErrNoRows)
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0)
				store.EXPECT().
					DeleteComment(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "Unauthorized User",
			commentID: comment.ID,
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, "unauthorized@example.com", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetComment(gomock.Any(), gomock.Eq(comment.ID)).
					Times(1).
					Return(db.GetCommentRow{
						ID:     comment.ID,
						UserID: int32(randomUser.ID),
					}, nil)
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq("unauthorized@example.com")).
					Times(1).
					Return(db.User{
						ID:       999,
						Username: "unauthorized",
						Email:    "unauthorized@example.com",
					}, nil)
				store.EXPECT().
					DeleteComment(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:      "Invalid Comment ID",
			commentID: 0,
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, randomUser.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetComment(gomock.Any(), gomock.Any()).
					Times(0)
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0)
				store.EXPECT().
					DeleteComment(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "Internal Server Error GetComment",
			commentID: comment.ID,
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, randomUser.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetComment(gomock.Any(), gomock.Eq(comment.ID)).
					Times(1).
					Return(db.GetCommentRow{}, sql.ErrConnDone)
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0)
				store.EXPECT().
					DeleteComment(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "Internal Server Error GetUser",
			commentID: comment.ID,
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, randomUser.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetComment(gomock.Any(), gomock.Eq(comment.ID)).
					Times(1).
					Return(db.GetCommentRow{
						ID:     comment.ID,
						UserID: int32(randomUser.ID),
					}, nil)
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(randomUser.Email)).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
				store.EXPECT().
					DeleteComment(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "Internal Server Error DeleteComment",
			commentID: comment.ID,
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, randomUser.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetComment(gomock.Any(), gomock.Eq(comment.ID)).
					Times(1).
					Return(db.GetCommentRow{
						ID:     comment.ID,
						UserID: int32(randomUser.ID),
					}, nil)
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(randomUser.Email)).
					Times(1).
					Return(randomUser, nil)
				store.EXPECT().
					DeleteComment(gomock.Any(), gomock.Eq(comment.ID)).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
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

			url := fmt.Sprintf("/comments/%d", tc.commentID)
			req, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, req, server.tokenMaker)

			server.router.ServeHTTP(recorder, req)

			tc.checkResponse(recorder)
		})
	}
}

func TestUpdateCommentAPI(t *testing.T) {
	randomUser, _ := generateRandomUser(t)
	comment := db.Comment{
		ID:      int64(utils.RandomInt(1, 100)),
		Content: "test comment",
		UserID:  int32(randomUser.ID),
		PostID:  utils.RandomInt(1, 100),
	}

	testCases := []struct {
		name          string
		commentID     int64
		body          gin.H
		setupAuth     func(t *testing.T, r *http.Request, maker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			commentID: comment.ID,
			body: gin.H{
				"content": "updated content",
			},
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, randomUser.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetComment(gomock.Any(), gomock.Eq(comment.ID)).
					Times(1).
					Return(db.GetCommentRow{
						ID:     comment.ID,
						UserID: int32(randomUser.ID),
					}, nil)
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(randomUser.Email)).
					Times(1).
					Return(randomUser, nil)
				params := db.UpdateCommentParams{
					ID:      comment.ID,
					Content: "updated content",
				}
				store.EXPECT().
					UpdateComment(gomock.Any(), gomock.Eq(params)).
					Times(1).
					Return(comment, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:      "Not Found",
			commentID: comment.ID,
			body: gin.H{
				"content": "updated content",
			},
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, randomUser.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetComment(gomock.Any(), gomock.Eq(comment.ID)).
					Times(1).
					Return(db.GetCommentRow{}, sql.ErrNoRows)
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0)
				store.EXPECT().
					UpdateComment(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "Unauthorized User",
			commentID: comment.ID,
			body: gin.H{
				"content": "updated content",
			},
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, "unauthorized@example.com", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetComment(gomock.Any(), gomock.Eq(comment.ID)).
					Times(1).
					Return(db.GetCommentRow{
						ID:     comment.ID,
						UserID: int32(randomUser.ID),
					}, nil)
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq("unauthorized@example.com")).
					Times(1).
					Return(db.User{
						ID:       999,
						Username: "unauthorized",
						Email:    "unauthorized@example.com",
					}, nil)
				store.EXPECT().
					DeleteComment(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:      "Invalid Comment ID",
			commentID: 0,
			body: gin.H{
				"content": "updated content",
			},
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, randomUser.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetComment(gomock.Any(), gomock.Any()).
					Times(0)
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0)
				store.EXPECT().
					DeleteComment(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "Internal Server Error GetComment",
			commentID: comment.ID,
			body: gin.H{
				"content": "updated content",
			},
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, randomUser.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetComment(gomock.Any(), gomock.Eq(comment.ID)).
					Times(1).
					Return(db.GetCommentRow{}, sql.ErrConnDone)
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0)
				store.EXPECT().
					UpdateComment(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "Internal Server Error GetUser",
			commentID: comment.ID,
			body: gin.H{
				"content": "updated content",
			},
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, randomUser.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetComment(gomock.Any(), gomock.Eq(comment.ID)).
					Times(1).
					Return(db.GetCommentRow{
						ID:     comment.ID,
						UserID: int32(randomUser.ID),
					}, nil)
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(randomUser.Email)).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
				store.EXPECT().
					UpdateComment(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "Internal Server Error UpdateComment",
			commentID: comment.ID,
			body: gin.H{
				"content": "updated content",
			},
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, randomUser.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetComment(gomock.Any(), gomock.Eq(comment.ID)).
					Times(1).
					Return(db.GetCommentRow{
						ID:     comment.ID,
						UserID: int32(randomUser.ID),
					}, nil)
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(randomUser.Email)).
					Times(1).
					Return(randomUser, nil)
				params := db.UpdateCommentParams{
					ID:      comment.ID,
					Content: "updated content",
				}
				store.EXPECT().
					UpdateComment(gomock.Any(), gomock.Eq(params)).
					Times(1).
					Return(db.Comment{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "Empty Content",
			commentID: comment.ID,
			body: gin.H{
				"content": "",
			},
			setupAuth: func(t *testing.T, r *http.Request, maker token.Maker) {
				addAuthorization(t, r, maker, authorizationTypeBearer, randomUser.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetComment(gomock.Any(), gomock.Eq(comment.ID)).
					Times(1).
					Return(db.GetCommentRow{
						ID:     comment.ID,
						UserID: int32(randomUser.ID),
					}, nil)
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(randomUser.Email)).
					Times(1).
					Return(randomUser, nil)
				store.EXPECT().
					UpdateComment(gomock.Any(), gomock.Any()).
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

			url := fmt.Sprintf("/comments/%d", tc.commentID)
			req, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, req, server.tokenMaker)

			server.router.ServeHTTP(recorder, req)

			tc.checkResponse(recorder)
		})
	}
}

func requireBodyMatchComment(t *testing.T, body *bytes.Buffer, comment db.Comment) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotComment db.Comment
	err = json.Unmarshal(data, &gotComment)
	require.NoError(t, err)

	require.Equal(t, comment.Content, gotComment.Content)
	require.Equal(t, comment.UserID, gotComment.UserID)
	require.Equal(t, comment.PostID, gotComment.PostID)
	require.WithinDuration(t, comment.CreatedAt, gotComment.CreatedAt, time.Second)
}

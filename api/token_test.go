package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	mockdb "github.com/aalug/blog-go/db/mock"
	db "github.com/aalug/blog-go/db/sqlc"
	"github.com/aalug/blog-go/token"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRenewAccessTokenAPI(t *testing.T) {
	user, _ := generateRandomUser(t)

	testCases := []struct {
		name          string
		body          func(server *Server) (gin.H, db.Session, *token.Payload)
		buildStubs    func(store *mockdb.MockStore, session db.Session, refreshPayload *token.Payload)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: func(server *Server) (gin.H, db.Session, *token.Payload) {
				refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
					user.Email,
					time.Minute,
				)
				require.NoError(t, err)
				require.NotEmpty(t, refreshToken)
				require.NotEmpty(t, refreshPayload)
				session := db.Session{
					ID:           uuid.UUID{},
					Email:        user.Email,
					RefreshToken: refreshToken,
					IsBlocked:    false,
					ExpiresAt:    time.Now().Add(time.Minute),
				}
				return gin.H{
					"refresh_token": refreshToken,
				}, session, refreshPayload
			},
			buildStubs: func(store *mockdb.MockStore, session db.Session, refreshPayload *token.Payload) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(refreshPayload.ID)).
					Times(1).
					Return(session, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Empty Body",
			body: func(server *Server) (gin.H, db.Session, *token.Payload) {
				refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
					user.Email,
					time.Minute,
				)
				require.NoError(t, err)
				require.NotEmpty(t, refreshToken)
				require.NotEmpty(t, refreshPayload)
				session := db.Session{
					ID:           uuid.UUID{},
					Email:        user.Email,
					RefreshToken: refreshToken,
					IsBlocked:    false,
					ExpiresAt:    time.Now().Add(time.Minute),
				}
				return gin.H{
					"refresh_token": "",
				}, session, refreshPayload
			},
			buildStubs: func(store *mockdb.MockStore, session db.Session, refreshPayload *token.Payload) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Not Found",
			body: func(server *Server) (gin.H, db.Session, *token.Payload) {
				refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
					user.Email,
					time.Minute,
				)
				require.NoError(t, err)
				require.NotEmpty(t, refreshToken)
				require.NotEmpty(t, refreshPayload)
				session := db.Session{
					ID:           uuid.UUID{},
					Email:        user.Email,
					RefreshToken: refreshToken,
					IsBlocked:    false,
					ExpiresAt:    time.Now().Add(time.Minute),
				}
				return gin.H{
					"refresh_token": refreshToken,
				}, session, refreshPayload
			},
			buildStubs: func(store *mockdb.MockStore, session db.Session, refreshPayload *token.Payload) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Session{}, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "Internal Server Error",
			body: func(server *Server) (gin.H, db.Session, *token.Payload) {
				refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
					user.Email,
					time.Minute,
				)
				require.NoError(t, err)
				require.NotEmpty(t, refreshToken)
				require.NotEmpty(t, refreshPayload)
				session := db.Session{
					ID:           uuid.UUID{},
					Email:        user.Email,
					RefreshToken: refreshToken,
					IsBlocked:    false,
					ExpiresAt:    time.Now().Add(time.Minute),
				}
				return gin.H{
					"refresh_token": refreshToken,
				}, session, refreshPayload
			},
			buildStubs: func(store *mockdb.MockStore, session db.Session, refreshPayload *token.Payload) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Session{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Invalid Token",
			body: func(server *Server) (gin.H, db.Session, *token.Payload) {
				refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
					user.Email,
					time.Minute,
				)
				require.NoError(t, err)
				require.NotEmpty(t, refreshToken)
				require.NotEmpty(t, refreshPayload)
				session := db.Session{
					ID:           uuid.UUID{},
					Email:        user.Email,
					RefreshToken: refreshToken,
					IsBlocked:    false,
					ExpiresAt:    time.Now().Add(time.Minute),
				}
				return gin.H{
					"refresh_token": "123",
				}, session, refreshPayload
			},
			buildStubs: func(store *mockdb.MockStore, session db.Session, refreshPayload *token.Payload) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Blocked Session",
			body: func(server *Server) (gin.H, db.Session, *token.Payload) {
				refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
					user.Email,
					time.Minute,
				)
				require.NoError(t, err)
				require.NotEmpty(t, refreshToken)
				require.NotEmpty(t, refreshPayload)
				session := db.Session{
					ID:           uuid.UUID{},
					Email:        user.Email,
					RefreshToken: refreshToken,
					IsBlocked:    true,
					ExpiresAt:    time.Now().Add(time.Minute),
				}
				return gin.H{
					"refresh_token": refreshToken,
				}, session, refreshPayload
			},
			buildStubs: func(store *mockdb.MockStore, session db.Session, refreshPayload *token.Payload) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(session, nil)
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

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			body, session, refreshPayload := tc.body(server)

			tc.buildStubs(store, session, refreshPayload)

			data, err := json.Marshal(body)
			require.NoError(t, err)

			url := "/tokens/renew"
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)

			tc.checkResponse(recorder)
		})
	}
}

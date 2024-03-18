package handlers

import (
	"bytes"
	"context"
	"errors"
	"github.com/jorgini/filmoteka/service"
	"github.com/jorgini/filmoteka/service/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_validation(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockUser, token string)

	testTable := []struct {
		name                 string
		headerName           string
		headerValue          string
		token                string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Ok",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(r *mock_service.MockUser, token string) {
				r.EXPECT().ParseToken(token).Return(1, nil)
				r.EXPECT().ValidateUser(1).Return(true, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "",
		},
		{
			name:                 "Invalid Header Name",
			headerName:           "",
			headerValue:          "Bearer token",
			token:                "token",
			mockBehavior:         func(r *mock_service.MockUser, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `empty header`,
		},
		{
			name:                 "Invalid Header Value",
			headerName:           "Authorization",
			headerValue:          "Bearr token",
			token:                "token",
			mockBehavior:         func(r *mock_service.MockUser, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `invalid auth header`,
		},
		{
			name:                 "Empty Token",
			headerName:           "Authorization",
			headerValue:          "Bearer ",
			token:                "token",
			mockBehavior:         func(r *mock_service.MockUser, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `token is empty`,
		},
		{
			name:        "Parse Error",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(r *mock_service.MockUser, token string) {
				r.EXPECT().ParseToken(token).Return(0, errors.New("invalid token"))
			},
			expectedStatusCode:   401,
			expectedResponseBody: `invalid token`,
		},
		{
			name:        "Not valid",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(r *mock_service.MockUser, token string) {
				r.EXPECT().ParseToken(token).Return(1, nil)
				r.EXPECT().ValidateUser(1).Return(false, nil)
			},
			expectedStatusCode:   423,
			expectedResponseBody: `this function locked for current user`,
		},
		{
			name:        "Service error",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(r *mock_service.MockUser, token string) {
				r.EXPECT().ParseToken(token).Return(1, nil)
				r.EXPECT().ValidateUser(1).Return(false, errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `something went wrong`,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockUser(c)
			test.mockBehavior(repo, test.token)

			services := &service.Service{User: repo}
			handler := Router{service: services}
			handler.AddEndPoint("GET", "/users", validation)

			// Init Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/users", nil)
			req.Header.Set(test.headerName, test.headerValue)

			handler.ServeHTTP(w, req)

			// Asserts
			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

func TestGetUserId(t *testing.T) {
	var getContext = func(id int) context.Context {
		ctx := context.WithValue(context.Background(), userCtx, id)
		return ctx
	}

	testTable := []struct {
		name       string
		ctx        context.Context
		id         int
		shouldFail bool
	}{
		{
			name: "Ok",
			ctx:  getContext(1),
			id:   1,
		},
		{
			ctx:        context.Background(),
			name:       "Empty",
			shouldFail: true,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			request, err := http.NewRequest("POST", "/smt", bytes.NewBufferString(""))
			if err != nil {
				t.Error(err)
			}

			req := request.Clone(test.ctx)
			id, err := getUserId(req)
			if test.shouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, id, test.id)
		})
	}
}

package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/jorgini/filmoteka"
	"github.com/jorgini/filmoteka/service"
	"github.com/jorgini/filmoteka/service/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_createNewUser(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockUser, user filmoteka.User)

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            filmoteka.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"login": "login", "password": "qwerty", "user_role": "regular"}`,
			inputUser: filmoteka.User{
				Login:    "login",
				Password: "qwerty",
				UserRole: "regular",
			},
			mockBehavior: func(r *mock_service.MockUser, user filmoteka.User) {
				r.EXPECT().CreateUser(user).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `"successful"`,
		},
		{
			name:                 "Wrong Input",
			inputBody:            `{"login": "login"}`,
			inputUser:            filmoteka.User{},
			mockBehavior:         func(r *mock_service.MockUser, user filmoteka.User) {},
			expectedStatusCode:   400,
			expectedResponseBody: `invalid state for required field(s)`,
		},
		{
			name:                 "Wrong role",
			inputBody:            `{"login": "login", "password": "qwerty", "user_role": "smt"}`,
			inputUser:            filmoteka.User{},
			mockBehavior:         func(r *mock_service.MockUser, user filmoteka.User) {},
			expectedStatusCode:   400,
			expectedResponseBody: `invalid state for required field(s)`,
		},
		{
			name:      "Service Error",
			inputBody: `{"login": "login", "password": "qwerty", "user_role": "regular"}`,
			inputUser: filmoteka.User{
				Login:    "login",
				Password: "qwerty",
				UserRole: "regular",
			},
			mockBehavior: func(r *mock_service.MockUser, user filmoteka.User) {
				r.EXPECT().CreateUser(user).Return(0, errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `something went wrong`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockUser(c)
			test.mockBehavior(repo, test.inputUser)

			services := &service.Service{User: repo}
			handler := Router{service: services}
			handler.AddEndPoint("POST", "/users", createNewUser)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/users",
				bytes.NewBufferString(test.inputBody))

			// Make Request
			handler.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, strings.ReplaceAll(w.Body.String(), "\n", ""))
		})
	}
}

func TestHandler_authUser(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockUser, input signInput)

	const token = "fmekwfmw"

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            signInput
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"login": "login", "password": "qwerty"}`,
			inputUser: signInput{
				Login:    "login",
				Password: "qwerty",
			},
			mockBehavior: func(r *mock_service.MockUser, input signInput) {
				r.EXPECT().GenerateToken(input.Login, input.Password).Return(token, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: fmt.Sprintf(`"%s"`, token),
		},
		{
			name:                 "Wrong Input",
			inputBody:            `{"login": "login"}`,
			inputUser:            signInput{},
			mockBehavior:         func(r *mock_service.MockUser, input signInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `missing required fields`,
		},
		{
			name:      "Service Error",
			inputBody: `{"login": "login", "password": "qwerty"}`,
			inputUser: signInput{
				Login:    "login",
				Password: "qwerty",
			},
			mockBehavior: func(r *mock_service.MockUser, input signInput) {
				r.EXPECT().GenerateToken(input.Login, input.Password).Return("", errors.New("something went wrong"))
			},
			expectedStatusCode:   401,
			expectedResponseBody: `something went wrong`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockUser(c)
			test.mockBehavior(repo, test.inputUser)

			services := &service.Service{User: repo}
			handler := Router{service: services}
			handler.AddEndPoint("GET", "/users", authUser)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/users",
				bytes.NewBufferString(test.inputBody))

			// Make Request
			handler.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, strings.ReplaceAll(w.Body.String(), "\n", ""))
		})
	}
}

func TestHandler_updateUser(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockUser, input updateInput)

	const token = "fmekwfmw"

	tests := []struct {
		name                 string
		headerName           string
		headerValue          string
		inputBody            string
		inputUser            updateInput
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Ok",
			headerName:  authorizationHeader,
			headerValue: fmt.Sprintf("Bearer %s", token),
			inputBody:   `{"login": "login", "user_role": "admin"}`,
			inputUser: updateInput{
				Login:    "login",
				UserRole: "admin",
			},
			mockBehavior: func(r *mock_service.MockUser, input updateInput) {
				r.EXPECT().ParseToken(token).Return(1, nil)
				r.EXPECT().ValidateUser(1).Return(true, nil)
				r.EXPECT().UpdateUser(input.Login, input.UserRole).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `"successfully update user role"`,
		},
		{
			name:        "Wrong Input",
			headerName:  authorizationHeader,
			headerValue: fmt.Sprintf("Bearer %s", token),
			inputBody:   `{"login": "login", "user_role": "avim"}`,
			inputUser: updateInput{
				Login:    "login",
				UserRole: "avim",
			},
			mockBehavior: func(r *mock_service.MockUser, input updateInput) {
				r.EXPECT().ParseToken(token).Return(1, nil)
				r.EXPECT().ValidateUser(1).Return(true, nil)
			},
			expectedStatusCode:   400,
			expectedResponseBody: `invalid state for required filed(s)`,
		},
		{
			name:        "Locked",
			headerName:  authorizationHeader,
			headerValue: fmt.Sprintf("Bearer %s", token),
			inputBody:   `{"login": "login", "user_role": "admin"}`,
			inputUser:   updateInput{},
			mockBehavior: func(r *mock_service.MockUser, input updateInput) {
				r.EXPECT().ParseToken(token).Return(1, nil)
				r.EXPECT().ValidateUser(1).Return(false, nil)
			},
			expectedStatusCode:   423,
			expectedResponseBody: `this function locked for current user`,
		},
		{
			name:        "Service Error",
			headerName:  authorizationHeader,
			headerValue: fmt.Sprintf("Bearer %s", token),
			inputBody:   `{"login": "login", "user_role": "admin"}`,
			inputUser: updateInput{
				Login:    "login",
				UserRole: "admin",
			},
			mockBehavior: func(r *mock_service.MockUser, input updateInput) {
				r.EXPECT().ParseToken(token).Return(1, nil)
				r.EXPECT().ValidateUser(1).Return(true, nil)
				r.EXPECT().UpdateUser(input.Login, input.UserRole).Return(errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `something went wrong`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockUser(c)
			test.mockBehavior(repo, test.inputUser)

			services := &service.Service{User: repo}
			handler := Router{service: services}
			handler.AddEndPoint("PUT", "/users", updateUser)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/users",
				bytes.NewBufferString(test.inputBody))
			req.Header.Set(test.headerName, test.headerValue)

			// Make Request
			handler.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, strings.ReplaceAll(w.Body.String(), "\n", ""))
		})
	}
}

func TestHandler_deleteUser(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockUser)

	const token = "fmekwfmw"

	tests := []struct {
		name                 string
		headerName           string
		headerValue          string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Ok",
			headerName:  authorizationHeader,
			headerValue: fmt.Sprintf("Bearer %s", token),
			mockBehavior: func(r *mock_service.MockUser) {
				r.EXPECT().ParseToken(token).Return(1, nil)
				r.EXPECT().DeleteUserById(1).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `"user successfully deleted"`,
		},
		{
			name:        "Wrong token",
			headerName:  authorizationHeader,
			headerValue: fmt.Sprintf("Bearer %s", token),
			mockBehavior: func(r *mock_service.MockUser) {
				r.EXPECT().ParseToken(token).Return(0, errors.New("invalid token"))
			},
			expectedStatusCode:   401,
			expectedResponseBody: `invalid token`,
		},
		{
			name:        "Service Error",
			headerName:  authorizationHeader,
			headerValue: fmt.Sprintf("Bearer %s", token),
			mockBehavior: func(r *mock_service.MockUser) {
				r.EXPECT().ParseToken(token).Return(1, nil)
				r.EXPECT().DeleteUserById(1).Return(errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `something went wrong`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockUser(c)
			test.mockBehavior(repo)

			services := &service.Service{User: repo}
			handler := Router{service: services}
			handler.AddEndPoint("DELETE", "/users", deleteUser)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/users",
				bytes.NewBufferString(""))
			req.Header.Set(test.headerName, test.headerValue)

			// Make Request
			handler.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, strings.ReplaceAll(w.Body.String(), "\n", ""))
		})
	}
}

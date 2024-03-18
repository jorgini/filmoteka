package handlers

import (
	"bytes"
	"errors"
	"github.com/jorgini/filmoteka"
	"github.com/jorgini/filmoteka/service"
	"github.com/jorgini/filmoteka/service/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestRouter_createNewActor(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r1 *mock_service.MockActor, r2 *mock_service.MockUser, actor filmoteka.Actor)

	var (
		date        = time.Time{}.AddDate(2022, 7, 10)
		headerName  = "Authorization"
		headerValue = "Bearer test"
		token       = "test"
		userId      = 1
	)

	tests := []struct {
		name                 string
		inputBody            string
		inputActor           filmoteka.Actor
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"name": "name", "surname": "surname", "sex": "male", "birthday" : "11-08-2023"}`,
			inputActor: filmoteka.Actor{
				Name:     "name",
				Surname:  "surname",
				Sex:      "male",
				Birthday: (*filmoteka.Date)(&date),
			},
			mockBehavior: func(r1 *mock_service.MockActor, r2 *mock_service.MockUser, actor filmoteka.Actor) {
				r2.EXPECT().ParseToken(token).Return(userId, nil)
				r2.EXPECT().ValidateUser(userId).Return(true, nil)
				r1.EXPECT().CreateActor(actor).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `"successfully create actor with id 1"`,
		},
		{
			name:       "Wrong Input",
			inputBody:  `{"name": "name"}`,
			inputActor: filmoteka.Actor{},
			mockBehavior: func(r1 *mock_service.MockActor, r2 *mock_service.MockUser, actor filmoteka.Actor) {
				r2.EXPECT().ParseToken(token).Return(userId, nil)
				r2.EXPECT().ValidateUser(userId).Return(true, nil)
			},
			expectedStatusCode:   400,
			expectedResponseBody: `invalid state for required field(s)`,
		},
		{
			name:      "Locked",
			inputBody: `{"name": "name", "surname": "surname", "sex": "female", "birthday": "11-08-2023"}`,
			inputActor: filmoteka.Actor{
				Name:     "name",
				Surname:  "surname",
				Sex:      "female",
				Birthday: (*filmoteka.Date)(&date),
			},
			mockBehavior: func(r1 *mock_service.MockActor, r2 *mock_service.MockUser, actor filmoteka.Actor) {
				r2.EXPECT().ParseToken(token).Return(userId, nil)
				r2.EXPECT().ValidateUser(userId).Return(false, nil)
			},
			expectedStatusCode:   423,
			expectedResponseBody: `this function locked for current user`,
		},
		{
			name:      "Service Error",
			inputBody: `{"name": "name", "surname": "surname", "sex": "female", "birthday": "11-08-2023"}`,
			inputActor: filmoteka.Actor{
				Name:     "name",
				Surname:  "surname",
				Sex:      "female",
				Birthday: (*filmoteka.Date)(&date),
			},
			mockBehavior: func(r1 *mock_service.MockActor, r2 *mock_service.MockUser, actor filmoteka.Actor) {
				r2.EXPECT().ParseToken(token).Return(userId, nil)
				r2.EXPECT().ValidateUser(userId).Return(true, nil)
				r1.EXPECT().CreateActor(actor).Return(0, errors.New("something went wrong"))
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

			repo1 := mock_service.NewMockActor(c)
			repo2 := mock_service.NewMockUser(c)
			test.mockBehavior(repo1, repo2, test.inputActor)

			services := &service.Service{Actor: repo1, User: repo2}
			handler := Router{service: services}
			handler.AddEndPoint("POST", "/actors", createNewActor)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/actors",
				bytes.NewBufferString(test.inputBody))
			req.Header.Set(headerName, headerValue)

			// Make Request
			handler.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, strings.ReplaceAll(w.Body.String(), "\n", ""))
		})
	}
}

func TestRouter_updateActor(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r1 *mock_service.MockActor, r2 *mock_service.MockUser, actor filmoteka.UpdateActorInput)

	var (
		date          = time.Time{}.AddDate(2022, 7, 10)
		headerName    = "Authorization"
		headerValue   = "Bearer test"
		token         = "test"
		userId        = 1
		actorId       = 1
		surnameUpdate = "surname"
	)

	tests := []struct {
		name                 string
		paramsName           string
		paramsValue          string
		inputBody            string
		inputActor           filmoteka.UpdateActorInput
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Ok",
			paramsName:  "id",
			paramsValue: "1",
			inputBody:   `{"surname": "surname", "birthday": "11-08-2023"}`,
			inputActor: filmoteka.UpdateActorInput{
				Id:       &actorId,
				Surname:  &surnameUpdate,
				Birthday: (*filmoteka.Date)(&date),
			},
			mockBehavior: func(r1 *mock_service.MockActor, r2 *mock_service.MockUser, actor filmoteka.UpdateActorInput) {
				r2.EXPECT().ParseToken(token).Return(userId, nil)
				r2.EXPECT().ValidateUser(userId).Return(true, nil)
				r1.EXPECT().UpdateActor(actor).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `"successfully update"`,
		},
		{
			name:        "Wrong Params",
			paramsName:  "id",
			paramsValue: "fafmek",
			inputBody:   `{"surname": "surname", "birthday": "11-08-2023"}`,
			inputActor: filmoteka.UpdateActorInput{
				Id:       &actorId,
				Surname:  &surnameUpdate,
				Birthday: (*filmoteka.Date)(&date),
			},
			mockBehavior: func(r1 *mock_service.MockActor, r2 *mock_service.MockUser, actor filmoteka.UpdateActorInput) {
				r2.EXPECT().ParseToken(token).Return(userId, nil)
				r2.EXPECT().ValidateUser(userId).Return(true, nil)
			},
			expectedStatusCode:   400,
			expectedResponseBody: `no id specified for update actor`,
		},
		{
			name:        "Wrong Input",
			paramsName:  "id",
			paramsValue: "1",
			inputBody:   `{}`,
			inputActor:  filmoteka.UpdateActorInput{},
			mockBehavior: func(r1 *mock_service.MockActor, r2 *mock_service.MockUser, actor filmoteka.UpdateActorInput) {
				r2.EXPECT().ParseToken(token).Return(userId, nil)
				r2.EXPECT().ValidateUser(userId).Return(true, nil)
			},
			expectedStatusCode:   400,
			expectedResponseBody: `invalid state for required filed id to update actor`,
		},
		{
			name:        "Locked",
			paramsName:  "id",
			paramsValue: "1",
			inputBody:   `{"surname": "surname", "birthday": "11-08-2023"}`,
			inputActor: filmoteka.UpdateActorInput{
				Id:       &actorId,
				Surname:  &surnameUpdate,
				Birthday: (*filmoteka.Date)(&date),
			},
			mockBehavior: func(r1 *mock_service.MockActor, r2 *mock_service.MockUser, actor filmoteka.UpdateActorInput) {
				r2.EXPECT().ParseToken(token).Return(userId, nil)
				r2.EXPECT().ValidateUser(userId).Return(false, nil)
			},
			expectedStatusCode:   423,
			expectedResponseBody: `this function locked for current user`,
		},
		{
			name:        "Service Error",
			paramsName:  "id",
			paramsValue: "1",
			inputBody:   `{"surname": "surname", "birthday": "11-08-2023"}`,
			inputActor: filmoteka.UpdateActorInput{
				Id:       &actorId,
				Surname:  &surnameUpdate,
				Birthday: (*filmoteka.Date)(&date),
			},
			mockBehavior: func(r1 *mock_service.MockActor, r2 *mock_service.MockUser, actor filmoteka.UpdateActorInput) {
				r2.EXPECT().ParseToken(token).Return(userId, nil)
				r2.EXPECT().ValidateUser(userId).Return(true, nil)
				r1.EXPECT().UpdateActor(actor).Return(errors.New("something went wrong"))
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

			repo1 := mock_service.NewMockActor(c)
			repo2 := mock_service.NewMockUser(c)
			test.mockBehavior(repo1, repo2, test.inputActor)

			services := &service.Service{Actor: repo1, User: repo2}
			handler := Router{service: services}
			handler.AddEndPoint("PUT", "/actors", updateActor)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/actors?"+test.paramsName+"="+test.paramsValue,
				bytes.NewBufferString(test.inputBody))
			req.Header.Set(headerName, headerValue)

			// Make Request
			handler.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, strings.ReplaceAll(w.Body.String(), "\n", ""))
		})
	}
}

func TestRouter_getActorsList(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockActor, page int)

	var (
		date  = time.Time{}.AddDate(2022, 7, 10)
		limit = 10
		actor = []filmoteka.ActorListItem{{
			Actor: filmoteka.Actor{
				Name:     "test",
				Surname:  "test",
				Sex:      "female",
				Birthday: (*filmoteka.Date)(&date),
			},
			Films: []filmoteka.Film{{
				Id:          1,
				Title:       "test",
				Description: "",
				IssueDate:   (*filmoteka.Date)(&date),
				Rating:      5,
			}},
		}}
	)

	tests := []struct {
		name                 string
		paramsName           string
		paramsValue          string
		page                 int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Ok",
			paramsName:  "page",
			paramsValue: "1",
			page:        1,
			mockBehavior: func(r *mock_service.MockActor, page int) {
				r.EXPECT().GetActorsList(page, limit).Return(actor, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"Actor":{"id":0,"name":"test","surname":"test","sex":"female","birthday":"11-08-2023"},"Films":[{"id":1,"title":"test","description":"","issue_date":"11-08-2023","rating":5}]}`,
		},
		{
			name:                 "Wrong Params",
			paramsName:           "paage",
			paramsValue:          "1",
			page:                 1,
			mockBehavior:         func(r *mock_service.MockActor, page int) {},
			expectedStatusCode:   400,
			expectedResponseBody: `page not selected`,
		},
		{
			name:                 "Wrong Input",
			paramsName:           "page",
			paramsValue:          "-1",
			page:                 -1,
			mockBehavior:         func(r *mock_service.MockActor, page int) {},
			expectedStatusCode:   400,
			expectedResponseBody: `page out of bounds`,
		},
		{
			name:        "Over page",
			paramsName:  "page",
			paramsValue: "2",
			page:        2,
			mockBehavior: func(r *mock_service.MockActor, page int) {
				r.EXPECT().GetActorsList(page, limit).Return([]filmoteka.ActorListItem{}, nil)
			},
			expectedStatusCode:   400,
			expectedResponseBody: `page out of bounds`,
		},
		{
			name:        "Service Error",
			paramsName:  "page",
			paramsValue: "1",
			page:        1,
			mockBehavior: func(r *mock_service.MockActor, page int) {
				r.EXPECT().GetActorsList(page, limit).Return(nil, errors.New("something went wrong"))
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

			repo := mock_service.NewMockActor(c)
			test.mockBehavior(repo, test.page)

			services := &service.Service{Actor: repo}
			handler := Router{service: services}
			handler.AddEndPoint("GET", "/actors/list", getActorsList)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/actors/list?"+test.paramsName+"="+test.paramsValue,
				bytes.NewBufferString(""))

			// Make Request
			handler.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, strings.ReplaceAll(w.Body.String(), "\n", ""))
		})
	}
}

func TestRouter_getActorById(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockActor)

	var (
		date    = time.Time{}.AddDate(2022, 7, 10)
		actorId = 1
		actor   = filmoteka.ActorListItem{
			Actor: filmoteka.Actor{
				Id:       actorId,
				Name:     "test",
				Surname:  "test",
				Sex:      "female",
				Birthday: (*filmoteka.Date)(&date),
			},
			Films: []filmoteka.Film{{
				Id:          1,
				Title:       "test",
				Description: "",
				IssueDate:   (*filmoteka.Date)(&date),
				Rating:      5,
			}},
		}
	)

	tests := []struct {
		name                 string
		paramsName           string
		paramsValue          string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Ok",
			paramsName:  "id",
			paramsValue: "1",
			mockBehavior: func(r *mock_service.MockActor) {
				r.EXPECT().GetActorById(actorId).Return(actor, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"Actor":{"id":1,"name":"test","surname":"test","sex":"female","birthday":"11-08-2023"},"Films":[{"id":1,"title":"test","description":"","issue_date":"11-08-2023","rating":5}]}`,
		},
		{
			name:                 "Wrong Params",
			paramsName:           "idd",
			paramsValue:          "1",
			mockBehavior:         func(r *mock_service.MockActor) {},
			expectedStatusCode:   400,
			expectedResponseBody: `id not selected`,
		},
		{
			name:                 "Wrong Input",
			paramsName:           "id",
			paramsValue:          "-1",
			mockBehavior:         func(r *mock_service.MockActor) {},
			expectedStatusCode:   400,
			expectedResponseBody: `id out of bounds`,
		},
		{
			name:        "Service Error",
			paramsName:  "id",
			paramsValue: "1",
			mockBehavior: func(r *mock_service.MockActor) {
				r.EXPECT().GetActorById(actorId).Return(filmoteka.ActorListItem{}, errors.New("something went wrong"))
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

			repo := mock_service.NewMockActor(c)
			test.mockBehavior(repo)

			services := &service.Service{Actor: repo}
			handler := Router{service: services}
			handler.AddEndPoint("GET", "/actors", getActorById)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/actors?"+test.paramsName+"="+test.paramsValue,
				bytes.NewBufferString(""))

			// Make Request
			handler.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, strings.ReplaceAll(w.Body.String(), "\n", ""))
		})
	}
}

func TestRouter_searchActor(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockActor, fragment filmoteka.ActorSearchFragment, page int)

	var (
		date  = time.Time{}.AddDate(2022, 7, 10)
		limit = 10
		actor = []filmoteka.ActorListItem{{
			Actor: filmoteka.Actor{
				Name:     "test",
				Surname:  "test",
				Sex:      "female",
				Birthday: (*filmoteka.Date)(&date),
			},
			Films: []filmoteka.Film{{
				Id:          1,
				Title:       "test",
				Description: "",
				IssueDate:   (*filmoteka.Date)(&date),
				Rating:      5,
			}},
		}}
		name    = "te"
		surname = "t"
	)

	tests := []struct {
		name                 string
		paramsName           string
		paramsValue          string
		page                 int
		inputBody            string
		fragment             filmoteka.ActorSearchFragment
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Ok",
			paramsName:  "page",
			paramsValue: "1",
			page:        1,
			inputBody:   `{"name":"te", "surname": "t"}`,
			fragment: filmoteka.ActorSearchFragment{
				Name:    &name,
				Surname: &surname,
			},
			mockBehavior: func(r *mock_service.MockActor, fragment filmoteka.ActorSearchFragment, page int) {
				r.EXPECT().SearchActor(page, limit, fragment).Return(actor, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"Actor":{"id":0,"name":"test","surname":"test","sex":"female","birthday":"11-08-2023"},"Films":[{"id":1,"title":"test","description":"","issue_date":"11-08-2023","rating":5}]}`,
		},
		{
			name:                 "Wrong Params",
			paramsName:           "paage",
			paramsValue:          "1",
			page:                 1,
			mockBehavior:         func(r *mock_service.MockActor, fragment filmoteka.ActorSearchFragment, page int) {},
			expectedStatusCode:   400,
			expectedResponseBody: `page not selected`,
		},
		{
			name:        "Wrong Input",
			paramsName:  "page",
			paramsValue: "1",
			page:        1,
			inputBody:   `{}`,
			fragment: filmoteka.ActorSearchFragment{
				Name:    nil,
				Surname: nil,
			},
			mockBehavior:         func(r *mock_service.MockActor, fragment filmoteka.ActorSearchFragment, page int) {},
			expectedStatusCode:   400,
			expectedResponseBody: `parameters for search not specified`,
		},
		{
			name:        "Bad request",
			paramsName:  "page",
			paramsValue: "1",
			page:        1,
			inputBody:   `{"name":"te", "surname": "t"}`,
			fragment: filmoteka.ActorSearchFragment{
				Name:    &name,
				Surname: &surname,
			},
			mockBehavior: func(r *mock_service.MockActor, fragment filmoteka.ActorSearchFragment, page int) {
				r.EXPECT().SearchActor(page, limit, fragment).Return([]filmoteka.ActorListItem{}, nil)
			},
			expectedStatusCode:   400,
			expectedResponseBody: `nothing find`,
		},
		{
			name:        "Service Error",
			paramsName:  "page",
			paramsValue: "1",
			page:        1,
			inputBody:   `{"name":"te", "surname": "t"}`,
			fragment: filmoteka.ActorSearchFragment{
				Name:    &name,
				Surname: &surname,
			},
			mockBehavior: func(r *mock_service.MockActor, fragment filmoteka.ActorSearchFragment, page int) {
				r.EXPECT().SearchActor(page, limit, fragment).Return(nil, errors.New("something went wrong"))
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

			repo := mock_service.NewMockActor(c)
			test.mockBehavior(repo, test.fragment, test.page)

			services := &service.Service{Actor: repo}
			handler := Router{service: services}
			handler.AddEndPoint("GET", "/actors/search", searchActor)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/actors/search?"+test.paramsName+"="+test.paramsValue,
				bytes.NewBufferString(test.inputBody))

			// Make Request
			handler.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, strings.ReplaceAll(w.Body.String(), "\n", ""))
		})
	}
}

func TestRouter_deleteActor(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r1 *mock_service.MockActor, r2 *mock_service.MockUser)

	var (
		headerName  = "Authorization"
		headerValue = "Bearer test"
		token       = "test"
		userId      = 1
		actorId     = 1
	)

	tests := []struct {
		name                 string
		paramsName           string
		paramsValue          string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Ok",
			paramsName:  "id",
			paramsValue: "1",
			mockBehavior: func(r1 *mock_service.MockActor, r2 *mock_service.MockUser) {
				r2.EXPECT().ParseToken(token).Return(userId, nil)
				r2.EXPECT().ValidateUser(userId).Return(true, nil)
				r1.EXPECT().DeleteActorById(actorId).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `"successful delete"`,
		},
		{
			name:        "Wrong Params",
			paramsName:  "id",
			paramsValue: "fafmek",
			mockBehavior: func(r1 *mock_service.MockActor, r2 *mock_service.MockUser) {
				r2.EXPECT().ParseToken(token).Return(userId, nil)
				r2.EXPECT().ValidateUser(userId).Return(true, nil)
			},
			expectedStatusCode:   400,
			expectedResponseBody: `id doesnt specified to delete actor`,
		},
		{
			name:        "Locked",
			paramsName:  "id",
			paramsValue: "1",
			mockBehavior: func(r1 *mock_service.MockActor, r2 *mock_service.MockUser) {
				r2.EXPECT().ParseToken(token).Return(userId, nil)
				r2.EXPECT().ValidateUser(userId).Return(false, nil)
			},
			expectedStatusCode:   423,
			expectedResponseBody: `this function locked for current user`,
		},
		{
			name:        "Service Error",
			paramsName:  "id",
			paramsValue: "1",
			mockBehavior: func(r1 *mock_service.MockActor, r2 *mock_service.MockUser) {
				r2.EXPECT().ParseToken(token).Return(userId, nil)
				r2.EXPECT().ValidateUser(userId).Return(true, nil)
				r1.EXPECT().DeleteActorById(actorId).Return(errors.New("something went wrong"))
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

			repo1 := mock_service.NewMockActor(c)
			repo2 := mock_service.NewMockUser(c)
			test.mockBehavior(repo1, repo2)

			services := &service.Service{Actor: repo1, User: repo2}
			handler := Router{service: services}
			handler.AddEndPoint("DELETE", "/actors", deleteActor)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/actors?"+test.paramsName+"="+test.paramsValue,
				bytes.NewBufferString(""))
			req.Header.Set(headerName, headerValue)

			// Make Request
			handler.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, strings.ReplaceAll(w.Body.String(), "\n", ""))
		})
	}
}

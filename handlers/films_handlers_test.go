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

func TestRouter_createNewFilm(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r1 *mock_service.MockFilm, r2 *mock_service.MockUser, film filmoteka.InputFilm)

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
		inputFilm            filmoteka.InputFilm
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"title": "title", "description": "test", "issue_date": "11-08-2023", "rating" : 5, "Cast": []}`,
			inputFilm: filmoteka.InputFilm{
				Film: filmoteka.Film{
					Title:       "title",
					Description: "test",
					IssueDate:   (*filmoteka.Date)(&date),
					Rating:      5,
				},
				Cast: nil,
			},
			mockBehavior: func(r1 *mock_service.MockFilm, r2 *mock_service.MockUser, film filmoteka.InputFilm) {
				r2.EXPECT().ParseToken(token).Return(userId, nil)
				r2.EXPECT().ValidateUser(userId).Return(true, nil)
				r1.EXPECT().CreateFilm(film).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `"successfully create film with id 1"`,
		},
		{
			name:      "Ok with cast",
			inputBody: `{"title": "title", "description": "test", "issue_date": "11-08-2023", "rating" : 5, "Cast": [{"name": "name", "surname":"surname"}]}`,
			inputFilm: filmoteka.InputFilm{
				Film: filmoteka.Film{
					Title:       "title",
					Description: "test",
					IssueDate:   (*filmoteka.Date)(&date),
					Rating:      5,
				},
				Cast: filmoteka.Cast{filmoteka.InputActor{Name: "name", Surname: "surname"}},
			},
			mockBehavior: func(r1 *mock_service.MockFilm, r2 *mock_service.MockUser, film filmoteka.InputFilm) {
				r2.EXPECT().ParseToken(token).Return(userId, nil)
				r2.EXPECT().ValidateUser(userId).Return(true, nil)
				r1.EXPECT().CreateFilm(film).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `"successfully create film with id 1"`,
		},
		{
			name:      "Wrong Input",
			inputBody: `{"title": "title"}`,
			inputFilm: filmoteka.InputFilm{},
			mockBehavior: func(r1 *mock_service.MockFilm, r2 *mock_service.MockUser, film filmoteka.InputFilm) {
				r2.EXPECT().ParseToken(token).Return(userId, nil)
				r2.EXPECT().ValidateUser(userId).Return(true, nil)
			},
			expectedStatusCode:   400,
			expectedResponseBody: `invalid state for required field(s)`,
		},
		{
			name:      "Locked",
			inputBody: `{"title": "title", "description": "", "issue_date": "11-08-2023", "rating" : "5", "cast": []}`,
			inputFilm: filmoteka.InputFilm{
				Film: filmoteka.Film{
					Title:       "title",
					Description: "",
					IssueDate:   (*filmoteka.Date)(&date),
					Rating:      5,
				},
				Cast: filmoteka.Cast{},
			},
			mockBehavior: func(r1 *mock_service.MockFilm, r2 *mock_service.MockUser, film filmoteka.InputFilm) {
				r2.EXPECT().ParseToken(token).Return(userId, nil)
				r2.EXPECT().ValidateUser(userId).Return(false, nil)
			},
			expectedStatusCode:   423,
			expectedResponseBody: `this function locked for current user`,
		},
		{
			name:      "Service Error",
			inputBody: `{"title": "title", "description": "", "issue_date": "11-08-2023", "rating" : 5, "cast": []}`,
			inputFilm: filmoteka.InputFilm{
				Film: filmoteka.Film{
					Title:       "title",
					Description: "",
					IssueDate:   (*filmoteka.Date)(&date),
					Rating:      5,
				},
				Cast: nil,
			},
			mockBehavior: func(r1 *mock_service.MockFilm, r2 *mock_service.MockUser, film filmoteka.InputFilm) {
				r2.EXPECT().ParseToken(token).Return(userId, nil)
				r2.EXPECT().ValidateUser(userId).Return(true, nil)
				r1.EXPECT().CreateFilm(film).Return(0, errors.New("something went wrong"))
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

			repo1 := mock_service.NewMockFilm(c)
			repo2 := mock_service.NewMockUser(c)
			test.mockBehavior(repo1, repo2, test.inputFilm)

			services := &service.Service{Film: repo1, User: repo2}
			handler := Router{service: services}
			handler.AddEndPoint("POST", "/films", createNewFilm)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/films",
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

func TestRouter_updateFilm(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r1 *mock_service.MockFilm, r2 *mock_service.MockUser, actor filmoteka.UpdateFilmInput)

	var (
		date        = time.Time{}.AddDate(2022, 7, 10)
		headerName  = "Authorization"
		headerValue = "Bearer test"
		token       = "test"
		userId      = 1
		filmId      = 1
		titleUpdate = "title"
	)

	tests := []struct {
		name                 string
		paramsName           string
		paramsValue          string
		inputBody            string
		inputFilm            filmoteka.UpdateFilmInput
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Ok",
			paramsName:  "id",
			paramsValue: "1",
			inputBody:   `{"title": "title", "issue_date": "11-08-2023"}`,
			inputFilm: filmoteka.UpdateFilmInput{
				Id:        &filmId,
				Title:     &titleUpdate,
				IssueDate: (*filmoteka.Date)(&date),
			},
			mockBehavior: func(r1 *mock_service.MockFilm, r2 *mock_service.MockUser, film filmoteka.UpdateFilmInput) {
				r2.EXPECT().ParseToken(token).Return(userId, nil)
				r2.EXPECT().ValidateUser(userId).Return(true, nil)
				r1.EXPECT().UpdateFilm(film).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `"successfully update"`,
		},
		{
			name:        "Ok with cast",
			paramsName:  "id",
			paramsValue: "1",
			inputBody:   `{"title": "title", "issue_date": "11-08-2023", "Cast": [{"name": "name", "surname": "surname"}]}`,
			inputFilm: filmoteka.UpdateFilmInput{
				Id:        &filmId,
				Title:     &titleUpdate,
				IssueDate: (*filmoteka.Date)(&date),
				Actors: &filmoteka.Cast{{
					Name:    "name",
					Surname: "surname",
				}},
			},
			mockBehavior: func(r1 *mock_service.MockFilm, r2 *mock_service.MockUser, film filmoteka.UpdateFilmInput) {
				r2.EXPECT().ParseToken(token).Return(userId, nil)
				r2.EXPECT().ValidateUser(userId).Return(true, nil)
				r1.EXPECT().UpdateFilm(film).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `"successfully update"`,
		},
		{
			name:        "Wrong Params",
			paramsName:  "id",
			paramsValue: "fafmek",
			inputBody:   `{"title": "title", "issue_date": "11-08-2023"}`,
			inputFilm: filmoteka.UpdateFilmInput{
				Id:        &filmId,
				Title:     &titleUpdate,
				IssueDate: (*filmoteka.Date)(&date),
			},
			mockBehavior: func(r1 *mock_service.MockFilm, r2 *mock_service.MockUser, film filmoteka.UpdateFilmInput) {
				r2.EXPECT().ParseToken(token).Return(userId, nil)
				r2.EXPECT().ValidateUser(userId).Return(true, nil)
			},
			expectedStatusCode:   400,
			expectedResponseBody: `no id specified to update film`,
		},
		{
			name:        "Wrong Input",
			paramsName:  "id",
			paramsValue: "1",
			inputBody:   `{"title": "title", "issue_date": "11-08-2023", "Cast": [{"name": "name"}]}`,
			inputFilm:   filmoteka.UpdateFilmInput{},
			mockBehavior: func(r1 *mock_service.MockFilm, r2 *mock_service.MockUser, film filmoteka.UpdateFilmInput) {
				r2.EXPECT().ParseToken(token).Return(userId, nil)
				r2.EXPECT().ValidateUser(userId).Return(true, nil)
			},
			expectedStatusCode:   400,
			expectedResponseBody: `invalid state for required filed(s)`,
		},
		{
			name:        "Locked",
			paramsName:  "id",
			paramsValue: "1",
			inputBody:   `{"title": "title", "issue_date": "11-08-2023"}`,
			inputFilm:   filmoteka.UpdateFilmInput{},
			mockBehavior: func(r1 *mock_service.MockFilm, r2 *mock_service.MockUser, film filmoteka.UpdateFilmInput) {
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
			inputBody:   `{"title": "title", "issue_date": "11-08-2023"}`,
			inputFilm: filmoteka.UpdateFilmInput{
				Id:        &filmId,
				Title:     &titleUpdate,
				IssueDate: (*filmoteka.Date)(&date),
			},
			mockBehavior: func(r1 *mock_service.MockFilm, r2 *mock_service.MockUser, film filmoteka.UpdateFilmInput) {
				r2.EXPECT().ParseToken(token).Return(userId, nil)
				r2.EXPECT().ValidateUser(userId).Return(true, nil)
				r1.EXPECT().UpdateFilm(film).Return(errors.New("something went wrong"))
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

			repo1 := mock_service.NewMockFilm(c)
			repo2 := mock_service.NewMockUser(c)
			test.mockBehavior(repo1, repo2, test.inputFilm)

			services := &service.Service{Film: repo1, User: repo2}
			handler := Router{service: services}
			handler.AddEndPoint("PUT", "/films", updateFilm)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/films?"+test.paramsName+"="+test.paramsValue,
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

func TestRouter_getSortedFilmsList(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockFilm, sortBy string, page int)

	var (
		date  = time.Time{}.AddDate(2022, 7, 10)
		limit = 10
		film  = []filmoteka.InputFilm{{
			Film: filmoteka.Film{
				Title:       "test",
				Description: "test",
				IssueDate:   (*filmoteka.Date)(&date),
				Rating:      5,
			},
			Cast: []filmoteka.InputActor{{
				Name:    "name",
				Surname: "surname",
			}},
		}}
	)

	tests := []struct {
		name                 string
		params               string
		page                 int
		sortBy               string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:   "Ok default",
			params: "page=1",
			page:   1,
			sortBy: "rating",
			mockBehavior: func(r *mock_service.MockFilm, sortBy string, page int) {
				r.EXPECT().GetSortedFilmList(sortBy, page, limit).Return(film, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `[{"id":0,"title":"test","description":"test","issue_date":"11-08-2023","rating":5,"Cast":[{"name":"name","surname":"surname"}]}]`,
		},
		{
			name:   "Ok with sort",
			params: "sort_by=issue_date&page=1",
			page:   1,
			sortBy: "issue_date",
			mockBehavior: func(r *mock_service.MockFilm, sortBy string, page int) {
				r.EXPECT().GetSortedFilmList(sortBy, page, limit).Return(film, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `[{"id":0,"title":"test","description":"test","issue_date":"11-08-2023","rating":5,"Cast":[{"name":"name","surname":"surname"}]}]`,
		},
		{
			name:                 "Wrong Params",
			params:               "paage=1",
			page:                 1,
			mockBehavior:         func(r *mock_service.MockFilm, sortBy string, page int) {},
			expectedStatusCode:   400,
			expectedResponseBody: `no page specified for sorted list`,
		},
		{
			name:                 "Wrong Input page",
			params:               "page=-1",
			page:                 -1,
			mockBehavior:         func(r *mock_service.MockFilm, sortBy string, page int) {},
			expectedStatusCode:   400,
			expectedResponseBody: `page out of bounds`,
		},
		{
			name:                 "Wrong Input sort",
			params:               "sort_by=smt&page=1",
			page:                 1,
			sortBy:               "smt",
			mockBehavior:         func(r *mock_service.MockFilm, sortBy string, page int) {},
			expectedStatusCode:   400,
			expectedResponseBody: `invalid parameter to sort films list`,
		},
		{
			name:   "Over page",
			params: "page=2",
			page:   2,
			sortBy: "rating",
			mockBehavior: func(r *mock_service.MockFilm, sortBy string, page int) {
				r.EXPECT().GetSortedFilmList(sortBy, page, limit).Return([]filmoteka.InputFilm{}, nil)
			},
			expectedStatusCode:   400,
			expectedResponseBody: `page out of bounds`,
		},
		{
			name:   "Service Error",
			params: "page=1",
			page:   1,
			sortBy: "rating",
			mockBehavior: func(r *mock_service.MockFilm, sortBy string, page int) {
				r.EXPECT().GetSortedFilmList(sortBy, page, limit).Return(nil, errors.New("something went wrong"))
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

			repo := mock_service.NewMockFilm(c)
			test.mockBehavior(repo, test.sortBy, test.page)

			services := &service.Service{Film: repo}
			handler := Router{service: services}
			handler.AddEndPoint("GET", "/films/list", getSortedFilmList)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/films/list?"+test.params, bytes.NewBufferString(""))

			// Make Request
			handler.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, strings.ReplaceAll(w.Body.String(), "\n", ""))
		})
	}
}

func TestRouter_getCurrentFilm(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockFilm)

	var (
		date   = time.Time{}.AddDate(2022, 7, 10)
		filmId = 1
		film   = filmoteka.InputFilm{
			Film: filmoteka.Film{
				Title:       "test",
				Description: "test",
				IssueDate:   (*filmoteka.Date)(&date),
				Rating:      5,
			},
			Cast: []filmoteka.InputActor{{
				Name:    "name",
				Surname: "surname",
			}},
		}
	)

	tests := []struct {
		name                 string
		params               string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:   "Ok",
			params: "id=1",
			mockBehavior: func(r *mock_service.MockFilm) {
				r.EXPECT().GetCurFilm(filmId).Return(film, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":0,"title":"test","description":"test","issue_date":"11-08-2023","rating":5,"Cast":[{"name":"name","surname":"surname"}]}`,
		},
		{
			name:                 "Wrong Params",
			params:               "idd=1",
			mockBehavior:         func(r *mock_service.MockFilm) {},
			expectedStatusCode:   400,
			expectedResponseBody: `id for get film not specified`,
		},
		{
			name:                 "Wrong Input",
			params:               "id=-1",
			mockBehavior:         func(r *mock_service.MockFilm) {},
			expectedStatusCode:   400,
			expectedResponseBody: `id out of bounds`,
		},
		{
			name:   "Service Error",
			params: "id=1",
			mockBehavior: func(r *mock_service.MockFilm) {
				r.EXPECT().GetCurFilm(filmId).Return(filmoteka.InputFilm{}, errors.New("something went wrong"))
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

			repo := mock_service.NewMockFilm(c)
			test.mockBehavior(repo)

			services := &service.Service{Film: repo}
			handler := Router{service: services}
			handler.AddEndPoint("GET", "/films", getCurrentFilm)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/films?"+test.params,
				bytes.NewBufferString(""))

			// Make Request
			handler.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, strings.ReplaceAll(w.Body.String(), "\n", ""))
		})
	}
}

func TestRouter_getSearchFilmList(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockFilm, fragment filmoteka.FilmSearchFragment, page int)

	var (
		date  = time.Time{}.AddDate(2022, 7, 10)
		limit = 10
		film  = []filmoteka.InputFilm{{
			Film: filmoteka.Film{
				Title:       "test",
				Description: "test",
				IssueDate:   (*filmoteka.Date)(&date),
				Rating:      5,
			},
			Cast: []filmoteka.InputActor{{
				Name:    "name",
				Surname: "surname",
			}},
		}}
		title   = "te"
		name    = "na"
		surname = "s"
	)

	tests := []struct {
		name                 string
		params               string
		page                 int
		inputBody            string
		fragment             filmoteka.FilmSearchFragment
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok search by title",
			params:    "page=1",
			page:      1,
			inputBody: `{"title": "te"}`,
			fragment: filmoteka.FilmSearchFragment{
				Title:   &title,
				Name:    nil,
				Surname: nil,
			},
			mockBehavior: func(r *mock_service.MockFilm, fragment filmoteka.FilmSearchFragment, page int) {
				r.EXPECT().GetSearchFilmList(page, limit, fragment).Return(film, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `[{"id":0,"title":"test","description":"test","issue_date":"11-08-2023","rating":5,"Cast":[{"name":"name","surname":"surname"}]}]`,
		},
		{
			name:      "Ok search by actor",
			params:    "page=1",
			page:      1,
			inputBody: `{"name": "na", "surname": "s"}`,
			fragment: filmoteka.FilmSearchFragment{
				Name:    &name,
				Surname: &surname,
			},
			mockBehavior: func(r *mock_service.MockFilm, fragment filmoteka.FilmSearchFragment, page int) {
				r.EXPECT().GetSearchFilmList(page, limit, fragment).Return(film, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `[{"id":0,"title":"test","description":"test","issue_date":"11-08-2023","rating":5,"Cast":[{"name":"name","surname":"surname"}]}]`,
		},
		{
			name:                 "Wrong Params",
			params:               "paage=1",
			page:                 1,
			mockBehavior:         func(r *mock_service.MockFilm, fragment filmoteka.FilmSearchFragment, page int) {},
			expectedStatusCode:   400,
			expectedResponseBody: `no page specified for search list`,
		},
		{
			name:                 "Wrong Input",
			params:               "page=1",
			page:                 1,
			inputBody:            `{}`,
			fragment:             filmoteka.FilmSearchFragment{},
			mockBehavior:         func(r *mock_service.MockFilm, fragment filmoteka.FilmSearchFragment, page int) {},
			expectedStatusCode:   400,
			expectedResponseBody: `parameters for search not specified`,
		},
		{
			name:      "Bad request",
			params:    "page=1",
			page:      1,
			inputBody: `{"name":"na", "surname": "s"}`,
			fragment: filmoteka.FilmSearchFragment{
				Name:    &name,
				Surname: &surname,
			},
			mockBehavior: func(r *mock_service.MockFilm, fragment filmoteka.FilmSearchFragment, page int) {
				r.EXPECT().GetSearchFilmList(page, limit, fragment).Return([]filmoteka.InputFilm{}, nil)
			},
			expectedStatusCode:   400,
			expectedResponseBody: `nothing find`,
		},
		{
			name:      "Service Error",
			params:    "page=1",
			page:      1,
			inputBody: `{"name":"na", "surname": "s"}`,
			fragment: filmoteka.FilmSearchFragment{
				Name:    &name,
				Surname: &surname,
			},
			mockBehavior: func(r *mock_service.MockFilm, fragment filmoteka.FilmSearchFragment, page int) {
				r.EXPECT().GetSearchFilmList(page, limit, fragment).Return(nil, errors.New("something went wrong"))
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

			repo := mock_service.NewMockFilm(c)
			test.mockBehavior(repo, test.fragment, test.page)

			services := &service.Service{Film: repo}
			handler := Router{service: services}
			handler.AddEndPoint("GET", "/films/search", getSearchFilmList)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/films/search?"+test.params,
				bytes.NewBufferString(test.inputBody))

			// Make Request
			handler.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, strings.ReplaceAll(w.Body.String(), "\n", ""))
		})
	}
}

func TestRouter_deleteFilm(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r1 *mock_service.MockFilm, r2 *mock_service.MockUser)

	var (
		headerName  = "Authorization"
		headerValue = "Bearer test"
		token       = "test"
		userId      = 1
		actorId     = 1
	)

	tests := []struct {
		name                 string
		params               string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:   "Ok",
			params: "id=1",
			mockBehavior: func(r1 *mock_service.MockFilm, r2 *mock_service.MockUser) {
				r2.EXPECT().ParseToken(token).Return(userId, nil)
				r2.EXPECT().ValidateUser(userId).Return(true, nil)
				r1.EXPECT().DeleteFilmById(actorId).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `"successfully delete"`,
		},
		{
			name:   "Wrong Params",
			params: "id=fksfm",
			mockBehavior: func(r1 *mock_service.MockFilm, r2 *mock_service.MockUser) {
				r2.EXPECT().ParseToken(token).Return(userId, nil)
				r2.EXPECT().ValidateUser(userId).Return(true, nil)
			},
			expectedStatusCode:   400,
			expectedResponseBody: `id doesnt specified to delete film`,
		},
		{
			name:   "Locked",
			params: "id=1",
			mockBehavior: func(r1 *mock_service.MockFilm, r2 *mock_service.MockUser) {
				r2.EXPECT().ParseToken(token).Return(userId, nil)
				r2.EXPECT().ValidateUser(userId).Return(false, nil)
			},
			expectedStatusCode:   423,
			expectedResponseBody: `this function locked for current user`,
		},
		{
			name:   "Service Error",
			params: "id=1",
			mockBehavior: func(r1 *mock_service.MockFilm, r2 *mock_service.MockUser) {
				r2.EXPECT().ParseToken(token).Return(userId, nil)
				r2.EXPECT().ValidateUser(userId).Return(true, nil)
				r1.EXPECT().DeleteFilmById(actorId).Return(errors.New("something went wrong"))
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

			repo1 := mock_service.NewMockFilm(c)
			repo2 := mock_service.NewMockUser(c)
			test.mockBehavior(repo1, repo2)

			services := &service.Service{Film: repo1, User: repo2}
			handler := Router{service: services}
			handler.AddEndPoint("DELETE", "/films", deleteFilm)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/films?"+test.params,
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

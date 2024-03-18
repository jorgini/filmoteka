// Code generated by MockGen. DO NOT EDIT.
// Source: service.go
//
// Generated by this command:
//
//	mockgen -source=service.go -destination=mocks/mock.go
//

// Package mock_service is a generated GoMock package.
package mock_service

import (
	"github.com/jorgini/filmoteka"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockUser is a mock of User interface.
type MockUser struct {
	ctrl     *gomock.Controller
	recorder *MockUserMockRecorder
}

// MockUserMockRecorder is the mock recorder for MockUser.
type MockUserMockRecorder struct {
	mock *MockUser
}

// NewMockUser creates a new mock instance.
func NewMockUser(ctrl *gomock.Controller) *MockUser {
	mock := &MockUser{ctrl: ctrl}
	mock.recorder = &MockUserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUser) EXPECT() *MockUserMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockUser) CreateUser(user filmoteka.User) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", user)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockUserMockRecorder) CreateUser(user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUser)(nil).CreateUser), user)
}

// DeleteUserById mocks base method.
func (m *MockUser) DeleteUserById(id int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUserById", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUserById indicates an expected call of DeleteUserById.
func (mr *MockUserMockRecorder) DeleteUserById(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUserById", reflect.TypeOf((*MockUser)(nil).DeleteUserById), id)
}

// GenerateToken mocks base method.
func (m *MockUser) GenerateToken(login, password string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateToken", login, password)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateToken indicates an expected call of GenerateToken.
func (mr *MockUserMockRecorder) GenerateToken(login, password any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateToken", reflect.TypeOf((*MockUser)(nil).GenerateToken), login, password)
}

// ParseToken mocks base method.
func (m *MockUser) ParseToken(token string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseToken", token)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseToken indicates an expected call of ParseToken.
func (mr *MockUserMockRecorder) ParseToken(token any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseToken", reflect.TypeOf((*MockUser)(nil).ParseToken), token)
}

// UpdateUser mocks base method.
func (m *MockUser) UpdateUser(login, userRole string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", login, userRole)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockUserMockRecorder) UpdateUser(login, userRole any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockUser)(nil).UpdateUser), login, userRole)
}

// ValidateUser mocks base method.
func (m *MockUser) ValidateUser(id int) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateUser", id)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateUser indicates an expected call of ValidateUser.
func (mr *MockUserMockRecorder) ValidateUser(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateUser", reflect.TypeOf((*MockUser)(nil).ValidateUser), id)
}

// MockActor is a mock of Actor interface.
type MockActor struct {
	ctrl     *gomock.Controller
	recorder *MockActorMockRecorder
}

// MockActorMockRecorder is the mock recorder for MockActor.
type MockActorMockRecorder struct {
	mock *MockActor
}

// NewMockActor creates a new mock instance.
func NewMockActor(ctrl *gomock.Controller) *MockActor {
	mock := &MockActor{ctrl: ctrl}
	mock.recorder = &MockActorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockActor) EXPECT() *MockActorMockRecorder {
	return m.recorder
}

// CreateActor mocks base method.
func (m *MockActor) CreateActor(actor filmoteka.Actor) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateActor", actor)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateActor indicates an expected call of CreateActor.
func (mr *MockActorMockRecorder) CreateActor(actor any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateActor", reflect.TypeOf((*MockActor)(nil).CreateActor), actor)
}

// DeleteActorById mocks base method.
func (m *MockActor) DeleteActorById(id int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteActorById", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteActorById indicates an expected call of DeleteActorById.
func (mr *MockActorMockRecorder) DeleteActorById(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteActorById", reflect.TypeOf((*MockActor)(nil).DeleteActorById), id)
}

// GetActorById mocks base method.
func (m *MockActor) GetActorById(id int) (filmoteka.ActorListItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetActorById", id)
	ret0, _ := ret[0].(filmoteka.ActorListItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetActorById indicates an expected call of GetActorById.
func (mr *MockActorMockRecorder) GetActorById(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetActorById", reflect.TypeOf((*MockActor)(nil).GetActorById), id)
}

// GetActorId mocks base method.
func (m *MockActor) GetActorId(name, surname string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetActorId", name, surname)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetActorId indicates an expected call of GetActorId.
func (mr *MockActorMockRecorder) GetActorId(name, surname any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetActorId", reflect.TypeOf((*MockActor)(nil).GetActorId), name, surname)
}

// GetActorsList mocks base method.
func (m *MockActor) GetActorsList(page, limit int) ([]filmoteka.ActorListItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetActorsList", page, limit)
	ret0, _ := ret[0].([]filmoteka.ActorListItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetActorsList indicates an expected call of GetActorsList.
func (mr *MockActorMockRecorder) GetActorsList(page, limit any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetActorsList", reflect.TypeOf((*MockActor)(nil).GetActorsList), page, limit)
}

// SearchActor mocks base method.
func (m *MockActor) SearchActor(page, limit int, fragment filmoteka.ActorSearchFragment) ([]filmoteka.ActorListItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchActor", page, limit, fragment)
	ret0, _ := ret[0].([]filmoteka.ActorListItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchActor indicates an expected call of SearchActor.
func (mr *MockActorMockRecorder) SearchActor(page, limit, fragment any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchActor", reflect.TypeOf((*MockActor)(nil).SearchActor), page, limit, fragment)
}

// UpdateActor mocks base method.
func (m *MockActor) UpdateActor(actor filmoteka.UpdateActorInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateActor", actor)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateActor indicates an expected call of UpdateActor.
func (mr *MockActorMockRecorder) UpdateActor(actor any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateActor", reflect.TypeOf((*MockActor)(nil).UpdateActor), actor)
}

// MockFilm is a mock of Film interface.
type MockFilm struct {
	ctrl     *gomock.Controller
	recorder *MockFilmMockRecorder
}

// MockFilmMockRecorder is the mock recorder for MockFilm.
type MockFilmMockRecorder struct {
	mock *MockFilm
}

// NewMockFilm creates a new mock instance.
func NewMockFilm(ctrl *gomock.Controller) *MockFilm {
	mock := &MockFilm{ctrl: ctrl}
	mock.recorder = &MockFilmMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFilm) EXPECT() *MockFilmMockRecorder {
	return m.recorder
}

// CreateActor mocks base method.
func (m *MockFilm) CreateActor(actor filmoteka.Actor) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateActor", actor)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateActor indicates an expected call of CreateActor.
func (mr *MockFilmMockRecorder) CreateActor(actor any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateActor", reflect.TypeOf((*MockFilm)(nil).CreateActor), actor)
}

// CreateFilm mocks base method.
func (m *MockFilm) CreateFilm(film filmoteka.InputFilm) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateFilm", film)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateFilm indicates an expected call of CreateFilm.
func (mr *MockFilmMockRecorder) CreateFilm(film any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateFilm", reflect.TypeOf((*MockFilm)(nil).CreateFilm), film)
}

// DeleteActorById mocks base method.
func (m *MockFilm) DeleteActorById(id int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteActorById", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteActorById indicates an expected call of DeleteActorById.
func (mr *MockFilmMockRecorder) DeleteActorById(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteActorById", reflect.TypeOf((*MockFilm)(nil).DeleteActorById), id)
}

// DeleteFilmById mocks base method.
func (m *MockFilm) DeleteFilmById(id int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteFilmById", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteFilmById indicates an expected call of DeleteFilmById.
func (mr *MockFilmMockRecorder) DeleteFilmById(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFilmById", reflect.TypeOf((*MockFilm)(nil).DeleteFilmById), id)
}

// GetActorById mocks base method.
func (m *MockFilm) GetActorById(id int) (filmoteka.ActorListItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetActorById", id)
	ret0, _ := ret[0].(filmoteka.ActorListItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetActorById indicates an expected call of GetActorById.
func (mr *MockFilmMockRecorder) GetActorById(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetActorById", reflect.TypeOf((*MockFilm)(nil).GetActorById), id)
}

// GetActorId mocks base method.
func (m *MockFilm) GetActorId(name, surname string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetActorId", name, surname)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetActorId indicates an expected call of GetActorId.
func (mr *MockFilmMockRecorder) GetActorId(name, surname any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetActorId", reflect.TypeOf((*MockFilm)(nil).GetActorId), name, surname)
}

// GetActorsList mocks base method.
func (m *MockFilm) GetActorsList(page, limit int) ([]filmoteka.ActorListItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetActorsList", page, limit)
	ret0, _ := ret[0].([]filmoteka.ActorListItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetActorsList indicates an expected call of GetActorsList.
func (mr *MockFilmMockRecorder) GetActorsList(page, limit any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetActorsList", reflect.TypeOf((*MockFilm)(nil).GetActorsList), page, limit)
}

// GetCurFilm mocks base method.
func (m *MockFilm) GetCurFilm(id int) (filmoteka.InputFilm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCurFilm", id)
	ret0, _ := ret[0].(filmoteka.InputFilm)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCurFilm indicates an expected call of GetCurFilm.
func (mr *MockFilmMockRecorder) GetCurFilm(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCurFilm", reflect.TypeOf((*MockFilm)(nil).GetCurFilm), id)
}

// GetSearchFilmList mocks base method.
func (m *MockFilm) GetSearchFilmList(page, limit int, fragment filmoteka.FilmSearchFragment) ([]filmoteka.InputFilm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSearchFilmList", page, limit, fragment)
	ret0, _ := ret[0].([]filmoteka.InputFilm)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSearchFilmList indicates an expected call of GetSearchFilmList.
func (mr *MockFilmMockRecorder) GetSearchFilmList(page, limit, fragment any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSearchFilmList", reflect.TypeOf((*MockFilm)(nil).GetSearchFilmList), page, limit, fragment)
}

// GetSortedFilmList mocks base method.
func (m *MockFilm) GetSortedFilmList(sortBy string, page, limit int) ([]filmoteka.InputFilm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSortedFilmList", sortBy, page, limit)
	ret0, _ := ret[0].([]filmoteka.InputFilm)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSortedFilmList indicates an expected call of GetSortedFilmList.
func (mr *MockFilmMockRecorder) GetSortedFilmList(sortBy, page, limit any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSortedFilmList", reflect.TypeOf((*MockFilm)(nil).GetSortedFilmList), sortBy, page, limit)
}

// SearchActor mocks base method.
func (m *MockFilm) SearchActor(page, limit int, fragment filmoteka.ActorSearchFragment) ([]filmoteka.ActorListItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchActor", page, limit, fragment)
	ret0, _ := ret[0].([]filmoteka.ActorListItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchActor indicates an expected call of SearchActor.
func (mr *MockFilmMockRecorder) SearchActor(page, limit, fragment any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchActor", reflect.TypeOf((*MockFilm)(nil).SearchActor), page, limit, fragment)
}

// UpdateActor mocks base method.
func (m *MockFilm) UpdateActor(actor filmoteka.UpdateActorInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateActor", actor)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateActor indicates an expected call of UpdateActor.
func (mr *MockFilmMockRecorder) UpdateActor(actor any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateActor", reflect.TypeOf((*MockFilm)(nil).UpdateActor), actor)
}

// UpdateFilm mocks base method.
func (m *MockFilm) UpdateFilm(film filmoteka.UpdateFilmInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateFilm", film)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateFilm indicates an expected call of UpdateFilm.
func (mr *MockFilmMockRecorder) UpdateFilm(film any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateFilm", reflect.TypeOf((*MockFilm)(nil).UpdateFilm), film)
}
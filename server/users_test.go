package server

import (
	"context"
	"io"
	"testing"

	"github.com/dhaifley/dlib"
	"github.com/dhaifley/dlib/ptypes"
	"github.com/dhaifley/dlib/dauth"
	"github.com/sirupsen/logrus/hooks/test"
	"google.golang.org/grpc"
)

type MockUserAccess struct {
	DBS dlib.SQLExecutor
}

func (m *MockUserAccess) GetUsers(opt *dauth.UserFind) <-chan dlib.Result {
	ch := make(chan dlib.Result, 256)
	go func() {
		defer close(ch)
		user := dauth.User{ID: 1, User: "test"}
		r := dlib.Result{Val: &user}
		ch <- r
	}()

	return ch
}

func (m *MockUserAccess) GetUserByID(id int64) <-chan dlib.Result {
	return m.GetUsers(nil)
}

func (m *MockUserAccess) DeleteUsers(opt *dauth.UserFind) <-chan dlib.Result {
	ch := make(chan dlib.Result, 256)
	go func() {
		defer close(ch)
		r := dlib.Result{Num: 1}
		ch <- r
	}()

	return ch
}

func (m *MockUserAccess) DeleteUserByID(id int64) <-chan dlib.Result {
	return m.DeleteUsers(nil)
}

func (m *MockUserAccess) SaveUser(a *dauth.User) <-chan dlib.Result {
	ch := make(chan dlib.Result, 256)
	go func() {
		defer close(ch)
		user := dauth.User{ID: 1, User: "test"}
		r := dlib.Result{Val: &user}
		ch <- r
	}()

	return ch
}

func (m *MockUserAccess) SaveUsers(a []dauth.User) <-chan dlib.Result {
	return m.SaveUser(nil)
}

type MockRFAuthGetUsersServer struct {
	grpc.ServerStream
	Results []ptypes.UserResponse
}

func (m *MockRFAuthGetUsersServer) Send(msg *ptypes.UserResponse) error {
	m.Results = append(m.Results, *msg)
	return nil
}

type MockRFAuthSaveUsersServer struct {
	grpc.ServerStream
	Results []ptypes.UserResponse
	Count   int16
}

func (m *MockRFAuthSaveUsersServer) Send(msg *ptypes.UserResponse) error {
	m.Results = append(m.Results, *msg)
	return nil
}

func (m *MockRFAuthSaveUsersServer) Recv() (*ptypes.UserRequest, error) {
	if m.Count < 1 {
		msg := ptypes.UserRequest{
			ID:   1,
			User: "test",
		}

		m.Count++
		return &msg, nil
	}

	return nil, io.EOF
}
func TestServerGetUsers(t *testing.T) {
	ma := MockUserAccess{}
	lm, _ := test.NewNullLogger()
	svr := Server{Users: &ma, Log: lm}
	var stream MockRFAuthGetUsersServer
	err := svr.GetUsers(&ptypes.UserRequest{User: "test"}, &stream)
	if err != nil {
		t.Error(err)
	}

	if stream.Results[0].ID != 1 {
		t.Errorf("ID expected: 1, got %v", stream.Results[0].ID)
	}
}

func TestServerSaveUsers(t *testing.T) {
	ma := MockUserAccess{}
	lm, _ := test.NewNullLogger()
	svr := Server{Users: &ma, Log: lm}
	var stream MockRFAuthSaveUsersServer
	err := svr.SaveUsers(&stream)
	if err != nil {
		t.Error(err)
	}

	if stream.Results[0].ID != 1 {
		t.Errorf("ID expected: 1, got %v", stream.Results[0].ID)
	}
}

func TestServerDeleteUsers(t *testing.T) {
	ma := MockUserAccess{}
	lm, _ := test.NewNullLogger()
	svr := Server{Users: &ma, Log: lm}
	res, err := svr.DeleteUsers(context.Background(), &ptypes.UserRequest{ID: 1})
	if err != nil {
		t.Error(err)
	}

	if res.Num != 1 {
		t.Errorf("Num expected: 1, got %v", res.Num)
	}
}

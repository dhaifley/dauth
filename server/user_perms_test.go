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

type MockUserPermAccess struct {
	DBS dlib.SQLExecutor
}

func (m *MockUserPermAccess) GetUserPerms(opt *dauth.UserPermFind) <-chan dlib.Result {
	ch := make(chan dlib.Result, 256)
	go func() {
		defer close(ch)
		up := dauth.UserPerm{ID: 1, UserID: 1, PermID: 1}
		r := dlib.Result{Val: &up}
		ch <- r
	}()

	return ch
}

func (m *MockUserPermAccess) GetUserPermByID(id int64) <-chan dlib.Result {
	return m.GetUserPerms(nil)
}

func (m *MockUserPermAccess) DeleteUserPerms(opt *dauth.UserPermFind) <-chan dlib.Result {
	ch := make(chan dlib.Result, 256)
	go func() {
		defer close(ch)
		r := dlib.Result{Num: 1}
		ch <- r
	}()

	return ch
}

func (m *MockUserPermAccess) DeleteUserPermByID(id int64) <-chan dlib.Result {
	return m.DeleteUserPerms(nil)
}

func (m *MockUserPermAccess) SaveUserPerm(a *dauth.UserPerm) <-chan dlib.Result {
	ch := make(chan dlib.Result, 256)
	go func() {
		defer close(ch)
		up := dauth.UserPerm{ID: 1, UserID: 1, PermID: 1}
		r := dlib.Result{Val: &up}
		ch <- r
	}()

	return ch
}

func (m *MockUserPermAccess) SaveUserPerms(a []dauth.UserPerm) <-chan dlib.Result {
	return m.SaveUserPerm(nil)
}

type MockRFAuthGetUserPermsServer struct {
	grpc.ServerStream
	Results []ptypes.UserPermResponse
}

func (m *MockRFAuthGetUserPermsServer) Send(msg *ptypes.UserPermResponse) error {
	m.Results = append(m.Results, *msg)
	return nil
}

type MockRFAuthSaveUserPermsServer struct {
	grpc.ServerStream
	Results []ptypes.UserPermResponse
	Count   int16
}

func (m *MockRFAuthSaveUserPermsServer) Send(msg *ptypes.UserPermResponse) error {
	m.Results = append(m.Results, *msg)
	return nil
}

func (m *MockRFAuthSaveUserPermsServer) Recv() (*ptypes.UserPermRequest, error) {
	if m.Count < 1 {
		msg := ptypes.UserPermRequest{ID: 1, UserID: 1, PermID: 1}
		m.Count++
		return &msg, nil
	}

	return nil, io.EOF
}
func TestServerGetUserPerms(t *testing.T) {
	ma := MockUserPermAccess{}
	lm, _ := test.NewNullLogger()
	svr := Server{UserPerms: &ma, Log: lm}
	var stream MockRFAuthGetUserPermsServer
	err := svr.GetUserPerms(&ptypes.UserPermRequest{ID: 1}, &stream)
	if err != nil {
		t.Error(err)
	}

	if stream.Results[0].ID != 1 {
		t.Errorf("ID expected: 1, got %v", stream.Results[0].ID)
	}
}

func TestServerSaveUserPerms(t *testing.T) {
	ma := MockUserPermAccess{}
	lm, _ := test.NewNullLogger()
	svr := Server{UserPerms: &ma, Log: lm}
	var stream MockRFAuthSaveUserPermsServer
	err := svr.SaveUserPerms(&stream)
	if err != nil {
		t.Error(err)
	}

	if stream.Results[0].ID != 1 {
		t.Errorf("ID expected: 1, got %v", stream.Results[0].ID)
	}
}

func TestServerDeleteUserPerms(t *testing.T) {
	ma := MockUserPermAccess{}
	lm, _ := test.NewNullLogger()
	svr := Server{UserPerms: &ma, Log: lm}
	res, err := svr.DeleteUserPerms(context.Background(), &ptypes.UserPermRequest{ID: 1})
	if err != nil {
		t.Error(err)
	}

	if res.Num != 1 {
		t.Errorf("Num expected: 1, got %v", res.Num)
	}
}

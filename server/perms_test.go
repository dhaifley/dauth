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

type MockPermAccess struct {
	DBS dlib.SQLExecutor
}

func (m *MockPermAccess) GetPerms(opt *dauth.PermFind) <-chan dlib.Result {
	ch := make(chan dlib.Result, 256)
	go func() {
		defer close(ch)
		perm := dauth.Perm{ID: 1, Service: "test", Name: "test"}
		r := dlib.Result{Val: &perm}
		ch <- r
	}()

	return ch
}

func (m *MockPermAccess) GetPermByID(id int64) <-chan dlib.Result {
	return m.GetPerms(nil)
}

func (m *MockPermAccess) DeletePerms(opt *dauth.PermFind) <-chan dlib.Result {
	ch := make(chan dlib.Result, 256)
	go func() {
		defer close(ch)
		r := dlib.Result{Num: 1}
		ch <- r
	}()

	return ch
}

func (m *MockPermAccess) DeletePermByID(id int64) <-chan dlib.Result {
	return m.DeletePerms(nil)
}

func (m *MockPermAccess) SavePerm(a *dauth.Perm) <-chan dlib.Result {
	ch := make(chan dlib.Result, 256)
	go func() {
		defer close(ch)
		perm := dauth.Perm{ID: 1, Service: "test", Name: "test"}
		r := dlib.Result{Val: &perm}
		ch <- r
	}()

	return ch
}

func (m *MockPermAccess) SavePerms(a []dauth.Perm) <-chan dlib.Result {
	return m.SavePerm(nil)
}

type MockRFAuthGetPermsServer struct {
	grpc.ServerStream
	Results []ptypes.PermResponse
}

func (m *MockRFAuthGetPermsServer) Send(msg *ptypes.PermResponse) error {
	m.Results = append(m.Results, *msg)
	return nil
}

type MockRFAuthSavePermsServer struct {
	grpc.ServerStream
	Results []ptypes.PermResponse
	Count   int16
}

func (m *MockRFAuthSavePermsServer) Send(msg *ptypes.PermResponse) error {
	m.Results = append(m.Results, *msg)
	return nil
}

func (m *MockRFAuthSavePermsServer) Recv() (*ptypes.PermRequest, error) {
	if m.Count < 1 {
		msg := ptypes.PermRequest{ID: 1, Service: "test", Name: "test"}
		m.Count++
		return &msg, nil
	}

	return nil, io.EOF
}
func TestServerGetPerms(t *testing.T) {
	ma := MockPermAccess{}
	lm, _ := test.NewNullLogger()
	svr := Server{Perms: &ma, Log: lm}
	var stream MockRFAuthGetPermsServer
	err := svr.GetPerms(&ptypes.PermRequest{Name: "test"}, &stream)
	if err != nil {
		t.Error(err)
	}

	if stream.Results[0].ID != 1 {
		t.Errorf("ID expected: 1, got %v", stream.Results[0].ID)
	}
}

func TestServerSavePerms(t *testing.T) {
	ma := MockPermAccess{}
	lm, _ := test.NewNullLogger()
	svr := Server{Perms: &ma, Log: lm}
	var stream MockRFAuthSavePermsServer
	err := svr.SavePerms(&stream)
	if err != nil {
		t.Error(err)
	}

	if stream.Results[0].ID != 1 {
		t.Errorf("ID expected: 1, got %v", stream.Results[0].ID)
	}
}

func TestServerDeletePerms(t *testing.T) {
	ma := MockPermAccess{}
	lm, _ := test.NewNullLogger()
	svr := Server{Perms: &ma, Log: lm}
	res, err := svr.DeletePerms(context.Background(), &ptypes.PermRequest{ID: 1})
	if err != nil {
		t.Error(err)
	}

	if res.Num != 1 {
		t.Errorf("Num expected: 1, got %v", res.Num)
	}
}

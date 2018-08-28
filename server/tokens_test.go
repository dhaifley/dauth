package server

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/dhaifley/dlib"
	"github.com/dhaifley/dlib/ptypes"
	"github.com/dhaifley/dlib/dauth"
	"github.com/sirupsen/logrus/hooks/test"
	"google.golang.org/grpc"
)

type MockTokenAccess struct {
	DBS dlib.SQLExecutor
}

func (m *MockTokenAccess) GetTokens(opt *dauth.TokenFind) <-chan dlib.Result {
	ch := make(chan dlib.Result, 256)
	go func() {
		defer close(ch)
		ct := time.Date(1983, 2, 2, 0, 0, 0, 0, time.Local)
		et := time.Date(2083, 2, 2, 0, 0, 0, 0, time.Local)
		token := dauth.Token{
			ID:      1,
			Token:   "test",
			UserID:  1,
			Created: &ct,
			Expires: &et,
		}

		r := dlib.Result{Val: &token}
		ch <- r
	}()

	return ch
}

func (m *MockTokenAccess) GetTokenByID(id int64) <-chan dlib.Result {
	return m.GetTokens(nil)
}

func (m *MockTokenAccess) DeleteTokens(opt *dauth.TokenFind) <-chan dlib.Result {
	ch := make(chan dlib.Result, 256)
	go func() {
		defer close(ch)
		r := dlib.Result{Num: 1}
		ch <- r
	}()

	return ch
}

func (m *MockTokenAccess) DeleteTokenByID(id int64) <-chan dlib.Result {
	return m.DeleteTokens(nil)
}

func (m *MockTokenAccess) SaveToken(a *dauth.Token) <-chan dlib.Result {
	ch := make(chan dlib.Result, 256)
	go func() {
		defer close(ch)
		token := dauth.Token{ID: 1, Token: "test"}
		r := dlib.Result{Val: &token}
		ch <- r
	}()

	return ch
}

func (m *MockTokenAccess) SaveTokens(a []dauth.Token) <-chan dlib.Result {
	return m.SaveToken(nil)
}

type MockRFAuthGetTokensServer struct {
	grpc.ServerStream
	Results []ptypes.TokenResponse
}

func (m *MockRFAuthGetTokensServer) Send(msg *ptypes.TokenResponse) error {
	m.Results = append(m.Results, *msg)
	return nil
}

type MockRFAuthSaveTokensServer struct {
	grpc.ServerStream
	Results []ptypes.TokenResponse
	Count   int16
}

func (m *MockRFAuthSaveTokensServer) Send(msg *ptypes.TokenResponse) error {
	m.Results = append(m.Results, *msg)
	return nil
}

func (m *MockRFAuthSaveTokensServer) Recv() (*ptypes.TokenRequest, error) {
	if m.Count < 1 {
		msg := ptypes.TokenRequest{
			ID:    1,
			Token: "test",
		}

		m.Count++
		return &msg, nil
	}

	return nil, io.EOF
}
func TestServerGetTokens(t *testing.T) {
	ma := MockTokenAccess{}
	lm, _ := test.NewNullLogger()
	svr := Server{Tokens: &ma, Log: lm}
	var stream MockRFAuthGetTokensServer
	err := svr.GetTokens(&ptypes.TokenRequest{Token: "test"}, &stream)
	if err != nil {
		t.Error(err)
	}

	if stream.Results[0].ID != 1 {
		t.Errorf("ID expected: 1, got %v", stream.Results[0].ID)
	}
}

func TestServerSaveTokens(t *testing.T) {
	ma := MockTokenAccess{}
	lm, _ := test.NewNullLogger()
	svr := Server{Tokens: &ma, Log: lm}
	var stream MockRFAuthSaveTokensServer
	err := svr.SaveTokens(&stream)
	if err != nil {
		t.Error(err)
	}

	if stream.Results[0].ID != 1 {
		t.Errorf("ID expected: 1, got %v", stream.Results[0].ID)
	}
}

func TestServerDeleteTokens(t *testing.T) {
	ma := MockTokenAccess{}
	lm, _ := test.NewNullLogger()
	svr := Server{Tokens: &ma, Log: lm}
	res, err := svr.DeleteTokens(context.Background(), &ptypes.TokenRequest{ID: 1})
	if err != nil {
		t.Error(err)
	}

	if res.Num != 1 {
		t.Errorf("Num expected: 1, got %v", res.Num)
	}
}

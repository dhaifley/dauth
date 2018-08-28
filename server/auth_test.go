package server

import (
	"context"
	"testing"

	"github.com/dhaifley/dlib"
	"github.com/dhaifley/dlib/ptypes"
	"github.com/sirupsen/logrus/hooks/test"
)

func TestServerAuth(t *testing.T) {
	mua := MockUserAccess{}
	mta := MockTokenAccess{}
	mpa := MockPermAccess{}
	mupa := MockUserPermAccess{}
	lm, _ := test.NewNullLogger()
	svr := Server{Users: &mua, Tokens: &mta, Perms: &mpa, UserPerms: &mupa, Log: lm}
	cases := []struct {
		req ptypes.AuthRequest
		exp bool
	}{
		{
			req: ptypes.AuthRequest{
				Token: &ptypes.TokenRequest{
					ID:     1,
					Token:  "test",
					UserID: 1,
				},
				Perm: &ptypes.PermRequest{
					ID:      1,
					Service: "test",
					Name:    "test",
				},
			},
			exp: true,
		},
		{
			req: ptypes.AuthRequest{
				Token: &ptypes.TokenRequest{
					ID:     1,
					Token:  "test",
					UserID: 1,
				},
				Perm: &ptypes.PermRequest{
					ID:      1,
					Service: "test",
					Name:    "wrong",
				},
			},
			exp: false,
		},
	}

	for _, c := range cases {
		res, err := svr.Auth(context.Background(), &c.req)
		if err != nil {
			t.Error(err)
		}

		if c.exp != res.Ok {
			t.Error("Failed to correctly authenticate")
		}
	}
}

func TestServerLogin(t *testing.T) {
	mua := MockUserAccess{}
	mta := MockTokenAccess{}
	lm, _ := test.NewNullLogger()
	svr := Server{Users: &mua, Tokens: &mta, Log: lm}
	req := ptypes.UserRequest{
		ID:   1,
		User: "test",
		Pass: dlib.EncodeBase64String("test"),
	}

	res, err := svr.Login(context.Background(), &req)
	if err != nil {
		t.Error(err)
	}

	if res.UserID != 1 {
		t.Errorf("User id expected: 1, got: %v", res.UserID)
	}
}

func TestServerLogout(t *testing.T) {
	mua := MockUserAccess{}
	mta := MockTokenAccess{}
	lm, _ := test.NewNullLogger()
	svr := Server{Users: &mua, Tokens: &mta, Log: lm}
	req := ptypes.TokenRequest{
		ID:     1,
		Token:  "test",
		UserID: 1,
	}

	res, err := svr.Logout(context.Background(), &req)
	if err != nil {
		t.Error(err)
	}

	if res.Token != "logout" {
		t.Errorf("Token expected: test, got: %v", res.Token)
	}

	if res.UserID != 0 {
		t.Errorf("User ID expected: 0, got: %v", res.UserID)
	}
}

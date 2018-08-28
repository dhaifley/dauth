package server

import (
	"context"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dhaifley/dlib"
	"github.com/dhaifley/dlib/dauth"
	"github.com/dhaifley/dlib/ptypes"
	"github.com/sirupsen/logrus"
)

// Auth authenticates a provided token and returns a user value.
func (s *Server) Auth(ctx context.Context,
	req *ptypes.AuthRequest) (*ptypes.AuthResponse, error) {
	if req.Token == nil {
		err := dlib.NewError(http.StatusBadRequest, "invalid token value")
		s.Log.WithFields(logrus.Fields{
			"rpc":     "Auth",
			"code":    http.StatusBadRequest,
			"context": ctx,
			"request": req,
		}).Error(err)
		return nil, err
	}

	q := dauth.TokenFind{}
	q.FromTokenRequest(req.Token)
	var t []dauth.Token
	var ch <-chan dlib.Result
	ch = s.Tokens.GetTokens(&q)
	for tr := range ch {
		if tr.Err != nil {
			switch err := tr.Err.(type) {
			case *dlib.Error:
				if err.Code == http.StatusNotFound {
					err := dlib.NewError(http.StatusUnauthorized, "unauthorized token")
					s.Log.WithFields(logrus.Fields{
						"rpc":     "Auth",
						"code":    http.StatusUnauthorized,
						"context": ctx,
						"request": req,
					}).Warning("unauthorized token")
					return nil, err
				}

				s.Log.WithFields(logrus.Fields{
					"rpc":     "Auth",
					"code":    http.StatusInternalServerError,
					"context": ctx,
					"request": req,
				}).Error(err)
				return nil, err
			default:
				s.Log.WithFields(logrus.Fields{
					"rpc":     "Auth",
					"code":    http.StatusInternalServerError,
					"context": ctx,
					"request": req,
				}).Error(err)
				return nil, err
			}
		}

		switch v := tr.Val.(type) {
		case *dauth.Token:
			t = append(t, *v)
		default:
			continue
		}
	}

	if len(t) == 0 {
		err := dlib.NewError(http.StatusUnauthorized, "unauthorized token")
		s.Log.WithFields(logrus.Fields{
			"rpc":     "Auth",
			"code":    http.StatusUnauthorized,
			"context": ctx,
			"request": req,
		}).Warning("unauthorized token")
		return nil, err
	}

	if t[0].Expires == nil || t[0].Expires.Unix() < time.Now().Unix() {
		err := dlib.NewError(http.StatusUnauthorized, "unauthorized token")
		s.Log.WithFields(logrus.Fields{
			"rpc":     "Auth",
			"code":    http.StatusUnauthorized,
			"context": ctx,
			"request": req,
		}).Warning("unauthorized token")
		return nil, err
	}

	qu := dauth.UserFind{ID: &t[0].UserID}
	var u []dauth.User
	ch = s.Users.GetUsers(&qu)
	for ur := range ch {
		if ur.Err != nil {
			switch err := ur.Err.(type) {
			case *dlib.Error:
				if err.Code == http.StatusNotFound {
					err := dlib.NewError(http.StatusUnauthorized, "unauthorized user")
					s.Log.WithFields(logrus.Fields{
						"rpc":     "Auth",
						"code":    http.StatusUnauthorized,
						"context": ctx,
						"request": req,
					}).Warning("unauthorized user")
					return nil, err
				}

				s.Log.WithFields(logrus.Fields{
					"rpc":     "Auth",
					"code":    http.StatusInternalServerError,
					"context": ctx,
					"request": req,
				}).Error(err)
				return nil, err
			default:
				s.Log.WithFields(logrus.Fields{
					"rpc":     "Auth",
					"code":    http.StatusInternalServerError,
					"context": ctx,
					"request": req,
				}).Error(err)
				return nil, err
			}
		}

		switch v := ur.Val.(type) {
		case *dauth.User:
			u = append(u, *v)
		default:
			continue
		}
	}

	if len(u) == 0 {
		err := dlib.NewError(http.StatusUnauthorized, "unknown user id")
		s.Log.WithFields(logrus.Fields{
			"rpc":     "Auth",
			"code":    http.StatusUnauthorized,
			"user_id": t[0].UserID,
			"context": ctx,
			"request": req,
		}).Warning("Unknown user id")
		return nil, err
	}

	u[0].Pass = ""
	ures := u[0].ToResponse()
	var pres ptypes.PermResponse
	qup := dauth.UserPermFind{UserID: &u[0].ID}
	ch = s.UserPerms.GetUserPerms(&qup)
	ok := false
	for upr := range ch {
		if upr.Err != nil {
			switch err := upr.Err.(type) {
			case *dlib.Error:
				if err.Code == http.StatusNotFound {
					err := dlib.NewError(http.StatusUnauthorized, "unauthorized user")
					s.Log.WithFields(logrus.Fields{
						"rpc":     "Auth",
						"code":    http.StatusUnauthorized,
						"context": ctx,
						"request": req,
					}).Warning("unauthorized user")
					return nil, err
				}

				s.Log.WithFields(logrus.Fields{
					"rpc":     "Auth",
					"code":    http.StatusInternalServerError,
					"context": ctx,
					"request": req,
				}).Error(err)
				return nil, err
			default:
				s.Log.WithFields(logrus.Fields{
					"rpc":     "Auth",
					"code":    http.StatusInternalServerError,
					"context": ctx,
					"request": req,
				}).Error(err)
				return nil, err
			}
		}

		switch v := upr.Val.(type) {
		case *dauth.UserPerm:
			qp := dauth.PermFind{ID: &v.PermID}
			chp := s.Perms.GetPerms(&qp)
			for up := range chp {
				if up.Err != nil {
					s.Log.WithFields(logrus.Fields{
						"rpc":     "Auth",
						"code":    http.StatusInternalServerError,
						"context": ctx,
						"request": req,
					}).Error(up.Err)
					continue
				}

				switch v := up.Val.(type) {
				case *dauth.Perm:
					if v.Service == "admin" || v.Service == req.Perm.Service {
						if v.Name == "admin" || v.Name == req.Perm.Name {
							pres = v.ToResponse()
							ok = true
							break
						}
					}
				default:
					continue
				}
			}
		default:
			continue
		}
	}

	res := ptypes.AuthResponse{
		Ok:   ok,
		User: &ures,
		Perm: &pres,
	}

	return &res, nil
}

// Login authenticates a provided user and creates a new token.
func (s *Server) Login(ctx context.Context,
	req *ptypes.UserRequest) (*ptypes.TokenResponse, error) {
	uq := dauth.User{}
	uq.FromRequest(req)
	if uq.User == "" || uq.Pass == "" {
		err := dlib.NewError(http.StatusUnauthorized, "unauthorized user")
		s.Log.WithFields(logrus.Fields{
			"rpc":     "Login",
			"code":    http.StatusUnauthorized,
			"context": ctx,
			"request": req,
		}).Warning("unauthorized user")
		return nil, err
	}

	pw, err := dlib.DecodeBase64String(uq.Pass)
	if err != nil {
		err := dlib.NewError(http.StatusUnauthorized, "unauthorized user")
		s.Log.WithFields(logrus.Fields{
			"rpc":     "Login",
			"code":    http.StatusUnauthorized,
			"context": ctx,
			"request": req,
		}).Warning("unauthorized user")
		return nil, err
	}

	pw, err = dlib.EncryptString(pw)
	if err != nil {
		s.Log.WithFields(logrus.Fields{
			"rpc":     "Login",
			"code":    http.StatusInternalServerError,
			"context": ctx,
			"request": req,
		}).Error(err)
		return nil, err
	}

	qu := dauth.UserFind{User: &uq.User, Pass: &pw}
	var u []dauth.User
	ch := s.Users.GetUsers(&qu)
	for ur := range ch {
		if ur.Err != nil {
			switch err := ur.Err.(type) {
			case *dlib.Error:
				if err.Code == http.StatusNotFound {
					err := dlib.NewError(http.StatusUnauthorized, "unauthorized user")
					s.Log.WithFields(logrus.Fields{
						"rpc":     "Login",
						"code":    http.StatusUnauthorized,
						"context": ctx,
						"request": req,
					}).Warning("unauthorized user")
					return nil, err
				}

				s.Log.WithFields(logrus.Fields{
					"rpc":     "Login",
					"code":    http.StatusInternalServerError,
					"context": ctx,
					"request": req,
				}).Error(err)
				return nil, err
			default:
				s.Log.WithFields(logrus.Fields{
					"rpc":     "Login",
					"code":    http.StatusInternalServerError,
					"context": ctx,
					"request": req,
				}).Error(err)
				return nil, err
			}
		}

		switch v := ur.Val.(type) {
		case *dauth.User:
			u = append(u, *v)
		case dauth.User:
			u = append(u, v)
		default:
			continue
		}
	}

	if len(u) == 0 {
		dlib.NewError(http.StatusUnauthorized, "unauthorized user")
		s.Log.WithFields(logrus.Fields{
			"rpc":     "Login",
			"code":    http.StatusUnauthorized,
			"context": ctx,
			"request": req,
		}).Warning("unauthorized user")
		return nil, err
	}

	ct := time.Now()
	et := ct.Add(time.Hour * 24)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":    u[0].User,
		"created": ct.Unix(),
		"expires": et.Unix(),
	})

	ts, err := token.SignedString([]byte(">^_^<"))
	if err != nil {
		s.Log.WithFields(logrus.Fields{
			"rpc":     "Login",
			"code":    http.StatusInternalServerError,
			"context": ctx,
			"request": req,
		}).Error(err)
		return nil, err
	}

	t := dauth.Token{
		Token:   ts,
		UserID:  u[0].ID,
		Created: &ct,
		Expires: &et,
	}

	ch = s.Tokens.SaveToken(&t)
	for tr := range ch {
		if tr.Err != nil {
			s.Log.WithFields(logrus.Fields{
				"rpc":     "Login",
				"code":    http.StatusInternalServerError,
				"context": ctx,
				"request": req,
			}).Error(err)
			return nil, err
		}
	}

	res := t.ToResponse()
	return &res, nil
}

// Logout destroys the provided token.
func (s *Server) Logout(ctx context.Context,
	req *ptypes.TokenRequest) (*ptypes.TokenResponse, error) {
	if req.Token == "" {
		err := dlib.NewError(http.StatusBadRequest, "invalid token value")
		s.Log.WithFields(logrus.Fields{
			"rpc":     "Logout",
			"code":    http.StatusBadRequest,
			"context": ctx,
			"request": req,
		}).Error("Invalid token value")

		return nil, err
	}

	q := dauth.TokenFind{Token: &req.Token}
	ch := s.Tokens.DeleteTokens(&q)
	for tr := range ch {
		if tr.Err != nil {
			s.Log.Error(tr.Err)
			return nil, tr.Err
		}

		if tr.Num == 0 {
			err := dlib.NewError(http.StatusNotFound, "token not found")
			s.Log.WithFields(logrus.Fields{
				"rpc":     "Logout",
				"code":    http.StatusNotFound,
				"context": ctx,
				"request": req,
			}).Error("Token not found")

			return nil, err
		}
	}

	res := ptypes.TokenResponse{
		Token:  "logout",
		UserID: 0,
	}

	return &res, nil
}

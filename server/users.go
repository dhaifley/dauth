package server

import (
	"context"
	"io"
	"net/http"

	"github.com/dhaifley/dlib"
	"github.com/dhaifley/dlib/dauth"
	"github.com/dhaifley/dlib/ptypes"
	"github.com/sirupsen/logrus"
)

// GetUsers returns a stream of users from the database.
func (s *Server) GetUsers(req *ptypes.UserRequest,
	stream ptypes.Auth_GetUsersServer) error {
	q := dauth.UserFind{}
	if err := q.FromUserRequest(req); err != nil {
		s.Log.WithFields(logrus.Fields{
			"rpc":     "GetUsers",
			"code":    http.StatusBadRequest,
			"request": req,
		}).Error(err)
		return err
	}

	ch := s.Users.GetUsers(&q)
	for r := range ch {
		if r.Err != nil {
			s.Log.WithFields(logrus.Fields{
				"rpc":     "GetUsers",
				"code":    http.StatusInternalServerError,
				"request": req,
			}).Error(r.Err)
			return r.Err
		}

		switch v := r.Val.(type) {
		case *dauth.User:
			v.Pass = ""
			res := v.ToResponse()
			if err := stream.Send(&res); err != nil {
				s.Log.WithFields(logrus.Fields{
					"rpc":      "GetUsers",
					"code":     http.StatusInternalServerError,
					"request":  req,
					"response": res,
				}).Error(err)
				return err
			}
		default:
			continue
		}
	}

	s.Log.WithFields(logrus.Fields{
		"rpc":     "GetUsers",
		"code":    http.StatusOK,
		"request": req,
	}).Info("GetUsers request processed")
	return nil
}

// SaveUsers serializes a stream of users to the database.
func (s *Server) SaveUsers(
	stream ptypes.Auth_SaveUsersServer) error {
	count := 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			s.Log.WithFields(logrus.Fields{
				"rpc":     "SaveUsers",
				"code":    http.StatusOK,
				"request": req,
				"count":   count,
			}).Info("SaveUsers request processed")
			return nil
		}

		if err != nil {
			s.Log.WithFields(logrus.Fields{
				"rpc":     "SaveUsers",
				"code":    http.StatusInternalServerError,
				"request": req,
				"count":   count,
			}).Error(err)
			return err
		}

		v := dauth.User{}
		v.FromRequest(req)
		if v.Pass != "" {
			dp, err := dlib.DecodeBase64String(v.Pass)
			if err != nil {
				s.Log.WithFields(logrus.Fields{
					"rpc":     "SaveUsers",
					"code":    http.StatusInternalServerError,
					"request": req,
					"count":   count,
				}).Error(err)
				return err
			}

			v.Pass, err = dlib.EncryptString(dp)
			if err != nil {
				s.Log.WithFields(logrus.Fields{
					"rpc":     "SaveUsers",
					"code":    http.StatusInternalServerError,
					"request": req,
					"count":   count,
				}).Error(err)
				return err
			}

		}

		ch := s.Users.SaveUser(&v)
		for r := range ch {
			if r.Err != nil {
				s.Log.WithFields(logrus.Fields{
					"rpc":     "SaveUsers",
					"code":    http.StatusInternalServerError,
					"request": req,
					"count":   count,
				}).Error(err)
				return err
			}
		}

		v.Pass = ""
		res := v.ToResponse()
		if err := stream.Send(&res); err != nil {
			s.Log.WithFields(logrus.Fields{
				"rpc":      "SaveUsers",
				"code":     http.StatusInternalServerError,
				"request":  req,
				"response": res,
				"count":    count,
			}).Error(err)
			return err
		}

		count++
	}
}

// DeleteUsers deletes users from the database.
func (s *Server) DeleteUsers(ctx context.Context, req *ptypes.UserRequest) (*ptypes.DeleteResponse, error) {
	q := dauth.UserFind{}
	if err := q.FromUserRequest(req); err != nil {
		s.Log.WithFields(logrus.Fields{
			"rpc":     "DeleteUsers",
			"code":    http.StatusBadRequest,
			"request": req,
		}).Error(err)
		return nil, err
	}

	ch := s.Users.DeleteUsers(&q)
	count := int64(0)
	for r := range ch {
		if r.Err != nil {
			s.Log.WithFields(logrus.Fields{
				"rpc":     "DeleteUsers",
				"code":    http.StatusInternalServerError,
				"request": req,
			}).Error(r.Err)
			return nil, r.Err
		}

		count += int64(r.Num)
	}

	s.Log.WithFields(logrus.Fields{
		"rpc":     "DeleteUsers",
		"code":    http.StatusOK,
		"request": req,
		"count":   count,
	}).Info("DeleteUsers request processed")
	return &ptypes.DeleteResponse{Num: count}, nil
}

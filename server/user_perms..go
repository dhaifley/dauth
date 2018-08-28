package server

import (
	"context"
	"io"
	"net/http"

	"github.com/dhaifley/dlib/dauth"
	"github.com/dhaifley/dlib/ptypes"
	"github.com/sirupsen/logrus"
)

// GetUserPerms returns a stream of user_perms from the database.
func (s *Server) GetUserPerms(req *ptypes.UserPermRequest,
	stream ptypes.Auth_GetUserPermsServer) error {
	q := dauth.UserPermFind{}
	if err := q.FromUserPermRequest(req); err != nil {
		s.Log.WithFields(logrus.Fields{
			"rpc":     "GetUserPerms",
			"code":    http.StatusBadRequest,
			"request": req,
		}).Error(err)
		return err
	}

	ch := s.UserPerms.GetUserPerms(&q)
	for r := range ch {
		if r.Err != nil {
			s.Log.WithFields(logrus.Fields{
				"rpc":     "GetUserPerms",
				"code":    http.StatusInternalServerError,
				"request": req,
			}).Error(r.Err)
			return r.Err
		}

		switch v := r.Val.(type) {
		case *dauth.UserPerm:
			res := v.ToResponse()
			if err := stream.Send(&res); err != nil {
				s.Log.WithFields(logrus.Fields{
					"rpc":      "GetUserPerms",
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
		"rpc":     "GetUserPerms",
		"code":    http.StatusOK,
		"request": req,
	}).Info("GetUserPerms request processed")
	return nil
}

// SaveUserPerms serializes a stream of user_perms to the database.
func (s *Server) SaveUserPerms(
	stream ptypes.Auth_SaveUserPermsServer) error {
	count := 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			s.Log.WithFields(logrus.Fields{
				"rpc":     "SaveUserPerms",
				"code":    http.StatusOK,
				"request": req,
				"count":   count,
			}).Info("SaveUserPerms request processed")
			return nil
		}

		if err != nil {
			s.Log.WithFields(logrus.Fields{
				"rpc":     "SaveUserPerms",
				"code":    http.StatusInternalServerError,
				"request": req,
				"count":   count,
			}).Error(err)
			return err
		}

		v := dauth.UserPerm{}
		v.FromRequest(req)
		ch := s.UserPerms.SaveUserPerm(&v)
		for r := range ch {
			if r.Err != nil {
				s.Log.WithFields(logrus.Fields{
					"rpc":     "SaveUserPerms",
					"code":    http.StatusInternalServerError,
					"request": req,
					"count":   count,
				}).Error(err)
				return err
			}
		}

		res := v.ToResponse()
		if err := stream.Send(&res); err != nil {
			s.Log.WithFields(logrus.Fields{
				"rpc":      "SaveUserPerms",
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

// DeleteUserPerms deletes user_perms from the database.
func (s *Server) DeleteUserPerms(ctx context.Context, req *ptypes.UserPermRequest) (*ptypes.DeleteResponse, error) {
	q := dauth.UserPermFind{}
	if err := q.FromUserPermRequest(req); err != nil {
		s.Log.WithFields(logrus.Fields{
			"rpc":     "DeleteUserPerms",
			"code":    http.StatusBadRequest,
			"request": req,
		}).Error(err)
		return nil, err
	}

	ch := s.UserPerms.DeleteUserPerms(&q)
	count := int64(0)
	for r := range ch {
		if r.Err != nil {
			s.Log.WithFields(logrus.Fields{
				"rpc":     "DeleteUserPerms",
				"code":    http.StatusInternalServerError,
				"request": req,
			}).Error(r.Err)
			return nil, r.Err
		}

		count += int64(r.Num)
	}

	s.Log.WithFields(logrus.Fields{
		"rpc":     "DeleteUserPerms",
		"code":    http.StatusOK,
		"request": req,
		"count":   count,
	}).Info("DeleteUserPerms request processed")
	return &ptypes.DeleteResponse{Num: count}, nil
}

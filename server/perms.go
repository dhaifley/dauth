package server

import (
	"context"
	"io"
	"net/http"

	"github.com/dhaifley/dlib/dauth"
	"github.com/dhaifley/dlib/ptypes"
	"github.com/sirupsen/logrus"
)

// GetPerms returns a stream of perms from the database.
func (s *Server) GetPerms(req *ptypes.PermRequest,
	stream ptypes.Auth_GetPermsServer) error {
	q := dauth.PermFind{}
	if err := q.FromPermRequest(req); err != nil {
		s.Log.WithFields(logrus.Fields{
			"rpc":     "GetPerms",
			"code":    http.StatusBadRequest,
			"request": req,
		}).Error(err)
		return err
	}

	ch := s.Perms.GetPerms(&q)
	for r := range ch {
		if r.Err != nil {
			s.Log.WithFields(logrus.Fields{
				"rpc":     "GetPerms",
				"code":    http.StatusInternalServerError,
				"request": req,
			}).Error(r.Err)
			return r.Err
		}

		switch v := r.Val.(type) {
		case *dauth.Perm:
			res := v.ToResponse()
			if err := stream.Send(&res); err != nil {
				s.Log.WithFields(logrus.Fields{
					"rpc":      "GetPerms",
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
		"rpc":     "GetPerms",
		"code":    http.StatusOK,
		"request": req,
	}).Info("GetPerms request processed")
	return nil
}

// SavePerms serializes a stream of perms to the database.
func (s *Server) SavePerms(
	stream ptypes.Auth_SavePermsServer) error {
	count := 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			s.Log.WithFields(logrus.Fields{
				"rpc":     "SavePerms",
				"code":    http.StatusOK,
				"request": req,
				"count":   count,
			}).Info("SavePerms request processed")
			return nil
		}

		if err != nil {
			s.Log.WithFields(logrus.Fields{
				"rpc":     "SavePerms",
				"code":    http.StatusInternalServerError,
				"request": req,
				"count":   count,
			}).Error(err)
			return err
		}

		v := dauth.Perm{}
		v.FromRequest(req)
		ch := s.Perms.SavePerm(&v)
		for r := range ch {
			if r.Err != nil {
				s.Log.WithFields(logrus.Fields{
					"rpc":     "SavePerms",
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
				"rpc":      "SavePerms",
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

// DeletePerms deletes perms from the database.
func (s *Server) DeletePerms(ctx context.Context, req *ptypes.PermRequest) (*ptypes.DeleteResponse, error) {
	q := dauth.PermFind{}
	if err := q.FromPermRequest(req); err != nil {
		s.Log.WithFields(logrus.Fields{
			"rpc":     "DeletePerms",
			"code":    http.StatusBadRequest,
			"request": req,
		}).Error(err)
		return nil, err
	}

	ch := s.Perms.DeletePerms(&q)
	count := int64(0)
	for r := range ch {
		if r.Err != nil {
			s.Log.WithFields(logrus.Fields{
				"rpc":     "DeletePerms",
				"code":    http.StatusInternalServerError,
				"request": req,
			}).Error(r.Err)
			return nil, r.Err
		}

		count += int64(r.Num)
	}

	s.Log.WithFields(logrus.Fields{
		"rpc":     "DeletePerms",
		"code":    http.StatusOK,
		"request": req,
		"count":   count,
	}).Info("DeletePerms request processed")
	return &ptypes.DeleteResponse{Num: count}, nil
}

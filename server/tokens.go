package server

import (
	"context"
	"io"
	"net/http"

	"github.com/dhaifley/dlib/dauth"
	"github.com/dhaifley/dlib/ptypes"
	"github.com/sirupsen/logrus"
)

// GetTokens returns a stream of tokens from the database.
func (s *Server) GetTokens(req *ptypes.TokenRequest,
	stream ptypes.Auth_GetTokensServer) error {
	q := dauth.TokenFind{}
	if err := q.FromTokenRequest(req); err != nil {
		s.Log.WithFields(logrus.Fields{
			"rpc":     "GetTokens",
			"code":    http.StatusInternalServerError,
			"context": stream.Context,
			"request": req,
		}).Error(err)
		return err
	}

	ch := s.Tokens.GetTokens(&q)
	for r := range ch {
		if r.Err != nil {
			s.Log.Error(r.Err)
			return r.Err
		}

		switch v := r.Val.(type) {
		case *dauth.Token:
			res := v.ToResponse()
			if err := stream.Send(&res); err != nil {
				s.Log.WithFields(logrus.Fields{
					"rpc":      "GetTokens",
					"code":     http.StatusInternalServerError,
					"context":  stream.Context,
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
		"rpc":     "GetTokens",
		"code":    http.StatusOK,
		"context": stream.Context,
		"request": req,
	}).Info("GetTokens request processed")
	return nil
}

// SaveTokens serializes a stream of tokens to the database.
func (s *Server) SaveTokens(
	stream ptypes.Auth_SaveTokensServer) error {
	count := 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			s.Log.WithFields(logrus.Fields{
				"rpc":     "SaveTokens",
				"code":    http.StatusOK,
				"context": stream.Context,
				"request": req,
				"count":   count,
			}).Info("SaveTokens request processed")
			return nil
		}

		if err != nil {
			s.Log.WithFields(logrus.Fields{
				"rpc":     "SaveTokens",
				"code":    http.StatusInternalServerError,
				"context": stream.Context,
				"request": req,
				"count":   count,
			}).Error(err)
			return err
		}

		v := dauth.Token{}
		v.FromRequest(req)
		ch := s.Tokens.SaveToken(&v)
		for r := range ch {
			if r.Err != nil {
				s.Log.WithFields(logrus.Fields{
					"rpc":     "SaveTokens",
					"code":    http.StatusInternalServerError,
					"context": stream.Context,
					"request": req,
					"count":   count,
				}).Error(err)
				return err
			}
		}

		res := v.ToResponse()
		if err := stream.Send(&res); err != nil {
			s.Log.WithFields(logrus.Fields{
				"rpc":      "SaveTokens",
				"code":     http.StatusInternalServerError,
				"context":  stream.Context,
				"request":  req,
				"response": res,
				"count":    count,
			}).Error(err)
			return err
		}

		count++
	}
}

// DeleteTokens deletes tokens from the database.
func (s *Server) DeleteTokens(ctx context.Context, req *ptypes.TokenRequest) (*ptypes.DeleteResponse, error) {
	q := dauth.TokenFind{}
	if err := q.FromTokenRequest(req); err != nil {
		s.Log.WithFields(logrus.Fields{
			"rpc":     "DeleteTokens",
			"code":    http.StatusBadRequest,
			"context": ctx,
			"request": req,
		}).Error(err)
		return nil, err
	}

	ch := s.Tokens.DeleteTokens(&q)
	count := int64(0)
	for r := range ch {
		if r.Err != nil {
			s.Log.WithFields(logrus.Fields{
				"rpc":     "DeleteTokens",
				"code":    http.StatusInternalServerError,
				"context": ctx,
				"request": req,
			}).Error(r.Err)
			return nil, r.Err
		}

		count += int64(r.Num)
	}

	s.Log.WithFields(logrus.Fields{
		"rpc":     "DeleteTokens",
		"code":    http.StatusOK,
		"context": ctx,
		"request": req,
		"count":   count,
	}).Info("DeleteTokens request processed")
	return &ptypes.DeleteResponse{Num: count}, nil
}

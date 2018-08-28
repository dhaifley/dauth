// Package server implements a server for the authentication services.
package server

import (
	"database/sql"

	"github.com/dhaifley/dauth/lib"
	"github.com/dhaifley/dlib"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Server values implement API server functionality.
type Server struct {
	SQL       dlib.SQLExecutor
	Tokens    lib.TokenAccessor
	Users     lib.UserAccessor
	Perms     lib.PermAccessor
	UserPerms lib.UserPermAccessor
	Log       logrus.FieldLogger
	Router    *mux.Router
}

// ConnectSQL connects to the cloud SQL database.
func (s *Server) ConnectSQL(dbs dlib.SQLExecutor) error {
	if s.Log != nil {
		s.Log.Info("Connecting to SQL database")
	}

	s.SQL = dbs
	if s.SQL == nil {
		db, err := sql.Open("postgres", viper.GetString("sql"))
		if err != nil {
			return err
		}

		db.SetMaxOpenConns(20)
		s.SQL = &dlib.SQLSession{DB: db}
	}

	s.Tokens = lib.NewTokenAccessor(s.SQL)
	s.Users = lib.NewUserAccessor(s.SQL)
	err := s.SQL.Ping()
	if err != nil {
		return err
	}

	if s.Log != nil {
		s.Log.WithField("database", dbs).Info("SQL database connected")
	}

	return nil
}

// Close releases all server resources for shutdown.
func (s *Server) Close() {
	if s.SQL != nil {
		s.SQL.Close()
	}
}

package server

import (
	"database/sql"
	"testing"

	"github.com/dhaifley/dlib"
)

type MockDBSession struct{}

func (m *MockDBSession) Close() error {
	return nil
}

func (m *MockDBSession) Exec(query string, args ...interface{}) (sql.Result, error) {
	return nil, nil
}

func (m *MockDBSession) Query(query string, args ...interface{}) (dlib.SQLRows, error) {
	return nil, nil
}

func (m *MockDBSession) Ping() error {
	return nil
}

func (m *MockDBSession) Stats() sql.DBStats {
	return sql.DBStats{OpenConnections: 1}
}

func TestServerConnectSQL(t *testing.T) {
	s := Server{}
	err := s.ConnectSQL(&MockDBSession{})
	if err != nil {
		t.Error(err)
	}

	defer s.Close()
}

func TestServerClose(t *testing.T) {
	s := Server{}
	defer s.Close()
}

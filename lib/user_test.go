package lib

import (
	"database/sql"
	"testing"

	"github.com/dhaifley/dlib"
	"github.com/dhaifley/dlib/dauth"
)

type MockUserResult struct{}

func (m *MockUserResult) LastInsertId() (int64, error) {
	return 1, nil
}

func (m *MockUserResult) RowsAffected() (int64, error) {
	return 1, nil
}

type MockUserRows struct {
	row int
}

func (m *MockUserRows) Close() error {
	return nil
}

func (m *MockUserRows) Next() bool {
	m.row++
	if m.row > 1 {
		return false
	}

	return true
}

func (m *MockUserRows) Scan(dest ...interface{}) error {
	switch v := dest[0].(type) {
	case *int64:
		*v = int64(1)
	case *int:
		*v = 1
	default:
		return dlib.NewError(500, "Invalid type")
	}

	if len(dest) > 1 {
		switch v := dest[1].(type) {
		case *string:
			*v = "test"
		default:
			return dlib.NewError(500, "Invalid type")
		}
	}

	if len(dest) > 2 {
		switch v := dest[2].(type) {
		case *string:
			*v = "test"
		default:
			return dlib.NewError(500, "Invalid type")
		}
	}

	if len(dest) > 3 {
		switch v := dest[3].(type) {
		case *sql.NullString:
			*v = sql.NullString{Valid: true, String: "test"}
		default:
			return dlib.NewError(500, "Invalid type")
		}
	}

	if len(dest) > 4 {
		switch v := dest[4].(type) {
		case *sql.NullString:
			*v = sql.NullString{Valid: true, String: "test"}
		default:
			return dlib.NewError(500, "Invalid type")
		}
	}

	return nil
}

type MockUserDBSession struct{}

func (m *MockUserDBSession) Close() error {
	return nil
}

func (m *MockUserDBSession) Exec(query string, args ...interface{}) (sql.Result, error) {
	mr := MockUserResult{}
	return &mr, nil
}

func (m *MockUserDBSession) Query(query string, args ...interface{}) (dlib.SQLRows, error) {
	mr := MockUserRows{}
	return &mr, nil
}

func (m *MockUserDBSession) Ping() error {
	return nil
}

func (m *MockUserDBSession) Stats() sql.DBStats {
	return sql.DBStats{OpenConnections: 0}
}

func TestNewUserAccessor(t *testing.T) {
	mdbs := MockUserDBSession{}
	ma := NewUserAccessor(&mdbs)
	_, ok := ma.(UserAccessor)
	if !ok {
		t.Errorf("Type expected: UserAccessor, got: %T", ma)
	}
}

func TestUserAccessGetUsers(t *testing.T) {
	mdbs := MockUserDBSession{}
	ma := NewUserAccessor(&mdbs)
	var as []dauth.User
	id := int64(1)
	c := ma.GetUsers(&dauth.UserFind{ID: &id})
	for r := range c {
		if r.Err != nil {
			t.Error(r.Err)
		}

		switch v := r.Val.(type) {
		case dauth.User:
			as = append(as, v)
		default:
			t.Errorf("Invalid data type returned")
		}
	}

	expected := "test"
	if as[0].User != expected {
		t.Errorf("User expected: %v, got: %v", expected, as[0].User)
	}
}

func TestUserAccessGetUserByID(t *testing.T) {
	mdbs := MockUserDBSession{}
	ma := NewUserAccessor(&mdbs)
	var a dauth.User
	c := ma.GetUserByID(1)
	for r := range c {
		if r.Err != nil {
			t.Error(r.Err)
		}

		switch v := r.Val.(type) {
		case dauth.User:
			a = v
		default:
			t.Errorf("Invalid data type returned")
		}
	}

	expected := "test"
	if a.User != expected {
		t.Errorf("User expected: %v, got: %v", expected, a.User)
	}
}

func TestUserAccessDeleteUserByID(t *testing.T) {
	mdbs := MockUserDBSession{}
	ma := NewUserAccessor(&mdbs)
	c := ma.DeleteUserByID(1)
	var n int
	for r := range c {
		if r.Err != nil {
			t.Error(r.Err)
		}

		n = r.Num
	}

	expected := 1
	if n != expected {
		t.Errorf("Delete count expected: %v, got: %v", expected, n)
	}
}

func TestUserAccessDeleteUsers(t *testing.T) {
	mdbs := MockUserDBSession{}
	ma := NewUserAccessor(&mdbs)
	id := "test"
	c := ma.DeleteUsers(&dauth.UserFind{User: &id})
	var n int
	for r := range c {
		if r.Err != nil {
			t.Error(r.Err)
		}

		n = r.Num
	}

	expected := 1
	if n != expected {
		t.Errorf("Delete count expected: %v, got: %v", expected, n)
	}
}

func TestUserAccessSaveUser(t *testing.T) {
	a := dauth.User{User: "test"}
	mdbs := MockUserDBSession{}
	ma := NewUserAccessor(&mdbs)
	c := ma.SaveUser(&a)
	for r := range c {
		if r.Err != nil {
			t.Error(r.Err)
		}

		switch v := r.Val.(type) {
		case dauth.User:
			a = v
		default:
			t.Errorf("Invalid data type returned")
		}
	}

	expected := "test"
	if a.User != expected {
		t.Errorf("User expected: %v, got: %v", expected, a.User)
	}
}

func TestUserAccessSaveUsers(t *testing.T) {
	a := []dauth.User{dauth.User{User: "test"}}
	mdbs := MockUserDBSession{}
	ma := NewUserAccessor(&mdbs)
	c := ma.SaveUsers(a)
	for r := range c {
		if r.Err != nil {
			t.Error(r.Err)
		}
	}

	expected := "test"
	if a[0].User != expected {
		t.Errorf("User expected: %v, got: %v", expected, a[0].User)
	}
}

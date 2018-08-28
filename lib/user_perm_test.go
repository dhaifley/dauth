package lib

import (
	"database/sql"
	"testing"

	"github.com/dhaifley/dlib"
	"github.com/dhaifley/dlib/dauth"
)

type MockUserPermResult struct{}

func (m *MockUserPermResult) LastInsertId() (int64, error) {
	return 1, nil
}

func (m *MockUserPermResult) RowsAffected() (int64, error) {
	return 1, nil
}

type MockUserPermRows struct {
	row int
}

func (m *MockUserPermRows) Close() error {
	return nil
}

func (m *MockUserPermRows) Next() bool {
	m.row++
	if m.row > 1 {
		return false
	}

	return true
}

func (m *MockUserPermRows) Scan(dest ...interface{}) error {
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
		case *int64:
			*v = int64(1)
		default:
			return dlib.NewError(500, "Invalid type")
		}
	}

	if len(dest) > 2 {
		switch v := dest[1].(type) {
		case *int64:
			*v = int64(1)
		default:
			return dlib.NewError(500, "Invalid type")
		}
	}

	return nil
}

type MockUserPermDBSession struct{}

func (m *MockUserPermDBSession) Close() error {
	return nil
}

func (m *MockUserPermDBSession) Exec(query string, args ...interface{}) (sql.Result, error) {
	mr := MockUserPermResult{}
	return &mr, nil
}

func (m *MockUserPermDBSession) Query(query string, args ...interface{}) (dlib.SQLRows, error) {
	mr := MockUserPermRows{}
	return &mr, nil
}

func (m *MockUserPermDBSession) Ping() error {
	return nil
}

func (m *MockUserPermDBSession) Stats() sql.DBStats {
	return sql.DBStats{OpenConnections: 0}
}

func TestNewUserPermAccessor(t *testing.T) {
	mdbs := MockUserPermDBSession{}
	ma := NewUserPermAccessor(&mdbs)
	_, ok := ma.(UserPermAccessor)
	if !ok {
		t.Errorf("Type expected: UserPermAccessor, got: %T", ma)
	}
}

func TestUserPermAccessGetUserPerms(t *testing.T) {
	mdbs := MockUserPermDBSession{}
	ma := NewUserPermAccessor(&mdbs)
	var as []dauth.UserPerm
	id := int64(1)
	c := ma.GetUserPerms(&dauth.UserPermFind{ID: &id})
	for r := range c {
		if r.Err != nil {
			t.Error(r.Err)
		}

		switch v := r.Val.(type) {
		case dauth.UserPerm:
			as = append(as, v)
		default:
			t.Errorf("Invalid data type returned")
		}
	}

	if as[0].ID != 1 {
		t.Errorf("ID expected: 1, got: %v", as[0].ID)
	}
}

func TestUserPermAccessGetUserPermByID(t *testing.T) {
	mdbs := MockUserPermDBSession{}
	ma := NewUserPermAccessor(&mdbs)
	var a dauth.UserPerm
	c := ma.GetUserPermByID(1)
	for r := range c {
		if r.Err != nil {
			t.Error(r.Err)
		}

		switch v := r.Val.(type) {
		case dauth.UserPerm:
			a = v
		default:
			t.Errorf("Invalid data type returned")
		}
	}

	if a.ID != 1 {
		t.Errorf("ID expected: 1, got: %v", a.ID)
	}
}

func TestUserPermAccessDeleteUserPermByID(t *testing.T) {
	mdbs := MockUserPermDBSession{}
	ma := NewUserPermAccessor(&mdbs)
	c := ma.DeleteUserPermByID(1)
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

func TestUserPermAccessDeleteUserPerms(t *testing.T) {
	mdbs := MockUserPermDBSession{}
	ma := NewUserPermAccessor(&mdbs)
	id := int64(1)
	c := ma.DeleteUserPerms(&dauth.UserPermFind{ID: &id})
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

func TestUserPermAccessSaveUserPerm(t *testing.T) {
	a := dauth.UserPerm{ID: 1}
	mdbs := MockUserPermDBSession{}
	ma := NewUserPermAccessor(&mdbs)
	c := ma.SaveUserPerm(&a)
	for r := range c {
		if r.Err != nil {
			t.Error(r.Err)
		}

		switch v := r.Val.(type) {
		case dauth.UserPerm:
			a = v
		default:
			t.Errorf("Invalid data type returned")
		}
	}

	if a.ID != 1 {
		t.Errorf("ID expected: 1, got: %v", a.ID)
	}
}

func TestUserPermAccessSaveUserPerms(t *testing.T) {
	a := []dauth.UserPerm{dauth.UserPerm{ID: 1}}
	mdbs := MockUserPermDBSession{}
	ma := NewUserPermAccessor(&mdbs)
	c := ma.SaveUserPerms(a)
	for r := range c {
		if r.Err != nil {
			t.Error(r.Err)
		}
	}

	if a[0].ID != 1 {
		t.Errorf("ID expected: 1, got: %v", a[0].ID)
	}
}

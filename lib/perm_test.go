package lib

import (
	"database/sql"
	"testing"

	"github.com/dhaifley/dlib"
	"github.com/dhaifley/dlib/dauth"
)

type MockPermResult struct{}

func (m *MockPermResult) LastInsertId() (int64, error) {
	return 1, nil
}

func (m *MockPermResult) RowsAffected() (int64, error) {
	return 1, nil
}

type MockPermRows struct {
	row int
}

func (m *MockPermRows) Close() error {
	return nil
}

func (m *MockPermRows) Next() bool {
	m.row++
	if m.row > 1 {
		return false
	}

	return true
}

func (m *MockPermRows) Scan(dest ...interface{}) error {
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

	return nil
}

type MockPermDBSession struct{}

func (m *MockPermDBSession) Close() error {
	return nil
}

func (m *MockPermDBSession) Exec(query string, args ...interface{}) (sql.Result, error) {
	mr := MockPermResult{}
	return &mr, nil
}

func (m *MockPermDBSession) Query(query string, args ...interface{}) (dlib.SQLRows, error) {
	mr := MockPermRows{}
	return &mr, nil
}

func (m *MockPermDBSession) Ping() error {
	return nil
}

func (m *MockPermDBSession) Stats() sql.DBStats {
	return sql.DBStats{OpenConnections: 0}
}

func TestNewPermAccessor(t *testing.T) {
	mdbs := MockPermDBSession{}
	ma := NewPermAccessor(&mdbs)
	_, ok := ma.(PermAccessor)
	if !ok {
		t.Errorf("Type expected: PermAccessor, got: %T", ma)
	}
}

func TestPermAccessGetPerms(t *testing.T) {
	mdbs := MockPermDBSession{}
	ma := NewPermAccessor(&mdbs)
	var as []dauth.Perm
	id := int64(1)
	c := ma.GetPerms(&dauth.PermFind{ID: &id})
	for r := range c {
		if r.Err != nil {
			t.Error(r.Err)
		}

		switch v := r.Val.(type) {
		case dauth.Perm:
			as = append(as, v)
		default:
			t.Errorf("Invalid data type returned")
		}
	}

	if as[0].ID != 1 {
		t.Errorf("ID expected: 1, got: %v", as[0].ID)
	}
}

func TestPermAccessGetPermByID(t *testing.T) {
	mdbs := MockPermDBSession{}
	ma := NewPermAccessor(&mdbs)
	var a dauth.Perm
	c := ma.GetPermByID(1)
	for r := range c {
		if r.Err != nil {
			t.Error(r.Err)
		}

		switch v := r.Val.(type) {
		case dauth.Perm:
			a = v
		default:
			t.Errorf("Invalid data type returned")
		}
	}

	if a.ID != 1 {
		t.Errorf("ID expected: 1, got: %v", a.ID)
	}
}

func TestPermAccessDeletePermByID(t *testing.T) {
	mdbs := MockPermDBSession{}
	ma := NewPermAccessor(&mdbs)
	c := ma.DeletePermByID(1)
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

func TestPermAccessDeletePerms(t *testing.T) {
	mdbs := MockPermDBSession{}
	ma := NewPermAccessor(&mdbs)
	id := int64(1)
	c := ma.DeletePerms(&dauth.PermFind{ID: &id})
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

func TestPermAccessSavePerm(t *testing.T) {
	a := dauth.Perm{ID: 1}
	mdbs := MockPermDBSession{}
	ma := NewPermAccessor(&mdbs)
	c := ma.SavePerm(&a)
	for r := range c {
		if r.Err != nil {
			t.Error(r.Err)
		}

		switch v := r.Val.(type) {
		case dauth.Perm:
			a = v
		default:
			t.Errorf("Invalid data type returned")
		}
	}

	if a.ID != 1 {
		t.Errorf("ID expected: 1, got: %v", a.ID)
	}
}

func TestPermAccessSavePerms(t *testing.T) {
	a := []dauth.Perm{dauth.Perm{ID: 1}}
	mdbs := MockPermDBSession{}
	ma := NewPermAccessor(&mdbs)
	c := ma.SavePerms(a)
	for r := range c {
		if r.Err != nil {
			t.Error(r.Err)
		}
	}

	if a[0].ID != 1 {
		t.Errorf("ID expected: 1, got: %v", a[0].ID)
	}
}

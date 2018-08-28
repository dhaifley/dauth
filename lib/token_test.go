package lib

import (
	"database/sql"
	"testing"
	"time"

	"github.com/dhaifley/dlib"
	"github.com/dhaifley/dlib/dauth"
)

type MockTokenResult struct{}

func (fr *MockTokenResult) LastInsertId() (int64, error) {
	return 1, nil
}

func (fr *MockTokenResult) RowsAffected() (int64, error) {
	return 1, nil
}

type MockTokenRows struct {
	row int
}

func (frs *MockTokenRows) Close() error {
	return nil
}

func (frs *MockTokenRows) Next() bool {
	frs.row++
	if frs.row > 1 {
		return false
	}

	return true
}

func (frs *MockTokenRows) Scan(dest ...interface{}) error {
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
		case *int64:
			*v = 1
		default:
			return dlib.NewError(500, "Invalid type")
		}
	}

	if len(dest) > 3 {
		switch v := dest[3].(type) {
		case *dlib.NullTime:
			dt := time.Date(1983, 2, 2, 0, 0, 0, 0, time.Local)
			*v = dlib.NullTime{Time: dt, Valid: true}
		default:
			return dlib.NewError(500, "Invalid type")
		}
	}

	if len(dest) > 4 {
		switch v := dest[4].(type) {
		case *dlib.NullTime:
			dt := time.Date(1983, 2, 2, 0, 0, 0, 0, time.Local)
			*v = dlib.NullTime{Time: dt, Valid: true}
		default:
			return dlib.NewError(500, "Invalid type")
		}
	}

	return nil
}

type MockTokenDBSession struct{}

func (m *MockTokenDBSession) Close() error {
	return nil
}

func (m *MockTokenDBSession) Exec(query string, args ...interface{}) (sql.Result, error) {
	fr := MockTokenResult{}
	return &fr, nil
}

func (m *MockTokenDBSession) Query(query string, args ...interface{}) (dlib.SQLRows, error) {
	fr := MockTokenRows{}
	return &fr, nil
}

func (m *MockTokenDBSession) Ping() error {
	return nil
}

func (m *MockTokenDBSession) Stats() sql.DBStats {
	return sql.DBStats{OpenConnections: 0}
}

func TestNewTokenAccessor(t *testing.T) {
	mdbs := MockTokenDBSession{}
	ma := NewTokenAccessor(&mdbs)
	_, ok := ma.(TokenAccessor)
	if !ok {
		t.Errorf("Type expected: TokenAccessor, got: %T", ma)
	}
}

func TestTokenAccessGetTokens(t *testing.T) {
	mdbs := MockTokenDBSession{}
	ma := NewTokenAccessor(&mdbs)
	var as []dauth.Token
	id := int64(1)
	c := ma.GetTokens(&dauth.TokenFind{ID: &id})
	for r := range c {
		if r.Err != nil {
			t.Error(r.Err)
		}

		switch v := r.Val.(type) {
		case dauth.Token:
			as = append(as, v)
		default:
			t.Errorf("Invalid data type returned")
		}
	}

	expected := int64(1)
	if as[0].ID != expected {
		t.Errorf("ID expected: %v, got: %v", expected, as[0].ID)
	}
}

func TestTokenAccessGetTokenByID(t *testing.T) {
	mdbs := MockTokenDBSession{}
	ma := NewTokenAccessor(&mdbs)
	var a dauth.Token
	c := ma.GetTokenByID(1)
	for r := range c {
		if r.Err != nil {
			t.Error(r.Err)
		}

		switch v := r.Val.(type) {
		case dauth.Token:
			a = v
		default:
			t.Errorf("Invalid data type returned")
		}
	}

	expected := int64(1)
	if a.ID != expected {
		t.Errorf("ID expected: %v, got: %v", expected, a.ID)
	}
}

func TestTokenAccessDeleteTokenByID(t *testing.T) {
	mdbs := MockTokenDBSession{}
	ma := NewTokenAccessor(&mdbs)
	c := ma.DeleteTokenByID(1)
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

func TestTokenAccessDeleteTokens(t *testing.T) {
	mdbs := MockTokenDBSession{}
	ma := NewTokenAccessor(&mdbs)
	id := int64(1)
	c := ma.DeleteTokens(&dauth.TokenFind{ID: &id})
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

func TestTokenAccessSaveToken(t *testing.T) {
	a := dauth.Token{ID: 1}
	mdbs := MockTokenDBSession{}
	ma := NewTokenAccessor(&mdbs)
	c := ma.SaveToken(&a)
	for r := range c {
		if r.Err != nil {
			t.Error(r.Err)
		}

		switch v := r.Val.(type) {
		case dauth.Token:
			a = v
		default:
			t.Errorf("Invalid data type returned")
		}
	}

	expected := int64(1)
	if a.ID != expected {
		t.Errorf("ID expected: %v, got: %v", expected, a.ID)
	}
}

func TestTokenAccessSaveTokens(t *testing.T) {
	a := []dauth.Token{dauth.Token{ID: 1}}
	mdbs := MockTokenDBSession{}
	ma := NewTokenAccessor(&mdbs)
	c := ma.SaveTokens(a)
	for r := range c {
		if r.Err != nil {
			t.Error(r.Err)
		}
	}

	expected := int64(1)
	if a[0].ID != expected {
		t.Errorf("ID expected: %v, got: %v", expected, a[0].ID)
	}
}

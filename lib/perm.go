package lib

import (
	"github.com/dhaifley/dlib"
	"github.com/dhaifley/dlib/dauth"
)

// PermAccess values are used to access perm records in the database.
type PermAccess struct {
	DBS dlib.SQLExecutor
}

// PermAccessor is an interface describing values capable of providing
// access to perm records in the database.
type PermAccessor interface {
	GetPerms(opt *dauth.PermFind) <-chan dlib.Result
	GetPermByID(id int64) <-chan dlib.Result
	DeletePerms(opt *dauth.PermFind) <-chan dlib.Result
	DeletePermByID(id int64) <-chan dlib.Result
	SavePerm(t *dauth.Perm) <-chan dlib.Result
	SavePerms(t []dauth.Perm) <-chan dlib.Result
}

// NewPermAccessor creates a new PermAccess instance and
// returns a pointer to it.
func NewPermAccessor(dbs dlib.SQLExecutor) PermAccessor {
	ua := PermAccess{DBS: dbs}
	return &ua
}

// GetPerms finds perm values in the database.
func (pa *PermAccess) GetPerms(opt *dauth.PermFind) <-chan dlib.Result {
	ch := make(chan dlib.Result, 256)

	go func() {
		defer close(ch)
		rows, err := pa.DBS.Query(`
			SELECT
				p.id,
				p.service,
				p.name
			FROM get_perms($1, $2, $3) AS p`,
			opt.ID,
			opt.Service,
			opt.Name)
		if err != nil {
			ch <- dlib.Result{Err: err}
			return
		}

		defer rows.Close()
		for rows.Next() {
			r := dauth.PermRow{}
			if err := rows.Scan(
				&r.ID,
				&r.Service,
				&r.Name,
			); err != nil {
				ch <- dlib.Result{Err: err}
				continue
			}

			v := r.ToPerm()
			ch <- dlib.Result{Val: v, Num: 1}
		}
	}()

	return ch
}

// GetPermByID finds a perm value in the database by ID.
func (pa *PermAccess) GetPermByID(id int64) <-chan dlib.Result {
	opt := dauth.PermFind{ID: &id}
	return pa.GetPerms(&opt)
}

// DeletePerms deletes perm values from the database.
func (pa *PermAccess) DeletePerms(opt *dauth.PermFind) <-chan dlib.Result {
	ch := make(chan dlib.Result, 256)

	go func() {
		defer close(ch)
		rows, err := pa.DBS.Query(
			"SELECT delete_perms($1, $2, $3) AS num",
			opt.ID,
			opt.Service,
			opt.Name)
		if err != nil {
			ch <- dlib.Result{Err: err}
			return
		}

		defer rows.Close()
		n := 0
		for rows.Next() {
			r := struct{ Num int }{Num: 0}
			if err := rows.Scan(&r.Num); err != nil {
				ch <- dlib.Result{Err: err}
				continue
			}

			n += r.Num
		}

		ch <- dlib.Result{Num: n, Err: nil}
	}()

	return ch
}

// DeletePermByID deletes a perm value from the database by ID.
func (pa *PermAccess) DeletePermByID(id int64) <-chan dlib.Result {
	opt := dauth.PermFind{ID: &id}
	return pa.DeletePerms(&opt)
}

// SavePerm saves a perm value to the database.
func (pa *PermAccess) SavePerm(u *dauth.Perm) <-chan dlib.Result {
	ch := make(chan dlib.Result, 256)

	go func() {
		defer close(ch)
		rows, err := pa.DBS.Query(
			"SELECT save_perm($1, $2, $3) AS id",
			u.ID,
			u.Service,
			u.Name)
		if err != nil {
			ch <- dlib.Result{Err: err}
			return
		}

		defer rows.Close()
		for rows.Next() {
			r := struct{ ID int64 }{ID: 0}
			if err := rows.Scan(&r.ID); err != nil {
				ch <- dlib.Result{Err: err}
				continue
			}

			u.ID = r.ID
		}

		ch <- dlib.Result{Val: *u, Err: nil}
	}()

	return ch
}

// SavePerms saves a slice of perm values to the database.
func (pa *PermAccess) SavePerms(u []dauth.Perm) <-chan dlib.Result {
	ch := make(chan dlib.Result, 256)

	go func() {
		defer close(ch)
		if u != nil {
			for _, v := range u {
				for sr := range pa.SavePerm(&v) {
					ch <- sr
				}
			}
		}
	}()

	return ch
}

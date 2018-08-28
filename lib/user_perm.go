package lib

import (
	"github.com/dhaifley/dlib"
	"github.com/dhaifley/dlib/dauth"
)

// UserPermAccess values are used to access user_perm records in the database.
type UserPermAccess struct {
	DBS dlib.SQLExecutor
}

// UserPermAccessor is an interface describing values capable of providing
// access to user_perm records in the database.
type UserPermAccessor interface {
	GetUserPerms(opt *dauth.UserPermFind) <-chan dlib.Result
	GetUserPermByID(id int64) <-chan dlib.Result
	DeleteUserPerms(opt *dauth.UserPermFind) <-chan dlib.Result
	DeleteUserPermByID(id int64) <-chan dlib.Result
	SaveUserPerm(t *dauth.UserPerm) <-chan dlib.Result
	SaveUserPerms(t []dauth.UserPerm) <-chan dlib.Result
}

// NewUserPermAccessor creates a new UserPermAccess instance and
// returns a pointer to it.
func NewUserPermAccessor(dbs dlib.SQLExecutor) UserPermAccessor {
	ua := UserPermAccess{DBS: dbs}
	return &ua
}

// GetUserPerms finds user_perm values in the database.
func (upa *UserPermAccess) GetUserPerms(opt *dauth.UserPermFind) <-chan dlib.Result {
	ch := make(chan dlib.Result, 256)

	go func() {
		defer close(ch)
		rows, err := upa.DBS.Query(`
			SELECT
				p.id,
				p.user_id,
				p.perm_id
			FROM get_user_perms($1, $2, $3) AS p`,
			opt.ID,
			opt.UserID,
			opt.PermID)
		if err != nil {
			ch <- dlib.Result{Err: err}
			return
		}

		defer rows.Close()
		for rows.Next() {
			r := dauth.UserPermRow{}
			if err := rows.Scan(
				&r.ID,
				&r.UserID,
				&r.PermID,
			); err != nil {
				ch <- dlib.Result{Err: err}
				continue
			}

			v := r.ToUserPerm()
			ch <- dlib.Result{Val: v, Num: 1}
		}
	}()

	return ch
}

// GetUserPermByID finds a user_perm value in the database by ID.
func (upa *UserPermAccess) GetUserPermByID(id int64) <-chan dlib.Result {
	opt := dauth.UserPermFind{ID: &id}
	return upa.GetUserPerms(&opt)
}

// DeleteUserPerms deletes user_perm values from the database.
func (upa *UserPermAccess) DeleteUserPerms(opt *dauth.UserPermFind) <-chan dlib.Result {
	ch := make(chan dlib.Result, 256)

	go func() {
		defer close(ch)
		rows, err := upa.DBS.Query(
			"SELECT delete_user_perms($1, $2, $3) AS num",
			opt.ID,
			opt.UserID,
			opt.PermID)
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

// DeleteUserPermByID deletes a user_perm value from the database by ID.
func (upa *UserPermAccess) DeleteUserPermByID(id int64) <-chan dlib.Result {
	opt := dauth.UserPermFind{ID: &id}
	return upa.DeleteUserPerms(&opt)
}

// SaveUserPerm saves a user_perm value to the database.
func (upa *UserPermAccess) SaveUserPerm(u *dauth.UserPerm) <-chan dlib.Result {
	ch := make(chan dlib.Result, 256)

	go func() {
		defer close(ch)
		rows, err := upa.DBS.Query(
			"SELECT save_user_perm($1, $2, $3) AS id",
			u.ID,
			u.UserID,
			u.PermID)
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

// SaveUserPerms saves a slice of user_perm values to the database.
func (upa *UserPermAccess) SaveUserPerms(u []dauth.UserPerm) <-chan dlib.Result {
	ch := make(chan dlib.Result, 256)

	go func() {
		defer close(ch)
		if u != nil {
			for _, v := range u {
				for sr := range upa.SaveUserPerm(&v) {
					ch <- sr
				}
			}
		}
	}()

	return ch
}

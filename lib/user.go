package lib

import (
	"github.com/dhaifley/dlib"
	"github.com/dhaifley/dlib/dauth"
)

// UserAccess values are used to access user records in the database.
type UserAccess struct {
	DBS dlib.SQLExecutor
}

// UserAccessor is an interface describing values capable of providing
// access to user records in the database.
type UserAccessor interface {
	GetUsers(opt *dauth.UserFind) <-chan dlib.Result
	GetUserByID(id int64) <-chan dlib.Result
	DeleteUsers(opt *dauth.UserFind) <-chan dlib.Result
	DeleteUserByID(id int64) <-chan dlib.Result
	SaveUser(t *dauth.User) <-chan dlib.Result
	SaveUsers(t []dauth.User) <-chan dlib.Result
}

// NewUserAccessor creates a new UserAccess instance and
// returns a pointer to it.
func NewUserAccessor(dbs dlib.SQLExecutor) UserAccessor {
	ua := UserAccess{DBS: dbs}
	return &ua
}

// GetUsers finds user values in the database.
func (ua *UserAccess) GetUsers(opt *dauth.UserFind) <-chan dlib.Result {
	ch := make(chan dlib.Result, 256)

	go func() {
		defer close(ch)
		rows, err := ua.DBS.Query(`
			SELECT
				u.id,
				u.user,
				u.pass,
				u.name,
				u.email
			FROM get_users($1, $2, $3, $4, $5) AS u`,
			opt.ID,
			opt.User,
			opt.Pass,
			opt.Name,
			opt.Email)
		if err != nil {
			ch <- dlib.Result{Err: err}
			return
		}

		defer rows.Close()
		for rows.Next() {
			r := dauth.UserRow{}
			if err := rows.Scan(
				&r.ID,
				&r.User,
				&r.Pass,
				&r.Name,
				&r.Email,
			); err != nil {
				ch <- dlib.Result{Err: err}
				continue
			}

			v := r.ToUser()
			ch <- dlib.Result{Val: v, Num: 1}
		}
	}()

	return ch
}

// GetUserByID finds a User value in the database by ID.
func (ua *UserAccess) GetUserByID(id int64) <-chan dlib.Result {
	opt := dauth.UserFind{ID: &id}
	return ua.GetUsers(&opt)
}

// DeleteUsers deletes User values from the database.
func (ua *UserAccess) DeleteUsers(opt *dauth.UserFind) <-chan dlib.Result {
	ch := make(chan dlib.Result, 256)

	go func() {
		defer close(ch)
		rows, err := ua.DBS.Query(
			"SELECT delete_users$1, $2, $3, $4, $5) AS num",
			opt.ID,
			opt.User,
			opt.Pass,
			opt.Name,
			opt.Email)
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

// DeleteUserByID deletes a User value from the database by ID.
func (ua *UserAccess) DeleteUserByID(id int64) <-chan dlib.Result {
	opt := dauth.UserFind{ID: &id}
	return ua.DeleteUsers(&opt)
}

// SaveUser saves a User value to the database.
func (ua *UserAccess) SaveUser(u *dauth.User) <-chan dlib.Result {
	ch := make(chan dlib.Result, 256)

	go func() {
		defer close(ch)
		rows, err := ua.DBS.Query(
			"SELECT save_user($1, $2, $3, $4, $5) AS id",
			u.ID,
			u.User,
			u.Pass,
			u.Name,
			u.Email)
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

// SaveUsers saves a slice of User values to the database.
func (ua *UserAccess) SaveUsers(u []dauth.User) <-chan dlib.Result {
	ch := make(chan dlib.Result, 256)

	go func() {
		defer close(ch)
		if u != nil {
			for _, v := range u {
				for sr := range ua.SaveUser(&v) {
					ch <- sr
				}
			}
		}
	}()

	return ch
}

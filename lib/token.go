package lib

import (
	"github.com/dhaifley/dlib"
	"github.com/dhaifley/dlib/dauth"
)

// TokenAccess values are used to access token records in the database.
type TokenAccess struct {
	DBS dlib.SQLExecutor
}

// TokenAccessor is an interface describing values capable of providing
// access to token records in the database.
type TokenAccessor interface {
	GetTokens(opt *dauth.TokenFind) <-chan dlib.Result
	GetTokenByID(id int64) <-chan dlib.Result
	DeleteTokens(opt *dauth.TokenFind) <-chan dlib.Result
	DeleteTokenByID(id int64) <-chan dlib.Result
	SaveToken(t *dauth.Token) <-chan dlib.Result
	SaveTokens(t []dauth.Token) <-chan dlib.Result
}

// NewTokenAccessor creates a new TokenAccess value for database access.
func NewTokenAccessor(dbs dlib.SQLExecutor) TokenAccessor {
	ta := TokenAccess{DBS: dbs}
	return &ta
}

// GetTokens finds token values in the database.
func (ta *TokenAccess) GetTokens(opt *dauth.TokenFind) <-chan dlib.Result {
	ch := make(chan dlib.Result, 256)

	go func() {
		defer close(ch)
		rows, err := ta.DBS.Query(`
			SELECT
				t.id,
				t.token,
				t.user,
				t.created,
				t.expires
			FROM get_tokens($1, $2, $3, $4, $5, $6, $7, $8) AS t`,
			opt.ID,
			opt.Token,
			opt.UserID,
			opt.Created,
			opt.Expires,
			opt.Start,
			opt.End,
			opt.Old)
		if err != nil {
			ch <- dlib.Result{Err: err}
			return
		}

		defer rows.Close()
		for rows.Next() {
			r := dauth.TokenRow{}
			if err := rows.Scan(
				&r.ID,
				&r.Token,
				&r.UserID,
				&r.Created,
				&r.Expires,
			); err != nil {
				ch <- dlib.Result{Err: err}
				continue
			}

			v := r.ToToken()
			ch <- dlib.Result{Val: v, Num: 1}
		}
	}()

	return ch
}

// GetTokenByID finds a Token value in the database by ID.
func (ta *TokenAccess) GetTokenByID(id int64) <-chan dlib.Result {
	opt := dauth.TokenFind{ID: &id}
	return ta.GetTokens(&opt)
}

// DeleteTokens deletes Token values from the database.
func (ta *TokenAccess) DeleteTokens(opt *dauth.TokenFind) <-chan dlib.Result {
	ch := make(chan dlib.Result, 256)

	go func() {
		defer close(ch)
		rows, err := ta.DBS.Query(
			"SELECT delete_tokens($1, $2, $3, $4, $5, $6, $7, $8) AS num",
			opt.ID,
			opt.Token,
			opt.UserID,
			opt.Created,
			opt.Expires,
			opt.Start,
			opt.End,
			opt.Old)
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

// DeleteTokenByID deletes a Token value from the database by ID.
func (ta *TokenAccess) DeleteTokenByID(id int64) <-chan dlib.Result {
	opt := dauth.TokenFind{ID: &id}
	return ta.DeleteTokens(&opt)
}

// SaveToken saves a Token value to the database.
func (ta *TokenAccess) SaveToken(t *dauth.Token) <-chan dlib.Result {
	ch := make(chan dlib.Result, 256)

	go func() {
		defer close(ch)
		rows, err := ta.DBS.Query(
			"SELECT save_token($1, $2, $3, $4, $5) AS id",
			t.ID,
			t.Token,
			t.UserID,
			t.Created,
			t.Expires)
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

			t.ID = r.ID
		}

		ch <- dlib.Result{Val: *t, Err: nil}
	}()

	return ch
}

// SaveTokens saves a slice of Token values to the database.
func (ta *TokenAccess) SaveTokens(t []dauth.Token) <-chan dlib.Result {
	ch := make(chan dlib.Result, 256)

	go func() {
		defer close(ch)
		if t != nil {
			for _, v := range t {
				for sr := range ta.SaveToken(&v) {
					ch <- sr
				}
			}
		}
	}()

	return ch
}

package dao

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/rand"

	customError "github.com/crosserclaws/Go-000/Week02/custom/error"
	pkgErrors "github.com/pkg/errors"
)

// User is a struct for user information.
type User struct {
	ID   int
	Name string
}

// QueryUserByID returns an user's information.
func QueryUserByID(id int) (*User, error) {
	row, err := fakeQueryUser(id)
	if err != nil {
		if customError.IsPlatformSpecific(err) {
			err = customError.StorageErrToCustomErr(err)
		}
		return nil, pkgErrors.WithStack(err)
	}
	return fakeRowToUser(row), nil
}

// QueryAllUsers returns an array of user information. The array may be empty.
func QueryAllUsers() ([]*User, error) {
	rows, err := fakeQueryUsers()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("No results when query all users.")
			return []*User{}, nil
		}
		return nil, pkgErrors.WithStack(err)
	}
	return fakeRowsToUsers(rows), nil
}

func fakeQueryUser(id int) (*sql.Row, error) {
	switch id {
	case 1:
		return &sql.Row{}, nil
	case 2:
		return nil, sql.ErrNoRows
	default:
		return nil, fmt.Errorf("An example error when querying an user with id=%d", id)
	}
}

func fakeRowToUser(rows *sql.Row) *User {
	return &User{1, "No.1"}
}

func fakeQueryUsers() (*sql.Rows, error) {
	v := rand.Int()
	switch v % 3 {
	case 0:
		return nil, sql.ErrNoRows
	case 1:
		return nil, errors.New("An example error when querying users")
	default:
		return &sql.Rows{}, nil
	}
}

func fakeRowsToUsers(*sql.Rows) []*User {
	return []*User{{1, "No.1"}, {2, "No.2"}}
}

package models

import (
	"fmt"
	"strings"
)

type User struct {
	ID          int64
	LoginName   string
	DisplayName string
	Email       string `xorm:"NOT NULL"`
	Mobile      string
}

// IsLoginNameExist checks if given user login_name exist,
// the user login_name should be noncased unique.
func IsLoginNameExist(name string) (bool, error) {
	if len(name) == 0 {
		return false, nil
	}
	return x.Get(&User{LoginName: strings.ToLower(name)})
}

// CreateUser creates record of a new user.
func CreateUser(u *User) (err error) {
	// if err = IsUsableUsername(u.Name); err != nil {
	// 	return err
	// }

	isExist, err := IsLoginNameExist(u.LoginName)
	if err != nil {
		return err
	} else if isExist {
		return ErrLoginNameAlreadyExist{u.LoginName}
	}

	u.Email = strings.ToLower(u.Email)

	sess := x.NewSession()
	defer sess.Close()
	if err = sess.Begin(); err != nil {
		return err
	}

	if _, err = sess.Insert(u); err != nil {
		return err
	}

	return sess.Commit()
}

func GetUserProfile(id int64) {
	user := &User{ID: id}

	x.Get(user)
	fmt.Printf("%+v\n", user)
}

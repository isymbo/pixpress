package models

import "fmt"

type ErrLoginNameAlreadyExist struct {
	LoginName string
}

func IsErrLoginNameAlreadyExist(err error) bool {
	_, ok := err.(ErrLoginNameAlreadyExist)
	return ok
}

func (err ErrLoginNameAlreadyExist) Error() string {
	return fmt.Sprintf("login_name already exists [login_name: %s]", err.LoginName)
}

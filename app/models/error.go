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

//    _____                                   ___________     __
//   /  _  \   ____  ____  ____   ______ _____\__    ___/___ |  | __ ____   ____
//  /  /_\  \_/ ___\/ ___\/ __ \ /  ___//  ___/ |    | /  _ \|  |/ // __ \ /    \
// /    |    \  \__\  \__\  ___/ \___ \ \___ \  |    |(  <_> )    <\  ___/|   |  \
// \____|__  /\___  >___  >___  >____  >____  > |____| \____/|__|_ \\___  >___|  /
//         \/     \/    \/    \/     \/     \/                    \/    \/     \/

type ErrAccessTokenNotExist struct {
	SHA string
}

func IsErrAccessTokenNotExist(err error) bool {
	_, ok := err.(ErrAccessTokenNotExist)
	return ok
}

func (err ErrAccessTokenNotExist) Error() string {
	return fmt.Sprintf("access token does not exist [sha: %s]", err.SHA)
}

type ErrAccessTokenEmpty struct {
}

func IsErrAccessTokenEmpty(err error) bool {
	_, ok := err.(ErrAccessTokenEmpty)
	return ok
}

func (err ErrAccessTokenEmpty) Error() string {
	return fmt.Sprintf("access token is empty")
}

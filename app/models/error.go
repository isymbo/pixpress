package models

import "fmt"

// .____                 .__         _______
// |    |    ____   ____ |__| ____   \      \ _____    _____   ____
// |    |   /  _ \ / ___\|  |/    \  /   |   \\__  \  /     \_/ __ \
// |    |__(  <_> ) /_/  >  |   |  \/    |    \/ __ \|  Y Y  \  ___/
// |_______ \____/\___  /|__|___|  /\____|__  (____  /__|_|  /\___  >
//         \/    /_____/         \/         \/     \/      \/     \/

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

type ErrEmailAlreadyUsed struct {
	Email string
}

func IsErrEmailAlreadyUsed(err error) bool {
	_, ok := err.(ErrEmailAlreadyUsed)
	return ok
}

func (err ErrEmailAlreadyUsed) Error() string {
	return fmt.Sprintf("e-mail has been used [email: %s]", err.Email)
}

// type ErrUserOwnRepos struct {
// 	UID int64
// }

// func IsErrUserOwnRepos(err error) bool {
// 	_, ok := err.(ErrUserOwnRepos)
// 	return ok
// }

// func (err ErrUserOwnRepos) Error() string {
// 	return fmt.Sprintf("user still has ownership of repositories [uid: %d]", err.UID)
// }

// type ErrUserHasOrgs struct {
// 	UID int64
// }

// func IsErrUserHasOrgs(err error) bool {
// 	_, ok := err.(ErrUserHasOrgs)
// 	return ok
// }

// func (err ErrUserHasOrgs) Error() string {
// 	return fmt.Sprintf("user still has membership of organizations [uid: %d]", err.UID)
// }

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

// .____                 .__           _________
// |    |    ____   ____ |__| ____    /   _____/ ____  __ _________   ____  ____
// |    |   /  _ \ / ___\|  |/    \   \_____  \ /  _ \|  |  \_  __ \_/ ___\/ __ \
// |    |__(  <_> ) /_/  >  |   |  \  /        (  <_> )  |  /|  | \/\  \__\  ___/
// |_______ \____/\___  /|__|___|  / /_______  /\____/|____/ |__|    \___  >___  >
//         \/    /_____/         \/          \/                          \/    \/

type ErrLoginSourceAlreadyExist struct {
	Name string
}

func IsErrLoginSourceAlreadyExist(err error) bool {
	_, ok := err.(ErrLoginSourceAlreadyExist)
	return ok
}

func (err ErrLoginSourceAlreadyExist) Error() string {
	return fmt.Sprintf("login source already exists [name: %s]", err.Name)
}

type ErrLoginSourceInUse struct {
	ID int64
}

func IsErrLoginSourceInUse(err error) bool {
	_, ok := err.(ErrLoginSourceInUse)
	return ok
}

func (err ErrLoginSourceInUse) Error() string {
	return fmt.Sprintf("login source is still used by some users [id: %d]", err.ID)
}

// __________               __
// \______   \____  _______/  |_
//  |     ___/  _ \/  ___/\   __\
//  |    |  (  <_> )___ \  |  |
//  |____|   \____/____  > |__|
//                     \/

type ErrPostNotExist struct {
	ID string
}

func IsErrPostNotExist(err error) bool {
	_, ok := err.(ErrPostNotExist)
	return ok
}

func (err ErrPostNotExist) Error() string {
	return fmt.Sprintf("post does not exist [id: %s]", err.ID)
}

//    _____   __    __                .__                           __
//   /  _  \_/  |__/  |______    ____ |  |__   _____   ____   _____/  |_
//  /  /_\  \   __\   __\__  \ _/ ___\|  |  \ /     \_/ __ \ /    \   __\
// /    |    \  |  |  |  / __ \\  \___|   Y  \  Y Y  \  ___/|   |  \  |
// \____|__  /__|  |__| (____  /\___  >___|  /__|_|  /\___  >___|  /__|
//         \/                \/     \/     \/      \/     \/     \/

type ErrAttachmentNotExist struct {
	ID   int64
	UUID string
}

func IsErrAttachmentNotExist(err error) bool {
	_, ok := err.(ErrAttachmentNotExist)
	return ok
}

func (err ErrAttachmentNotExist) Error() string {
	return fmt.Sprintf("attachment does not exist [id: %d, uuid: %s]", err.ID, err.UUID)
}

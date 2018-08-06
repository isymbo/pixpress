package errors

import "fmt"

type PostNotExist struct {
	ID       int64
	AuthorID int64
}

func IsPostNotExist(err error) bool {
	_, ok := err.(PostNotExist)
	return ok
}

func (err PostNotExist) Error() string {
	return fmt.Sprintf("post does not exist [id: %d, author_id: %d]", err.ID, err.AuthorID)
}

package errors

import "fmt"

type CoverImgNotExist struct {
	ID   int64
	UUID string
}

func IsCoverImgNotExist(err error) bool {
	_, ok := err.(CoverImgNotExist)
	return ok
}

func (err CoverImgNotExist) Error() string {
	return fmt.Sprintf("cover image does not exist [id: %d, uuid: %d]", err.ID, err.UUID)
}

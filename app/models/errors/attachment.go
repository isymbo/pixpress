package errors

import "fmt"

type AttachmentNotExist struct {
	ID   int64
	UUID string
}

func IsAttachmentNotExist(err error) bool {
	_, ok := err.(AttachmentNotExist)
	return ok
}

func (err AttachmentNotExist) Error() string {
	return fmt.Sprintf("attachment does not exist [id: %d, uuid: %d]", err.ID, err.UUID)
}

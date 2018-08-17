package models

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"time"

	"github.com/go-xorm/xorm"
	gouuid "github.com/satori/go.uuid"

	"github.com/isymbo/pixpress/setting"
)

// Attachment represent a attachment of issue/comment/release.
type Attachment struct {
	ID     int64
	UUID   string `xorm:"uuid UNIQUE"`
	PostID int64  `xorm:"INDEX"`
	Name   string

	Created     time.Time `xorm:"-" json:"-"`
	CreatedUnix int64
}

func (a *Attachment) BeforeInsert() {
	a.CreatedUnix = time.Now().Unix()
}

func (a *Attachment) AfterSet(colName string, _ xorm.Cell) {
	switch colName {
	case "created_unix":
		a.Created = time.Unix(a.CreatedUnix, 0).Local()
	}
}

// AttachmentLocalPath returns where attachment is stored in local file system based on given UUID.
func AttachmentLocalPath(uuid string) string {
	return path.Join(setting.Attachment.Path, uuid[0:1], uuid[1:2], uuid)
}

// LocalPath returns where attachment is stored in local file system.
func (attach *Attachment) LocalPath() string {
	return AttachmentLocalPath(attach.UUID)
}

// NewAttachment creates a new attachment object.
func NewAttachment(name string, buf []byte, file multipart.File) (_ *Attachment, err error) {
	uuid, err := gouuid.NewV4()
	if err != nil {
		return nil, fmt.Errorf("UUID NewV4: %v", err)
	}
	attach := &Attachment{
		UUID: uuid.String(),
		Name: name,
	}

	localPath := attach.LocalPath()
	if err = os.MkdirAll(path.Dir(localPath), os.ModePerm); err != nil {
		return nil, fmt.Errorf("MkdirAll: %v", err)
	}

	fw, err := os.Create(localPath)
	if err != nil {
		return nil, fmt.Errorf("Create: %v", err)
	}
	defer fw.Close()

	if _, err = fw.Write(buf); err != nil {
		return nil, fmt.Errorf("Write: %v", err)
	} else if _, err = io.Copy(fw, file); err != nil {
		return nil, fmt.Errorf("Copy: %v", err)
	}

	if _, err := x.Insert(attach); err != nil {
		return nil, err
	}

	return attach, nil
}

func getAttachmentByUUID(e Engine, uuid string) (*Attachment, error) {
	attach := &Attachment{UUID: uuid}
	has, err := x.Get(attach)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, ErrAttachmentNotExist{0, uuid}
	}
	return attach, nil
}

func getAttachmentsByUUIDs(e Engine, uuids []string) ([]*Attachment, error) {
	if len(uuids) == 0 {
		return []*Attachment{}, nil
	}

	// Silently drop invalid uuids.
	attachments := make([]*Attachment, 0, len(uuids))
	return attachments, e.In("uuid", uuids).Find(&attachments)
}

// GetAttachmentByUUID returns attachment by given UUID.
func GetAttachmentByUUID(uuid string) (*Attachment, error) {
	return getAttachmentByUUID(x, uuid)
}

func getAttachmentsByPostID(e Engine, postID int64) ([]*Attachment, error) {
	attachments := make([]*Attachment, 0, 5)
	return attachments, e.Where("post_id = ?", postID).Find(&attachments)
}

// GetAttachmentsByPostID returns all attachments of a post.
func GetAttachmentsByPostID(postID int64) ([]*Attachment, error) {
	return getAttachmentsByPostID(x, postID)
}

// DeleteAttachment deletes the given attachment and optionally the associated file.
func DeleteAttachment(a *Attachment, remove bool) error {
	_, err := DeleteAttachments([]*Attachment{a}, remove)
	return err
}

// DeleteAttachments deletes the given attachments and optionally the associated files.
func DeleteAttachments(attachments []*Attachment, remove bool) (int, error) {
	for i, a := range attachments {
		if remove {
			if err := os.Remove(a.LocalPath()); err != nil {
				return i, err
			}
		}

		if _, err := x.Delete(a); err != nil {
			return i, err
		}
	}

	return len(attachments), nil
}

// DeleteAttachmentsByPost deletes all attachments associated with the given post.
func DeleteAttachmentsByPost(postId int64, remove bool) (int, error) {
	attachments, err := GetAttachmentsByPostID(postId)

	if err != nil {
		return 0, err
	}

	return DeleteAttachments(attachments, remove)
}

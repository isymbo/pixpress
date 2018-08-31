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
	log "gopkg.in/clog.v1"

	"github.com/isymbo/pixpress/setting"
)

// CoverImg represent a attachment of post.
type CoverImg struct {
	ID     int64
	UUID   string `xorm:"uuid UNIQUE"`
	PostID int64  `xorm:"INDEX"`
	Name   string

	Created     time.Time `xorm:"-" json:"-"`
	CreatedUnix int64
}

func (a *CoverImg) BeforeInsert() {
	a.CreatedUnix = time.Now().Unix()
}

func (a *CoverImg) AfterSet(colName string, _ xorm.Cell) {
	switch colName {
	case "created_unix":
		a.Created = time.Unix(a.CreatedUnix, 0).Local()
	}
}

// CoverImgLocalPath returns where attachment is stored in local file system based on given UUID.
func CoverImgLocalPath(uuid string) string {
	return path.Join(setting.Cover.Path, uuid[0:1], uuid[1:2], uuid)
}

// LocalPath returns where attachment is stored in local file system.
func (attach *CoverImg) LocalPath() string {
	return CoverImgLocalPath(attach.UUID)
}

// NewCoverImg creates a new cover image object.
func NewCoverImg(name string, buf []byte, file multipart.File) (_ *CoverImg, err error) {
	uuid, err := gouuid.NewV4()
	if err != nil {
		return nil, fmt.Errorf("UUID NewV4: %v", err)
	}
	attach := &CoverImg{
		UUID: uuid.String(),
		Name: name,
	}

	log.Trace("cover image: %+v", attach)

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

func getCoverImgByUUID(e Engine, uuid string) (*CoverImg, error) {
	attach := &CoverImg{UUID: uuid}
	has, err := x.Get(attach)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, ErrCoverImgNotExist{0, uuid}
	}
	return attach, nil
}

func getCoverImgsByUUIDs(e Engine, uuids []string) ([]*CoverImg, error) {
	if len(uuids) == 0 {
		return []*CoverImg{}, nil
	}

	// Silently drop invalid uuids.
	attachments := make([]*CoverImg, 0, len(uuids))
	return attachments, e.In("uuid", uuids).Find(&attachments)
}

// GetCoverImgByUUID returns cover image by given UUID.
func GetCoverImgByUUID(uuid string) (*CoverImg, error) {
	return getCoverImgByUUID(x, uuid)
}

func getCoverImgsByPostID(e Engine, postID int64) ([]*CoverImg, error) {
	attachments := make([]*CoverImg, 0, 5)

	// return attachments, e.Where("post_id = ?", postID).Find(&attachments)
	err := e.Where("post_id = ?", postID).Find(&attachments)
	log.Trace("attachments: %+v", attachments)

	return attachments, err
}

// GetAttachmentsByPostID returns all attachments of a post.
func GetCoverImgsByPostID(postID int64) ([]*CoverImg, error) {
	return getCoverImgsByPostID(x, postID)
}

// DeleteAttachment deletes the given attachment and optionally the associated file.
func DeleteCoverImg(a *CoverImg, remove bool) error {
	_, err := DeleteCoverImgs([]*CoverImg{a}, remove)
	return err
}

// DeleteAttachments deletes the given attachments and optionally the associated files.
func DeleteCoverImgs(attachments []*CoverImg, remove bool) (int, error) {
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
func DeleteCoverImgsByPost(postId int64, remove bool) (int, error) {
	attachments, err := GetCoverImgsByPostID(postId)

	if err != nil {
		return 0, err
	}

	return DeleteCoverImgs(attachments, remove)
}

package post

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	log "gopkg.in/clog.v1"

	"github.com/isymbo/pixpress/app/controllers/context"
	"github.com/isymbo/pixpress/app/models"
	"github.com/isymbo/pixpress/setting"
)

var (
	ErrFileTypeForbidden = errors.New("File type is not allowed")
	ErrTooManyFiles      = errors.New("Maximum number of files to upload exceeded")
)

func renderAttachmentSettings(c *context.Context) {
	c.Data["RequireDropzone"] = true

	c.Data["IsAttachmentEnabled"] = setting.Attachment.Enabled
	c.Data["AttachmentAllowedTypes"] = setting.Attachment.AllowedTypes
	c.Data["AttachmentMaxSize"] = setting.Attachment.MaxSize
	c.Data["AttachmentMaxFiles"] = setting.Attachment.MaxFiles
}

func uploadAttachment(c *context.Context, allowedTypes []string) {
	file, header, err := c.Req.FormFile("file")
	log.Trace("Attachment form file: %+v", file)
	if err != nil {
		c.Error(500, fmt.Sprintf("FormFile: %v", err))
		return
	}
	defer file.Close()

	buf := make([]byte, 1024)
	n, _ := file.Read(buf)
	if n > 0 {
		buf = buf[:n]
	}
	fileType := http.DetectContentType(buf)

	allowed := false
	for _, t := range allowedTypes {
		t := strings.Trim(t, " ")
		if t == "*/*" || t == fileType {
			allowed = true
			break
		}
	}

	if !allowed {
		c.Error(400, ErrFileTypeForbidden.Error())
		return
	}

	attach, err := models.NewAttachment(header.Filename, buf, file)
	if err != nil {
		c.Error(500, fmt.Sprintf("NewAttachment: %v", err))
		return
	}

	log.Trace("New attachment uploaded: %s", attach.UUID)
	c.JSON(200, map[string]string{
		"uuid": attach.UUID,
	})
}

func UploadPixAttachment(c *context.Context) {
	if !setting.Attachment.Enabled {
		c.NotFound()
		return
	}

	uploadAttachment(c, strings.Split(setting.Attachment.AllowedTypes, ","))
}

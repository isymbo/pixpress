package post

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Unknwon/com"

	"github.com/isymbo/pixpress/app/controllers/context"
	"github.com/isymbo/pixpress/app/models"
	"github.com/isymbo/pixpress/setting"
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
	// log.Trace("Attachment form file: %+v", file)
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

	// log.Trace("New attachment uploaded: %s", attach.UUID)
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

func RenderPixAttachment(c *context.Context) {
	attach, err := models.GetAttachmentByUUID(c.Params(":uuid"))
	if err != nil {
		c.NotFoundOrServerError("GetAttachmentByUUID", models.IsErrAttachmentNotExist, err)
		return
	} else if !com.IsFile(attach.LocalPath()) {
		c.NotFound()
		return
	}

	fr, err := os.Open(attach.LocalPath())
	if err != nil {
		c.Handle(500, "Open", err)
		return
	}
	defer fr.Close()

	c.Header().Set("Cache-Control", "public,max-age=86400")
	fmt.Println("attach.Name:", attach.Name)
	c.Header().Set("Content-Disposition", fmt.Sprintf(`inline; filename="%s"`, attach.Name))
	if err = ServeData(c, attach.Name, fr); err != nil {
		c.Handle(500, "ServeData", err)
		return
	}
}

func DownloadPixAttachment(c *context.Context) {
	attach, err := models.GetAttachmentByUUID(c.Params(":uuid"))
	if err != nil {
		c.NotFoundOrServerError("GetAttachmentByUUID", models.IsErrAttachmentNotExist, err)
		return
	} else if !com.IsFile(attach.LocalPath()) {
		c.NotFound()
		return
	}

	fr, err := os.Open(attach.LocalPath())
	if err != nil {
		c.Handle(500, "Open", err)
		return
	}
	defer fr.Close()

	c.Header().Set("Cache-Control", "public,max-age=86400")
	fmt.Println("attach.Name:", attach.Name)
	c.Header().Set("Content-Disposition", fmt.Sprintf(`inline; filename="%s"`, attach.Name))
	if err = ServeDownloadData(c, attach.Name, fr); err != nil {
		c.Handle(500, "ServeDownloadData", err)
		return
	}
}

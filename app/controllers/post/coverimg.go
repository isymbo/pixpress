package post

import (
	"fmt"
	"net/http"
	"strings"

	log "gopkg.in/clog.v1"

	"github.com/isymbo/pixpress/app/controllers/context"
	"github.com/isymbo/pixpress/app/models"
	"github.com/isymbo/pixpress/setting"
)

// var (
// 	ErrFileTypeForbidden = errors.New("File type is not allowed")
// 	ErrTooManyFiles      = errors.New("Maximum number of files to upload exceeded")
// )

func renderCoverSettings(c *context.Context) {
	c.Data["RequireDropzone"] = true

	c.Data["CoverAllowedTypes"] = setting.Cover.AllowedTypes
	c.Data["CoverMaxSize"] = setting.Cover.MaxSize
	c.Data["CoverMaxFiles"] = setting.Cover.MaxFiles
}

func uploadCoverImg(c *context.Context, allowedTypes []string) {
	file, header, err := c.Req.FormFile("file")

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

	attach, err := models.NewCoverImg(header.Filename, buf, file)
	if err != nil {
		c.Error(500, fmt.Sprintf("NewCoverImg: %v", err))
		return
	}

	log.Trace("New cover image uploaded: %s", attach.UUID)
	c.JSON(200, map[string]string{
		"uuid": attach.UUID,
	})
}

func UploadPixCoverImg(c *context.Context) {
	uploadCoverImg(c, strings.Split(setting.Cover.AllowedTypes, ","))
}

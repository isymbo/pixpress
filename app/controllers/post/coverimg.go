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

	// log.Trace("New cover image uploaded: %s", attach.UUID)
	c.JSON(200, map[string]string{
		"uuid": attach.UUID,
	})
}

func UploadPixCoverImg(c *context.Context) {
	uploadCoverImg(c, strings.Split(setting.Cover.AllowedTypes, ","))
}

func RenderPixCoverImg(c *context.Context) {
	cover, err := models.GetCoverImgByUUID(c.Params(":uuid"))
	if err != nil {
		c.NotFoundOrServerError("GetCoverImgByUUID", models.IsErrCoverImgNotExist, err)
		return
	} else if !com.IsFile(cover.LocalPath()) {
		c.NotFound()
		return
	}

	fr, err := os.Open(cover.LocalPath())
	if err != nil {
		c.Handle(500, "Open", err)
		return
	}
	defer fr.Close()

	c.Header().Set("Cache-Control", "public,max-age=86400")
	fmt.Println("cover.Name:", cover.Name)
	c.Header().Set("Content-Disposition", fmt.Sprintf(`inline; filename="%s"`, cover.Name))
	if err = ServeData(c, cover.Name, fr); err != nil {
		c.Handle(500, "ServeData", err)
		return
	}
}

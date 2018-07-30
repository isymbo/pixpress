package post

import (
	"github.com/go-macaron/binding"
	"gopkg.in/macaron.v1"

	"github.com/isymbo/pixpress/app/controllers/context"
)

const (
	PIXNEW    = "pix/new"
	PIXEDIT   = "pix/edit"
	PIXDELETE = "pix/delete"
)

type Post struct {
	LoginName string `form:"LoginName" binding:"Required"`
	Password  string `form:"Password" binding:"Required"`
	Email     string
	CName     string
}

func InitRoutes(m *macaron.Macaron) {

	reqSignIn := context.ReqSignIn
	bindIgnErr := binding.BindIgnErr

	m.Group("/pix", func() {
		m.Get("/:pixid", reqSignIn, Show)
		m.Combo("/new").
			Get(NewPix).
			Post(bindIgnErr(Post{}), NewPixPost)
	})
}

func NewPix(c *context.Context) {
	c.Data["Title"] = "新建作品"

	c.Success(PIXNEW)
}

func Show() {

}

func NewPixPost() {

}

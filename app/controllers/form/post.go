package form

import (
	"mime/multipart"

	"github.com/go-macaron/binding"
	macaron "gopkg.in/macaron.v1"
)

type CreatePost struct {
	Title string `form:"Title" binding:"Required;MaxSize(100)"`
	// Content MaxSize is " 1 << 16 - 4"
	Content string `form:"Content" binding:"Required;MaxSize(65532)"`
}

func (f *CreatePost) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return validate(errs, ctx.Data, f, ctx.Locale)
}

const (
	COVERIMG_LOCAL string = "local"
	COVERIMG_AWS   string = "aws"
)

type CoverImg struct {
	Source   string
	CoverImg *multipart.FileHeader
	URL      string `binding:"OmitEmpty;MaxSize(254)"`
}

func (f *CoverImg) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return validate(errs, ctx.Data, f, ctx.Locale)
}

package form

import (
	"github.com/go-macaron/binding"
	macaron "gopkg.in/macaron.v1"
)

type CreatePost struct {
	Title   string `binding:"Required;MaxSize(100)"`
	Content string `binding:"Required;MaxSize(255)"`
	// Content string `binding:"Required;MaxSize( 1 << 16 - 4)"`
}

func (f *CreatePost) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return validate(errs, ctx.Data, f, ctx.Locale)
}

package user

import (
	"net/http"

	"github.com/isymbo/pixpress/app/controllers/context"
)

const (
	LOGIN  = "user/login"
	LOGOUT = "user/logout"
)

func Login(c *context.Context) {
	// c.Title("sign_in")

	//c.Success(LOGIN)
	c.HTML(http.StatusOK, LOGIN)
}

func Logout(c *context.Context) {
	// c.Title("sign_out")

	//c.Success(LOGOUT)
	c.HTML(http.StatusOK, LOGOUT)
}

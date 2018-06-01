package user

import (
	"github.com/isymbo/pixpress/app/controllers/context"
	"gopkg.in/macaron.v1"
)

const (
	LOGIN  = "user/login"
	LOGOUT = "user/logout"
)

func InitRoutes(m *macaron.Macaron) {
	m.Group("/user", func() {
		m.Get("/login", Login)
		m.Get("/logout", Logout)
	})
}

func Login(c *context.Context) {
	c.Title("sign_in")

	c.Success(LOGIN)
	// c.HTML(http.StatusOK, LOGIN)
}

func Logout(c *context.Context) {
	c.Title("sign_out")

	c.Success(LOGOUT)
}

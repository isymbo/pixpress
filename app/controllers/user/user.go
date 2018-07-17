package user

import (
	"fmt"
	"net/url"

	"gopkg.in/macaron.v1"

	"github.com/go-macaron/binding"
	// ldap "github.com/jtblin/go-ldap-client"
	log "gopkg.in/clog.v1"

	"github.com/isymbo/pixpress/app/controllers/context"
	"github.com/isymbo/pixpress/app/models"
	"github.com/isymbo/pixpress/app/models/errors"
	"github.com/isymbo/pixpress/setting"
)

const (
	LOGIN  = "user/login"
	LOGOUT = "user/logout"
	HOME   = "user/home"
)

type User struct {
	LoginName string `form:"LoginName" binding:"Required"`
	Password  string `form:"Password" binding:"Required"`
	Email     string
	CName     string
}

func InitRoutes(m *macaron.Macaron) {

	reqSignIn := context.ReqSignIn
	//reqSignOut := context.ReqSignOut

	bindIgnErr := binding.BindIgnErr

	m.Group("/user", func() {
		m.Get("/:id", reqSignIn, Home)
		m.Combo("/login").
			Get(Login).
			Post(bindIgnErr(User{}), LoginPost)
		m.Get("/logout", reqSignIn, Logout)
	})

	// m.Group("/user", func() {
	// 	m.Group("/login", func() {
	// 		m.Combo("").Get(user.Login).
	// 			Post(bindIgnErr(form.SignIn{}), user.LoginPost)
	// 		m.Combo("/two_factor").Get(user.LoginTwoFactor).Post(user.LoginTwoFactorPost)
	// 		m.Combo("/two_factor_recovery_code").Get(user.LoginTwoFactorRecoveryCode).Post(user.LoginTwoFactorRecoveryCodePost)
	// 	})

	// 	m.Get("/sign_up", user.SignUp)
	// 	m.Post("/sign_up", bindIgnErr(form.Register{}), user.SignUpPost)
	// 	m.Get("/reset_password", user.ResetPasswd)
	// 	m.Post("/reset_password", user.ResetPasswdPost)
	// }, reqSignOut)
}

// AutoLogin reads cookie and try to auto-login.
func AutoLogin(c *context.Context) (bool, error) {
	if !models.HasEngine {
		return false, nil
	}

	uname := c.GetCookie(setting.Security.CookieUserName)
	if len(uname) == 0 {
		return false, nil
	}

	isSucceed := false
	defer func() {
		if !isSucceed {
			log.Trace("auto-login cookie cleared: %s", uname)
			c.SetCookie(setting.Security.CookieUserName, "", -1, setting.AppSubURL)
			c.SetCookie(setting.Security.CookieRememberName, "", -1, setting.AppSubURL)
			c.SetCookie(setting.Security.LoginStatusCookieName, "", -1, setting.AppSubURL)
		}
	}()

	u, err := models.GetUserByName(uname)
	if err != nil {
		if !errors.IsUserNotExist(err) {
			return false, fmt.Errorf("GetUserByName: %v", err)
		}
		return false, nil
	}

	// TODO, Uncommment for now since not use it.
	// if val, ok := c.GetSuperSecureCookie(u.Rands+u.Passwd, setting.Security.CookieRememberName); !ok || val != u.LoginName {
	// 	return false, nil
	// }

	isSucceed = true
	c.Session.Set("uid", u.ID)
	c.Session.Set("uname", u.LoginName)
	c.SetCookie(setting.Session.CSRFCookieName, "", -1, setting.AppSubURL)
	if setting.Security.EnableLoginStatusCookie {
		c.SetCookie(setting.Security.LoginStatusCookieName, "true", 0, setting.AppSubURL)
	}
	return true, nil
}

func Login(c *context.Context) {
	c.Title("sign_in")

	// Check auto-login
	isSucceed, err := AutoLogin(c)
	log.Trace("AutoLogin Succeed: %+v", isSucceed)
	if err != nil {
		c.ServerError("AutoLogin", err)
		return
	}

	redirectTo := c.Query("redirect_to")
	if len(redirectTo) > 0 {
		c.SetCookie("redirect_to", redirectTo, 0, setting.AppSubURL)
	} else {
		redirectTo, _ = url.QueryUnescape(c.GetCookie("redirect_to"))
	}

	if isSucceed {
		if isValidRedirect(redirectTo) {
			c.Redirect(redirectTo)
		} else {
			c.SubURLRedirect("/")
		}
		c.SetCookie("redirect_to", "", -1, setting.AppSubURL)
		return
	}

	// Display normal login page
	loginSources, err := models.ActivatedLoginSources()
	if err != nil {
		c.ServerError("ActivatedLoginSources", err)
		return
	}
	c.Data["LoginSources"] = loginSources

	c.Success(LOGIN)
	// c.HTML(http.StatusOK, LOGIN)
}

func Logout(c *context.Context) {
	c.Title("sign_out")

	c.Success(LOGOUT)
}

func Home(c *context.Context) {
	c.Title("home")

	id := c.ParamsInt64("id")

	models.GetUserProfile(id)

	c.Success(HOME)
}

func afterLogin(c *context.Context, u *models.User, remember bool) {
	if remember {
		days := 86400 * setting.Security.LoginRememberDays
		c.SetCookie(setting.Security.CookieUserName, u.LoginName, days, setting.AppSubURL, "", setting.Security.CookieSecure, true)
		//c.SetSuperSecureCookie(u.Rands+u.Passwd, setting.CookieRememberName, u.Name, days, setting.AppSubURL, "", setting.CookieSecure, true)
	}

	c.Session.Set("uid", u.ID)
	c.Session.Set("uname", u.LoginName)

	// Clear whatever CSRF has right now, force to generate a new one
	c.SetCookie(setting.Session.CSRFCookieName, "", -1, setting.AppSubURL)
	if setting.Security.EnableLoginStatusCookie {
		c.SetCookie(setting.Security.LoginStatusCookieName, "true", 0, setting.AppSubURL)
	}

	redirectTo, _ := url.QueryUnescape(c.GetCookie("redirect_to"))
	c.SetCookie("redirect_to", "", -1, setting.AppSubURL)
	if isValidRedirect(redirectTo) {
		c.Redirect(redirectTo)
		return
	}

	c.SubURLRedirect("/")
}

// isValidRedirect returns false if the URL does not redirect to same site.
// False: //url, http://url
// True: /url
func isValidRedirect(url string) bool {
	return len(url) >= 2 && url[0] == '/' && url[1] != '/'
}

func LoginPost(c *context.Context, u User) {
	c.Title("sign_in")

	// loginSources, err := models.ActivatedLoginSources()
	// if err != nil {
	// 	c.ServerError("ActivatedLoginSources", err)
	// 	return
	// }
	// c.Data["LoginSources"] = loginSources

	// if c.HasError() {
	// 	c.Success(LOGIN)
	// 	return
	// }

	// u, err := models.UserLogin(u.LoginName, u.Password, f.LoginSource)
	// if err != nil {
	// 	switch err.(type) {
	// 	// case errors.UserNotExist:
	// 	// 	c.FormErr("UserName", "Password")
	// 	// 	c.RenderWithErr(c.Tr("form.username_password_incorrect"), LOGIN, &f)
	// 	// case errors.LoginSourceMismatch:
	// 	// 	c.FormErr("LoginSource")
	// 	// 	c.RenderWithErr(c.Tr("form.auth_source_mismatch"), LOGIN, &f)

	// 	default:
	// 		c.ServerError("UserLogin", err)
	// 	}
	// 	return
	// }

	user, _ := models.UserLogin(u.LoginName, u.Password, 101)

	afterLogin(c, user, true)

	c.Success(HOME)
}

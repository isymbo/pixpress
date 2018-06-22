package user

import (
	"log"

	"gopkg.in/macaron.v1"

	"github.com/go-macaron/binding"
	ldap "github.com/jtblin/go-ldap-client"

	//"github.com/isymbo/pixpress/app/controllers/auth/ldap"
	"github.com/isymbo/pixpress/app/controllers/context"
	"github.com/isymbo/pixpress/app/models"
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

	bindIgnErr := binding.BindIgnErr

	m.Group("/user", func() {
		m.Get("/:id", Home)
		m.Combo("/login").
			Get(Login).
			Post(bindIgnErr(User{}), LoginPost)
		m.Get("/logout", Logout)
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

func Login(c *context.Context) {
	c.Title("sign_in")

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

func LoginPost(c *context.Context, u User) {
	c.Title("sign_in")

	//return fmt.Sprintf("LoginName: %s\nPassword: %v", u.LoginName, u.Password)
	client := &ldap.LDAPClient{
		Base:         setting.Ldap.Base,
		Host:         setting.Ldap.Host,
		Port:         setting.Ldap.Port,
		UseSSL:       false,
		SkipTLS:      true,
		BindDN:       setting.Ldap.BindDn,
		BindPassword: setting.Ldap.Password,
		UserFilter:   "(sAMAccountName=%s)",
		Attributes:   []string{"displayName", "mail", "mobile", "sAMAccountName"},
	}
	// It is the responsibility of the caller to close the connection
	defer client.Close()

	ok, rmap, err := client.Authenticate(u.LoginName, u.Password)
	if err != nil {
		log.Printf("Error authenticating user %s: %+v", u.LoginName, err)
	}
	if ok {
		newUser := &models.User{
			LoginName:   rmap["sAMAccountName"],
			DisplayName: rmap["displayName"],
			Email:       rmap["mail"],
			Mobile:      rmap["mobile"],
		}
		if err = models.CreateUser(newUser); err != nil {
			if models.IsErrLoginNameAlreadyExist(err) {
				c.Success(HOME)
				return
			}
		}
	}

	c.Success(LOGIN)
}

package context

import (
	"net/url"

	log "gopkg.in/clog.v1"
	"gopkg.in/macaron.v1"

	"github.com/isymbo/pixpress/app/controllers/auth"
	"github.com/isymbo/pixpress/setting"
)

type ToggleOptions struct {
	SignInRequired  bool
	SignOutRequired bool
	AdminRequired   bool
	DisableCSRF     bool
}

var (
	ReqSignIn        = Toggle(&ToggleOptions{SignInRequired: true})
	IgnSignIn        = Toggle(&ToggleOptions{SignInRequired: setting.Service.RequireSignInView})
	IgnSignInAndCsrf = Toggle(&ToggleOptions{DisableCSRF: true})
	ReqSignOut       = Toggle(&ToggleOptions{SignOutRequired: true})
	AdminReq         = Toggle(&ToggleOptions{SignInRequired: true, AdminRequired: true})
)

func Toggle(options *ToggleOptions) macaron.Handler {
	return func(c *Context) {
		// Cannot view any page before installation.
		// if !setting.InstallLock {
		// 	c.Redirect(setting.AppSubURL + "/install")
		// 	return
		// }

		// Check prohibit login users.
		// if c.IsLogged && c.User.ProhibitLogin {
		// 	c.Data["Title"] = c.Tr("auth.prohibit_login")
		// 	c.HTML(200, "user/auth/prohibit_login")
		// 	return
		// }

		// Check non-logged users landing page.
		if !c.IsLogged && c.Req.RequestURI == "/" && setting.Server.LandingPageURL != setting.LANDING_PAGE_HOME {
			c.Redirect(setting.AppSubURL + string(setting.Server.LandingPageURL))
			log.Trace(setting.AppSubURL + string(setting.Server.LandingPageURL))
			return
		}

		// Redirect to dashboard if user tries to visit any non-login page.
		if options.SignOutRequired && c.IsLogged && c.Req.RequestURI != "/" {
			c.Redirect(setting.AppSubURL + "/")
			return
		}

		// if !options.SignOutRequired && !options.DisableCSRF && c.Req.Method == "POST" && !auth.IsAPIPath(c.Req.URL.Path) {
		// 	csrf.Validate(c.Context, c.csrf)
		// 	if c.Written() {
		// 		return
		// 	}
		// }

		if options.SignInRequired {
			if !c.IsLogged {
				// Restrict API calls with error message.
				if auth.IsAPIPath(c.Req.URL.Path) {
					c.JSON(403, map[string]string{
						"message": "Only signed in user is allowed to call APIs.",
					})
					return
				}

				c.SetCookie("redirect_to", url.QueryEscape(setting.AppSubURL+c.Req.RequestURI), 0, setting.AppSubURL)
				c.Redirect(setting.AppSubURL + "/user/login")
				return
			}
			// else if !c.User.IsActive && setting.Service.RegisterEmailConfirm {
			// 	c.Data["Title"] = c.Tr("auth.active_your_account")
			// 	c.HTML(200, "user/auth/activate")
			// 	return
			// }
		}

		// Redirect to log in page if auto-signin info is provided and has not signed in.
		if !options.SignOutRequired && !c.IsLogged && !auth.IsAPIPath(c.Req.URL.Path) &&
			len(c.GetCookie(setting.Security.CookieUserName)) > 0 {
			c.SetCookie("redirect_to", url.QueryEscape(setting.AppSubURL+c.Req.RequestURI), 0, setting.AppSubURL)
			c.Redirect(setting.AppSubURL + "/user/login")
			return
		}

		if options.AdminRequired {
			if !c.User.IsAdmin {
				c.Error(403)
				return
			}
			c.Data["PageIsAdmin"] = true
		}
	}
}

package auth

import (
	"strings"
	"time"

	"github.com/go-macaron/session"
	gouuid "github.com/satori/go.uuid"
	log "gopkg.in/clog.v1"
	"gopkg.in/macaron.v1"

	"github.com/isymbo/pixpress/app/models"
	"github.com/isymbo/pixpress/app/models/errors"
	"github.com/isymbo/pixpress/setting"
	"github.com/isymbo/pixpress/util"
)

func IsAPIPath(url string) bool {
	return strings.HasPrefix(url, "/api/")
}

// SignedInID returns the id of signed in user.
func SignedInID(c *macaron.Context, sess session.Store) int64 {
	if !models.HasEngine {
		return 0
	}

	// Check access token.
	if IsAPIPath(c.Req.URL.Path) {
		tokenSHA := c.Query("token")
		if len(tokenSHA) <= 0 {
			tokenSHA = c.Query("access_token")
		}
		if len(tokenSHA) == 0 {
			// Well, check with header again.
			auHead := c.Req.Header.Get("Authorization")
			if len(auHead) > 0 {
				auths := strings.Fields(auHead)
				if len(auths) == 2 && auths[0] == "token" {
					tokenSHA = auths[1]
				}
			}
		}

		// Let's see if token is valid.
		if len(tokenSHA) > 0 {
			t, err := models.GetAccessTokenBySHA(tokenSHA)
			if err != nil {
				if !models.IsErrAccessTokenNotExist(err) && !models.IsErrAccessTokenEmpty(err) {
					log.Error(2, "GetAccessTokenBySHA: %v", err)
				}
				return 0
			}
			t.Updated = time.Now()
			if err = models.UpdateAccessToken(t); err != nil {
				log.Error(2, "UpdateAccessToken: %v", err)
			}
			return t.UID
		}
	}

	uid := sess.Get("uid")

	if uid == nil {
		return 0
	}
	if id, ok := uid.(int64); ok {
		if _, err := models.GetUserByID(id); err != nil {
			if !errors.IsUserNotExist(err) {
				log.Error(2, "GetUserByID: %v", err)
			}
			return 0
		}
		return id
	}
	return 0
}

// SignedInUser returns the user object of signed user.
// It returns a bool value to indicate whether user uses basic auth or not.
// FIXME, TODO
func SignedInUser(ctx *macaron.Context, sess session.Store) (*models.User, bool) {
	if !models.HasEngine {
		return nil, false
	}

	uid := SignedInID(ctx, sess)

	if uid <= 0 {
		if setting.Service.EnableReverseProxyAuth {
			webAuthUser := ctx.Req.Header.Get(setting.Security.ReverseProxyAuthUser)
			if len(webAuthUser) > 0 {
				u, err := models.GetUserByName(webAuthUser)
				if err != nil {
					if !errors.IsUserNotExist(err) {
						log.Error(4, "GetUserByName: %v", err)
						return nil, false
					}

					// Check if enabled auto-registration.
					if setting.Service.EnableReverseProxyAutoRegister {

						// TODO, FIXME
						t, err := gouuid.NewV4()
						if err != nil {
							return nil, false
						}

						u := &models.User{
							LoginName: webAuthUser,
							Email:     t.String() + "@localhost",
							Passwd:    webAuthUser,
							IsActive:  true,
						}
						if err = models.CreateUser(u); err != nil {
							// FIXME: should I create a system notice?
							log.Error(4, "CreateUser: %v", err)
							return nil, false
						} else {
							return u, false
						}
					}
				}
				return u, false
			}
		}

		// Check with basic auth.
		baHead := ctx.Req.Header.Get("Authorization")
		if len(baHead) > 0 {
			auths := strings.Fields(baHead)
			if len(auths) == 2 && auths[0] == "Basic" {
				uname, passwd, _ := util.BasicAuthDecode(auths[1])

				u, err := models.UserLogin(uname, passwd, -1)
				if err != nil {
					if !errors.IsUserNotExist(err) {
						log.Error(4, "UserLogin: %v", err)
					}
					return nil, false
				}

				return u, true
			}
		}
		return nil, false
	}

	u, err := models.GetUserByID(uid)
	if err != nil {
		log.Error(4, "GetUserById: %v", err)
		return nil, false
	}
	return u, false
}

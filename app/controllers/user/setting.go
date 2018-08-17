package user

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/Unknwon/com"
	log "gopkg.in/clog.v1"

	"github.com/isymbo/pixpress/app/controllers/context"
	"github.com/isymbo/pixpress/app/controllers/form"
	"github.com/isymbo/pixpress/app/models"
	"github.com/isymbo/pixpress/util"
)

func Settings(c *context.Context) {
	// uname := c.Params(":loginname")
	// log.Trace("loginname: %+v", uname)
	// // Special handle for FireFox requests favicon.ico.
	// if uname == "favicon.ico" {
	// 	c.ServeFile(path.Join(setting.Server.StaticRootPath, "public/img/favicon.png"))
	// 	return
	// } else if strings.HasSuffix(uname, ".png") {
	// 	c.Error(404)
	// 	return
	// }

	// ctxUser := GetUserByName(c, strings.TrimSuffix(uname, ".keys"))
	// log.Trace("ctxUser: %+v", ctxUser)
	// if c.Written() {
	// 	return
	// }

	// c.Data["Title"] = ctxUser.DisplayName
	// c.Data["PageIsUserProfile"] = true
	// c.Data["Owner"] = ctxUser

	// c.Success(PROFILE)

	c.Title("用户头像")
	c.PageIs("SettingsAvatar")

	c.Success(SETTINGS_AVATAR)

}

// FIXME: limit upload size
func UpdateAvatarSetting(c *context.Context, f form.Avatar, ctxUser *models.User) error {
	ctxUser.UseCustomAvatar = f.Source == form.AVATAR_LOCAL
	log.Trace("UseCustomAvatar: %+v", ctxUser.UseCustomAvatar)

	if len(f.Gravatar) > 0 {
		ctxUser.Avatar = util.MD5(f.Gravatar)
		//ctxUser.AvatarEmail = f.Gravatar
	}

	if f.Avatar != nil && f.Avatar.Filename != "" {
		r, err := f.Avatar.Open()
		if err != nil {
			return fmt.Errorf("open avatar reader: %v", err)
		}
		defer r.Close()

		data, err := ioutil.ReadAll(r)
		if err != nil {
			return fmt.Errorf("read avatar content: %v", err)
		}
		if !util.IsImageFile(data) {
			return errors.New("上传的文件不是一张图片！")
		}
		if err = ctxUser.UploadAvatar(data); err != nil {
			return fmt.Errorf("upload avatar: %v", err)
		}
	} else {
		// // No avatar is uploaded but setting has been changed to enable,
		// // generate a random one when needed.
		// if ctxUser.UseCustomAvatar && !com.IsFile(ctxUser.CustomAvatarPath()) {
		// 	if err := ctxUser.GenerateRandomAvatar(); err != nil {
		// 		log.Error(2, "generate random avatar [%d]: %v", ctxUser.ID, err)
		// 	}
		// }

		// No avatar is uploaded but setting has been changed to enable,
		// use default avatar.
		if ctxUser.UseCustomAvatar && !com.IsFile(ctxUser.CustomAvatarPath()) {
			// if err := ctxUser.GenerateRandomAvatar(); err != nil {
			// 	log.Error(2, "generate random avatar [%d]: %v", ctxUser.ID, err)
			// }
		}
	}

	if err := models.UpdateUser(ctxUser); err != nil {
		return fmt.Errorf("update user: %v", err)
	}

	return nil
}

func SettingsAvatar(c *context.Context) {
	c.Title("用户头像")
	c.PageIs("SettingsAvatar")

	c.Success(SETTINGS_AVATAR)
}

func SettingsAvatarPost(c *context.Context, f form.Avatar) {
	if err := UpdateAvatarSetting(c, f, c.User); err != nil {
		c.Flash.Error(err.Error())
	} else {
		c.Flash.Success("您的头像设置更新成功！")
	}

	c.SubURLRedirect("/user/settings/avatar")
}

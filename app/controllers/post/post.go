package post

import (
	"errors"
	"io"
	"time"

	"github.com/go-macaron/binding"
	log "gopkg.in/clog.v1"
	"gopkg.in/macaron.v1"

	"github.com/isymbo/pixpress/app/controllers/context"
	"github.com/isymbo/pixpress/app/controllers/form"
	"github.com/isymbo/pixpress/app/controllers/routes"
	"github.com/isymbo/pixpress/app/models"
	"github.com/isymbo/pixpress/setting"
	"github.com/isymbo/pixpress/util"
)

const (
	PIXNEW    = "pix/new"
	PIXEDIT   = "pix/edit"
	PIXDELETE = "pix/delete"
	PIXES     = "pix/list"
	PIXVIEW   = "pix/view"
	PIXHOME   = "/pix"

	EXPLORE_PIX = "explore/pix"
)

var (
	ErrFileTypeForbidden = errors.New("File type is not allowed")
	ErrTooManyFiles      = errors.New("Maximum number of files to upload exceeded")
)

func InitRoutes(m *macaron.Macaron) {

	reqSignIn := context.ReqSignIn
	bindIgnErr := binding.BindIgnErr

	m.Group("/pix", func() {
		m.Get("", reqSignIn, ListPix)
		m.Get("/:pixid", reqSignIn, ViewPix)
		m.Combo("/:pixid/edit").
			Get(reqSignIn, EditPix).
			Post(bindIgnErr(form.CreatePost{}), EditPixPost)
		m.Get("/:pixid/delete", reqSignIn, DeletePix)
		m.Combo("/new").
			Get(reqSignIn, NewPix).
			Post(bindIgnErr(form.CreatePost{}), NewPixPost)
		m.Post("/attachments", UploadPixAttachment)
		m.Post("/coverimgs", UploadPixCoverImg)
	})

	m.Get("/attachments/:uuid", RenderPixAttachment)
	m.Get("/attachments/:uuid/action/:action", NumDownloadsUpdate, DownloadPixAttachment)
	m.Get("/covers/:uuid", RenderPixCoverImg)
}

func ServeData(c *context.Context, name string, reader io.Reader) error {
	buf := make([]byte, 1024)
	n, _ := reader.Read(buf)
	if n >= 0 {
		buf = buf[:n]
	}

	if !util.IsImageFile(buf) {
		c.Resp.Header().Set("Content-Disposition", "attachment; filename=\""+name+"\"")
		c.Resp.Header().Set("Content-Transfer-Encoding", "binary")
	}
	c.Resp.Write(buf)
	_, err := io.Copy(c.Resp, reader)
	return err
}

func ServeDownloadData(c *context.Context, name string, reader io.Reader) error {
	buf := make([]byte, 1024)
	n, _ := reader.Read(buf)
	if n >= 0 {
		buf = buf[:n]
	}

	c.Resp.Header().Set("Content-Disposition", "attachment; filename=\""+name+"\"")
	c.Resp.Header().Set("Content-Transfer-Encoding", "binary")
	c.Resp.Write(buf)
	_, err := io.Copy(c.Resp, reader)
	return err
}

func NewPix(c *context.Context) {
	c.Data["Title"] = "新建作品"
	c.Data["PageIsPixes"] = true
	c.Data["PageIsNewPix"] = true
	renderAttachmentSettings(c)
	renderCoverSettings(c)

	if c.Written() {
		return
	}

	c.Success(PIXNEW)
}

func EditPix(c *context.Context) {
	c.Data["Title"] = "作品信息"
	c.Data["PageIsPixes"] = true
	c.Data["PageIsEditPix"] = true
	renderAttachmentSettings(c)
	renderCoverSettings(c)

	post, err := models.GetPostByID(c.ParamsInt64(":pixid"))
	if err != nil {
		if models.IsErrPostNotExist(err) {
			c.Handle(404, "", nil)
		} else {
			c.Handle(500, "GetPostByID", err)
		}
		return
	}

	post.Attachments, err = models.GetAttachmentsByPostID(c.ParamsInt64(":pixid"))
	if err != nil {
		if models.IsErrAttachmentNotExist(err) {
			c.Handle(404, "", nil)
		} else {
			c.Handle(500, "GetAttachmentsByPostID", err)
		}
		return
	}

	post.CoverImg, err = models.GetCoverImgsByPostID(c.ParamsInt64(":pixid"))
	if err != nil {
		if models.IsErrCoverImgNotExist(err) {
			c.Handle(404, "", nil)
		} else {
			c.Handle(500, "GetCoverImgsByPostID", err)
		}
		return
	}

	c.Data["title"] = post.Title
	c.Data["content"] = post.Content
	c.Data["Post"] = post

	// log.Trace("Post: %+v", post)
	// log.Trace("Post.Attachments: %+v", post.Attachments)
	// log.Trace("Post.CoverImg: %+v", post.CoverImg)

	c.Success(PIXEDIT)
}

func EditPixPost(c *context.Context, f form.CreatePost) {
	c.Data["Title"] = "作品信息"
	c.Data["PageIsPixes"] = true
	c.Data["PageIsEditPixPost"] = true

	if c.HasError() {
		c.Flash.Error(c.Data["ErrorMsg"].(string))
		c.HTML(200, PIXEDIT)
		return
	}

	post, err := models.GetPostByID(c.ParamsInt64(":pixid"))
	if err != nil {
		if models.IsErrPostNotExist(err) {
			c.Handle(404, "", nil)
		} else {
			c.Handle(500, "GetPostByID", err)
		}
		return
	}

	c.Data["Post"] = post

	post.Title = f.Title
	post.Content = f.Content

	if err = models.UpdatePost(post); err != nil {
		c.Handle(500, "UpdatePost", err)
		return
	}

	c.SubURLRedirect(PIXHOME)
}

func ListPix(c *context.Context) {
	c.Data["Title"] = "我的作品列表"
	c.Data["PageIsPixes"] = true
	c.Data["PageIsListPix"] = true

	if c.HasError() {
		c.Flash.Error(c.Data["ErrorMsg"].(string))
	}

	routes.RenderPostSearchByAuthorID(c, &routes.PostSearchByAuthorIDOptions{
		Type:     models.POST_TYPE_PIX,
		Counter:  models.CountPostsByAuthorID,
		Ranger:   models.PostsByAuthorID,
		PageSize: setting.UI.User.PostPagingNum,
		OrderBy:  "updated_unix DESC",
		TplName:  PIXES,
		AuthorID: c.User.ID,
	})
}

// func NewPixPost(c *context.Context, f form.CreatePost) {
// 	c.Data["Title"] = "新建作品"
// 	c.Data["PageIsPixes"] = true
// 	c.Data["PageIsNewPix"] = true
// 	renderAttachmentSettings(c)

// 	if c.HasError() {
// 		c.Success(PIXNEW)
// 		return
// 	}

// 	post := &models.Post{
// 		AuthorID: c.User.ID,
// 		Author:   c.User,
// 		Title:    f.Title,
// 		Content:  f.Content,
// 		PostType: models.POST_TYPE_PIX,
// 	}

// 	if err := models.NewPost(post); err != nil {
// 		c.Handle(500, "NewPost", err)
// 		return
// 	}

// 	c.SubURLRedirect(PIXHOME)
// }

func DeletePix(c *context.Context) {
	post, err := models.GetPostByID(c.ParamsInt64(":pixid"))
	if err != nil {
		if models.IsErrPostNotExist(err) {
			c.Handle(404, "", nil)
		} else {
			c.Handle(500, "GetPostByID", err)
		}
		return
	}

	if err = models.DeletePost(post); err != nil {
		c.Handle(500, "DeletePost", err)
		return
	}

	c.SubURLRedirect(PIXHOME)
}

// func PixCoverImgPost(c *context.Context, f form.CoverImg) {
// 	if err := UpdatePixCoverImg(c, f, c.Post.Post); err != nil {
// 		c.Flash.Error(err.Error())
// 	} else {
// 		c.Flash.Success("封面图像上传成功！")
// 	}
// 	// c.SubURLRedirect(c.Repo.RepoLink + "/settings")

// 	// FIXME, TODO, fix this when no need to request url path
// 	log.Trace("c.Link: %+v", c.Req.URL.Path)
// 	c.Req.URL.Path = strings.TrimRight(c.Req.URL.Path, "/pix/new/cover")
// 	log.Trace("afterTrim c.Link: %+v", c.Req.URL.Path)

// 	c.SubURLRedirect(PIXNEW)
// }

// // FIXME: limit upload size
// func UpdatePixCoverImg(c *context.Context, f form.CoverImg, ctxPost *models.Post) error {
// 	if f.CoverImg != nil {
// 		r, err := f.CoverImg.Open()
// 		if err != nil {
// 			return fmt.Errorf("open coverimg reader: %v", err)
// 		}
// 		defer r.Close()

// 		data, err := ioutil.ReadAll(r)
// 		if err != nil {
// 			return fmt.Errorf("read coverimg content: %v", err)
// 		}
// 		if !util.IsImageFile(data) {
// 			return errors.New("上传的文件不是一张图片！")
// 		}
// 		if err = ctxPost.UploadCoverImg(data); err != nil {
// 			return fmt.Errorf("upload avatar: %v", err)
// 		}
// 	} else {
// 		// No avatar is uploaded and reset setting back.
// 		if !com.IsFile(ctxPost.CoverImgPath()) {
// 			// ctxRepo.UseCustomAvatar = false
// 		}
// 	}

// 	if err := models.UpdatePost(ctxPost); err != nil {
// 		return fmt.Errorf("update post: %v", err)
// 	}

// 	return nil
// }

func NewPixPost(c *context.Context, f form.CreatePost) {
	c.Data["Title"] = "新建作品"
	c.Data["PageIsPixes"] = true
	c.Data["PageIsNewPix"] = true
	renderAttachmentSettings(c)
	renderCoverSettings(c)

	if c.Written() {
		return
	}

	var attachments []string
	if setting.Attachment.Enabled {
		attachments = f.Files
	}

	var cover []string
	cover = f.CoverImg

	if c.HasError() {
		c.Flash.Error(c.Data["ErrorMsg"].(string))
		// c.Redirect(fmt.Sprintf("%s/issues/%d", c.Repo.RepoLink, issue.Index))
		return
	}

	var err error
	// var post *models.CreatePostOptions
	// defer func() {
	// 	// Check if issue admin/poster changes the status of issue.
	// 	if (c.Repo.IsWriter() || (c.IsLogged && issue.IsPoster(c.User.ID))) &&
	// 		(f.Status == "reopen" || f.Status == "close") &&
	// 		!(issue.IsPull && issue.PullRequest.HasMerged) {

	// 		// Duplication and conflict check should apply to reopen pull request.
	// 		var pr *models.PullRequest

	// 		if f.Status == "reopen" && issue.IsPull {
	// 			pull := issue.PullRequest
	// 			pr, err = models.GetUnmergedPullRequest(pull.HeadRepoID, pull.BaseRepoID, pull.HeadBranch, pull.BaseBranch)
	// 			if err != nil {
	// 				if !models.IsErrPullRequestNotExist(err) {
	// 					c.ServerError("GetUnmergedPullRequest", err)
	// 					return
	// 				}
	// 			}

	// 			// Regenerate patch and test conflict.
	// 			if pr == nil {
	// 				if err = issue.PullRequest.UpdatePatch(); err != nil {
	// 					c.ServerError("UpdatePatch", err)
	// 					return
	// 				}

	// 				issue.PullRequest.AddToTaskQueue()
	// 			}
	// 		}

	// 		if pr != nil {
	// 			c.Flash.Info(c.Tr("repo.pulls.open_unmerged_pull_exists", pr.Index))
	// 		} else {
	// 			if err = issue.ChangeStatus(c.User, c.Repo.Repository, f.Status == "close"); err != nil {
	// 				log.Error(2, "ChangeStatus: %v", err)
	// 			} else {
	// 				log.Trace("Issue [%d] status changed to closed: %v", issue.ID, issue.IsClosed)
	// 			}
	// 		}
	// 	}

	// 	// Redirect to comment hashtag if there is any actual content.
	// 	typeName := "issues"
	// 	if issue.IsPull {
	// 		typeName = "pulls"
	// 	}
	// 	if comment != nil {
	// 		c.Redirect(fmt.Sprintf("%s/%s/%d#%s", c.Repo.RepoLink, typeName, issue.Index, comment.HashTag()))
	// 	} else {
	// 		c.Redirect(fmt.Sprintf("%s/%s/%d", c.Repo.RepoLink, typeName, issue.Index))
	// 	}
	// }()

	// Fix #321: Allow empty comments, as long as we have title, content and coverimg
	if len(f.Title) == 0 && len(f.Content) == 0 && len(cover) == 0 {
		return
	}

	po := &models.CreatePostOptions{
		PostType:    models.POST_TYPE_PIX,
		Doer:        c.User,
		Title:       f.Title,
		Content:     f.Content,
		Attachments: attachments,
		CoverImg:    cover,
	}

	post, err := models.CreatePostByOption(po)
	if err != nil {
		c.ServerError("CreatePost", err)
		return
	}

	log.Info("Post created: %d", post.ID)

	c.SubURLRedirect(PIXHOME)
}

func AnonViewPix(c *context.Context) {
	c.Data["Title"] = "作品信息"
	c.Data["PageIsExplore"] = true
	c.Data["PageIsAnonViewPix"] = true
	renderAttachmentSettings(c)
	renderCoverSettings(c)

	post, err := models.GetPostByID(c.ParamsInt64(":pixid"))
	if err != nil {
		if models.IsErrPostNotExist(err) {
			c.Handle(404, "", nil)
		} else {
			c.Handle(500, "GetPostByID", err)
		}
		return
	}
	post.Updated = time.Unix(post.UpdatedUnix, 0).Local()

	post.Author, err = models.GetUserByPostID(c.ParamsInt64(":pixid"))
	if err != nil {
		if models.IsErrUserNotExist(err) {
			c.Handle(404, "", nil)
		} else {
			c.Handle(500, "GetUserByPostID", err)
		}
		return
	}

	post.Attachments, err = models.GetAttachmentsByPostID(c.ParamsInt64(":pixid"))
	if err != nil {
		if models.IsErrAttachmentNotExist(err) {
			c.Handle(404, "", nil)
		} else {
			c.Handle(500, "GetAttachmentsByPostID", err)
		}
		return
	}

	post.CoverImg, err = models.GetCoverImgsByPostID(c.ParamsInt64(":pixid"))
	if err != nil {
		if models.IsErrCoverImgNotExist(err) {
			c.Handle(404, "", nil)
		} else {
			c.Handle(500, "GetCoverImgsByPostID", err)
		}
		return
	}

	c.Data["title"] = post.Title
	c.Data["content"] = post.Content
	c.Data["Post"] = post

	c.Success(EXPLORE_PIX)
}

func NumViewsUpdate(c *context.Context) {
	post, err := models.GetPostByID(c.ParamsInt64(":pixid"))
	if err != nil {
		if models.IsErrPostNotExist(err) {
			c.Handle(404, "", nil)
		} else {
			c.Handle(500, "GetPostByID", err)
		}
		return
	}

	models.PostIncNumViews(post)
}

func NumDownloadsUpdate(c *context.Context) {
	attachment, err := models.GetAttachmentByUUID(c.Params(":uuid"))
	if err != nil {
		c.NotFoundOrServerError("GetAttachmentByUUID", models.IsErrAttachmentNotExist, err)
		return
	}

	post, err := models.GetPostByID(attachment.PostID)
	if err != nil {
		if models.IsErrPostNotExist(err) {
			c.Handle(404, "", nil)
		} else {
			c.Handle(500, "GetPostByID", err)
		}
		return
	}

	models.PostIncNumDownloads(post)
}

func ViewPix(c *context.Context) {
	c.Data["Title"] = "作品信息"
	c.Data["PageIsPix"] = true
	c.Data["PageIsViewPix"] = true
	renderAttachmentSettings(c)
	renderCoverSettings(c)

	post, err := models.GetPostByID(c.ParamsInt64(":pixid"))
	if err != nil {
		if models.IsErrPostNotExist(err) {
			c.Handle(404, "", nil)
		} else {
			c.Handle(500, "GetPostByID", err)
		}
		return
	}

	post.Author, err = models.GetUserByPostID(c.ParamsInt64(":pixid"))
	if err != nil {
		if models.IsErrUserNotExist(err) {
			c.Handle(404, "", nil)
		} else {
			c.Handle(500, "GetUserByPostID", err)
		}
		return
	}

	post.Attachments, err = models.GetAttachmentsByPostID(c.ParamsInt64(":pixid"))
	if err != nil {
		if models.IsErrAttachmentNotExist(err) {
			c.Handle(404, "", nil)
		} else {
			c.Handle(500, "GetAttachmentsByPostID", err)
		}
		return
	}

	post.CoverImg, err = models.GetCoverImgsByPostID(c.ParamsInt64(":pixid"))
	if err != nil {
		if models.IsErrCoverImgNotExist(err) {
			c.Handle(404, "", nil)
		} else {
			c.Handle(500, "GetCoverImgsByPostID", err)
		}
		return
	}

	c.Data["title"] = post.Title
	c.Data["content"] = post.Content
	c.Data["Post"] = post

	c.Success(PIXVIEW)
}

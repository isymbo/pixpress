package routes

import (
	"github.com/Unknwon/paginater"
	log "gopkg.in/clog.v1"

	"github.com/isymbo/pixpress/app/controllers/context"
	"github.com/isymbo/pixpress/app/models"
	"github.com/isymbo/pixpress/setting"
)

const (
	HOME                  = "home"
	EXPLORE_WORKS         = "explore/works"
	EXPLORE_USERS         = "explore/users"
	EXPLORE_ORGANIZATIONS = "explore/organizations"
)

func Home(c *context.Context) {
	if c.IsLogged {
		// if !c.User.IsActive && setting.Service.RegisterEmailConfirm {
		// 	c.Data["Title"] = c.Tr("auth.active_your_account")
		// 	c.Success(user.ACTIVATE)
		// } else {
		// 	user.Dashboard(c)
		// }
		if !c.User.IsActive {
			c.Redirect(setting.AppSubURL + "/fixme")
		}
		log.Trace("HOME: c.IsLogged: %+v", c.IsLogged)
		c.SubURLRedirect("user/home")
		return
		// return
	}
	log.Trace("HOME: c.IsLogged: %+v", c.IsLogged)

	// Check auto-login.
	uname := c.GetCookie(setting.Security.CookieUserName)
	log.Trace("uname: %s", uname)
	if len(uname) != 0 {
		c.Redirect(setting.AppSubURL + "/user/login")
		return
	}

	c.Data["PageIsHome"] = true
	c.Success(HOME)
}

func ExploreWorks(c *context.Context) {
	c.Data["Title"] = "发现"
	c.Data["PageIsExplore"] = true
	c.Data["PageIsExploreWorks"] = true

	RenderPostSearch(c, &PostSearchOptions{
		Type:     models.POST_TYPE_PIX,
		Counter:  models.CountPosts,
		Ranger:   models.Posts,
		PageSize: setting.UI.ExplorePagingNum,
		OrderBy:  "updated_unix DESC",
		TplName:  EXPLORE_WORKS,
	})
}

type UserSearchOptions struct {
	Type     models.UserType
	Counter  func() int64
	Ranger   func(int, int) ([]*models.User, error)
	PageSize int
	OrderBy  string
	TplName  string
}

func RenderUserSearch(c *context.Context, opts *UserSearchOptions) {
	page := c.QueryInt("page")
	if page <= 1 {
		page = 1
	}

	var (
		users []*models.User
		count int64
		err   error
	)

	keyword := c.Query("q")
	if len(keyword) == 0 {
		users, err = opts.Ranger(page, opts.PageSize)
		if err != nil {
			c.ServerError("Ranger", err)
			return
		}
		count = opts.Counter()
	} else {
		users, count, err = models.SearchUserByName(&models.SearchUserOptions{
			Keyword:  keyword,
			Type:     opts.Type,
			OrderBy:  opts.OrderBy,
			Page:     page,
			PageSize: opts.PageSize,
		})
		if err != nil {
			c.ServerError("SearchUserByName", err)
			return
		}
	}
	c.Data["Keyword"] = keyword
	c.Data["Total"] = count
	c.Data["Page"] = paginater.New(int(count), opts.PageSize, page, 5)
	c.Data["Users"] = users

	c.Success(opts.TplName)
}

func ExploreUsers(c *context.Context) {
	// c.Data["Title"] = c.Tr("explore")
	// c.Data["PageIsExplore"] = true
	// c.Data["PageIsExploreUsers"] = true

	// RenderUserSearch(c, &UserSearchOptions{
	// 	Type:     models.USER_TYPE_INDIVIDUAL,
	// 	Counter:  models.CountUsers,
	// 	Ranger:   models.Users,
	// 	PageSize: setting.UI.ExplorePagingNum,
	// 	OrderBy:  "updated_unix DESC",
	// 	TplName:  EXPLORE_USERS,
	// })
}

func ExploreOrganizations(c *context.Context) {
	// c.Data["Title"] = c.Tr("explore")
	// c.Data["PageIsExplore"] = true
	// c.Data["PageIsExploreOrganizations"] = true

	// RenderUserSearch(c, &UserSearchOptions{
	// 	Type:     models.USER_TYPE_ORGANIZATION,
	// 	Counter:  models.CountOrganizations,
	// 	Ranger:   models.Organizations,
	// 	PageSize: setting.UI.ExplorePagingNum,
	// 	OrderBy:  "updated_unix DESC",
	// 	TplName:  EXPLORE_ORGANIZATIONS,
	// })
}

func NotFound(c *context.Context) {
	c.Data["Title"] = "Page Not Found"
	c.NotFound()
}

type PostSearchOptions struct {
	Type     models.PostType
	Counter  func() int64
	Ranger   func(int, int) ([]*models.Post, error)
	PageSize int
	OrderBy  string
	TplName  string
}

func RenderPostSearch(c *context.Context, opts *PostSearchOptions) {
	page := c.QueryInt("page")
	if page <= 1 {
		page = 1
	}

	var (
		posts []*models.Post
		count int64
		err   error
	)

	keyword := c.Query("q")
	if len(keyword) == 0 {
		posts, err = opts.Ranger(page, opts.PageSize)
		if err != nil {
			c.ServerError("Ranger", err)
			return
		}
		count = opts.Counter()
	} else {
		posts, count, err = models.SearchPostByName(&models.SearchPostOptions{
			Keyword:  keyword,
			Type:     opts.Type,
			OrderBy:  opts.OrderBy,
			Page:     page,
			PageSize: opts.PageSize,
		})
		if err != nil {
			c.ServerError("SearchPostByName", err)
			return
		}
	}

	for _, p := range posts {
		u, _ := models.GetUserByID(p.AuthorID)
		p.Author = u

		cover, _ := models.GetCoverImgsByPostID(p.ID)
		p.CoverImg = cover
	}

	c.Data["Keyword"] = keyword
	c.Data["Total"] = count
	c.Data["Page"] = paginater.New(int(count), opts.PageSize, page, 5)
	c.Data["Posts"] = posts

	c.Success(opts.TplName)
}

// func ExplorePixes(c *context.Context) {
// 	c.Data["Title"] = "作品列表"
// 	c.Data["PageIsExplore"] = true
// 	c.Data["PageIsExplorePixes"] = true

// 	RenderPostSearch(c, &PostSearchOptions{
// 		Type:     models.POST_TYPE_PIX,
// 		Counter:  models.CountPixes,
// 		Ranger:   models.Pixes,
// 		PageSize: setting.UI.ExplorePagingNum,
// 		OrderBy:  "updated_unix DESC",
// 		TplName:  EXPLORE_PIXES,
// 	})
// }

type PostSearchByAuthorIDOptions struct {
	Type     models.PostType
	Counter  func(int64) int64
	Ranger   func(int, int, int64) ([]*models.Post, error)
	PageSize int
	OrderBy  string
	TplName  string
	AuthorID int64
}

func RenderPostSearchByAuthorID(c *context.Context, opts *PostSearchByAuthorIDOptions) {
	page := c.QueryInt("page")
	if page <= 1 {
		page = 1
	}

	var (
		posts []*models.Post
		count int64
		err   error
	)

	keyword := c.Query("q")
	if len(keyword) == 0 {
		posts, err = opts.Ranger(page, opts.PageSize, opts.AuthorID)
		if err != nil {
			c.ServerError("Ranger", err)
			return
		}
		count = opts.Counter(opts.AuthorID)
	} else {
		posts, count, err = models.SearchPostByAuthorID(&models.SearchPostByAuthorIDOptions{
			AuthorID: opts.AuthorID,
			Type:     opts.Type,
			OrderBy:  opts.OrderBy,
			Page:     page,
			PageSize: opts.PageSize,
		})
		if err != nil {
			c.ServerError("SearchPostByAuthorID", err)
			return
		}
	}

	for _, p := range posts {
		u, _ := models.GetUserByID(p.AuthorID)
		p.Author = u

		cover, _ := models.GetCoverImgsByPostID(p.ID)
		p.CoverImg = cover
	}

	c.Data["Keyword"] = keyword
	c.Data["Total"] = count
	c.Data["Page"] = paginater.New(int(count), opts.PageSize, page, 5)
	c.Data["Posts"] = posts

	c.Success(opts.TplName)
}

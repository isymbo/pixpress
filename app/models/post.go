package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-xorm/xorm"
	log "gopkg.in/clog.v1"

	"github.com/isymbo/pixpress/app/models/errors"
	"github.com/isymbo/pixpress/setting"
)

//
type PostType int

const (
	// Post Type
	POST_TYPE_PIX PostType = iota
	POST_TYPE_ARTICLE
)

// type ContentType string

type Post struct {
	ID          int64
	AuthorID    int64  `xorm:"NOT NULL"`
	Author      *User  `xorm:"-" json:"-"`
	Title       string `xorm:"VARCHAR(100)"`
	Content     string `xorm:"TEXT"`
	PostType    PostType
	CoverImg    string
	NumComments int `xorm:"NOT NULL DEFAULT 0"`

	Created     time.Time `xorm:"-" json:"-"`
	CreatedUnix int64
	Updated     time.Time `xorm:"-" json:"-"`
	UpdatedUnix int64

	Attachments []*Attachment `xorm:"-" json:"-"`
}

type PostsOptions struct {
	UserID   int64
	AuthorID int64
	Page     int
	Labels   string
	// SortType    string
}

func (po *Post) BeforeInsert() {
	po.CreatedUnix = time.Now().Unix()
	po.UpdatedUnix = po.CreatedUnix
}

func (po *Post) BeforeUpdate() {
	po.UpdatedUnix = time.Now().Unix()
}

func (po *Post) AfterSet(colName string, _ xorm.Cell) {
	switch colName {
	case "created_unix":
		po.Created = time.Unix(po.CreatedUnix, 0).Local()
	case "updated_unix":
		po.Updated = time.Unix(po.UpdatedUnix, 0).Local()
	}
}

func NewPost(p *Post) (err error) {
	sess := x.NewSession()
	defer sess.Close()
	if err = sess.Begin(); err != nil {
		return err
	}

	if _, err = sess.Insert(p); err != nil {
		return err
	}

	log.Trace("New post %d is created", p.ID)
	return sess.Commit()
}

// func Posts(opts *PostsOptions) ([]*Post, error) {
// 	sess := buildPostsQuery(opts)
// 	if sess == nil {
// 		return make([]*Post, 0), nil
// 	}

// 	sess.Limit(setting.UI.IssuePagingNum, (opts.Page-1)*setting.UI.IssuePagingNum)

// 	posts := make([]*Post, 0, setting.UI.IssuePagingNum)
// 	if err := sess.Find(&posts); err != nil {
// 		return nil, fmt.Errorf("Find: %v", err)
// 	}

// 	// // FIXME: use IssueList to improve performance.
// 	// for i := range issues {
// 	// 	if err := issues[i].LoadAttributes(); err != nil {
// 	// 		return nil, fmt.Errorf("LoadAttributes [%d]: %v", issues[i].ID, err)
// 	// 	}
// 	// }

// 	return posts, nil
// }

type SearchPostOptions struct {
	Keyword  string
	Type     PostType
	OrderBy  string
	Page     int
	PageSize int // Can be smaller than or equal to setting.UI.ExplorePagingNum
}

// SearchUserByName takes keyword and part of user name to search,
// it returns results in given range and number of total results.
func SearchPostByName(opts *SearchPostOptions) (posts []*Post, _ int64, _ error) {
	if len(opts.Keyword) == 0 {
		return posts, 0, nil
	}

	opts.Keyword = strings.ToLower(opts.Keyword)

	if opts.PageSize <= 0 || opts.PageSize > setting.UI.ExplorePagingNum {
		opts.PageSize = setting.UI.ExplorePagingNum
	}
	if opts.Page <= 0 {
		opts.Page = 1
	}

	searchQuery := "%" + opts.Keyword + "%"
	posts = make([]*Post, 0, opts.PageSize)
	// Append conditions
	sess := x.Where("LOWER(lower_name) LIKE ?", searchQuery).
		Or("LOWER(full_name) LIKE ?", searchQuery).
		And("type = ?", opts.Type)

	var countSess xorm.Session
	countSess = *sess
	count, err := countSess.Count(new(Post))
	if err != nil {
		return nil, 0, fmt.Errorf("Count: %v", err)
	}

	if len(opts.OrderBy) > 0 {
		sess.OrderBy(opts.OrderBy)
	}
	return posts, count, sess.Limit(opts.PageSize, (opts.Page-1)*opts.PageSize).Find(&posts)
}

func countPosts(e Engine) int64 {
	count, _ := e.Where("type=0").Count(new(Post))
	return count
}

// CountPosts returns number of posts.
func CountPosts() int64 {
	return countPosts(x)
}

// Posts returns number of posts in given page.
func Posts(page, pageSize int) ([]*Post, error) {
	posts := make([]*Post, 0, pageSize)
	return posts, x.Limit(pageSize, (page-1)*pageSize).Where("post_type=0").Asc("id").Find(&posts)
}

func getRawPostByID(e Engine, id int64) (*Post, error) {
	post := new(Post)
	has, err := e.Id(id).Get(post)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, errors.PostNotExist{id, 0}
	}
	return post, nil
}

// GetPostByID returns a post by given ID.
func GetPostByID(id int64) (*Post, error) {
	return getRawPostByID(x, id)
}

func updatePost(e Engine, p *Post) error {
	_, err := e.Id(p.ID).AllCols().Update(p)
	return err
}

// UpdatePost updates information of given post.
func UpdatePost(p *Post) error {
	return updatePost(x, p)
}

func deletePost(e Engine, p *Post) error {
	_, err := e.Id(p.ID).AllCols().Delete(p)
	return err
}

// DeletePost deletes information of given post.
func DeletePost(p *Post) error {
	return deletePost(x, p)
}

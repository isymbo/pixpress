package models

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Unknwon/com"
	"github.com/go-xorm/xorm"
	"github.com/nfnt/resize"
	log "gopkg.in/clog.v1"

	"github.com/isymbo/pixpress/app/models/errors"
	"github.com/isymbo/pixpress/setting"
)

type PostType int

const (
	// Post Type
	POST_TYPE_PIX PostType = iota
	POST_TYPE_ARTICLE
)

// type ContentType string

type Post struct {
	ID           int64
	AuthorID     int64  `xorm:"NOT NULL"`
	Author       *User  `xorm:"-" json:"-"`
	Title        string `xorm:"VARCHAR(100)"`
	Content      string `xorm:"TEXT"`
	PostType     PostType
	NumComments  int `xorm:"NOT NULL DEFAULT 0"`
	NumViews     int `xorm:"NOT NULL DEFAULT 0"`
	NumLikes     int `xorm:"NOT NULL DEFAULT 0"`
	NumDownloads int `xorm:"NOT NULL DEFAULT 0"`

	Created     time.Time `xorm:"-" json:"-"`
	CreatedUnix int64
	Updated     time.Time `xorm:"-" json:"-"`
	UpdatedUnix int64

	CoverImg    []*CoverImg   `xorm:"-" json:"-"`
	Attachments []*Attachment `xorm:"-" json:"-"`
}

type CreatePostOptions struct {
	PostType    PostType
	Doer        *User
	Title       string
	Content     string
	Attachments []string // UUIDs of attachments
	CoverImg    []string // UUIDs of coverimg
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

// SearchPostByName takes keyword and part of user name to search,
// it returns results in given range and number of total results.
func SearchPostByName(opts *SearchPostOptions) (posts []*Post, _ int64, _ error) {
	if len(opts.Keyword) == 0 {
		return posts, 0, nil
	}

	opts.Keyword = strings.ToLower(opts.Keyword)

	// if opts.PageSize <= 0 || opts.PageSize > setting.UI.User.PostPagingNum {
	// 	opts.PageSize = setting.UI.User.PostPagingNum
	// }
	if opts.Page <= 0 {
		opts.Page = 1
	}

	searchQuery := "%" + opts.Keyword + "%"
	posts = make([]*Post, 0, opts.PageSize)

	// Append conditions
	// sess := x.Where("LOWER(lower_name) LIKE ?", searchQuery).
	// 	Or("LOWER(full_name) LIKE ?", searchQuery).
	// 	And("type = ?", opts.Type)
	sess := x.Where("title LIKE ?", searchQuery).Or("title LIKE ?", searchQuery)
	if DbCfg.DbType == "mysql" {
		sess = x.Where("MATCH (title, content) AGAINST (? IN NATURAL LANGUAGE MODE)", opts.Keyword)
	}

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

type SearchPostByAuthorIDOptions struct {
	AuthorID int64
	Type     PostType
	OrderBy  string
	Page     int
	PageSize int // Can be smaller than or equal to setting.UI.ExplorePagingNum
}

// SearchPostByName takes keyword and part of user name to search,
// it returns results in given range and number of total results.
func SearchPostByAuthorID(opts *SearchPostByAuthorIDOptions) (posts []*Post, _ int64, _ error) {
	if opts.AuthorID == 0 {
		return posts, 0, nil
	}

	if opts.PageSize <= 0 || opts.PageSize > setting.UI.User.PostPagingNum {
		opts.PageSize = setting.UI.User.PostPagingNum
	}
	if opts.Page <= 0 {
		opts.Page = 1
	}

	// searchQuery := "%" + opts.Keyword + "%"
	searchQuery := strconv.FormatInt(opts.AuthorID, 10)
	posts = make([]*Post, 0, opts.PageSize)
	// Append conditions
	sess := x.Where("author_id = ?", searchQuery).
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
	count, _ := e.Where("post_type=0").Count(new(Post))
	return count
}

// CountPosts returns number of posts.
func CountPosts() int64 {
	return countPosts(x)
}

// Posts returns number of posts in given page.
func Posts(page, pageSize int) ([]*Post, error) {
	posts := make([]*Post, 0, pageSize)
	return posts, x.Limit(pageSize, (page-1)*pageSize).Where("post_type=0").Desc("id").Find(&posts)
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

// UploadCoverImg saves cover image for post.
// FIXME: split uploads to different subdirs in case we have massive number of posts.
func (post *Post) UploadCoverImg(data []byte) error {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("decode image: %v", err)
	}

	os.MkdirAll(setting.Cover.Path, os.ModePerm)
	fw, err := os.Create(post.CoverImgPath())
	if err != nil {
		return fmt.Errorf("create post cover image directory: %v", err)
	}
	defer fw.Close()

	// m := resize.Resize(avatar.AVATAR_SIZE, avatar.AVATAR_SIZE, img, resize.NearestNeighbor)
	m := resize.Resize(300, 300, img, resize.NearestNeighbor)
	if err = png.Encode(fw, m); err != nil {
		return fmt.Errorf("encode image: %v", err)
	}

	return nil
}

// DeleteCoverImg deletes the post cover image.
func (post *Post) DeleteCoverImg() error {
	log.Trace("DeleteCoverImg [%d]: %s", post.ID, post.CoverImgPath())
	if err := os.Remove(post.CoverImgPath()); err != nil {
		return err
	}

	return UpdatePost(post)
}

// CoverImgPath returns post cover image file path.
func (post *Post) CoverImgPath() string {
	return filepath.Join(setting.Cover.Path, com.ToStr(post.ID))
}

// CreatePost creates a post.
func CreatePostByOption(opts *CreatePostOptions) (post *Post, err error) {
	sess := x.NewSession()
	defer sess.Close()
	if err = sess.Begin(); err != nil {
		return nil, err
	}

	post, err = createPostByOption(sess, opts)
	if err != nil {
		return nil, err
	}

	return post, sess.Commit()
}

// type Post struct {
// 	ID          int64
// 	AuthorID    int64  `xorm:"NOT NULL"`
// 	Author      *User  `xorm:"-" json:"-"`
// 	Title       string `xorm:"VARCHAR(100)"`
// 	Content     string `xorm:"TEXT"`
// 	PostType    PostType
// 	CoverImg    string
// 	NumComments int `xorm:"NOT NULL DEFAULT 0"`

// 	Created     time.Time `xorm:"-" json:"-"`
// 	CreatedUnix int64
// 	Updated     time.Time `xorm:"-" json:"-"`
// 	UpdatedUnix int64

// 	Attachments []*Attachment `xorm:"-" json:"-"`
// }

func createPostByOption(e *xorm.Session, opts *CreatePostOptions) (_ *Post, err error) {
	post := &Post{
		AuthorID: opts.Doer.ID,
		Author:   opts.Doer,
		PostType: opts.PostType,
		Title:    opts.Title,
		Content:  opts.Content,
	}
	if _, err = e.Insert(post); err != nil {
		return nil, err
	}

	// Check attachments
	attachments := make([]*Attachment, 0, len(opts.Attachments))
	for _, uuid := range opts.Attachments {
		attach, err := getAttachmentByUUID(e, uuid)
		if err != nil {
			if IsErrAttachmentNotExist(err) {
				continue
			}
			return nil, fmt.Errorf("getAttachmentByUUID [%s]: %v", uuid, err)
		}
		attachments = append(attachments, attach)
	}

	for i := range attachments {
		attachments[i].PostID = post.ID

		// No assign value could be 0, so ignore AllCols().
		if _, err = e.Id(attachments[i].ID).Update(attachments[i]); err != nil {
			return nil, fmt.Errorf("update attachment [%d]: %v", attachments[i].ID, err)
		}
	}

	// Check coverimg
	coverimg := make([]*CoverImg, 0, len(opts.CoverImg))
	for _, uuid := range opts.CoverImg {
		cover, err := getCoverImgByUUID(e, uuid)
		if err != nil {
			if IsErrCoverImgNotExist(err) {
				continue
			}
			return nil, fmt.Errorf("getCoverImgByUUID [%s]: %v", uuid, err)
		}
		coverimg = append(coverimg, cover)
	}

	for i := range coverimg {
		coverimg[i].PostID = post.ID

		// No assign value could be 0, so ignore AllCols().
		if _, err = e.Id(coverimg[i].ID).Update(coverimg[i]); err != nil {
			return nil, fmt.Errorf("update coverimg [%d]: %v", coverimg[i].ID, err)
		}
	}

	return post, nil
}

func PostIncNumViews(p *Post) error {
	p.NumViews += 1
	return updatePost(x, p)
}

func PostIncNumDownloads(p *Post) error {
	p.NumDownloads += 1
	return updatePost(x, p)
}

func countPostsByAuthorID(e Engine, authorID int64) int64 {
	count, _ := e.Where("author_id = ?", authorID).And("post_type=0").Count(new(Post))
	return count
}

// CountPosts returns number of posts by authorID.
func CountPostsByAuthorID(authorID int64) int64 {
	return countPostsByAuthorID(x, authorID)
}

// Posts returns number of posts in given page.
func PostsByAuthorID(page, pageSize int, authorID int64) ([]*Post, error) {
	posts := make([]*Post, 0, pageSize)
	return posts, x.Limit(pageSize, (page-1)*pageSize).Where("author_id = ?", authorID).And("post_type=0").Desc("id").Find(&posts)
}

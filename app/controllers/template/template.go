// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package template

import (
	"container/list"
	"fmt"
	"html/template"
	"mime"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"

	"github.com/isymbo/pixpress/setting"
	"github.com/isymbo/pixpress/util"
)

// TODO: only initialize map once and save to a local variable to reduce copies.
func NewFuncMap() []template.FuncMap {
	return []template.FuncMap{map[string]interface{}{
		"GoVer": func() string {
			return strings.Title(runtime.Version())
		},
		// "UseHTTPS": func() bool {
		// 	return strings.HasPrefix(setting.AppURL, "https")
		// },
		"AppName": func() string {
			return setting.AppName
		},
		"AppSubURL": func() string {
			return setting.AppSubURL
		},
		"AppURL": func() string {
			return setting.AppURL
		},
		"AppVer": func() string {
			return setting.AppVer
		},
		"ShowFooterTemplateLoadTime": func() bool {
			return setting.Other.ShowFooterTemplateLoadTime
		},
		"LoadTimes": func(startTime time.Time) string {
			return fmt.Sprint(time.Since(startTime).Nanoseconds()/1e6) + "ms"
		},
		// "AvatarLink":       tool.AvatarLink,
		// "AppendAvatarSize": tool.AppendAvatarSize,
		"Safe":     Safe,
		"Sanitize": bluemonday.UGCPolicy().Sanitize,
		//"Str2html": Str2html,
		"TimeSince":    util.TimeSince,
		"RawTimeSince": util.RawTimeSince,
		"FileSize":     util.FileSize,
		"Subtract":     util.Subtract,
		"Add": func(a, b int) int {
			return a + b
		},
		// "ActionIcon": ActionIcon,
		"DateFmtLong": func(t time.Time) string {
			return t.Format(time.RFC1123Z)
		},
		"DateFmtShort": func(t time.Time) string {
			return t.Format("Jan 02, 2006")
		},
		"List": List,
		"SubStr": func(str string, start, length int) string {
			if len(str) == 0 {
				return ""
			}
			end := start + length
			if length == -1 {
				end = len(str)
			}
			if len(str) < end {
				return str
			}
			return str[start:end]
		},
		"Join":           strings.Join,
		"EllipsisString": util.EllipsisString,
		"ShortSHA1":      util.ShortSHA1,
		"MD5":            util.MD5,
		"FilenameIsImage": func(filename string) bool {
			mimeType := mime.TypeByExtension(filepath.Ext(filename))
			return strings.HasPrefix(mimeType, "image/")
		},
	}}
}

func Safe(raw string) template.HTML {
	return template.HTML(raw)
}

// func Str2html(raw string) template.HTML {
// 	return template.HTML(markup.Sanitize(raw))
// }

func List(l *list.List) chan interface{} {
	e := l.Front()
	c := make(chan interface{})
	go func() {
		for e != nil {
			c <- e.Value
			e = e.Next()
		}
		close(c)
	}()
	return c
}

// Replaces all prefixes 'old' in 's' with 'new'.
func ReplaceLeft(s, old, new string) string {
	old_len, new_len, i, n := len(old), len(new), 0, 0
	for ; i < len(s) && strings.HasPrefix(s[i:], old); n += 1 {
		i += old_len
	}

	// simple optimization
	if n == 0 {
		return s
	}

	// allocating space for the new string
	newLen := n*new_len + len(s[i:])
	replacement := make([]byte, newLen, newLen)

	j := 0
	for ; j < n*new_len; j += new_len {
		copy(replacement[j:j+new_len], new)
	}

	copy(replacement[j:], s[i:])
	return string(replacement)
}

// NewLine2br simply replaces "\n" to "<br>".
func NewLine2br(raw string) string {
	return strings.Replace(raw, "\n", "<br>", -1)
}

func Sha1(str string) string {
	return util.SHA1(str)
}

// func ToUTF8WithErr(content []byte) (error, string) {
// 	charsetLabel, err := util.DetectEncoding(content)
// 	if err != nil {
// 		return err, ""
// 	} else if charsetLabel == "UTF-8" {
// 		return nil, string(content)
// 	}

// 	encoding, _ := charset.Lookup(charsetLabel)
// 	if encoding == nil {
// 		return fmt.Errorf("Unknown encoding: %s", charsetLabel), string(content)
// 	}

// 	// If there is an error, we concatenate the nicely decoded part and the
// 	// original left over. This way we won't loose data.
// 	result, n, err := transform.String(encoding.NewDecoder(), string(content))
// 	if err != nil {
// 		result = result + string(content[n:])
// 	}

// 	return err, result
// }

// func ToUTF8(content string) string {
// 	_, res := ToUTF8WithErr([]byte(content))
// 	return res
// }

func EscapePound(str string) string {
	return strings.NewReplacer("%", "%25", "#", "%23", " ", "%20", "?", "%3F").Replace(str)
}

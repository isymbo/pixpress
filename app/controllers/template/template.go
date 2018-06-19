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
		// "TimeSince":    tool.TimeSince,
		// "RawTimeSince": tool.RawTimeSince,
		// "FileSize":     tool.FileSize,
		// "Subtract":     tool.Subtract,
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
		"Join": strings.Join,
		// "EllipsisString":    tool.EllipsisString,
		// "Sha1": Sha1,
		// "ShortSHA1":         tool.ShortSHA1,
		// "MD5":               tool.MD5,
		// "RenderCommitMessage":   RenderCommitMessage,
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

// func Str2html(raw string) template.HTML {
// 	return template.HTML(markup.Sanitize(raw))
// }

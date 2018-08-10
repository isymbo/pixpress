package form

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/Unknwon/com"
	"github.com/go-macaron/binding"
	macaron "gopkg.in/macaron.v1"
)

const ERR_ALPHA_DASH_DOT_SLASH = "AlphaDashDotSlashError"

var AlphaDashDotSlashPattern = regexp.MustCompile("[^\\d\\w-_\\./]")

func init() {
	binding.SetNameMapper(com.ToSnakeCase)
	binding.AddRule(&binding.Rule{
		IsMatch: func(rule string) bool {
			return rule == "AlphaDashDotSlash"
		},
		IsValid: func(errs binding.Errors, name string, v interface{}) (bool, binding.Errors) {
			if AlphaDashDotSlashPattern.MatchString(fmt.Sprintf("%v", v)) {
				errs.Add([]string{name}, ERR_ALPHA_DASH_DOT_SLASH, "AlphaDashDotSlash")
				return false, errs
			}
			return true, errs
		},
	})
}

type Form interface {
	binding.Validator
}

// Assign assign form values back to the template data.
func Assign(form interface{}, data map[string]interface{}) {
	typ := reflect.TypeOf(form)
	val := reflect.ValueOf(form)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		fieldName := field.Tag.Get("form")
		// Allow ignored fields in the struct
		if fieldName == "-" {
			continue
		} else if len(fieldName) == 0 {
			fieldName = com.ToSnakeCase(field.Name)
		}

		data[fieldName] = val.Field(i).Interface()
	}
}

func getRuleBody(field reflect.StructField, prefix string) string {
	for _, rule := range strings.Split(field.Tag.Get("binding"), ";") {
		if strings.HasPrefix(rule, prefix) {
			return rule[len(prefix) : len(rule)-1]
		}
	}
	return ""
}

func getSize(field reflect.StructField) string {
	return getRuleBody(field, "Size(")
}

func getMinSize(field reflect.StructField) string {
	return getRuleBody(field, "MinSize(")
}

func getMaxSize(field reflect.StructField) string {
	return getRuleBody(field, "MaxSize(")
}

func getInclude(field reflect.StructField) string {
	return getRuleBody(field, "Include(")
}

func validate(errs binding.Errors, data map[string]interface{}, f Form, l macaron.Locale) binding.Errors {
	if errs.Len() == 0 {
		return errs
	}

	data["HasError"] = true
	Assign(f, data)

	typ := reflect.TypeOf(f)
	val := reflect.ValueOf(f)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		fieldName := field.Tag.Get("form")
		// Allow ignored fields in the struct
		if fieldName == "-" {
			continue
		}

		if errs[0].FieldNames[0] == field.Name {
			data["Err_"+field.Name] = true

			// FIXME, TODO, will use locale in the future
			// trName := field.Tag.Get("locale")
			// if len(trName) == 0 {
			// 	trName = l.Tr("form." + field.Name)
			// } else {
			// 	trName = l.Tr(trName)
			// }

			trName := field.Tag.Get("locale")
			if len(trName) == 0 {
				trName = field.Name
			}

			switch errs[0].Classification {
			case binding.ERR_REQUIRED:
				// data["ErrorMsg"] = trName + l.Tr("form.require_error")
				data["ErrorMsg"] = trName + "不能为空。"
			case binding.ERR_ALPHA_DASH:
				// data["ErrorMsg"] = trName + l.Tr("form.alpha_dash_error")
				data["ErrorMsg"] = trName + "必须为英文字母、阿拉伯数字或横线（-_）。"
			case binding.ERR_ALPHA_DASH_DOT:
				// data["ErrorMsg"] = trName + l.Tr("form.alpha_dash_dot_error")
				data["ErrorMsg"] = trName + "必须为英文字母、阿拉伯数字、横线（-_）或点。"
			case ERR_ALPHA_DASH_DOT_SLASH:
				// data["ErrorMsg"] = trName + l.Tr("form.alpha_dash_dot_slash_error")
				data["ErrorMsg"] = trName + "必须为英文字母、阿拉伯数字、横线（-_）、点或斜线。"
			case binding.ERR_SIZE:
				// data["ErrorMsg"] = trName + l.Tr("form.size_error", getSize(field))
				data["ErrorMsg"] = trName + "长度必须为 " + getSize(field)
			case binding.ERR_MIN_SIZE:
				// data["ErrorMsg"] = trName + l.Tr("form.min_size_error", getMinSize(field))
				data["ErrorMsg"] = trName + "长度最小为 " + getMinSize(field) + " 个字符。"
			case binding.ERR_MAX_SIZE:
				// data["ErrorMsg"] = trName + l.Tr("form.max_size_error", getMaxSize(field))
				data["ErrorMsg"] = trName + "长度最大为 " + getMaxSize(field) + " 个字符。"
			case binding.ERR_EMAIL:
				// data["ErrorMsg"] = trName + l.Tr("form.email_error")
				data["ErrorMsg"] = trName + "不是一个有效的邮箱地址。"
			case binding.ERR_URL:
				// data["ErrorMsg"] = trName + l.Tr("form.url_error")
				data["ErrorMsg"] = trName + "不是一个有效的 URL。"
			case binding.ERR_INCLUDE:
				// data["ErrorMsg"] = trName + l.Tr("form.include_error", getInclude(field))
				data["ErrorMsg"] = trName + "必须包含子字符串 " + getInclude(field) + " 。"
			default:
				// data["ErrorMsg"] = l.Tr("form.unknown_error") + " " + errs[0].Classification
				data["ErrorMsg"] = "未知错误：" + errs[0].Classification
			}

			return errs
		}
	}

	return errs
}

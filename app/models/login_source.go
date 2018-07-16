package models

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Unknwon/com"
	"github.com/go-macaron/binding"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	jsoniter "github.com/json-iterator/go"
	log "gopkg.in/clog.v1"
	"gopkg.in/ini.v1"

	"github.com/isymbo/pixpress/app/controllers/auth/ldap"
	"github.com/isymbo/pixpress/app/models/errors"
)

type LoginType int

// Note: new type must append to the end of list to maintain compatibility.
const (
	LOGIN_NOTYPE LoginType = iota
	LOGIN_PLAIN            // 1
	LOGIN_LDAP             // 2
	LOGIN_DLDAP            // 3
)

var LoginNames = map[LoginType]string{
	LOGIN_LDAP:  "LDAP (via BindDN)",  // Direct Bind
	LOGIN_DLDAP: "LDAP (simple auth)", // Simple Auth
}

var SecurityProtocolNames = map[ldap.SecurityProtocol]string{
	ldap.SECURITY_PROTOCOL_UNENCRYPTED: "Unencrypted",
	ldap.SECURITY_PROTOCOL_LDAPS:       "LDAPS",
	ldap.SECURITY_PROTOCOL_START_TLS:   "StartTLS",
}

// Ensure structs implemented interface.
var (
	_ core.Conversion = &LDAPConfig{}
)

type LDAPConfig struct {
	*ldap.Source `ini:"config"`
}

func (cfg *LDAPConfig) FromDB(bs []byte) error {
	return jsoniter.Unmarshal(bs, &cfg)
}

func (cfg *LDAPConfig) ToDB() ([]byte, error) {
	return jsoniter.Marshal(cfg)
}

func (cfg *LDAPConfig) SecurityProtocolName() string {
	return SecurityProtocolNames[cfg.SecurityProtocol]
}

// AuthSourceFile contains information of an authentication source file.
type AuthSourceFile struct {
	abspath string
	file    *ini.File
}

// SetGeneral sets new value to the given key in the general (default) section.
func (f *AuthSourceFile) SetGeneral(name, value string) {
	f.file.Section("").Key(name).SetValue(value)
}

// SetConfig sets new values to the "config" section.
func (f *AuthSourceFile) SetConfig(cfg core.Conversion) error {
	return f.file.Section("config").ReflectFrom(cfg)
}

// Save writes updates into file system.
func (f *AuthSourceFile) Save() error {
	return f.file.SaveTo(f.abspath)
}

// LoginSource represents an external way for authorizing users.
type LoginSource struct {
	ID        int64
	Type      LoginType
	Name      string          `xorm:"UNIQUE"`
	IsActived bool            `xorm:"NOT NULL DEFAULT false"`
	Cfg       core.Conversion `xorm:"TEXT"`

	Created     time.Time `xorm:"-" json:"-"`
	CreatedUnix int64
	Updated     time.Time `xorm:"-" json:"-"`
	UpdatedUnix int64

	LocalFile *AuthSourceFile `xorm:"-" json:"-"`
}

func (s *LoginSource) BeforeInsert() {
	s.CreatedUnix = time.Now().Unix()
	s.UpdatedUnix = s.CreatedUnix
}

func (s *LoginSource) BeforeUpdate() {
	s.UpdatedUnix = time.Now().Unix()
}

// Cell2Int64 converts a xorm.Cell type to int64,
// and handles possible irregular cases.
func Cell2Int64(val xorm.Cell) int64 {
	switch (*val).(type) {
	case []uint8:
		log.Trace("Cell2Int64 ([]uint8): %v", *val)
		return com.StrTo(string((*val).([]uint8))).MustInt64()
	}
	return (*val).(int64)
}

func (s *LoginSource) BeforeSet(colName string, val xorm.Cell) {
	switch colName {
	case "type":
		switch LoginType(Cell2Int64(val)) {
		case LOGIN_LDAP, LOGIN_DLDAP:
			s.Cfg = new(LDAPConfig)
		default:
			panic("unrecognized login source type: " + com.ToStr(*val))
		}
	}
}

func (s *LoginSource) AfterSet(colName string, _ xorm.Cell) {
	switch colName {
	case "created_unix":
		s.Created = time.Unix(s.CreatedUnix, 0).Local()
	case "updated_unix":
		s.Updated = time.Unix(s.UpdatedUnix, 0).Local()
	}
}

func (s *LoginSource) TypeName() string {
	return LoginNames[s.Type]
}

func (s *LoginSource) IsLDAP() bool {
	return s.Type == LOGIN_LDAP
}

func (s *LoginSource) IsDLDAP() bool {
	return s.Type == LOGIN_DLDAP
}

func (s *LoginSource) HasTLS() bool {
	return ((s.IsLDAP() || s.IsDLDAP()) &&
		s.LDAP().SecurityProtocol > ldap.SECURITY_PROTOCOL_UNENCRYPTED)
}

func (s *LoginSource) UseTLS() bool {
	switch s.Type {
	case LOGIN_LDAP, LOGIN_DLDAP:
		return s.LDAP().SecurityProtocol != ldap.SECURITY_PROTOCOL_UNENCRYPTED
	}

	return false
}

func (s *LoginSource) SkipVerify() bool {
	switch s.Type {
	case LOGIN_LDAP, LOGIN_DLDAP:
		return s.LDAP().SkipVerify
	}

	return false
}

func (s *LoginSource) LDAP() *LDAPConfig {
	return s.Cfg.(*LDAPConfig)
}

func CreateLoginSource(source *LoginSource) error {
	has, err := x.Get(&LoginSource{Name: source.Name})
	if err != nil {
		return err
	} else if has {
		return ErrLoginSourceAlreadyExist{source.Name}
	}

	_, err = x.Insert(source)
	return err
}

// LoginSources returns all login sources defined.
func LoginSources() ([]*LoginSource, error) {
	sources := make([]*LoginSource, 0, 2)
	if err := x.Find(&sources); err != nil {
		return nil, err
	}

	return append(sources, localLoginSources.List()...), nil
}

// ActivatedLoginSources returns login sources that are currently activated.
func ActivatedLoginSources() ([]*LoginSource, error) {
	sources := make([]*LoginSource, 0, 2)
	if err := x.Where("is_actived = ?", true).Find(&sources); err != nil {
		return nil, fmt.Errorf("find activated login sources: %v", err)
	}
	return append(sources, localLoginSources.ActivatedList()...), nil
}

// GetLoginSourceByID returns login source by given ID.
func GetLoginSourceByID(id int64) (*LoginSource, error) {
	source := new(LoginSource)
	has, err := x.Id(id).Get(source)
	if err != nil {
		return nil, err
	} else if !has {
		return localLoginSources.GetLoginSourceByID(id)
	}
	return source, nil
}

// UpdateLoginSource updates information of login source to database or local file.
func UpdateLoginSource(source *LoginSource) error {
	if source.LocalFile == nil {
		_, err := x.Id(source.ID).AllCols().Update(source)

		return err
	}

	source.LocalFile.SetGeneral("name", source.Name)
	source.LocalFile.SetGeneral("is_activated", com.ToStr(source.IsActived))
	if err := source.LocalFile.SetConfig(source.Cfg); err != nil {
		return fmt.Errorf("LocalFile.SetConfig: %v", err)
	} else if err = source.LocalFile.Save(); err != nil {
		return fmt.Errorf("LocalFile.Save: %v", err)
	}

	localLoginSources.UpdateLoginSource(source)
	return nil
}

func DeleteSource(source *LoginSource) error {
	count, err := x.Count(&User{LoginSource: source.ID})
	if err != nil {
		return err
	} else if count > 0 {
		return ErrLoginSourceInUse{source.ID}
	}
	_, err = x.Id(source.ID).Delete(new(LoginSource))
	return err
}

// CountLoginSources returns total number of login sources.
func CountLoginSources() int64 {
	count, _ := x.Count(new(LoginSource))
	return count + int64(localLoginSources.Len())
}

// LocalLoginSources contains authentication sources configured and loaded from local files.
// Calling its methods is thread-safe; otherwise, please maintain the mutex accordingly.
type LocalLoginSources struct {
	sync.RWMutex
	sources []*LoginSource
}

func (s *LocalLoginSources) Len() int {
	return len(s.sources)
}

// List returns full clone of login sources.
func (s *LocalLoginSources) List() []*LoginSource {
	s.RLock()
	defer s.RUnlock()

	list := make([]*LoginSource, s.Len())
	for i := range s.sources {
		list[i] = &LoginSource{}
		*list[i] = *s.sources[i]
	}
	return list
}

// ActivatedList returns clone of activated login sources.
func (s *LocalLoginSources) ActivatedList() []*LoginSource {
	s.RLock()
	defer s.RUnlock()

	list := make([]*LoginSource, 0, 2)
	for i := range s.sources {
		if !s.sources[i].IsActived {
			continue
		}

		source := &LoginSource{}
		*source = *s.sources[i]
		list = append(list, source)
	}
	return list
}

// GetLoginSourceByID returns a clone of login source by given ID.
func (s *LocalLoginSources) GetLoginSourceByID(id int64) (*LoginSource, error) {
	s.RLock()
	defer s.RUnlock()

	for i := range s.sources {
		if s.sources[i].ID == id {
			source := &LoginSource{}
			*source = *s.sources[i]
			return source, nil
		}
	}

	return nil, errors.LoginSourceNotExist{id}
}

// UpdateLoginSource updates in-memory copy of the authentication source.
func (s *LocalLoginSources) UpdateLoginSource(source *LoginSource) {
	s.Lock()
	defer s.Unlock()

	source.Updated = time.Now()
	for i := range s.sources {
		if s.sources[i].ID == source.ID {
			*s.sources[i] = *source
			break
		}
	}
}

var localLoginSources = &LocalLoginSources{}

// LoadAuthSources loads authentication sources from local files
// and converts them into login sources.
func LoadAuthSources(p string) {
	authdPath := p
	if !com.IsDir(authdPath) {
		return
	}

	paths, err := com.GetFileListBySuffix(authdPath, ".conf")
	if err != nil {
		log.Fatal(2, "Failed to list authentication sources: %v", err)
	}

	localLoginSources.sources = make([]*LoginSource, 0, len(paths))

	for _, fpath := range paths {
		authSource, err := ini.Load(fpath)
		if err != nil {
			log.Fatal(2, "Failed to load authentication source: %v", err)
		}
		//authSource.NameMapper = ini.TitleUnderscore

		// Set general attributes
		s := authSource.Section("")
		loginSource := &LoginSource{
			// ID:        s.Key("id").MustInt64(),
			// Name:      s.Key("name").String(),
			// IsActived: s.Key("is_activated").MustBool(),
			ID:        s.Key("ID").MustInt64(),
			Name:      s.Key("NAME").String(),
			IsActived: s.Key("IS_ACTIVATED").MustBool(),
			LocalFile: &AuthSourceFile{
				abspath: fpath,
				file:    authSource,
			},
		}

		fi, err := os.Stat(fpath)
		if err != nil {
			log.Fatal(2, "Failed to load authentication source: %v", err)
		}
		loginSource.Updated = fi.ModTime()

		// Parse authentication source file
		// authType := s.Key("type").String()
		authType := s.Key("TYPE").String()
		switch authType {
		case "ldap_bind_dn":
			loginSource.Type = LOGIN_LDAP
			loginSource.Cfg = &LDAPConfig{}
		case "ldap_simple_auth":
			loginSource.Type = LOGIN_DLDAP
			loginSource.Cfg = &LDAPConfig{}
		default:
			log.Fatal(2, "Failed to load authentication source: unknown type '%s'", authType)
		}

		if err = authSource.Section("config").MapTo(loginSource.Cfg); err != nil {
			log.Fatal(2, "Failed to parse authentication source 'config': %v", err)
		}

		localLoginSources.sources = append(localLoginSources.sources, loginSource)
	}
}

// .____     ________      _____ __________
// |    |    \______ \    /  _  \\______   \
// |    |     |    |  \  /  /_\  \|     ___/
// |    |___  |    `   \/    |    \    |
// |_______ \/_______  /\____|__  /____|
//         \/        \/         \/

func composeFullName(firstname, surname, username string) string {
	switch {
	case len(firstname) == 0 && len(surname) == 0:
		return username
	case len(firstname) == 0:
		return surname
	case len(surname) == 0:
		return firstname
	default:
		return firstname + " " + surname
	}
}

// LoginViaLDAP queries if login/password is valid against the LDAP directory pool,
// and create a local user if success when enabled.
func LoginViaLDAP(user *User, login, password string, source *LoginSource, autoRegister bool) (*User, error) {
	displayname, mobile, mail, uid, isAdmin, succeed := source.Cfg.(*LDAPConfig).SearchEntry(login, password, source.Type == LOGIN_LDAP)
	if !succeed {
		// User not in LDAP, do nothing
		return nil, errors.UserNotExist{0, login}
	}

	if !autoRegister {
		return user, nil
	}

	// // Fallback.
	// if len(username) == 0 {
	// 	username = login
	// }

	// Force Fallback
	username := login

	// Validate username make sure it satisfies requirement.
	if binding.AlphaDashDotPattern.MatchString(username) {
		return nil, fmt.Errorf("Invalid pattern for attribute 'username' [%s]: must be valid alpha or numeric or dash(-_) or dot characters", username)
	}

	if len(mail) == 0 {
		mail = fmt.Sprintf("%s@localhost", username)
	}

	user = &User{
		Email:       mail,
		Mobile:      mobile,
		DisplayName: displayname,
		LoginType:   source.Type,
		LoginSource: source.ID,
		LoginName:   login,
		IsActive:    true,
		IsAdmin:     isAdmin,
		LDAPUID:     uid,
	}

	ok, err := IsUserExist(0, user.LoginName)
	if err != nil {
		return user, err
	}

	if ok {
		return user, UpdateUser(user)
	}

	return user, CreateUser(user)
}

func remoteUserLogin(user *User, login, password string, source *LoginSource, autoRegister bool) (*User, error) {
	if !source.IsActived {
		return nil, errors.LoginSourceNotActivated{source.ID}
	}

	switch source.Type {
	case LOGIN_LDAP, LOGIN_DLDAP:
		return LoginViaLDAP(user, login, password, source, autoRegister)
	}

	return nil, errors.InvalidLoginSourceType{source.Type}
}

// UserLogin validates user name and password via given login source ID.
// If the loginSourceID is negative, it will abort login process if user is not found.
func UserLogin(username, password string, loginSourceID int64) (*User, error) {
	var user *User
	if strings.Contains(username, "@") {
		user = &User{Email: strings.ToLower(username)}
	} else {
		user = &User{LoginName: strings.ToLower(username)}
	}

	hasUser, err := x.Get(user)
	if err != nil {
		return nil, fmt.Errorf("get user record: %v", err)
	}

	if hasUser {
		// Note: This check is unnecessary but to reduce user confusion at login page
		// and make it more consistent at user's perspective.
		if loginSourceID >= 0 && user.LoginSource != loginSourceID {

			return nil, errors.LoginSourceMismatch{loginSourceID, user.LoginSource}
		}

		// Validate password hash fetched from database for local accounts
		if user.LoginType == LOGIN_NOTYPE ||
			user.LoginType == LOGIN_PLAIN {
			if user.ValidatePassword(password) {
				return user, nil
			}

			return nil, errors.UserNotExist{user.ID, user.LoginName}
		}

		// Remote login to the login source the user is associated with
		source, err := GetLoginSourceByID(user.LoginSource)
		if err != nil {
			return nil, err
		}

		return remoteUserLogin(user, user.LoginName, password, source, false)
	}

	// Non-local login source is always greater than 0
	if loginSourceID <= 0 {
		return nil, errors.UserNotExist{-1, username}
	}

	source, err := GetLoginSourceByID(loginSourceID)
	if err != nil {
		return nil, err
	}

	return remoteUserLogin(nil, username, password, source, true)
}

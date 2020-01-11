// package user contains User business object (BO) and Mysql-DAO implementations.
package user

import (
	"encoding/json"
	"github.com/btnguyen2k/consu/reddo"
	"github.com/btnguyen2k/consu/semita"
	"strings"
)

const (
	fieldUserUsername = "uname"
	fieldUserData     = "data"
	attrUserPassword  = "pwd"
	attrUserAesKey    = "akey"
	attrUserName      = "name"
	attrUserGroupId   = "gid"
)

// User is the business object
type User struct {
	username string `json:"uname"`
	data     string `json:"udata"`
	root     interface{}
	s        *semita.Semita
}

func (u *User) GetUsername() string {
	return u.username
}

func (u *User) SetUsername(username string) *User {
	u.username = strings.TrimSpace(strings.ToLower(username))
	return u
}

func (u *User) GetData() string {
	return u.data
}

func (u *User) SetData(data string) *User {
	u.data = strings.TrimSpace(data)
	var jsData interface{}
	if err := json.Unmarshal([]byte(u.data), &jsData); err == nil {
		u.root = jsData
	} else {
		u.root = make(map[string]interface{})
	}
	u.s = semita.NewSemita(u.root)
	return u
}

func (u *User) setAttr(attr string, value interface{}) *User {
	u.s.SetValue(attr, value)
	data, _ := json.Marshal(u.s.Unwrap())
	return u.SetData(string(data))
}

func (u *User) GetPassword() string {
	if v, e := u.s.GetValueOfType(attrUserPassword, reddo.TypeString); e == nil {
		return v.(string)
	}
	return ""
}

func (u *User) SetPassword(value string) *User {
	return u.setAttr(attrUserPassword, strings.TrimSpace(value))
}

func (u *User) GetName() string {
	if v, e := u.s.GetValueOfType(attrUserName, reddo.TypeString); e == nil {
		return v.(string)
	}
	return ""
}

func (u *User) SetName(value string) *User {
	return u.setAttr(attrUserName, strings.TrimSpace(value))
}

func (u *User) GetAesKey() string {
	if v, e := u.s.GetValueOfType(attrUserAesKey, reddo.TypeString); e == nil {
		return v.(string)
	}
	return ""
}

func (u *User) SetAesKey(value string) *User {
	return u.setAttr(attrUserAesKey, strings.TrimSpace(value))
}

func (u *User) GetGroupId() string {
	if v, e := u.s.GetValueOfType(attrUserGroupId, reddo.TypeString); e == nil {
		return v.(string)
	}
	return ""
}

func (u *User) SetGroupId(value string) *User {
	return u.setAttr(attrUserGroupId, strings.TrimSpace(strings.ToLower(value)))
}

func NewUserBo(username, data string) *User {
	user := &User{}
	return user.SetUsername(username).SetData(data)
}

// UserDao defines API to access User storage
type UserDao interface {
	// Delete removes the specified business object from storage
	Delete(bo *User) (bool, error)

	// Create persists a new business object to storage
	Create(bo *User) (bool, error)

	// Get retrieves a business object from storage
	Get(username string) (*User, error)

	// GetN retrieves N business objects from storage
	GetN(fromOffset, maxNumRows int) ([]*User, error)

	// GetAll retrieves all availables business objects from storage
	GetAll() ([]*User, error)

	// Update modifies an existing business object
	Update(bo *User) (bool, error)
}

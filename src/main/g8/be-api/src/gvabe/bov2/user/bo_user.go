package user

import (
	"encoding/json"
	"strings"

	"github.com/btnguyen2k/consu/reddo"

	"main/src/henge"
)

// NewUser is helper function to create new User bo
//
// available since template-v0.2.0
func NewUser(appVersion uint64, id, maskId string) *User {
	user := &User{
		UniversalBo: *henge.NewUniversalBo(id, appVersion),
	}
	return user.SetMaskId(maskId).sync()
}

// NewUserFromUbo is helper function to create User bo from a universal bo
//
// available since template-v0.2.0
func NewUserFromUbo(ubo *henge.UniversalBo) *User {
	if ubo == nil {
		return nil
	}
	user := User{UniversalBo: *ubo.Clone()}
	if v, err := user.GetExtraAttrAs(UserField_MaskId, reddo.TypeString); err != nil {
		return nil
	} else {
		user.maskId = v.(string)
	}
	if v, err := user.GetDataAttrAs(UserAttr_DisplayName, reddo.TypeString); err != nil {
		return nil
	} else {
		user.displayName = v.(string)
	}
	if v, err := user.GetDataAttrAs(UserAttr_IsAdmin, reddo.TypeBool); err != nil {
		return nil
	} else {
		user.isAdmin = v.(bool)
	}
	if v, err := user.GetDataAttrAs(UserAttr_Password, reddo.TypeString); err != nil {
		return nil
	} else {
		user.password = v.(string)
	}
	return (&user).sync()
}

const (
	// mask-id is also a unique id, used when we do not wish to expose user's id
	// (if we do not wish to use mask-id, simply set id as its value
	UserField_MaskId = "mid"

	// 'password' is used for authentication
	UserAttr_Password = "pwd"

	// 'display-name' is used for displaying purpose
	UserAttr_DisplayName = "dname"

	// 'is-admin' is to flag if user has administrative privilege
	UserAttr_IsAdmin = "isadm"

	userAttr_Ubo = "_ubo"
)

// User is the business object
//	- User inherits unique id from bo.UniversalBo
//
// available since template-v0.2.0
type User struct {
	henge.UniversalBo `json:"_ubo"`
	maskId            string `json:"mid"`
	password          string `json:"pwd"`
	displayName       string `json:"dname"`
	isAdmin           bool   `json:"isadm"`
}

// MarshalJSON implements json.encode.Marshaler.MarshalJSON
//	TODO: lock for read?
func (user *User) MarshalJSON() ([]byte, error) {
	user.sync()
	m := map[string]interface{}{
		userAttr_Ubo: user.UniversalBo.Clone(),
		"_cols": map[string]interface{}{
			UserField_MaskId: user.maskId,
		},
		"_attrs": map[string]interface{}{
			UserAttr_DisplayName: user.displayName,
			UserAttr_IsAdmin:     user.isAdmin,
			UserAttr_Password:    user.password,
		},
	}
	return json.Marshal(m)
}

// UnmarshalJSON implements json.decode.Unmarshaler.UnmarshalJSON
//	TODO: lock for write?
func (user *User) UnmarshalJSON(data []byte) error {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	var err error
	if m[userAttr_Ubo] != nil {
		js, _ := json.Marshal(m[userAttr_Ubo])
		if err = json.Unmarshal(js, &user.UniversalBo); err != nil {
			return err
		}
	}
	if _cols, ok := m["_cols"].(map[string]interface{}); ok {
		if user.maskId, err = reddo.ToString(_cols[UserField_MaskId]); err != nil {
			return err
		}
	}
	if _attrs, ok := m["_attrs"].(map[string]interface{}); ok {
		if user.displayName, err = reddo.ToString(_attrs[UserAttr_DisplayName]); err != nil {
			return err
		}
		if user.isAdmin, err = reddo.ToBool(_attrs[UserAttr_IsAdmin]); err != nil {
			return err
		}
		if user.password, err = reddo.ToString(_attrs[UserAttr_Password]); err != nil {
			return err
		}
	}
	user.sync()
	return nil
}

// GetMaskUniqueId returns value of user's 'mask-id' attribute
func (user *User) GetMaskId() string {
	return user.maskId
}

// SetMaskId sets value of user's 'mask-uid' attribute
func (user *User) SetMaskId(v string) *User {
	user.maskId = strings.TrimSpace(strings.ToLower(v))
	return user
}

// GetPassword returns value of user's 'password' attribute
func (user *User) GetPassword() string {
	return user.password
}

// SetPassword sets value of user's 'password' attribute
func (user *User) SetPassword(v string) *User {
	user.password = strings.TrimSpace(v)
	return user
}

// GetDisplayName returns value of user's 'display-name' attribute
func (user *User) GetDisplayName() string {
	return user.displayName
}

// SetDisplayName sets value of user's 'display-name' attribute
func (user *User) SetDisplayName(v string) *User {
	user.displayName = strings.TrimSpace(v)
	return user
}

// IsAdmin returns value of user's 'is-admin' attribute
func (user *User) IsAdmin() bool {
	return user.isAdmin
}

// SetAdmin sets value of user's 'is-admin' attribute
func (user *User) SetAdmin(v bool) *User {
	user.isAdmin = v
	return user
}

// sync is called to synchronize BO's attributes to its UniversalBo
func (user *User) sync() *User {
	user.SetDataAttr(UserAttr_Password, user.password)
	user.SetDataAttr(UserAttr_DisplayName, user.displayName)
	user.SetDataAttr(UserAttr_IsAdmin, user.isAdmin)
	user.SetExtraAttr(UserField_MaskId, user.maskId)
	user.UniversalBo.Sync()
	return user
}

// UserDao defines API to access User storage
//
// available since template-v0.2.0
type UserDao interface {
	// Delete removes the specified business object from storage
	Delete(bo *User) (bool, error)

	// Create persists a new business object to storage
	Create(bo *User) (bool, error)

	// Get retrieves a business object from storage
	Get(username string) (*User, error)

	// GetN retrieves N business objects from storage
	GetN(fromOffset, maxNumRows int) ([]*User, error)

	// GetAll retrieves all available business objects from storage
	GetAll() ([]*User, error)

	// Update modifies an existing business object
	Update(bo *User) (bool, error)
}

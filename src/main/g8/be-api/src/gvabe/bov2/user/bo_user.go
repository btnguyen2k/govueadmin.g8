package user

import (
	"encoding/json"
	"strings"

	"github.com/btnguyen2k/consu/reddo"
	"github.com/btnguyen2k/henge"
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

func (u *User) ToMap(postFunc henge.FuncPostUboToMap) map[string]interface{} {
	result := map[string]interface{}{
		henge.FieldId:        u.GetId(),
		UserField_MaskId:     u.maskId,
		UserAttr_IsAdmin:     u.isAdmin,
		UserAttr_DisplayName: u.displayName,
	}
	if postFunc != nil {
		result = postFunc(result)
	}
	return result
}

// MarshalJSON implements json.encode.Marshaler.MarshalJSON
//	TODO: lock for read?
func (u *User) MarshalJSON() ([]byte, error) {
	u.sync()
	m := map[string]interface{}{
		userAttr_Ubo: u.UniversalBo.Clone(),
		"_cols": map[string]interface{}{
			UserField_MaskId: u.maskId,
		},
		"_attrs": map[string]interface{}{
			UserAttr_DisplayName: u.displayName,
			UserAttr_IsAdmin:     u.isAdmin,
			UserAttr_Password:    u.password,
		},
	}
	return json.Marshal(m)
}

// UnmarshalJSON implements json.decode.Unmarshaler.UnmarshalJSON
//	TODO: lock for write?
func (u *User) UnmarshalJSON(data []byte) error {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	var err error
	if m[userAttr_Ubo] != nil {
		js, _ := json.Marshal(m[userAttr_Ubo])
		if err = json.Unmarshal(js, &u.UniversalBo); err != nil {
			return err
		}
	}
	if _cols, ok := m["_cols"].(map[string]interface{}); ok {
		if u.maskId, err = reddo.ToString(_cols[UserField_MaskId]); err != nil {
			return err
		}
	}
	if _attrs, ok := m["_attrs"].(map[string]interface{}); ok {
		if u.displayName, err = reddo.ToString(_attrs[UserAttr_DisplayName]); err != nil {
			return err
		}
		if u.isAdmin, err = reddo.ToBool(_attrs[UserAttr_IsAdmin]); err != nil {
			return err
		}
		if u.password, err = reddo.ToString(_attrs[UserAttr_Password]); err != nil {
			return err
		}
	}
	u.sync()
	return nil
}

// GetMaskUniqueId returns value of user's 'mask-id' attribute
func (u *User) GetMaskId() string {
	return u.maskId
}

// SetMaskId sets value of user's 'mask-uid' attribute
func (u *User) SetMaskId(v string) *User {
	u.maskId = strings.TrimSpace(strings.ToLower(v))
	return u
}

// GetPassword returns value of user's 'password' attribute
func (u *User) GetPassword() string {
	return u.password
}

// SetPassword sets value of user's 'password' attribute
func (u *User) SetPassword(v string) *User {
	u.password = strings.TrimSpace(v)
	return u
}

// GetDisplayName returns value of user's 'display-name' attribute
func (u *User) GetDisplayName() string {
	return u.displayName
}

// SetDisplayName sets value of user's 'display-name' attribute
func (u *User) SetDisplayName(v string) *User {
	u.displayName = strings.TrimSpace(v)
	return u
}

// IsAdmin returns value of user's 'is-admin' attribute
func (u *User) IsAdmin() bool {
	return u.isAdmin
}

// SetAdmin sets value of user's 'is-admin' attribute
func (u *User) SetAdmin(v bool) *User {
	u.isAdmin = v
	return u
}

// sync is called to synchronize BO's attributes to its UniversalBo
func (u *User) sync() *User {
	u.SetDataAttr(UserAttr_Password, u.password)
	u.SetDataAttr(UserAttr_DisplayName, u.displayName)
	u.SetDataAttr(UserAttr_IsAdmin, u.isAdmin)
	u.SetExtraAttr(UserField_MaskId, u.maskId)
	u.UniversalBo.Sync()
	return u
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

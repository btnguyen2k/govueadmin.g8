package user

import (
	"encoding/json"
	"strings"

	"github.com/btnguyen2k/consu/reddo"
	"github.com/btnguyen2k/henge"
)

// NewUser is helper function to create new User bo
//
// Available since template-v0.2.0
func NewUser(appVersion uint64, id, maskId string) *User {
	user := &User{
		UniversalBo: henge.NewUniversalBo(id, appVersion),
	}
	return user.SetMaskId(maskId).sync()
}

// NewUserFromUbo is helper function to create User bo from a universal bo
//
// Available since template-v0.2.0
func NewUserFromUbo(ubo *henge.UniversalBo) *User {
	if ubo == nil {
		return nil
	}
	ubo = ubo.Clone()
	user := &User{UniversalBo: ubo}
	if v, err := ubo.GetExtraAttrAs(UserFieldMaskId, reddo.TypeString); err != nil {
		return nil
	} else {
		user.maskId, _ = v.(string)
	}
	if v, err := ubo.GetDataAttrAs(UserAttrDisplayName, reddo.TypeString); err != nil {
		return nil
	} else {
		user.displayName, _ = v.(string)
	}
	if v, err := ubo.GetDataAttrAs(UserAttrIsAdmin, reddo.TypeBool); err != nil {
		return nil
	} else {
		user.isAdmin, _ = v.(bool)
	}
	if v, err := ubo.GetDataAttrAs(UserAttrPassword, reddo.TypeString); err != nil {
		return nil
	} else {
		user.password, _ = v.(string)
	}
	return user.sync()
}

const (
	// UserFieldMaskId is a also unique id for BO user, used when we do not wish to expose user's id
	// (if we do not wish to use mask-id, simply set id as its value.
	UserFieldMaskId = "mid"

	// UserAttrPassword is user's password, used for authentication.
	UserAttrPassword = "pwd"

	// UserAttrDisplayName is used for displaying purpose.
	UserAttrDisplayName = "dname"

	// UserAttrIsAdmin is a flag to mark if user has administrative privilege.
	UserAttrIsAdmin = "isadm"

	// userAttr_Ubo is for internal use only!
	userAttr_Ubo = "_ubo"
)

// User is the business object
//	- User inherits unique id from bo.UniversalBo
//
// available since template-v0.2.0
type User struct {
	*henge.UniversalBo `json:"_ubo"`
	maskId             string `json:"mid"`
	password           string `json:"pwd"`
	displayName        string `json:"dname"`
	isAdmin            bool   `json:"isadm"`
}

// ToMap transforms user's attributes to a map.
func (u *User) ToMap(postFunc henge.FuncPostUboToMap) map[string]interface{} {
	result := map[string]interface{}{
		henge.FieldId:       u.GetId(),
		UserFieldMaskId:     u.maskId,
		UserAttrIsAdmin:     u.isAdmin,
		UserAttrDisplayName: u.displayName,
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
			UserFieldMaskId: u.maskId,
		},
		"_attrs": map[string]interface{}{
			UserAttrDisplayName: u.displayName,
			UserAttrIsAdmin:     u.isAdmin,
			UserAttrPassword:    u.password,
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
		if u.maskId, err = reddo.ToString(_cols[UserFieldMaskId]); err != nil {
			return err
		}
	}
	if _attrs, ok := m["_attrs"].(map[string]interface{}); ok {
		if u.displayName, err = reddo.ToString(_attrs[UserAttrDisplayName]); err != nil {
			return err
		}
		if u.isAdmin, err = reddo.ToBool(_attrs[UserAttrIsAdmin]); err != nil {
			return err
		}
		if u.password, err = reddo.ToString(_attrs[UserAttrPassword]); err != nil {
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
	u.SetDataAttr(UserAttrPassword, u.password)
	u.SetDataAttr(UserAttrDisplayName, u.displayName)
	u.SetDataAttr(UserAttrIsAdmin, u.isAdmin)
	u.SetExtraAttr(UserFieldMaskId, u.maskId)
	u.UniversalBo.Sync()
	return u
}

package user

import (
	"encoding/json"
	"fmt"
	"log"
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

// NewUserFromUbo is helper function to create User App bo from a universal bo
//
// available since template-v0.2.0
func NewUserFromUbo(ubo *henge.UniversalBo) *User {
	if ubo == nil {
		return nil
	}
	user := User{}
	if err := json.Unmarshal([]byte(ubo.GetDataJson()), &user); err != nil {
		log.Print(fmt.Sprintf("[WARN] NewUserFromUbo - error unmarshalling JSON data: %e", err))
		log.Print(err)
		return nil
	}
	user.UniversalBo = *ubo.Clone()
	if maskId, err := user.GetExtraAttrAs(UserField_MaskId, reddo.TypeString); err == nil {
		user.SetMaskId(maskId.(string))
	}
	return &user
}

const (
	// mask-id is also a unique id, used when we do not wish to expose user's id
	// (if we do not wish to use mask-id, simply set id as its value
	UserField_MaskId = "mid"

	// 'display-name' is used for displaying purpose
	UserAttr_DisplayName = "dname"

	userAttr_Ubo = "_ubo"
)

// User is the business object
//	- User inherits unique id from bo.UniversalBo
//
// available since template-v0.2.0
type User struct {
	henge.UniversalBo `json:"_ubo"`
	maskId            string `json:"mid"`
	displayName       string `json:""`
}

// MarshalJSON implements json.encode.Marshaler.MarshalJSON
//	TODO: lock for read?
func (user *User) MarshalJSON() ([]byte, error) {
	user.sync()
	m := map[string]interface{}{
		userAttr_Ubo:         user.UniversalBo.Clone(),
		UserField_MaskId:     user.maskId,
		UserAttr_DisplayName: user.displayName,
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
		if err := json.Unmarshal(js, &user.UniversalBo); err != nil {
			return err
		}
	}
	if user.displayName, err = reddo.ToString(m[UserAttr_DisplayName]); err != nil {
		return err
	}
	if user.maskId, err = reddo.ToString(m[UserField_MaskId]); err != nil {
		return err
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

// GetDisplayName returns value of user's 'display-name' attribute
func (user *User) GetDisplayName() string {
	return user.displayName
}

// SetDisplayName sets value of user's 'display-name' attribute
func (user *User) SetDisplayName(v string) *User {
	user.displayName = strings.TrimSpace(v)
	return user
}

func (user *User) sync() *User {
	user.SetDataAttr(UserAttr_DisplayName, user.displayName)
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

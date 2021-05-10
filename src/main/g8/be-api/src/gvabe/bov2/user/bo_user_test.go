package user

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/btnguyen2k/henge"
)

func TestNewUser(t *testing.T) {
	name := "TestNewUser"
	_tagVersion := uint64(1337)
	_id := "admin@local"
	_maskId := "admin"
	_pwd := "mypassword"
	_displayName := "Administrator"
	_isAdmin := true
	user := NewUser(_tagVersion, _id, _maskId)
	if user == nil {
		t.Fatalf("%s failed: nil", name)
	}
	user.SetPassword(_pwd).SetDisplayName(_displayName).SetAdmin(_isAdmin)
	if tagVersion := user.GetTagVersion(); tagVersion != _tagVersion {
		t.Fatalf("%s failed: expected tag-version to be %#v but received %#v", name, _tagVersion, tagVersion)
	}
	if id := user.GetId(); id != _id {
		t.Fatalf("%s failed: expected bo's id to be %#v but received %#v", name, _id, id)
	}
	if maskId := user.GetMaskId(); maskId != _maskId {
		t.Fatalf("%s failed: expected bo's mask-id to be %#v but received %#v", name, _maskId, maskId)
	}
	if password := user.GetPassword(); password != _pwd {
		t.Fatalf("%s failed: expected bo's password to be %#v but received %#v", name, _pwd, password)
	}
	if displayName := user.GetDisplayName(); displayName != _displayName {
		t.Fatalf("%s failed: expected bo's display-name to be %#v but received %#v", name, _displayName, displayName)
	}
	if isAdmin := user.IsAdmin(); isAdmin != _isAdmin {
		t.Fatalf("%s failed: expected bo's display-name to be %#v but received %#v", name, _isAdmin, isAdmin)
	}
}

func TestNewUserFromUbo(t *testing.T) {
	name := "TestNewUserFromUbo"

	if NewUserFromUbo(nil) != nil {
		t.Fatalf("%s failed: NewUserFromUbo(nil) should return nil", name)
	}
	_tagVersion := uint64(1337)
	_id := "admin@local"
	_maskId := "admin"
	_pwd := "mypassword"
	_displayName := "Administrator"
	_isAdmin := true
	ubo := henge.NewUniversalBo(_id, _tagVersion)
	ubo.SetExtraAttr(UserFieldMaskId, _maskId)
	ubo.SetDataAttr(UserAttrDisplayName, _displayName)
	ubo.SetDataAttr(UserAttrPassword, _pwd)
	ubo.SetDataAttr(UserAttrIsAdmin, _isAdmin)

	user := NewUserFromUbo(ubo)
	if tagVersion := user.GetTagVersion(); tagVersion != _tagVersion {
		t.Fatalf("%s failed: expected tag-version to be %#v but received %#v", name, _tagVersion, tagVersion)
	}
	if id := user.GetId(); id != _id {
		t.Fatalf("%s failed: expected bo's id to be %#v but received %#v", name, _id, id)
	}
	if maskId := user.GetMaskId(); maskId != _maskId {
		t.Fatalf("%s failed: expected bo's mask-id to be %#v but received %#v", name, _maskId, maskId)
	}
	if password := user.GetPassword(); password != _pwd {
		t.Fatalf("%s failed: expected bo's password to be %#v but received %#v", name, _pwd, password)
	}
	if displayName := user.GetDisplayName(); displayName != _displayName {
		t.Fatalf("%s failed: expected bo's display-name to be %#v but received %#v", name, _displayName, displayName)
	}
	if isAdmin := user.IsAdmin(); isAdmin != _isAdmin {
		t.Fatalf("%s failed: expected bo's display-name to be %#v but received %#v", name, _isAdmin, isAdmin)
	}
}

func TestUser_ToMap(t *testing.T) {
	name := "TestUser_ToMap"
	_tagVersion := uint64(1337)
	_id := "admin@local"
	_maskId := "admin"
	_pwd := "mypassword"
	_displayName := "Administrator"
	_isAdmin := true
	user := NewUser(_tagVersion, _id, _maskId)
	if user == nil {
		t.Fatalf("%s failed: nil", name)
	}
	user.SetPassword(_pwd).SetDisplayName(_displayName).SetAdmin(_isAdmin)

	m := user.ToMap(nil)
	expected := map[string]interface{}{
		henge.FieldId:       _id,
		UserFieldMaskId:     _maskId,
		UserAttrIsAdmin:     _isAdmin,
		UserAttrDisplayName: _displayName,
	}
	if !reflect.DeepEqual(m, expected) {
		t.Fatalf("%s failed: expected %#v but received %#v", name, expected, m)
	}

	m = user.ToMap(func(input map[string]interface{}) map[string]interface{} {
		return map[string]interface{}{
			"FieldId":              input[henge.FieldId],
			"UserFieldMaskId":     input[UserFieldMaskId],
			"UserAttrIsAdmin":     input[UserAttrIsAdmin],
			"UserAttrDisplayName": input[UserAttrDisplayName],
		}
	})
	expected = map[string]interface{}{
		"FieldId":              _id,
		"UserFieldMaskId":     _maskId,
		"UserAttrIsAdmin":     _isAdmin,
		"UserAttrDisplayName": _displayName,
	}
	if !reflect.DeepEqual(m, expected) {
		t.Fatalf("%s failed: expected %#v but received %#v", name, expected, m)
	}
}

func TestUser_json(t *testing.T) {
	name := "TestUser_json"
	_tagVersion := uint64(1337)
	_id := "admin@local"
	_maskId := "admin"
	_pwd := "mypassword"
	_displayName := "Administrator"
	_isAdmin := true
	user1 := NewUser(_tagVersion, _id, _maskId)
	if user1 == nil {
		t.Fatalf("%s failed: nil", name)
	}
	user1.SetPassword(_pwd).SetDisplayName(_displayName).SetAdmin(_isAdmin)
	js1, _ := json.Marshal(user1)

	var user2 *User
	err := json.Unmarshal(js1, &user2)
	if err != nil {
		t.Fatalf("%s failed: %e", name, err)
	}
	if user1.GetTagVersion() != user2.GetTagVersion() {
		t.Fatalf("%s failed: expected %#v but received %#v", name, user1.GetTagVersion(), user2.GetTagVersion())
	}
	if user1.GetId() != user2.GetId() {
		t.Fatalf("%s failed: expected %#v but received %#v", name, user1.GetId(), user2.GetId())
	}
	if user1.GetMaskId() != user2.GetMaskId() {
		t.Fatalf("%s failed: expected %#v but received %#v", name, user1.GetMaskId(), user2.GetMaskId())
	}
	if user1.GetPassword() != user2.GetPassword() {
		t.Fatalf("%s failed: expected %#v but received %#v", name, user1.GetPassword(), user2.GetPassword())
	}
	if user1.GetDisplayName() != user2.GetDisplayName() {
		t.Fatalf("%s failed: expected %#v but received %#v", name, user1.GetDisplayName(), user2.GetDisplayName())
	}
	if user1.IsAdmin() != user2.IsAdmin() {
		t.Fatalf("%s failed: expected %#v but received %#v", name, user1.IsAdmin(), user2.IsAdmin())
	}
	if user1.GetChecksum() != user2.GetChecksum() {
		t.Fatalf("%s failed: expected %#v but received %#v", name, user1.GetChecksum(), user2.GetChecksum())
	}
}

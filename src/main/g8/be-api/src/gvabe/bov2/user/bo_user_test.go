package user

import (
	"encoding/json"
	"testing"
)

func TestNewUser(t *testing.T) {
	name := "TestNewUser"
	_appVersion := uint64(1337)
	_id := "admin@local"
	_maskId := "admin"
	_displayName := "Administrator"
	user := NewUser(_appVersion, _id, _maskId)
	if user == nil {
		t.Fatalf("%s failed: nil", name)
	}
	user.SetDisplayName(_displayName)
	if appVersion := user.GetAppVersion(); appVersion != _appVersion {
		t.Fatalf("%s failed: expected app-version to be %#v but received %#v", name, _appVersion, appVersion)
	}
	if id := user.GetId(); id != _id {
		t.Fatalf("%s failed: expected bo's id to be %#v but received %#v", name, _id, id)
	}
	if maskId := user.GetMaskId(); maskId != _maskId {
		t.Fatalf("%s failed: expected bo's mask-id to be %#v but received %#v", name, _maskId, maskId)
	}
	if displayName := user.GetDisplayName(); displayName != _displayName {
		t.Fatalf("%s failed: expected bo's display-name to be %#v but received %#v", name, _displayName, displayName)
	}
}

func TestUser_json(t *testing.T) {
	name := "TestUser_json"
	_appVersion := uint64(1337)
	_id := "admin@local"
	_maskId := "admin"
	_displayName := "Administrator"
	user1 := NewUser(_appVersion, _id, _maskId)
	if user1 == nil {
		t.Fatalf("%s failed: nil", name)
	}
	user1.SetDisplayName(_displayName)
	js1, _ := json.Marshal(user1)

	var user2 *User
	err := json.Unmarshal(js1, &user2)
	if err != nil {
  		t.Fatalf("%s failed: %e", name, err)
	}
	if user1.GetAppVersion() != user2.GetAppVersion() {
		t.Fatalf("%s failed: expected %#v but received %#v", name, user1.GetAppVersion(), user2.GetAppVersion())
	}
	if user1.GetId() != user2.GetId() {
		t.Fatalf("%s failed: expected %#v but received %#v", name, user1.GetId(), user2.GetId())
	}
	if user1.GetMaskId() != user2.GetMaskId() {
		t.Fatalf("%s failed: expected %#v but received %#v", name, user1.GetMaskId(), user2.GetMaskId())
	}
	if user1.GetDisplayName() != user2.GetDisplayName() {
		t.Fatalf("%s failed: expected %#v but received %#v", name, user1.GetDisplayName(), user2.GetDisplayName())
	}
	if user1.GetChecksum() != user2.GetChecksum() {
		t.Fatalf("%s failed: expected %#v but received %#v", name, user1.GetChecksum(), user2.GetChecksum())
	}
}

package blog

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/btnguyen2k/henge"
	"main/src/gvabe/bov2/user"
	"main/src/utils"
)

func TestNewBlogVote(t *testing.T) {
	name := "TestNewBlogVote"
	_tagVersion := uint64(1337)
	_userId := "admin@local"
	_userMaskId := "admin"
	_user := user.NewUser(_tagVersion, _userId, _userMaskId)

	_targetId := utils.UniqueId()
	_value := 123
	vote := NewBlogVote(_tagVersion, _user, _targetId, _value)
	if vote == nil {
		t.Fatalf("%s failed: nil", name)
	}
	if tagVersion := vote.GetTagVersion(); tagVersion != _tagVersion {
		t.Fatalf("%s failed: expected tag-version to be %#v but received %#v", name, _tagVersion, tagVersion)
	}
	if ownerId := vote.GetOwnerId(); ownerId != _userId {
		t.Fatalf("%s failed: expected owner-id to be %#v but received %#v", name, _userId, ownerId)
	}
	if targetId := vote.GetTargetId(); targetId != _targetId {
		t.Fatalf("%s failed: expected target-id to be %#v but received %#v", name, _targetId, targetId)
	}
	if value := vote.GetValue(); value != _value {
		t.Fatalf("%s failed: expected value to be %#v but received %#v", name, _value, value)
	}
}

func TestNewBlogVoteFromUbo(t *testing.T) {
	name := "TestNewBlogVoteFromUbo"

	if NewBlogVoteFromUbo(nil) != nil {
		t.Fatalf("%s failed: NewBlogVoteFromUbo(nil) should return nil", name)
	}
	_tagVersion := uint64(1337)
	_id := utils.UniqueId()
	_userId := utils.UniqueId()
	_targetId := utils.UniqueId()
	_value := 123
	ubo := henge.NewUniversalBo(_id, _tagVersion)
	ubo.SetExtraAttr(VoteFieldOwnerId, _userId)
	ubo.SetExtraAttr(VoteFieldValue, _value)
	ubo.SetExtraAttr(VoteFieldTargetId, _targetId)

	vote := NewBlogVoteFromUbo(ubo)
	if vote == nil {
		t.Fatalf("%s failed: nil", name)
	}
	if tagVersion := vote.GetTagVersion(); tagVersion != _tagVersion {
		t.Fatalf("%s failed: expected tag-version to be %#v but received %#v", name, _tagVersion, tagVersion)
	}
	if id := vote.GetId(); id != _id {
		t.Fatalf("%s failed: expected bo's id to be %#v but received %#v", name, _id, id)
	}
	if ownerId := vote.GetOwnerId(); ownerId != _userId {
		t.Fatalf("%s failed: expected owner-id to be %#v but received %#v", name, _userId, ownerId)
	}
	if targetId := vote.GetTargetId(); targetId != _targetId {
		t.Fatalf("%s failed: expected target-id to be %#v but received %#v", name, _targetId, targetId)
	}
	if value := vote.GetValue(); value != _value {
		t.Fatalf("%s failed: expected value to be %#v but received %#v", name, _value, value)
	}
}

func TestBlogVote_ToMap(t *testing.T) {
	name := "TestBlogVote_ToMap"
	_tagVersion := uint64(1337)
	_userId := "admin@local"
	_userMaskId := "admin"
	_user := user.NewUser(_tagVersion, _userId, _userMaskId)

	_targetId := utils.UniqueId()
	_value := 123
	vote := NewBlogVote(_tagVersion, _user, _targetId, _value)
	if vote == nil {
		t.Fatalf("%s failed: nil", name)
	}

	m := vote.ToMap(nil)
	expected := map[string]interface{}{
		henge.FieldId:          vote.GetId(),
		henge.FieldTimeCreated: vote.GetTimeCreated(),
		VoteFieldTargetId:      _targetId,
		VoteFieldOwnerId:       _userId,
		VoteFieldValue:         _value,
	}
	if !reflect.DeepEqual(m, expected) {
		t.Fatalf("%s failed: expected %#v but received %#v", name, expected, m)
	}

	m = vote.ToMap(func(input map[string]interface{}) map[string]interface{} {
		return map[string]interface{}{
			"FieldId":           input[henge.FieldId],
			"FieldTimeCreated":  input[henge.FieldTimeCreated],
			"VoteFieldTargetId": input[VoteFieldTargetId],
			"VoteFieldOwnerId":  input[VoteFieldOwnerId],
			"VoteFieldValue":    input[VoteFieldValue],
		}
	})
	expected = map[string]interface{}{
		"FieldId":           vote.GetId(),
		"FieldTimeCreated":  vote.GetTimeCreated(),
		"VoteFieldTargetId": _targetId,
		"VoteFieldOwnerId":  _userId,
		"VoteFieldValue":    _value,
	}
	if !reflect.DeepEqual(m, expected) {
		t.Fatalf("%s failed: expected %#v but received %#v", name, expected, m)
	}
}

func TestBlogVote_json(t *testing.T) {
	name := "TestBlogVote_json"
	_tagVersion := uint64(1337)
	_userId := "admin@local"
	_userMaskId := "admin"
	_user := user.NewUser(_tagVersion, _userId, _userMaskId)

	_targetId := utils.UniqueId()
	_value := 123
	vote := NewBlogVote(_tagVersion, _user, _targetId, _value)
	if vote == nil {
		t.Fatalf("%s failed: nil", name)
	}
	js1, _ := json.Marshal(vote)

	var vote2 *BlogVote
	err := json.Unmarshal(js1, &vote2)
	if err != nil {
		t.Fatalf("%s failed: %e", name, err)
	}
	if tagVersion := vote2.GetTagVersion(); tagVersion != _tagVersion {
		t.Fatalf("%s failed: expected tag-version to be %#v but received %#v", name, _tagVersion, tagVersion)
	}
	if ownerId := vote2.GetOwnerId(); ownerId != _userId {
		t.Fatalf("%s failed: expected owner-id to be %#v but received %#v", name, _userId, ownerId)
	}
	if targetId := vote2.GetTargetId(); targetId != _targetId {
		t.Fatalf("%s failed: expected target-id to be %#v but received %#v", name, _targetId, targetId)
	}
	if value := vote2.GetValue(); value != _value {
		t.Fatalf("%s failed: expected value to be %#v but received %#v", name, _value, value)
	}
	if t2, t1 := vote2.GetTimeCreated(), vote.GetTimeCreated(); !t2.Equal(t1) {
		t.Fatalf("%s failed: expected %#v but received %#v", name, t1.Format(time.RFC3339), t2.Format(time.RFC3339))
	}
	if vote2.GetChecksum() != vote.GetChecksum() {
		t.Fatalf("%s failed: expected %#v but received %#v", name, vote2.GetChecksum(), vote.GetChecksum())
	}
}

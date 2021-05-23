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

func TestNewBlogPost(t *testing.T) {
	name := "TestNewBlogPost"
	_tagVersion := uint64(1337)
	_userId := "admin@local"
	_userMaskId := "admin"
	_user := user.NewUser(_tagVersion, _userId, _userMaskId)

	_postIsPublic := true
	_postTitle := "Blog post title"
	_postContent := "Blog post content"
	post := NewBlogPost(_tagVersion, _user, _postIsPublic, _postTitle, _postContent)
	if post == nil {
		t.Fatalf("%s failed: nil", name)
	}
	_postNumComments := 12
	_postVotesUp := 34
	_postVotesDown := 56
	post.SetNumComments(_postNumComments).SetNumVotesUp(_postVotesUp).SetNumVotesDown(_postVotesDown)
	if tagVersion := post.GetTagVersion(); tagVersion != _tagVersion {
		t.Fatalf("%s failed: expected tag-version to be %#v but received %#v", name, _tagVersion, tagVersion)
	}
	if ownerId := post.GetOwnerId(); ownerId != _userId {
		t.Fatalf("%s failed: expected owner-id to be %#v but received %#v", name, _userId, ownerId)
	}
	if title := post.GetTitle(); title != _postTitle {
		t.Fatalf("%s failed: expected title to be %#v but received %#v", name, _postTitle, title)
	}
	if content := post.GetContent(); content != _postContent {
		t.Fatalf("%s failed: expected content to be %#v but received %#v", name, _postContent, content)
	}
	if numComments := post.GetNumComments(); numComments != _postNumComments {
		t.Fatalf("%s failed: expected num-comments to be %#v but received %#v", name, _postNumComments, numComments)
	}
	if numVotesUp := post.GetNumVotesUp(); numVotesUp != _postVotesUp {
		t.Fatalf("%s failed: expected num-votes-up to be %#v but received %#v", name, _postVotesUp, numVotesUp)
	}
	if numVotesDown := post.GetNumVotesDown(); numVotesDown != _postVotesDown {
		t.Fatalf("%s failed: expected num-votes-down to be %#v but received %#v", name, _postVotesDown, numVotesDown)
	}
}

func TestNewBlogPostFromUbo(t *testing.T) {
	name := "TestNewBlogPostFromUbo"

	if NewBlogPostFromUbo(nil) != nil {
		t.Fatalf("%s failed: NewBlogPostFromUbo(nil) should return nil", name)
	}
	_tagVersion := uint64(1337)
	_id := utils.UniqueId()
	_isPublic := true
	_userId := utils.UniqueId()
	_content := "Blog content"
	_title := "Blog title"
	_numComments := 12
	_numVotesUp := 34
	_numVotesDown := 56
	ubo := henge.NewUniversalBo(_id, _tagVersion)
	ubo.SetExtraAttr(PostFieldIsPublic, _isPublic)
	ubo.SetExtraAttr(PostFieldOwnerId, _userId)
	ubo.SetDataAttr(PostAttrContent, _content)
	ubo.SetDataAttr(PostAttrTitle, _title)
	ubo.SetDataAttr(PostAttrNumComments, _numComments)
	ubo.SetDataAttr(PostAttrNumVotesUp, _numVotesUp)
	ubo.SetDataAttr(PostAttrNumVotesDown, _numVotesDown)

	post := NewBlogPostFromUbo(ubo)
	if post == nil {
		t.Fatalf("%s failed: nil", name)
	}
	if tagVersion := post.GetTagVersion(); tagVersion != _tagVersion {
		t.Fatalf("%s failed: expected tag-version to be %#v but received %#v", name, _tagVersion, tagVersion)
	}
	if id := post.GetId(); id != _id {
		t.Fatalf("%s failed: expected bo's id to be %#v but received %#v", name, _id, id)
	}
	if ownerId := post.GetOwnerId(); ownerId != _userId {
		t.Fatalf("%s failed: expected owner-id to be %#v but received %#v", name, _userId, ownerId)
	}
	if title := post.GetTitle(); title != _title {
		t.Fatalf("%s failed: expected title to be %#v but received %#v", name, _title, title)
	}
	if content := post.GetContent(); content != _content {
		t.Fatalf("%s failed: expected content to be %#v but received %#v", name, _content, content)
	}
	if numComments := post.GetNumComments(); numComments != _numComments {
		t.Fatalf("%s failed: expected num-comments to be %#v but received %#v", name, _numComments, numComments)
	}
	if numVotesUp := post.GetNumVotesUp(); numVotesUp != _numVotesUp {
		t.Fatalf("%s failed: expected num-votes-up to be %#v but received %#v", name, _numVotesUp, numVotesUp)
	}
	if numVotesDown := post.GetNumVotesDown(); numVotesDown != _numVotesDown {
		t.Fatalf("%s failed: expected num-votes-down to be %#v but received %#v", name, _numVotesDown, numVotesDown)
	}
}

func TestBlogPost_ToMap(t *testing.T) {
	name := "TestBlogPost_ToMap"
	_tagVersion := uint64(1337)
	_userId := "admin@local"
	_userMaskId := "admin"
	_user := user.NewUser(_tagVersion, _userId, _userMaskId)

	_postIsPublic := true
	_postTitle := "Blog post title"
	_postContent := "Blog post content"
	post := NewBlogPost(_tagVersion, _user, _postIsPublic, _postTitle, _postContent)
	if post == nil {
		t.Fatalf("%s failed: nil", name)
	}
	_postNumComments := 12
	_postVotesUp := 34
	_postVotesDown := 56
	post.SetNumComments(_postNumComments).SetNumVotesUp(_postVotesUp).SetNumVotesDown(_postVotesDown)

	m := post.ToMap(nil)
	expected := map[string]interface{}{
		henge.FieldId:          post.GetId(),
		henge.FieldTimeCreated: post.GetTimeCreated(),
		PostFieldOwnerId:       _userId,
		PostFieldIsPublic:      _postIsPublic,
		PostAttrTitle:          _postTitle,
		PostAttrContent:        _postContent,
		PostAttrNumComments:    _postNumComments,
		PostAttrNumVotesUp:     _postVotesUp,
		PostAttrNumVotesDown:   _postVotesDown,
	}
	if !reflect.DeepEqual(m, expected) {
		t.Fatalf("%s failed: expected %#v but received %#v", name, expected, m)
	}

	m = post.ToMap(func(input map[string]interface{}) map[string]interface{} {
		return map[string]interface{}{
			"FieldId":              input[henge.FieldId],
			"FieldTimeCreated":     input[henge.FieldTimeCreated],
			"PostFieldOwnerId":     input[PostFieldOwnerId],
			"PostFieldIsPublic":    input[PostFieldIsPublic],
			"PostAttrTitle":        input[PostAttrTitle],
			"PostAttrContent":      input[PostAttrContent],
			"PostAttrNumComments":  input[PostAttrNumComments],
			"PostAttrNumVotesUp":   input[PostAttrNumVotesUp],
			"PostAttrNumVotesDown": input[PostAttrNumVotesDown],
		}
	})
	expected = map[string]interface{}{
		"FieldId":              post.GetId(),
		"FieldTimeCreated":     post.GetTimeCreated(),
		"PostFieldOwnerId":     _userId,
		"PostFieldIsPublic":    _postIsPublic,
		"PostAttrTitle":        _postTitle,
		"PostAttrContent":      _postContent,
		"PostAttrNumComments":  _postNumComments,
		"PostAttrNumVotesUp":   _postVotesUp,
		"PostAttrNumVotesDown": _postVotesDown,
	}
	if !reflect.DeepEqual(m, expected) {
		t.Fatalf("%s failed: expected %#v but received %#v", name, expected, m)
	}
}

func TestBlogPost_json(t *testing.T) {
	name := "TestBlogPost_json"
	_tagVersion := uint64(1337)
	_userId := "admin@local"
	_userMaskId := "admin"
	_user := user.NewUser(_tagVersion, _userId, _userMaskId)

	_postIsPublic := true
	_postTitle := "Blog post title"
	_postContent := "Blog post content"
	post := NewBlogPost(_tagVersion, _user, _postIsPublic, _postTitle, _postContent)
	if post == nil {
		t.Fatalf("%s failed: nil", name)
	}
	_postNumComments := 12
	_postVotesUp := 34
	_postVotesDown := 56
	post.SetNumComments(_postNumComments).SetNumVotesUp(_postVotesUp).SetNumVotesDown(_postVotesDown)
	js1, _ := json.Marshal(post)

	var post2 *BlogPost
	err := json.Unmarshal(js1, &post2)
	if err != nil {
		t.Fatalf("%s failed: %e", name, err)
	}
	if tagVersion := post2.GetTagVersion(); tagVersion != _tagVersion {
		t.Fatalf("%s failed: expected tag-version to be %#v but received %#v", name, _tagVersion, tagVersion)
	}
	if ownerId := post2.GetOwnerId(); ownerId != _userId {
		t.Fatalf("%s failed: expected owner-id to be %#v but received %#v", name, _userId, ownerId)
	}
	if title := post2.GetTitle(); title != _postTitle {
		t.Fatalf("%s failed: expected title to be %#v but received %#v", name, _postTitle, title)
	}
	if content := post2.GetContent(); content != _postContent {
		t.Fatalf("%s failed: expected content to be %#v but received %#v", name, _postContent, content)
	}
	if numComments := post2.GetNumComments(); numComments != _postNumComments {
		t.Fatalf("%s failed: expected num-comments to be %#v but received %#v", name, _postNumComments, numComments)
	}
	if numVotesUp := post2.GetNumVotesUp(); numVotesUp != _postVotesUp {
		t.Fatalf("%s failed: expected num-votes-up to be %#v but received %#v", name, _postVotesUp, numVotesUp)
	}
	if numVotesDown := post2.GetNumVotesDown(); numVotesDown != _postVotesDown {
		t.Fatalf("%s failed: expected num-votes-down to be %#v but received %#v", name, _postVotesDown, numVotesDown)
	}
	if t2, t1 := post2.GetTimeCreated(), post.GetTimeCreated(); !t2.Equal(t1) {
		t.Fatalf("%s failed: expected %#v but received %#v", name, t1.Format(time.RFC3339), t2.Format(time.RFC3339))
	}
	if post2.GetChecksum() != post.GetChecksum() {
		t.Fatalf("%s failed: expected %#v but received %#v", name, post2.GetChecksum(), post.GetChecksum())
	}
}

package blog

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/btnguyen2k/henge"
	"main/src/gvabe/bov2/user"
	"main/src/utils"
)

func TestNewBlogComment(t *testing.T) {
	name := "TestNewBlogComment"
	_tagVersion := uint64(1337)
	_userId := "admin@local"
	_userMaskId := "admin"
	_user := user.NewUser(_tagVersion, _userId, _userMaskId)

	_postIsPublic := true
	_postTitle := "Blog post title"
	_postContent := "Blog post content"
	_post := NewBlogPost(_tagVersion, _user, _postIsPublic, _postTitle, _postContent)

	_commentContent := "Blog comment content"
	comment := NewBlogComment(_tagVersion, _user, _post, nil, _commentContent)

	if comment == nil {
		t.Fatalf("%s failed: nil", name)
	}
	if tagVersion := comment.GetTagVersion(); tagVersion != _tagVersion {
		t.Fatalf("%s failed: expected tag-version to be %#v but received %#v", name, _tagVersion, tagVersion)
	}
	if ownerId := comment.GetOwnerId(); ownerId != _userId {
		t.Fatalf("%s failed: expected owner-id to be %#v but received %#v", name, _userId, ownerId)
	}
	if postId := comment.GetPostId(); postId != _post.GetId() {
		t.Fatalf("%s failed: expected post-id to be %#v but received %#v", name, _post.GetId(), postId)
	}
	if parentId := comment.GetParentId(); parentId != "" {
		t.Fatalf("%s failed: expected parent-id to be %#v but received %#v", name, "", parentId)
	}
	if content := comment.GetContent(); content != _commentContent {
		t.Fatalf("%s failed: expected content to be %#v but received %#v", name, _commentContent, content)
	}
}

func TestNewBlogCommentFromUbo(t *testing.T) {
	name := "TestNewBlogCommentFromUbo"

	if NewBlogCommentFromUbo(nil) != nil {
		t.Fatalf("%s failed: NewBlogCommentFromUbo(nil) should return nil", name)
	}
	_tagVersion := uint64(1337)
	_id := utils.UniqueId()
	_postId := utils.UniqueId()
	_userId := utils.UniqueId()
	_parentId := utils.UniqueId()
	_content := "Blog content"
	ubo := henge.NewUniversalBo(_id, _tagVersion)
	ubo.SetExtraAttr(CommentFieldPostId, _postId)
	ubo.SetExtraAttr(CommentFieldOwnerId, _userId)
	ubo.SetExtraAttr(CommentFieldParentId, _parentId)
	ubo.SetDataAttr(CommentAttrContent, _content)

	comment := NewBlogCommentFromUbo(ubo)
	if tagVersion := comment.GetTagVersion(); tagVersion != _tagVersion {
		t.Fatalf("%s failed: expected tag-version to be %#v but received %#v", name, _tagVersion, tagVersion)
	}
	if id := comment.GetId(); id != _id {
		t.Fatalf("%s failed: expected bo's id to be %#v but received %#v", name, _id, id)
	}
	if ownerId := comment.GetOwnerId(); ownerId != _userId {
		t.Fatalf("%s failed: expected owner-id to be %#v but received %#v", name, _userId, ownerId)
	}
	if postId := comment.GetPostId(); postId != _postId {
		t.Fatalf("%s failed: expected post-id to be %#v but received %#v", name, _postId, postId)
	}
	if parentId := comment.GetParentId(); parentId != _parentId {
		t.Fatalf("%s failed: expected parent-id to be %#v but received %#v", name, _parentId, parentId)
	}
	if content := comment.GetContent(); content != _content {
		t.Fatalf("%s failed: expected content to be %#v but received %#v", name, _content, content)
	}
}

func TestBlogComment_ToMap(t *testing.T) {
	name := "TestBlogComment_ToMap"
	_tagVersion := uint64(1337)
	_userId := "admin@local"
	_userMaskId := "admin"
	_user := user.NewUser(_tagVersion, _userId, _userMaskId)

	_postIsPublic := true
	_postTitle := "Blog post title"
	_postContent := "Blog post content"
	_post := NewBlogPost(_tagVersion, _user, _postIsPublic, _postTitle, _postContent)

	_commentContent := "Blog comment content"
	comment := NewBlogComment(_tagVersion, _user, _post, nil, _commentContent)

	m := comment.ToMap(nil)
	expected := map[string]interface{}{
		henge.FieldId:        comment.GetId(),
		CommentFieldPostId:   _post.GetId(),
		CommentFieldOwnerId:  _userId,
		CommentFieldParentId: "",
		CommentAttrContent:   _commentContent,
	}
	if !reflect.DeepEqual(m, expected) {
		t.Fatalf("%s failed: expected %#v but received %#v", name, expected, m)
	}

	m = comment.ToMap(func(input map[string]interface{}) map[string]interface{} {
		return map[string]interface{}{
			"FieldId":              input[henge.FieldId],
			"CommentFieldPostId":   input[CommentFieldPostId],
			"CommentFieldOwnerId":  input[CommentFieldOwnerId],
			"CommentFieldParentId": input[CommentFieldParentId],
			"CommentAttrContent":   input[CommentAttrContent],
		}
	})
	expected = map[string]interface{}{
		"FieldId":              comment.GetId(),
		"CommentFieldPostId":   _post.GetId(),
		"CommentFieldOwnerId":  _userId,
		"CommentFieldParentId": "",
		"CommentAttrContent":   _commentContent,
	}
	if !reflect.DeepEqual(m, expected) {
		t.Fatalf("%s failed: expected %#v but received %#v", name, expected, m)
	}
}

func TestBlogComment_json(t *testing.T) {
	name := "TestBlogComment_json"
	_tagVersion := uint64(1337)
	_userId := "admin@local"
	_userMaskId := "admin"
	_user := user.NewUser(_tagVersion, _userId, _userMaskId)

	_postIsPublic := true
	_postTitle := "Blog post title"
	_postContent := "Blog post content"
	_post := NewBlogPost(_tagVersion, _user, _postIsPublic, _postTitle, _postContent)

	_commentContent := "Blog comment content"
	comment := NewBlogComment(_tagVersion, _user, _post, nil, _commentContent)
	js1, _ := json.Marshal(comment)

	var comment2 *BlogComment
	err := json.Unmarshal(js1, &comment2)
	if err != nil {
		t.Fatalf("%s failed: %e", name, err)
	}
	if tagVersion := comment2.GetTagVersion(); tagVersion != _tagVersion {
		t.Fatalf("%s failed: expected tag-version to be %#v but received %#v", name, _tagVersion, tagVersion)
	}
	if ownerId := comment2.GetOwnerId(); ownerId != _userId {
		t.Fatalf("%s failed: expected owner-id to be %#v but received %#v", name, _userId, ownerId)
	}
	if postId := comment2.GetPostId(); postId != _post.GetId() {
		t.Fatalf("%s failed: expected post-id to be %#v but received %#v", name, _post.GetId(), postId)
	}
	if parentId := comment2.GetParentId(); parentId != "" {
		t.Fatalf("%s failed: expected parent-id to be %#v but received %#v", name, "", parentId)
	}
	if content := comment2.GetContent(); content != _commentContent {
		t.Fatalf("%s failed: expected content to be %#v but received %#v", name, _commentContent, content)
	}
}

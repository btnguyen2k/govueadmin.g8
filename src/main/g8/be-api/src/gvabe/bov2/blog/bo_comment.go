package blog

import (
	"encoding/json"
	"strings"

	"github.com/btnguyen2k/consu/reddo"
	"github.com/btnguyen2k/henge"

	"main/src/gvabe/bov2/user"
	"main/src/utils"
)

// NewBlogComment is helper function to create new BlogComment bo.
//
// Available since template-v0.2.0
func NewBlogComment(appVersion uint64, owner *user.User, post *BlogPost, parent *BlogComment, content string) *BlogComment {
	comment := &BlogComment{
		UniversalBo: henge.NewUniversalBo(utils.UniqueId(), appVersion),
		ownerId:     owner.GetId(),
		postId:      post.GetId(),
		content:     strings.TrimSpace(content),
	}
	if parent != nil {
		comment.parentId = parent.GetId()
	}
	return comment.sync()
}

// NewBlogCommentFromUbo is helper function to create BlogComment bo from a universal bo.
//
// Available since template-v0.2.0
func NewBlogCommentFromUbo(ubo *henge.UniversalBo) *BlogComment {
	if ubo == nil {
		return nil
	}
	ubo = ubo.Clone()
	comment := &BlogComment{UniversalBo: ubo}
	if v, err := ubo.GetExtraAttrAs(CommentFieldOwnerId, reddo.TypeString); err != nil {
		return nil
	} else {
		comment.ownerId = v.(string)
	}
	if v, err := ubo.GetExtraAttrAs(CommentFieldPostId, reddo.TypeString); err != nil {
		return nil
	} else {
		comment.postId = v.(string)
	}
	if v, err := ubo.GetExtraAttrAs(CommentFieldParentId, reddo.TypeString); err != nil {
		return nil
	} else {
		comment.parentId = v.(string)
	}
	if v, err := ubo.GetDataAttrAs(CommentAttrContent, reddo.TypeString); err != nil {
		return nil
	} else {
		comment.content = v.(string)
	}
	return comment.sync()
}

const (
	// CommentFieldOwnerId is id of the user who created the comment.
	CommentFieldOwnerId = "oid"

	// CommentFieldPostId is id of the blog post the comment belongs to.
	CommentFieldPostId = "pid"

	// CommentFieldParentId is id of the parent comment.
	CommentFieldParentId = "prid"

	// CommentAttrContent is comment's content.
	CommentAttrContent = "cont"

	// commentAttr_Ubo is for internal use only!
	commentAttr_Ubo = "_ubo"
)

// BlogComment is the business object.
//   - BlogComment inherits unique id from bo.UniversalBo
//
// Available since template-v0.2.0
type BlogComment struct {
	*henge.UniversalBo `json:"_ubo"`
	ownerId            string `json:"oid"`
	postId             string `json:"pid"`
	parentId           string `json:"prid"`
	content            string `json:"cont"`
}

// ToMap transforms comment's attributes to a map.
//
// Available since template-v0.3.0
func (c *BlogComment) ToMap(postFunc henge.FuncPostUboToMap) map[string]interface{} {
	result := map[string]interface{}{
		henge.FieldId:          c.GetId(),
		henge.FieldTimeCreated: c.GetTimeCreated(),
		CommentFieldPostId:     c.GetPostId(),
		CommentFieldOwnerId:    c.GetOwnerId(),
		CommentFieldParentId:   c.GetParentId(),
		CommentAttrContent:     c.GetContent(),
	}
	if postFunc != nil {
		result = postFunc(result)
	}
	return result
}

// MarshalJSON implements json.encode.Marshaler.MarshalJSON.
// TODO: lock for read?
func (c *BlogComment) MarshalJSON() ([]byte, error) {
	c.sync()
	m := map[string]interface{}{
		commentAttr_Ubo: c.UniversalBo.Clone(),
		"_cols": map[string]interface{}{
			CommentFieldOwnerId:  c.ownerId,
			CommentFieldPostId:   c.postId,
			CommentFieldParentId: c.parentId,
		},
		"_attrs": map[string]interface{}{
			CommentAttrContent: c.content,
		},
	}
	return json.Marshal(m)
}

// UnmarshalJSON implements json.decode.Unmarshaler.UnmarshalJSON.
// TODO: lock for write?
func (c *BlogComment) UnmarshalJSON(data []byte) error {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	var err error
	if m[commentAttr_Ubo] != nil {
		js, _ := json.Marshal(m[commentAttr_Ubo])
		if err = json.Unmarshal(js, &c.UniversalBo); err != nil {
			return err
		}
	}
	if _cols, ok := m["_cols"].(map[string]interface{}); ok {
		if c.ownerId, err = reddo.ToString(_cols[CommentFieldOwnerId]); err != nil {
			return err
		}
		if c.postId, err = reddo.ToString(_cols[CommentFieldPostId]); err != nil {
			return err
		}
		if c.parentId, err = reddo.ToString(_cols[CommentFieldParentId]); err != nil {
			return err
		}
	}
	if _attrs, ok := m["_attrs"].(map[string]interface{}); ok {
		if c.content, err = reddo.ToString(_attrs[CommentAttrContent]); err != nil {
			return err
		}
	}
	c.sync()
	return nil
}

// GetOwnerId returns value of blog comment's 'owner-id' attribute
func (c *BlogComment) GetOwnerId() string {
	return c.ownerId
}

// SetOwnerId sets value of blog comment's 'owner-id' attribute
func (c *BlogComment) SetOwnerId(v string) *BlogComment {
	c.ownerId = strings.TrimSpace(strings.ToLower(v))
	return c
}

// GetPostId returns value of blog comment's 'post-id' attribute
func (c *BlogComment) GetPostId() string {
	return c.postId
}

// SetPostId sets value of blog comment's 'post-id' attribute
func (c *BlogComment) SetPostId(v string) *BlogComment {
	c.postId = strings.TrimSpace(strings.ToLower(v))
	return c
}

// GetParentId returns value of blog comment's 'parent-id' attribute
func (c *BlogComment) GetParentId() string {
	return c.parentId
}

// SetParentId sets value of blog comment's 'parent-id' attribute
func (c *BlogComment) SetParentId(v string) *BlogComment {
	c.parentId = strings.TrimSpace(strings.ToLower(v))
	return c
}

// GetContent returns value of blog comment's 'content' attribute
func (c *BlogComment) GetContent() string {
	return c.content
}

// SetContent sets value of blog comment's 'content' attribute
func (c *BlogComment) SetContent(v string) *BlogComment {
	c.content = strings.TrimSpace(v)
	return c
}

// sync is called to synchronize BO's attributes to its UniversalBo
func (c *BlogComment) sync() *BlogComment {
	c.SetDataAttr(CommentAttrContent, c.content)
	c.SetExtraAttr(CommentFieldOwnerId, c.ownerId)
	c.SetExtraAttr(CommentFieldPostId, c.postId)
	c.SetExtraAttr(CommentFieldParentId, c.parentId)
	c.UniversalBo.Sync()
	return c
}

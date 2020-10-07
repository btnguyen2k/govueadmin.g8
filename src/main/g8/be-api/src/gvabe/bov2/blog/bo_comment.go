package user

import (
	"encoding/json"
	"strings"

	"github.com/btnguyen2k/consu/reddo"

	userv2 "main/src/gvabe/bov2/user"
	"main/src/henge"
	"main/src/utils"
)

// NewBlogComment is helper function to create new BlogComment bo
//
// available since template-v0.2.0
func NewBlogComment(appVersion uint64, owner *userv2.User, post *BlogPost, parent *BlogComment, content string) *BlogComment {
	comment := &BlogComment{
		UniversalBo: *henge.NewUniversalBo(utils.UniqueId(), appVersion),
		ownerId:     strings.TrimSpace(strings.ToLower(owner.GetId())),
		postId:      strings.TrimSpace(strings.ToLower(post.GetId())),
		content:     strings.TrimSpace(content),
	}
	if parent != nil {
		comment.parentId = strings.TrimSpace(strings.ToLower(parent.GetId()))
	}
	return comment.sync()
}

// NewBlogCommentFromUbo is helper function to create BlogComment bo from a universal bo
//
// available since template-v0.2.0
func NewBlogCommentFromUbo(ubo *henge.UniversalBo) *BlogComment {
	if ubo == nil {
		return nil
	}
	comment := BlogComment{UniversalBo: *ubo.Clone()}
	if v, err := comment.GetExtraAttrAs(CommentField_OwnerId, reddo.TypeString); err != nil {
		return nil
	} else {
		comment.ownerId = v.(string)
	}
	if v, err := comment.GetExtraAttrAs(CommentField_PostId, reddo.TypeString); err != nil {
		return nil
	} else {
		comment.postId = v.(string)
	}
	if v, err := comment.GetExtraAttrAs(CommentField_ParentId, reddo.TypeString); err != nil {
		return nil
	} else {
		comment.parentId = v.(string)
	}
	if v, err := comment.GetDataAttrAs(CommentAttr_Content, reddo.TypeString); err != nil {
		return nil
	} else {
		comment.content = v.(string)
	}
	return (&comment).sync()
}

const (
	// id of user who is owner of the blog comment
	CommentField_OwnerId = "oid"

	// id of the blog post the comment belongs to
	CommentField_PostId = "pid"

	// id of the parent comment
	CommentField_ParentId = "prid"

	// content of blog comment
	CommentAttr_Content = "cont"

	commentAttr_Ubo = "_ubo"
)

// BlogComment is the business object
//	- BlogComment inherits unique id from bo.UniversalBo
//
// available since template-v0.2.0
type BlogComment struct {
	henge.UniversalBo `json:"_ubo"`
	ownerId           string `json:"oid"`
	postId            string `json:"pid"`
	parentId          string `json:"prid"`
	content           string `json:"cont"`
}

// MarshalJSON implements json.encode.Marshaler.MarshalJSON
//	TODO: lock for read?
func (c *BlogComment) MarshalJSON() ([]byte, error) {
	c.sync()
	m := map[string]interface{}{
		commentAttr_Ubo: c.UniversalBo.Clone(),
		"_cols": map[string]interface{}{
			CommentField_OwnerId:  c.ownerId,
			CommentField_PostId:   c.postId,
			CommentField_ParentId: c.parentId,
		},
		"_attrs": map[string]interface{}{
			CommentAttr_Content: c.content,
		},
	}
	return json.Marshal(m)
}

// UnmarshalJSON implements json.decode.Unmarshaler.UnmarshalJSON
//	TODO: lock for write?
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
		if c.ownerId, err = reddo.ToString(_cols[CommentField_OwnerId]); err != nil {
			return err
		}
		if c.postId, err = reddo.ToString(_cols[CommentField_PostId]); err != nil {
			return err
		}
		if c.parentId, err = reddo.ToString(_cols[CommentField_ParentId]); err != nil {
			return err
		}
	}
	if _attrs, ok := m["_attrs"].(map[string]interface{}); ok {
		if c.content, err = reddo.ToString(_attrs[CommentAttr_Content]); err != nil {
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
	c.SetDataAttr(CommentAttr_Content, c.content)
	c.SetExtraAttr(CommentField_OwnerId, c.ownerId)
	c.SetExtraAttr(CommentField_PostId, c.postId)
	c.SetExtraAttr(CommentField_ParentId, c.parentId)
	c.UniversalBo.Sync()
	return c
}

// BlogCommentDao defines API to access BlogComment storage
//
// available since template-v0.2.0
type BlogCommentDao interface {
	// Delete removes the specified business object from storage
	Delete(bo *BlogComment) (bool, error)

	// Create persists a new business object to storage
	Create(bo *BlogComment) (bool, error)

	// Get retrieves a business object from storage
	Get(id string) (*BlogComment, error)

	// GetN retrieves N business objects from storage
	GetN(fromOffset, maxNumRows int) ([]*BlogComment, error)

	// GetAll retrieves all available business objects from storage
	GetAll() ([]*BlogComment, error)

	// Update modifies an existing business object
	Update(bo *BlogComment) (bool, error)
}

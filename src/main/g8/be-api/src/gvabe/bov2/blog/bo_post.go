package blog

import (
	"encoding/json"
	"strings"

	"github.com/btnguyen2k/consu/reddo"

	userv2 "main/src/gvabe/bov2/user"
	"main/src/henge"
	"main/src/utils"
)

// NewBlogPost is helper function to create new BlogPost bo
//
// available since template-v0.2.0
func NewBlogPost(appVersion uint64, owner *userv2.User, isPublic bool, title, content string) *BlogPost {
	post := &BlogPost{
		UniversalBo:  *henge.NewUniversalBo(utils.UniqueId(), appVersion),
		ownerId:      strings.TrimSpace(strings.ToLower(owner.GetId())),
		isPublic:     isPublic,
		title:        strings.TrimSpace(title),
		content:      strings.TrimSpace(content),
		numComments:  0,
		numVotesUp:   0,
		numVotesDown: 0,
	}
	return post.sync()
}

// NewBlogPostFromUbo is helper function to create BlogPost bo from a universal bo
//
// available since template-v0.2.0
func NewBlogPostFromUbo(ubo *henge.UniversalBo) *BlogPost {
	if ubo == nil {
		return nil
	}
	post := BlogPost{UniversalBo: *ubo.Clone()}
	if v, err := post.GetExtraAttrAs(PostField_OwnerId, reddo.TypeString); err != nil {
		return nil
	} else {
		post.ownerId = v.(string)
	}
	if v, err := post.GetExtraAttrAs(PostField_IsPublic, reddo.TypeBool); err != nil {
		return nil
	} else {
		post.isPublic = v.(bool)
	}
	if v, err := post.GetDataAttrAs(PostAttr_Title, reddo.TypeString); err != nil {
		return nil
	} else {
		post.title = v.(string)
	}
	if v, err := post.GetDataAttrAs(PostAttr_Content, reddo.TypeString); err != nil {
		return nil
	} else {
		post.content = v.(string)
	}
	if v, err := post.GetDataAttrAs(PostAttr_NumComments, reddo.TypeInt); err != nil {
		return nil
	} else {
		post.numComments = int(v.(int64))
	}
	if v, err := post.GetDataAttrAs(PostAttr_NumVotesUp, reddo.TypeInt); err != nil {
		return nil
	} else {
		post.numVotesUp = int(v.(int64))
	}
	if v, err := post.GetDataAttrAs(PostAttr_NumVotesDown, reddo.TypeInt); err != nil {
		return nil
	} else {
		post.numVotesDown = int(v.(int64))
	}
	return (&post).sync()
}

const (
	// id of user who is owner of the blog post
	PostField_OwnerId = "oid"

	// flag to mark if the blog post is public or private
	PostField_IsPublic = "ispub"

	// title of blog post
	PostAttr_Title = "title"

	// content of blog post
	PostAttr_Content = "cont"

	// number of comments
	PostAttr_NumComments = "ncmts"

	// number of votes up
	PostAttr_NumVotesUp = "vup"

	// number of votes down
	PostAttr_NumVotesDown = "vdown"

	postAttr_Ubo = "_ubo"
)

// BlogPost is the business object
//	- BlogPost inherits unique id from bo.UniversalBo
//
// available since template-v0.2.0
type BlogPost struct {
	henge.UniversalBo `json:"_ubo"`
	ownerId           string `json:"oid"`
	isPublic          bool   `json:"ispub"`
	title             string `json:"title"`
	content           string `json:"cont"`
	numComments       int    `json:"ncmts"`
	numVotesUp        int    `json:"vup"`
	numVotesDown      int    `json:"vdown"`
}

// MarshalJSON implements json.encode.Marshaler.MarshalJSON
//	TODO: lock for read?
func (p *BlogPost) MarshalJSON() ([]byte, error) {
	p.sync()
	m := map[string]interface{}{
		postAttr_Ubo: p.UniversalBo.Clone(),
		"_cols": map[string]interface{}{
			PostField_OwnerId:  p.ownerId,
			PostField_IsPublic: p.isPublic,
		},
		"_attrs": map[string]interface{}{
			PostAttr_Title:        p.title,
			PostAttr_Content:      p.content,
			PostAttr_NumComments:  p.numComments,
			PostAttr_NumVotesUp:   p.numVotesUp,
			PostAttr_NumVotesDown: p.numVotesDown,
		},
	}
	return json.Marshal(m)
}

// UnmarshalJSON implements json.decode.Unmarshaler.UnmarshalJSON
//	TODO: lock for write?
func (p *BlogPost) UnmarshalJSON(data []byte) error {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	var err error
	if m[postAttr_Ubo] != nil {
		js, _ := json.Marshal(m[postAttr_Ubo])
		if err = json.Unmarshal(js, &p.UniversalBo); err != nil {
			return err
		}
	}
	if _cols, ok := m["_cols"].(map[string]interface{}); ok {
		if p.ownerId, err = reddo.ToString(_cols[PostField_OwnerId]); err != nil {
			return err
		}
		if p.isPublic, err = reddo.ToBool(_cols[PostField_IsPublic]); err != nil {
			return err
		}
	}
	if _attrs, ok := m["_attrs"].(map[string]interface{}); ok {
		if p.title, err = reddo.ToString(_attrs[PostAttr_Title]); err != nil {
			return err
		}
		if p.content, err = reddo.ToString(_attrs[PostAttr_Content]); err != nil {
			return err
		}
		if v, err := reddo.ToInt(_attrs[PostAttr_NumComments]); err != nil {
			return err
		} else {
			p.numComments = int(v)
		}
		if v, err := reddo.ToInt(_attrs[PostAttr_NumVotesUp]); err != nil {
			return err
		} else {
			p.numVotesUp = int(v)
		}
		if v, err := reddo.ToInt(_attrs[PostAttr_NumVotesDown]); err != nil {
			return err
		} else {
			p.numVotesDown = int(v)
		}
	}
	p.sync()
	return nil
}

// GetOwnerId returns value of blog post's 'owner-id' attribute
func (p *BlogPost) GetOwnerId() string {
	return p.ownerId
}

// SetOwnerId sets value of blog post's 'owner-id' attribute
func (p *BlogPost) SetOwnerId(v string) *BlogPost {
	p.ownerId = strings.TrimSpace(strings.ToLower(v))
	return p
}

// IsPublic returns value of blog post's 'is-public' attribute
func (p *BlogPost) IsPublic() bool {
	return p.isPublic
}

// SetPublic sets value of blog post's 'is-public' attribute
func (p *BlogPost) SetPublic(v bool) *BlogPost {
	p.isPublic = v
	return p
}

// GetTitle returns value of blog post's 'title' attribute
func (p *BlogPost) GetTitle() string {
	return p.title
}

// SetTitle sets value of blog post's 'title' attribute
func (p *BlogPost) SetTitle(v string) *BlogPost {
	p.title = strings.TrimSpace(v)
	return p
}

// GetContent returns value of blog post's 'content' attribute
func (p *BlogPost) GetContent() string {
	return p.content
}

// SetContent sets value of blog post's 'content' attribute
func (p *BlogPost) SetContent(v string) *BlogPost {
	p.content = strings.TrimSpace(v)
	return p
}

// GetNumComments returns value of blog post's 'num-comments' attribute
func (p *BlogPost) NumComments() int {
	return p.numComments
}

// SetNumComments sets value of blog post's 'num-comments' attribute
func (p *BlogPost) SetNumComments(v int) *BlogPost {
	p.numComments = v
	return p
}

// IncNumComments increases value of blog post's 'num-comments' attribute
func (p *BlogPost) IncNumComments(delta int) *BlogPost {
	p.numComments += delta
	return p
}

// GetNumVotesUp returns value of blog post's 'num-votes-up' attribute
func (p *BlogPost) GetNumVotesUp() int {
	return p.numVotesUp
}

// SetNumVotesUp sets value of blog post's 'num-votes-up' attribute
func (p *BlogPost) SetNumVotesUp(v int) *BlogPost {
	p.numVotesUp = v
	return p
}

// IncNumVotesUp increases value of blog post's 'num-votes-up' attribute
func (p *BlogPost) IncNumVotesUp(delta int) *BlogPost {
	p.numVotesUp += delta
	return p
}

// GetNumVotesDown returns value of blog post's 'num-votes-down' attribute
func (p *BlogPost) GetNumVotesDown() int {
	return p.numVotesDown
}

// SetNumVotesDown sets value of blog post's 'num-votes-down' attribute
func (p *BlogPost) SetNumVotesDown(v int) *BlogPost {
	p.numVotesDown = v
	return p
}

// IncNumVotesDown increases value of blog post's 'num-votes-down' attribute
func (p *BlogPost) IncNumVotesDown(delta int) *BlogPost {
	p.numVotesDown += delta
	return p
}

// sync is called to synchronize BO's attributes to its UniversalBo
func (p *BlogPost) sync() *BlogPost {
	p.SetDataAttr(PostAttr_Title, p.title)
	p.SetDataAttr(PostAttr_Content, p.content)
	p.SetDataAttr(PostAttr_NumComments, p.numComments)
	p.SetDataAttr(PostAttr_NumVotesUp, p.numVotesUp)
	p.SetDataAttr(PostAttr_NumVotesDown, p.numVotesDown)
	p.SetExtraAttr(PostField_OwnerId, p.ownerId)
	p.SetExtraAttr(PostField_IsPublic, p.isPublic)
	p.UniversalBo.Sync()
	return p
}

// BlogPostDao defines API to access BlogPost storage
//
// available since template-v0.2.0
type BlogPostDao interface {
	// Delete removes the specified business object from storage
	Delete(bo *BlogPost) (bool, error)

	// Create persists a new business object to storage
	Create(bo *BlogPost) (bool, error)

	// Get retrieves a business object from storage
	Get(id string) (*BlogPost, error)

	// GetN retrieves N business objects from storage
	GetN(fromOffset, maxNumRows int) ([]*BlogPost, error)

	// GetAll retrieves all available business objects from storage
	GetAll() ([]*BlogPost, error)

	// Update modifies an existing business object
	Update(bo *BlogPost) (bool, error)
}

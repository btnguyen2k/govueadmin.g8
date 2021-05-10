package blog

import (
	"encoding/json"
	"strings"

	"github.com/btnguyen2k/consu/reddo"
	"github.com/btnguyen2k/henge"

	"main/src/gvabe/bov2/user"
	"main/src/utils"
)

// NewBlogPost is helper function to create new BlogPost bo
//
// Available since template-v0.2.0
func NewBlogPost(appVersion uint64, owner *user.User, isPublic bool, title, content string) *BlogPost {
	post := &BlogPost{
		UniversalBo:  henge.NewUniversalBo(utils.UniqueId(), appVersion),
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
// Available since template-v0.2.0
func NewBlogPostFromUbo(ubo *henge.UniversalBo) *BlogPost {
	if ubo == nil {
		return nil
	}
	ubo = ubo.Clone()
	post := &BlogPost{UniversalBo: ubo}
	if v, err := ubo.GetExtraAttrAs(PostFieldOwnerId, reddo.TypeString); err != nil {
		return nil
	} else {
		post.ownerId = v.(string)
	}
	if v, err := ubo.GetExtraAttrAs(PostFieldIsPublic, reddo.TypeInt); err != nil {
		return nil
	} else {
		post.isPublic = v.(int64) != 0
	}
	if v, err := ubo.GetDataAttrAs(PostAttrTitle, reddo.TypeString); err != nil {
		return nil
	} else {
		post.title = v.(string)
	}
	if v, err := ubo.GetDataAttrAs(PostAttrContent, reddo.TypeString); err != nil {
		return nil
	} else {
		post.content = v.(string)
	}
	if v, err := ubo.GetDataAttrAs(PostAttrNumComments, reddo.TypeInt); err != nil {
		return nil
	} else {
		post.numComments = int(v.(int64))
	}
	if v, err := ubo.GetDataAttrAs(PostAttrNumVotesUp, reddo.TypeInt); err != nil {
		return nil
	} else {
		post.numVotesUp = int(v.(int64))
	}
	if v, err := ubo.GetDataAttrAs(PostAttrNumVotesDown, reddo.TypeInt); err != nil {
		return nil
	} else {
		post.numVotesDown = int(v.(int64))
	}
	return post.sync()
}

const (
	// PostFieldOwnerId is id of the user who made the blog post.
	PostFieldOwnerId = "oid"

	// PostFieldIsPublic is a flag to mark if the blog post is public or private.
	PostFieldIsPublic = "ispub"

	// PostAttrTitle is blog post's title.
	PostAttrTitle = "title"

	// PostAttrContent is blog post's content.
	PostAttrContent = "cont"

	// PostAttrNumComments is blog post's number of comments.
	PostAttrNumComments = "ncmts"

	// PostAttrNumVotesUp is blog post's number of votes up.
	PostAttrNumVotesUp = "vup"

	// PostAttrNumVotesDown is blog post's number of votes down.
	PostAttrNumVotesDown = "vdown"

	// postAttr_Ubo is for internal use only!
	postAttr_Ubo = "_ubo"
)

// BlogPost is the business object.
//   - BlogPost inherits unique id from bo.UniversalBo
//
// Available since template-v0.2.0
type BlogPost struct {
	*henge.UniversalBo `json:"_ubo"`
	ownerId            string `json:"oid"`
	isPublic           bool   `json:"ispub"`
	title              string `json:"title"`
	content            string `json:"cont"`
	numComments        int    `json:"ncmts"`
	numVotesUp         int    `json:"vup"`
	numVotesDown       int    `json:"vdown"`
}

// ToMap transforms post's attributes to a map.
func (p *BlogPost) ToMap(postFunc henge.FuncPostUboToMap) map[string]interface{} {
	result := map[string]interface{}{
		henge.FieldId:        p.GetId(),
		PostFieldOwnerId:     p.ownerId,
		PostFieldIsPublic:    p.isPublic,
		PostAttrTitle:        p.title,
		PostAttrContent:      p.content,
		PostAttrNumComments:  p.numComments,
		PostAttrNumVotesUp:   p.numVotesUp,
		PostAttrNumVotesDown: p.numVotesDown,
	}
	if postFunc != nil {
		result = postFunc(result)
	}
	return result
}

// MarshalJSON implements json.encode.Marshaler.MarshalJSON.
// TODO: lock for read?
func (p *BlogPost) MarshalJSON() ([]byte, error) {
	p.sync()
	m := map[string]interface{}{
		postAttr_Ubo: p.UniversalBo.Clone(),
		"_cols": map[string]interface{}{
			PostFieldOwnerId:  p.ownerId,
			PostFieldIsPublic: p.isPublic,
		},
		"_attrs": map[string]interface{}{
			PostAttrTitle:        p.title,
			PostAttrContent:      p.content,
			PostAttrNumComments:  p.numComments,
			PostAttrNumVotesUp:   p.numVotesUp,
			PostAttrNumVotesDown: p.numVotesDown,
		},
	}
	return json.Marshal(m)
}

// UnmarshalJSON implements json.decode.Unmarshaler.UnmarshalJSON.
// TODO: lock for write?
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
		if p.ownerId, err = reddo.ToString(_cols[PostFieldOwnerId]); err != nil {
			return err
		}
		if p.isPublic, err = reddo.ToBool(_cols[PostFieldIsPublic]); err != nil {
			return err
		}
	}
	if _attrs, ok := m["_attrs"].(map[string]interface{}); ok {
		if p.title, err = reddo.ToString(_attrs[PostAttrTitle]); err != nil {
			return err
		}
		if p.content, err = reddo.ToString(_attrs[PostAttrContent]); err != nil {
			return err
		}
		if v, err := reddo.ToInt(_attrs[PostAttrNumComments]); err != nil {
			return err
		} else {
			p.numComments = int(v)
		}
		if v, err := reddo.ToInt(_attrs[PostAttrNumVotesUp]); err != nil {
			return err
		} else {
			p.numVotesUp = int(v)
		}
		if v, err := reddo.ToInt(_attrs[PostAttrNumVotesDown]); err != nil {
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
func (p *BlogPost) GetNumComments() int {
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
	vIsPublic := 1
	if !p.isPublic {
		vIsPublic = 0
	}
	p.SetDataAttr(PostAttrTitle, p.title)
	p.SetDataAttr(PostAttrContent, p.content)
	p.SetDataAttr(PostAttrNumComments, p.numComments)
	p.SetDataAttr(PostAttrNumVotesUp, p.numVotesUp)
	p.SetDataAttr(PostAttrNumVotesDown, p.numVotesDown)
	p.SetExtraAttr(PostFieldOwnerId, p.ownerId)
	p.SetExtraAttr(PostFieldIsPublic, vIsPublic)
	p.UniversalBo.Sync()
	return p
}

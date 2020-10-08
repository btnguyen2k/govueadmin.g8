package blog

import (
	"encoding/json"
	"strings"

	"github.com/btnguyen2k/consu/reddo"

	userv2 "main/src/gvabe/bov2/user"
	"main/src/henge"
	"main/src/utils"
)

// NewBlogVote is helper function to create new BlogVote bo
//
// available since template-v0.2.0
func NewBlogVote(appVersion uint64, owner *userv2.User, targetId string, value int) *BlogVote {
	vote := &BlogVote{
		UniversalBo: *henge.NewUniversalBo(utils.UniqueId(), appVersion),
		ownerId:     strings.TrimSpace(strings.ToLower(owner.GetId())),
		targetId:    strings.TrimSpace(strings.ToLower(targetId)),
		value:       value,
	}
	return vote.sync()
}

// NewBlogVoteFromUbo is helper function to create BlogVote bo from a universal bo
//
// available since template-v0.2.0
func NewBlogVoteFromUbo(ubo *henge.UniversalBo) *BlogVote {
	if ubo == nil {
		return nil
	}
	vote := BlogVote{UniversalBo: *ubo.Clone()}
	if v, err := vote.GetExtraAttrAs(VoteField_OwnerId, reddo.TypeString); err != nil {
		return nil
	} else {
		vote.ownerId = v.(string)
	}
	if v, err := vote.GetExtraAttrAs(VoteField_TargetId, reddo.TypeString); err != nil {
		return nil
	} else {
		vote.targetId = v.(string)
	}
	if v, err := vote.GetExtraAttrAs(VoteField_Value, reddo.TypeInt); err != nil {
		return nil
	} else {
		vote.value = int(v.(int64))
	}
	return (&vote).sync()
}

const (
	// id of user who is owner of the vote
	VoteField_OwnerId = "oid"

	// id of the target this vote is for
	VoteField_TargetId = "tid"

	// value of the vote (-1 or 1)
	VoteField_Value = "v"

	voteAttr_Ubo = "_ubo"
)

// BlogVote is the business object
//	- BlogVote inherits unique id from bo.UniversalBo
//
// available since template-v0.2.0
type BlogVote struct {
	henge.UniversalBo `json:"_ubo"`
	ownerId           string `json:"oid"`
	targetId          string `json:"tid"`
	value             int    `json:"v"`
}

// MarshalJSON implements json.encode.Marshaler.MarshalJSON
//	TODO: lock for read?
func (v *BlogVote) MarshalJSON() ([]byte, error) {
	v.sync()
	m := map[string]interface{}{
		voteAttr_Ubo: v.UniversalBo.Clone(),
		"_cols": map[string]interface{}{
			VoteField_OwnerId:  v.ownerId,
			VoteField_TargetId: v.targetId,
			VoteField_Value:    v.value,
		},
	}
	return json.Marshal(m)
}

// UnmarshalJSON implements json.decode.Unmarshaler.UnmarshalJSON
//	TODO: lock for write?
func (v *BlogVote) UnmarshalJSON(data []byte) error {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	var err error
	if m[voteAttr_Ubo] != nil {
		js, _ := json.Marshal(m[voteAttr_Ubo])
		if err = json.Unmarshal(js, &v.UniversalBo); err != nil {
			return err
		}
	}
	if _cols, ok := m["_cols"].(map[string]interface{}); ok {
		if v.ownerId, err = reddo.ToString(_cols[VoteField_OwnerId]); err != nil {
			return err
		}
		if v.targetId, err = reddo.ToString(_cols[VoteField_TargetId]); err != nil {
			return err
		}
		if _v, err := reddo.ToInt(_cols[VoteField_Value]); err != nil {
			return err
		} else {
			v.value = int(_v)
		}
	}
	v.sync()
	return nil
}

// GetOwnerId returns value of vote's 'owner-id' attribute
func (v *BlogVote) GetOwnerId() string {
	return v.ownerId
}

// SetOwnerId sets value of vote's 'owner-id' attribute
func (v *BlogVote) SetOwnerId(_v string) *BlogVote {
	v.ownerId = strings.TrimSpace(strings.ToLower(_v))
	return v
}

// GetTargetId returns value of vote's 'target-id' attribute
func (v *BlogVote) GetTargetId() string {
	return v.targetId
}

// SetTargetId sets value of vote's 'target-id' attribute
func (v *BlogVote) SetPostId(_v string) *BlogVote {
	v.targetId = strings.TrimSpace(strings.ToLower(_v))
	return v
}

// GetValue returns value of vote's 'value' attribute
func (v *BlogVote) GetValue() int {
	return v.value
}

// SetValue sets value of vote's 'value' attribute
func (v *BlogVote) SetValue(_v int) *BlogVote {
	v.value = _v
	return v
}

// sync is called to synchronize BO's attributes to its UniversalBo
func (v *BlogVote) sync() *BlogVote {
	v.SetExtraAttr(VoteField_OwnerId, v.ownerId)
	v.SetExtraAttr(VoteField_TargetId, v.targetId)
	v.SetExtraAttr(VoteField_Value, v.value)
	v.UniversalBo.Sync()
	return v
}

// BlogVoteDao defines API to access BlogVote storage
//
// available since template-v0.2.0
type BlogVoteDao interface {
	// Delete removes the specified business object from storage
	Delete(bo *BlogVote) (bool, error)

	// Create persists a new business object to storage
	Create(bo *BlogVote) (bool, error)

	// Get retrieves a business object from storage
	Get(id string) (*BlogVote, error)

	// GetN retrieves N business objects from storage
	GetN(fromOffset, maxNumRows int) ([]*BlogVote, error)

	// GetAll retrieves all available business objects from storage
	GetAll() ([]*BlogVote, error)

	// Update modifies an existing business object
	Update(bo *BlogVote) (bool, error)
}

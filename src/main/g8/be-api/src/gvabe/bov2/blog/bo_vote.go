package blog

import (
	"encoding/json"
	"strings"

	"github.com/btnguyen2k/consu/reddo"
	"github.com/btnguyen2k/henge"

	"main/src/gvabe/bov2/user"
	"main/src/utils"
)

// NewBlogVote is helper function to create new BlogVote bo.
//
// Available since template-v0.2.0
func NewBlogVote(appVersion uint64, owner *user.User, targetId string, value int) *BlogVote {
	vote := &BlogVote{
		UniversalBo: henge.NewUniversalBo(utils.UniqueId(), appVersion),
		ownerId:     strings.TrimSpace(strings.ToLower(owner.GetId())),
		targetId:    strings.TrimSpace(strings.ToLower(targetId)),
		value:       value,
	}
	return vote.sync()
}

// NewBlogVoteFromUbo is helper function to create BlogVote bo from a universal bo.
//
// Available since template-v0.2.0
func NewBlogVoteFromUbo(ubo *henge.UniversalBo) *BlogVote {
	if ubo == nil {
		return nil
	}
	ubo = ubo.Clone()
	vote := &BlogVote{UniversalBo: ubo}
	if v, err := ubo.GetExtraAttrAs(VoteFieldOwnerId, reddo.TypeString); err != nil {
		return nil
	} else {
		vote.ownerId = v.(string)
	}
	if v, err := ubo.GetExtraAttrAs(VoteFieldTargetId, reddo.TypeString); err != nil {
		return nil
	} else {
		vote.targetId = v.(string)
	}
	if v, err := ubo.GetExtraAttrAs(VoteFieldValue, reddo.TypeInt); err != nil {
		return nil
	} else {
		vote.value = int(v.(int64))
	}
	return vote.sync()
}

const (
	// VoteFieldOwnerId is id of the user who made the vote.
	VoteFieldOwnerId = "oid"

	// VoteFieldTargetId is id of the target this vote is for.
	VoteFieldTargetId = "tid"

	// VoteFieldValue is value of the vote (-1 or 1)
	VoteFieldValue = "v"

	// voteAttr_Ubo is for internal use only!
	voteAttr_Ubo = "_ubo"
)

// BlogVote is the business object.
//   - BlogVote inherits unique id from bo.UniversalBo
//
// Available since template-v0.2.0
type BlogVote struct {
	*henge.UniversalBo `json:"_ubo"`
	ownerId            string `json:"oid"`
	targetId           string `json:"tid"`
	value              int    `json:"v"`
}

// ToMap transforms vote's attributes to a map.
//
// Available since template-v0.3.0
func (v *BlogVote) ToMap(postFunc henge.FuncPostUboToMap) map[string]interface{} {
	result := map[string]interface{}{
		henge.FieldId:     v.GetId(),
		VoteFieldTargetId: v.GetTargetId(),
		VoteFieldOwnerId:  v.GetOwnerId(),
		VoteFieldValue:    v.GetValue(),
	}
	if postFunc != nil {
		result = postFunc(result)
	}
	return result
}

// MarshalJSON implements json.encode.Marshaler.MarshalJSON.
// TODO: lock for read?
func (v *BlogVote) MarshalJSON() ([]byte, error) {
	v.sync()
	m := map[string]interface{}{
		voteAttr_Ubo: v.UniversalBo.Clone(),
		"_cols": map[string]interface{}{
			VoteFieldOwnerId:  v.ownerId,
			VoteFieldTargetId: v.targetId,
			VoteFieldValue:    v.value,
		},
	}
	return json.Marshal(m)
}

// UnmarshalJSON implements json.decode.Unmarshaler.UnmarshalJSON.
// TODO: lock for write?
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
		if v.ownerId, err = reddo.ToString(_cols[VoteFieldOwnerId]); err != nil {
			return err
		}
		if v.targetId, err = reddo.ToString(_cols[VoteFieldTargetId]); err != nil {
			return err
		}
		if _v, err := reddo.ToInt(_cols[VoteFieldValue]); err != nil {
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
func (v *BlogVote) SetTargetId(_v string) *BlogVote {
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
	v.SetExtraAttr(VoteFieldOwnerId, v.ownerId)
	v.SetExtraAttr(VoteFieldTargetId, v.targetId)
	v.SetExtraAttr(VoteFieldValue, v.value)
	v.UniversalBo.Sync()
	return v
}

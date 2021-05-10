package blog

import (
	"github.com/btnguyen2k/godal"
	"github.com/btnguyen2k/godal/sql"
	"github.com/btnguyen2k/henge"
	"main/src/gvabe/bov2/user"
)

const (
	TableBlogComment   = "gva_blog_comment"
	CommentColOwnerId  = "zownid"
	CommentColPostId   = "zpid"
	CommentColParentId = "zprid"
)

// BlogCommentDao defines API to access BlogComment storage.
//
// Available since template-v0.2.0
type BlogCommentDao interface {
	// Delete removes the specified business object from storage.
	Delete(bo *BlogComment) (bool, error)

	// Create persists a new business object to storage.
	Create(bo *BlogComment) (bool, error)

	// Get retrieves a business object from storage.
	Get(id string) (*BlogComment, error)

	// GetN retrieves N business objects from storage.
	GetN(fromOffset, maxNumRows int, filter godal.FilterOpt, sorting *godal.SortingOpt) ([]*BlogComment, error)

	// GetAll retrieves all available business objects from storage.
	GetAll(filter godal.FilterOpt, sorting *godal.SortingOpt) ([]*BlogComment, error)

	// Update modifies an existing business object.
	Update(bo *BlogComment) (bool, error)
}

// BaseBlogCommentDaoImpl is a generic implementation of BlogCommentDao.
//
// Available since template-v0.3.0
type BaseBlogCommentDaoImpl struct {
	henge.UniversalDao
}

// // GdaoCreateFilter implements IGenericDao.GdaoCreateFilter
// func (dao *BaseBlogCommentDaoImpl) GdaoCreateFilter(_ string, gbo godal.IGenericBo) godal.FilterOpt {
// 	return godal.MakeFilter(map[string]interface{}{henge.FieldId: gbo.GboGetAttrUnsafe(henge.FieldId, reddo.TypeString)})
// }

// Delete implements BlogCommentDao.Delete
func (dao *BaseBlogCommentDaoImpl) Delete(comment *BlogComment) (bool, error) {
	return dao.UniversalDao.Delete(comment.sync().UniversalBo)
}

// Create implements BlogCommentDao.Create
func (dao *BaseBlogCommentDaoImpl) Create(comment *BlogComment) (bool, error) {
	return dao.UniversalDao.Create(comment.sync().UniversalBo)
}

// Get implements BlogCommentDao.Get
func (dao *BaseBlogCommentDaoImpl) Get(id string) (*BlogComment, error) {
	ubo, err := dao.UniversalDao.Get(id)
	if err != nil {
		return nil, err
	}
	return NewBlogCommentFromUbo(ubo), nil
}

// GetN implements BlogCommentDao.GetN
func (dao *BaseBlogCommentDaoImpl) GetN(fromOffset, maxNumRows int, filter godal.FilterOpt, sorting *godal.SortingOpt) ([]*BlogComment, error) {
	uboList, err := dao.UniversalDao.GetN(fromOffset, maxNumRows, filter, sorting)
	if err != nil {
		return nil, err
	}
	result := make([]*BlogComment, 0)
	for _, ubo := range uboList {
		app := NewBlogCommentFromUbo(ubo)
		result = append(result, app)
	}
	return result, nil
}

// GetAll implements BlogCommentDao.GetAll
func (dao *BaseBlogCommentDaoImpl) GetAll(filter godal.FilterOpt, sorting *godal.SortingOpt) ([]*BlogComment, error) {
	return dao.GetN(0, 0, filter, sorting)
}

// Update implements BlogCommentDao.Update
func (dao *BaseBlogCommentDaoImpl) Update(comment *BlogComment) (bool, error) {
	return dao.UniversalDao.Update(comment.sync().UniversalBo)
}

/*----------------------------------------------------------------------*/

const (
	TableBlogPost   = "gva_blog_post"
	PostColOwnerId  = "zownid"
	PostColIsPublic = "zispub"
)

// BlogPostDao defines API to access BlogPost storage.
//
// Available since template-v0.2.0
type BlogPostDao interface {
	// Delete removes the specified business object from storage.
	Delete(bo *BlogPost) (bool, error)

	// Create persists a new business object to storage.
	Create(bo *BlogPost) (bool, error)

	// Get retrieves a business object from storage.
	Get(id string) (*BlogPost, error)

	// GetUserPostsN retrieves first N user's blog posts of a user, latest posts first.
	GetUserPostsN(user *user.User, fromOffset, maxNumRows int) ([]*BlogPost, error)

	// GetUserPostsAll retrieves all available user's blog posts, latest posts first.
	GetUserPostsAll(user *user.User) ([]*BlogPost, error)

	// GetUserFeedN retrieves first N blog posts for user's feed, latest posts first.
	GetUserFeedN(user *user.User, fromOffset, maxNumRows int) ([]*BlogPost, error)

	// GetUserFeedAll retrieves all available blog posts for user's feed, latest posts first.
	GetUserFeedAll(user *user.User) ([]*BlogPost, error)

	// Update modifies an existing business object.
	Update(bo *BlogPost) (bool, error)
}

// BaseBlogPostDaoImpl is a generic implementation of BlogPostDao.
//
// Available since template-v0.3.0
type BaseBlogPostDaoImpl struct {
	henge.UniversalDao
}

// // GdaoCreateFilter implements IGenericDao.GdaoCreateFilter
// func (dao *BaseBlogPostDaoImpl) GdaoCreateFilter(_ string, gbo godal.IGenericBo) interface{} {
// 	return map[string]interface{}{henge.SqlColId: gbo.GboGetAttrUnsafe(henge.FieldId, reddo.TypeString)}
// }

// Delete implements BlogPostDao.Delete
func (dao *BaseBlogPostDaoImpl) Delete(post *BlogPost) (bool, error) {
	return dao.UniversalDao.Delete(post.UniversalBo.Clone())
}

// Create implements BlogPostDao.Create
func (dao *BaseBlogPostDaoImpl) Create(post *BlogPost) (bool, error) {
	return dao.UniversalDao.Create(post.sync().UniversalBo.Clone())
}

// Get implements BlogPostDao.Get
func (dao *BaseBlogPostDaoImpl) Get(id string) (*BlogPost, error) {
	ubo, err := dao.UniversalDao.Get(id)
	if err != nil {
		return nil, err
	}
	return NewBlogPostFromUbo(ubo), nil
}

// GetUserPostsN implements BlogPostDao.GetUserPostsN
func (dao *BaseBlogPostDaoImpl) GetUserPostsN(user *user.User, fromOffset, maxNumRows int) ([]*BlogPost, error) {
	filter := &godal.FilterOptFieldOpValue{FieldName: PostFieldOwnerId, Operator: godal.FilterOpEqual, Value: user.GetId()}
	sorting := (&godal.SortingField{FieldName: henge.FieldTimeCreated, Descending: true}).ToSortingOpt()
	uboList, err := dao.UniversalDao.GetN(fromOffset, maxNumRows, filter, sorting)
	if err != nil {
		return nil, err
	}
	result := make([]*BlogPost, 0)
	for _, ubo := range uboList {
		app := NewBlogPostFromUbo(ubo)
		result = append(result, app)
	}
	return result, nil
}

// GetUserPostsAll implements BlogPostDao.GetAll
func (dao *BaseBlogPostDaoImpl) GetUserPostsAll(user *user.User) ([]*BlogPost, error) {
	return dao.GetUserPostsN(user, 0, 0)
}

// GetUserFeedN implements BlogPostDao.GetUserPostsN
func (dao *BaseBlogPostDaoImpl) GetUserFeedN(user *user.User, fromOffset, maxNumRows int) ([]*BlogPost, error) {
	filter := (&godal.FilterOptOr{}).
		Add(&godal.FilterOptFieldOpValue{FieldName: PostFieldOwnerId, Operator: godal.FilterOpEqual, Value: user.GetId()}).
		Add(&godal.FilterOptFieldOpValue{FieldName: PostFieldIsPublic, Operator: godal.FilterOpEqual, Value: 1})
	sorting := (&godal.SortingField{FieldName: henge.FieldTimeCreated, Descending: true}).ToSortingOpt()
	uboList, err := dao.UniversalDao.GetN(fromOffset, maxNumRows, filter, sorting)
	if err != nil {
		return nil, err
	}
	result := make([]*BlogPost, 0)
	for _, ubo := range uboList {
		app := NewBlogPostFromUbo(ubo)
		result = append(result, app)
	}
	return result, nil
}

// GetUserFeedAll implements BlogPostDao.GetUserFeedAll
func (dao *BaseBlogPostDaoImpl) GetUserFeedAll(user *user.User) ([]*BlogPost, error) {
	return dao.GetUserFeedN(user, 0, 0)
}

// Update implements BlogPostDao.Update
func (dao *BaseBlogPostDaoImpl) Update(post *BlogPost) (bool, error) {
	return dao.UniversalDao.Update(post.sync().UniversalBo.Clone())
}

/*----------------------------------------------------------------------*/

const (
	TableBlogVote   = "gva_blog_vote"
	VoteColOwnerId  = "zownid"
	VoteColTargetId = "ztid"
	VoteColValue    = "zval"
)

// BlogVoteDao defines API to access BlogVote storage.
//
// Available since template-v0.2.0
type BlogVoteDao interface {
	// GetUserVoteForTarget retrieves a user's vote against a target.
	GetUserVoteForTarget(user *user.User, targetId string) (*BlogVote, error)

	// Delete removes the specified business object from storage.
	Delete(bo *BlogVote) (bool, error)

	// Create persists a new business object to storage.
	Create(bo *BlogVote) (bool, error)

	// Get retrieves a business object from storage.
	Get(id string) (*BlogVote, error)

	// GetN retrieves N business objects from storage.
	GetN(fromOffset, maxNumRows int, filter godal.FilterOpt, sorting *godal.SortingOpt) ([]*BlogVote, error)

	// GetAll retrieves all available business objects from storage.
	GetAll(filter godal.FilterOpt, sorting *godal.SortingOpt) ([]*BlogVote, error)

	// Update modifies an existing business object.
	Update(bo *BlogVote) (bool, error)
}

// BaseBlogVoteDaoImpl is a generic implementation of BlogVoteDao.
//
// Available since template-v0.3.0
type BaseBlogVoteDaoImpl struct {
	henge.UniversalDao
}

// // GdaoCreateFilter implements IGenericDao.GdaoCreateFilter
// func (dao *BaseBlogVoteDaoImpl) GdaoCreateFilter(_ string, gbo godal.IGenericBo) interface{} {
// 	return map[string]interface{}{henge.SqlColId: gbo.GboGetAttrUnsafe(henge.FieldId, reddo.TypeString)}
// }

// GetUserVoteForTarget implements BlogVoteDao.GetUserVoteForTarget
func (dao *BaseBlogVoteDaoImpl) GetUserVoteForTarget(user *user.User, targetId string) (*BlogVote, error) {
	if user == nil || targetId == "" {
		return nil, nil
	}
	filter := &sql.FilterAnd{FilterAndOr: sql.FilterAndOr{
		Filters: []sql.IFilter{
			&sql.FilterFieldValue{Field: VoteColOwnerId, Operator: "=", Value: user.GetId()},
			&sql.FilterFieldValue{Field: VoteColTargetId, Operator: "=", Value: targetId},
		}},
	}
	uboList, err := dao.UniversalDao.GetAll(filter, nil)
	if err != nil {
		return nil, err
	}
	if uboList == nil || len(uboList) == 0 {
		return nil, nil
	}
	return NewBlogVoteFromUbo(uboList[0]), nil
}

// Delete implements BlogVoteDao.Delete
func (dao *BaseBlogVoteDaoImpl) Delete(vote *BlogVote) (bool, error) {
	return dao.UniversalDao.Delete(vote.UniversalBo.Clone())
}

// Create implements BlogVoteDao.Create
func (dao *BaseBlogVoteDaoImpl) Create(vote *BlogVote) (bool, error) {
	return dao.UniversalDao.Create(vote.sync().UniversalBo.Clone())
}

// Get implements BlogVoteDao.Get
func (dao *BaseBlogVoteDaoImpl) Get(id string) (*BlogVote, error) {
	ubo, err := dao.UniversalDao.Get(id)
	if err != nil {
		return nil, err
	}
	return NewBlogVoteFromUbo(ubo), nil
}

// GetN implements BlogVoteDao.GetN
func (dao *BaseBlogVoteDaoImpl) GetN(fromOffset, maxNumRows int, filter godal.FilterOpt, sorting *godal.SortingOpt) ([]*BlogVote, error) {
	uboList, err := dao.UniversalDao.GetN(fromOffset, maxNumRows, filter, sorting)
	if err != nil {
		return nil, err
	}
	result := make([]*BlogVote, 0)
	for _, ubo := range uboList {
		app := NewBlogVoteFromUbo(ubo)
		result = append(result, app)
	}
	return result, nil
}

// GetAll implements BlogVoteDao.GetAll
func (dao *BaseBlogVoteDaoImpl) GetAll(filter godal.FilterOpt, sorting *godal.SortingOpt) ([]*BlogVote, error) {
	return dao.GetN(0, 0, filter, sorting)
}

// Update implements BlogVoteDao.Update
func (dao *BaseBlogVoteDaoImpl) Update(vote *BlogVote) (bool, error) {
	return dao.UniversalDao.Update(vote.sync().UniversalBo.Clone())
}

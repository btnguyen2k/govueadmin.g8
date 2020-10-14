package blog

import (
	"github.com/btnguyen2k/consu/reddo"
	"github.com/btnguyen2k/godal"
	"github.com/btnguyen2k/godal/sql"
	"github.com/btnguyen2k/prom"

	userv2 "main/src/gvabe/bov2/user"
	"main/src/henge"
)

const TableBlogPost = "gva_blog_post"
const (
	PostCol_OwnerId  = "zownid"
	PostCol_IsPublic = "zispub"
)

// NewBlogPostDaoSql is helper method to create SQL-implementation of BlogPostDao
//
// available since template-v0.2.0
func NewBlogPostDaoSql(sqlc *prom.SqlConnect, tableName string) BlogPostDao {
	dao := &BlogPostDaoSql{}
	dao.UniversalDao = henge.NewUniversalDaoSql(
		sqlc, tableName,
		map[string]string{
			PostCol_OwnerId:  PostField_OwnerId,
			PostCol_IsPublic: PostField_IsPublic,
		})
	return dao
}

// BlogPostDaoSql is SQL-implementation of BlogPostDao
//
// available since template-v0.2.0
type BlogPostDaoSql struct {
	henge.UniversalDao
}

// GdaoCreateFilter implements IGenericDao.GdaoCreateFilter
func (dao *BlogPostDaoSql) GdaoCreateFilter(_ string, gbo godal.IGenericBo) interface{} {
	return map[string]interface{}{henge.ColId: gbo.GboGetAttrUnsafe(henge.FieldId, reddo.TypeString)}
}

// Delete implements BlogPostDao.Delete
func (dao *BlogPostDaoSql) Delete(post *BlogPost) (bool, error) {
	return dao.UniversalDao.Delete(post.UniversalBo.Clone())
}

// Create implements BlogPostDao.Create
func (dao *BlogPostDaoSql) Create(post *BlogPost) (bool, error) {
	return dao.UniversalDao.Create(post.sync().UniversalBo.Clone())
}

// Get implements BlogPostDao.Get
func (dao *BlogPostDaoSql) Get(id string) (*BlogPost, error) {
	ubo, err := dao.UniversalDao.Get(id)
	if err != nil {
		return nil, err
	}
	return NewBlogPostFromUbo(ubo), nil
}

// GetUserPostsN implements BlogPostDao.GetUserPostsN
func (dao *BlogPostDaoSql) GetUserPostsN(user *userv2.User, fromOffset, maxNumRows int) ([]*BlogPost, error) {
	filter := &sql.FilterFieldValue{Field: PostCol_OwnerId, Operation: "=", Value: user.GetId()}
	uboList, err := dao.UniversalDao.GetN(fromOffset, maxNumRows, filter, nil)
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
func (dao *BlogPostDaoSql) GetUserPostsAll(user *userv2.User) ([]*BlogPost, error) {
	return dao.GetUserPostsN(user, 0, 0)
}

// Update implements BlogPostDao.Update
func (dao *BlogPostDaoSql) Update(post *BlogPost) (bool, error) {
	return dao.UniversalDao.Update(post.sync().UniversalBo.Clone())
}

package blog

import (
	"github.com/btnguyen2k/consu/reddo"
	"github.com/btnguyen2k/godal"
	"github.com/btnguyen2k/prom"

	"main/src/henge"
)

const TableBlogComment = "gva_blog_comment"
const (
	CommentCol_OwnerId  = "zownid"
	CommentCol_PostId   = "zpid"
	CommentCol_ParentId = "zprid"
)

// NewBlogCommentDaoSql is helper method to create SQL-implementation of BlogCommentDao
//
// available since template-v0.2.0
func NewBlogCommentDaoSql(sqlc *prom.SqlConnect, tableName string) BlogCommentDao {
	dao := &BlogCommentDaoSql{}
	dao.UniversalDao = henge.NewUniversalDaoSql(
		sqlc, tableName,
		map[string]string{
			CommentCol_OwnerId:  CommentField_OwnerId,
			CommentCol_PostId:   CommentField_PostId,
			CommentCol_ParentId: CommentField_ParentId,
		})
	return dao
}

// BlogCommentDaoSql is SQL-implementation of BlogCommentDao
//
// available since template-v0.2.0
type BlogCommentDaoSql struct {
	henge.UniversalDao
}

// GdaoCreateFilter implements IGenericDao.GdaoCreateFilter
func (dao *BlogCommentDaoSql) GdaoCreateFilter(_ string, gbo godal.IGenericBo) interface{} {
	return map[string]interface{}{henge.ColId: gbo.GboGetAttrUnsafe(henge.FieldId, reddo.TypeString)}
}

// Delete implements BlogCommentDao.Delete
func (dao *BlogCommentDaoSql) Delete(comment *BlogComment) (bool, error) {
	return dao.UniversalDao.Delete(comment.UniversalBo.Clone())
}

// Create implements BlogCommentDao.Create
func (dao *BlogCommentDaoSql) Create(comment *BlogComment) (bool, error) {
	return dao.UniversalDao.Create(comment.sync().UniversalBo.Clone())
}

// Get implements BlogCommentDao.Get
func (dao *BlogCommentDaoSql) Get(id string) (*BlogComment, error) {
	ubo, err := dao.UniversalDao.Get(id)
	if err != nil {
		return nil, err
	}
	return NewBlogCommentFromUbo(ubo), nil
}

// GetN implements BlogCommentDao.GetN
func (dao *BlogCommentDaoSql) GetN(fromOffset, maxNumRows int) ([]*BlogComment, error) {
	uboList, err := dao.UniversalDao.GetN(fromOffset, maxNumRows)
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
func (dao *BlogCommentDaoSql) GetAll() ([]*BlogComment, error) {
	return dao.GetN(0, 0)
}

// Update implements BlogCommentDao.Update
func (dao *BlogCommentDaoSql) Update(comment *BlogComment) (bool, error) {
	return dao.UniversalDao.Update(comment.sync().UniversalBo.Clone())
}

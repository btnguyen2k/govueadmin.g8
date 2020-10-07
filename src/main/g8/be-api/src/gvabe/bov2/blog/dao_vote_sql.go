package user

import (
	"github.com/btnguyen2k/consu/reddo"
	"github.com/btnguyen2k/godal"
	"github.com/btnguyen2k/prom"

	"main/src/henge"
)

const TableBlogVote = "gva_blog_vote"
const (
	VoteCol_OwnerId  = "zownid"
	VoteCol_TargetId = "ztid"
	VoteCol_Value    = "zval"
)

// NewBlogVoteDaoSql is helper method to create SQL-implementation of BlogVoteDao
//
// available since template-v0.2.0
func NewBlogVoteDaoSql(sqlc *prom.SqlConnect, tableName string) BlogVoteDao {
	dao := &BlogVoteDaoSql{}
	dao.UniversalDao = henge.NewUniversalDaoSql(
		sqlc, tableName,
		map[string]string{
			VoteCol_OwnerId:  VoteField_OwnerId,
			VoteCol_TargetId: VoteField_TargetId,
			VoteCol_Value:    VoteField_Value,
		})
	return dao
}

// BlogVoteDaoSql is SQL-implementation of BlogVoteDao
//
// available since template-v0.2.0
type BlogVoteDaoSql struct {
	henge.UniversalDao
}

// GdaoCreateFilter implements IGenericDao.GdaoCreateFilter
func (dao *BlogVoteDaoSql) GdaoCreateFilter(_ string, gbo godal.IGenericBo) interface{} {
	return map[string]interface{}{henge.ColId: gbo.GboGetAttrUnsafe(henge.FieldId, reddo.TypeString)}
}

// Delete implements BlogVoteDao.Delete
func (dao *BlogVoteDaoSql) Delete(vote *BlogVote) (bool, error) {
	return dao.UniversalDao.Delete(vote.UniversalBo.Clone())
}

// Create implements BlogVoteDao.Create
func (dao *BlogVoteDaoSql) Create(vote *BlogVote) (bool, error) {
	return dao.UniversalDao.Create(vote.sync().UniversalBo.Clone())
}

// Get implements BlogVoteDao.Get
func (dao *BlogVoteDaoSql) Get(id string) (*BlogVote, error) {
	ubo, err := dao.UniversalDao.Get(id)
	if err != nil {
		return nil, err
	}
	return NewBlogVoteFromUbo(ubo), nil
}

// GetN implements BlogVoteDao.GetN
func (dao *BlogVoteDaoSql) GetN(fromOffset, maxNumRows int) ([]*BlogVote, error) {
	uboList, err := dao.UniversalDao.GetN(fromOffset, maxNumRows)
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
func (dao *BlogVoteDaoSql) GetAll() ([]*BlogVote, error) {
	return dao.GetN(0, 0)
}

// Update implements BlogVoteDao.Update
func (dao *BlogVoteDaoSql) Update(vote *BlogVote) (bool, error) {
	return dao.UniversalDao.Update(vote.sync().UniversalBo.Clone())
}

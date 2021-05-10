package blog

import (
	"github.com/btnguyen2k/henge"
	"github.com/btnguyen2k/prom"
)

// NewBlogCommentDaoSql is helper method to create SQL-implementation of BlogCommentDao.
//
// Available since template-v0.2.0
func NewBlogCommentDaoSql(sqlc *prom.SqlConnect, tableName string, txModeOnWrite bool) BlogCommentDao {
	dao := &BaseBlogCommentDaoImpl{}
	dao.UniversalDao = henge.NewUniversalDaoSql(
		sqlc, tableName, txModeOnWrite,
		map[string]string{
			CommentColOwnerId:  CommentFieldOwnerId,
			CommentColPostId:   CommentFieldPostId,
			CommentColParentId: CommentFieldParentId,
		})
	return dao
}

// NewBlogPostDaoSql is helper method to create SQL-implementation of BlogPostDao.
//
// Available since template-v0.2.0
func NewBlogPostDaoSql(sqlc *prom.SqlConnect, tableName string, txModeOnWrite bool) BlogPostDao {
	dao := &BaseBlogPostDaoImpl{}
	dao.UniversalDao = henge.NewUniversalDaoSql(
		sqlc, tableName, txModeOnWrite,
		map[string]string{
			PostColOwnerId:  PostFieldOwnerId,
			PostColIsPublic: PostFieldIsPublic,
		})
	return dao
}

// NewBlogVoteDaoSql is helper method to create SQL-implementation of BlogVoteDao.
//
// Available since template-v0.2.0
func NewBlogVoteDaoSql(sqlc *prom.SqlConnect, tableName string, txModeOnWrite bool) BlogVoteDao {
	dao := &BaseBlogVoteDaoImpl{}
	dao.UniversalDao = henge.NewUniversalDaoSql(
		sqlc, tableName, txModeOnWrite,
		map[string]string{
			VoteColOwnerId:  VoteFieldOwnerId,
			VoteColTargetId: VoteFieldTargetId,
			VoteColValue:    VoteFieldValue,
		})
	return dao
}

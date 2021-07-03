package blog

import (
	"github.com/btnguyen2k/henge"
	"github.com/btnguyen2k/prom"
)

// NewBlogCommentDaoDynamodb is helper method to create AWS DynamoDB-implementation of BlogCommentDao.
//
// Available since template-v0.3.0
func NewBlogCommentDaoDynamodb(adc *prom.AwsDynamodbConnect, tableName string) BlogCommentDao {
	dao := &BaseBlogCommentDaoImpl{}
	spec := &henge.DynamodbDaoSpec{}
	dao.UniversalDao = henge.NewUniversalDaoDynamodb(adc, tableName, spec)
	return dao
}

// NewBlogPostDaoDynamodb is helper method to create AWS DynamoDB-implementation of BlogPostDao.
//
// Available since template-v0.3.0
func NewBlogPostDaoDynamodb(adc *prom.AwsDynamodbConnect, tableName string) BlogPostDao {
	dao := &BaseBlogPostDaoImpl{}
	spec := &henge.DynamodbDaoSpec{}
	dao.UniversalDao = henge.NewUniversalDaoDynamodb(adc, tableName, spec)
	return dao
}

// NewBlogVoteDaoDynamodb is helper method to create AWS DynamoDB-implementation of BlogVoteDao.
//
// Available since template-v0.3.0
func NewBlogVoteDaoDynamodb(adc *prom.AwsDynamodbConnect, tableName string) BlogVoteDao {
	dao := &BaseBlogVoteDaoImpl{}
	spec := &henge.DynamodbDaoSpec{}
	dao.UniversalDao = henge.NewUniversalDaoDynamodb(adc, tableName, spec)
	return dao
}

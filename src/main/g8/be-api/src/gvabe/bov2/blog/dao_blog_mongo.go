package blog

import (
	"github.com/btnguyen2k/henge"
	"github.com/btnguyen2k/prom"
)

// NewBlogCommentDaoMongo is helper method to create MongoDB-implementation of BlogCommentDao.
//
// Available since template-v0.3.0
func NewBlogCommentDaoMongo(mc *prom.MongoConnect, collectionName string, txModeOnWrite bool) BlogCommentDao {
	dao := &BaseBlogCommentDaoImpl{}
	dao.UniversalDao = henge.NewUniversalDaoMongo(mc, collectionName, txModeOnWrite)
	return dao
}

// NewBlogPostDaoMongo is helper method to create MongoDB-implementation of BlogPostDao.
//
// Available since template-v0.3.0
func NewBlogPostDaoMongo(mc *prom.MongoConnect, collectionName string, txModeOnWrite bool) BlogPostDao {
	dao := &BaseBlogPostDaoImpl{}
	dao.UniversalDao = henge.NewUniversalDaoMongo(mc, collectionName, txModeOnWrite)
	return dao
}

// NewBlogVoteDaoMongo is helper method to create MongoDB-implementation of BlogVoteDao.
//
// Available since template-v0.3.0
func NewBlogVoteDaoMongo(mc *prom.MongoConnect, collectionName string, txModeOnWrite bool) BlogVoteDao {
	dao := &BaseBlogVoteDaoImpl{}
	dao.UniversalDao = henge.NewUniversalDaoMongo(mc, collectionName, txModeOnWrite)
	return dao
}

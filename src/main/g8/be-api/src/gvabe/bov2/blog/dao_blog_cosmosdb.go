package blog

import (
	"github.com/btnguyen2k/henge"
	promsql "github.com/btnguyen2k/prom/sql"
)

// NewBlogCommentDaoCosmosdb is helper method to create Azure Cosmos DB-implementation of BlogCommentDao.
//
// Note: txModeOnWrite is not currently used!
//
// Available since template-v0.3.0
func NewBlogCommentDaoCosmosdb(sqlc *promsql.SqlConnect, tableName string, txModeOnWrite bool) BlogCommentDao {
	dao := &BaseBlogCommentDaoImpl{}
	spec := &henge.CosmosdbDaoSpec{
		PkName:        henge.CosmosdbColId,
		TxModeOnWrite: txModeOnWrite,
	}
	dao.UniversalDao = henge.NewUniversalDaoCosmosdbSql(sqlc, tableName, spec)
	return dao
}

// NewBlogPostDaoCosmosdb is helper method to create Azure Cosmos DB-implementation of BlogPostDao.
//
// Note: txModeOnWrite is not currently used!
//
// Available since template-v0.3.0
func NewBlogPostDaoCosmosdb(sqlc *promsql.SqlConnect, tableName string, txModeOnWrite bool) BlogPostDao {
	dao := &BaseBlogPostDaoImpl{}
	spec := &henge.CosmosdbDaoSpec{
		PkName:        henge.CosmosdbColId,
		TxModeOnWrite: txModeOnWrite,
	}
	dao.UniversalDao = henge.NewUniversalDaoCosmosdbSql(sqlc, tableName, spec)
	return dao
}

// NewBlogVoteDaoCosmosdb is helper method to create Azure Cosmos DB-implementation of BlogVoteDao.
//
// Note: txModeOnWrite is not currently used!
//
// Available since template-v0.3.0
func NewBlogVoteDaoCosmosdb(sqlc *promsql.SqlConnect, tableName string, txModeOnWrite bool) BlogVoteDao {
	dao := &BaseBlogVoteDaoImpl{}
	spec := &henge.CosmosdbDaoSpec{
		PkName:        henge.CosmosdbColId,
		TxModeOnWrite: txModeOnWrite,
	}
	dao.UniversalDao = henge.NewUniversalDaoCosmosdbSql(sqlc, tableName, spec)
	return dao
}

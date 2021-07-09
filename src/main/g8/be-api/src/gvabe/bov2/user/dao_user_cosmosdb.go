package user

import (
	"github.com/btnguyen2k/henge"
	"github.com/btnguyen2k/prom"
)

// NewUserDaoCosmosdb is helper method to create Azure Cosmos DB-implementation of UserDao
//
// Note: txModeOnWrite is not currently used!
//
// Available since template-v0.3.0
func NewUserDaoCosmosdb(sqlc *prom.SqlConnect, tableName string, txModeOnWrite bool) UserDao {
	dao := &BaseUserDaoImpl{}
	spec := &henge.CosmosdbDaoSpec{
		PkName:        henge.CosmosdbColId,
		TxModeOnWrite: txModeOnWrite,
	}
	dao.UniversalDao = henge.NewUniversalDaoCosmosdbSql(sqlc, tableName, spec)
	return dao
}

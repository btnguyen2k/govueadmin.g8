package user

import (
	"github.com/btnguyen2k/henge"
	"github.com/btnguyen2k/prom"
)

// NewUserDaoSql is helper method to create SQL-implementation of UserDao
//
// Available since template-v0.2.0
func NewUserDaoSql(sqlc *prom.SqlConnect, tableName string, txModeOnWrite bool) UserDao {
	dao := &BaseUserDaoImpl{}
	dao.UniversalDao = henge.NewUniversalDaoSql(
		sqlc, tableName, txModeOnWrite,
		map[string]string{UserColMaskUid: UserFieldMaskId})
	return dao
}

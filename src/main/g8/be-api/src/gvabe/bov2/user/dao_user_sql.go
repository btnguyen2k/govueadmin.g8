package user

import (
	"github.com/btnguyen2k/henge"
	promsql "github.com/btnguyen2k/prom/sql"
)

// NewUserDaoSql is helper method to create SQL-implementation of UserDao
//
// Available since template-v0.2.0
func NewUserDaoSql(sqlc *promsql.SqlConnect, tableName string, txModeOnWrite bool) UserDao {
	dao := &BaseUserDaoImpl{}
	dao.UniversalDao = henge.NewUniversalDaoSql(
		sqlc, tableName, txModeOnWrite,
		map[string]string{UserColMaskUid: UserFieldMaskId})
	return dao
}
